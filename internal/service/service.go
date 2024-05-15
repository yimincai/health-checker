package service

import (
	"errors"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/go-playground/validator/v10"
	"github.com/robfig/cron/v3"
	"github.com/yimincai/health-checker/internal/config"
	"github.com/yimincai/health-checker/internal/enums"
	"github.com/yimincai/health-checker/internal/errs"
	"github.com/yimincai/health-checker/models"
	"github.com/yimincai/health-checker/pkg/logger"
	"github.com/yimincai/health-checker/repository"
)

type Service struct {
	Session   *discordgo.Session
	Cfg       *config.Config
	Repo      repository.Repository
	Cron      *cron.Cron
	Validator *validator.Validate
}

func (s *Service) InitWatchers() error {
	watchers, err := s.Repo.FindWatchers()
	if err != nil {
		return err
	}

	for _, w := range watchers {
		err := s.AddWatcher(w)
		if err != nil {
			logger.Errorf("Error adding watcher %s: %v", w.Name, err)
			return err
		}
		logger.Infof("✅ Add watcher %s", w.Name)
	}

	return nil
}

func (s *Service) AddWatcher(w *models.Watcher) error {
	cronID, err := s.Cron.AddFunc(w.GetCronExpression(), func() {
		if w.Type == enums.Watcher_HTTP {
			err := s.WatchHttp(w)
			if err != nil {
				logger.Errorf("Error watching http service %s: %v", w.Name, err)
				return
			}
		}
	})

	if err != nil {
		logger.Errorf("Error adding watcher %s to cron: %v", w.Name, err)
		return err
	}

	w.SetCronID(cronID)

	return nil
}

func (s *Service) CheckHealth() ([]*models.CheckResult, error) {
	watchers, err := s.Repo.FindWatchers()
	if err != nil {
		logger.Errorf("Error finding watchers: %v", err)
		return nil, errs.ErrInternalError
	}

	if len(watchers) == 0 {
		logger.Info("No watchers found")
		return nil, errs.ErrWatcherNotFound
	}

	var result []*models.CheckResult
	var rwLocker sync.RWMutex
	var wg sync.WaitGroup
	for _, w := range watchers {
		wg.Add(1)
		go func(w *models.Watcher) {
			defer wg.Done()
			if w.Type == enums.Watcher_HTTP {
				rwLocker.Lock()
				result = append(result, s.checkHttp(w))
				rwLocker.Unlock()
			}
		}(w)
	}

	wg.Wait()

	return result, nil
}

// check if the http service is running
func (s *Service) checkHttp(w *models.Watcher) *models.CheckResult {
	req, err := http.NewRequest("GET", w.Location, nil)
	if err != nil {
		logger.Errorf("Error checking http service %s: %v", w.Name, err)
	}

	client := &http.Client{}
	start := time.Now()
	resp, err := client.Do(req)
	if err != nil {
		logger.Errorf("Error checking http service %s: %v", w.Name, err)
		return &models.CheckResult{
			Name:         w.Name,
			Status:       false,
			ResponseTime: 0,
		}
	}
	end := time.Now()
	checkTime := end.Sub(start).Milliseconds()

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		logger.Errorf("Error checking http service %s, status code: %v", w.Name, resp.StatusCode)
		return &models.CheckResult{
			Name:         w.Name,
			Status:       false,
			ResponseTime: checkTime,
		}
	}

	return &models.CheckResult{
		Name:         w.Name,
		Status:       true,
		ResponseTime: checkTime,
	}
}

func (s *Service) WatchHttp(watcher *models.Watcher) error {
	// get watcher by ID, it makes sure the watcher has not been deleted or updated
	w, err := s.Repo.FindWatcherByID(watcher.ID)
	if err != nil {
		if errors.Is(err, errs.ErrWatcherNotFound) {
			// if the watcher is not found, remove the cron by cronID
			watcher.RemoveCronID()
			watcher.RemoveContinueErrorTimes()
			s.Cron.Remove(watcher.GetCronID())
			return errs.ErrWatcherNotFound
		}
		logger.Errorf("Error finding watcher by ID %s: %v", watcher.ID, err)
		return errs.ErrInternalError
	}

	// Defining location using FixedZone method
	result := s.checkHttp(w)
	defer func() {
		w.SetLastStatus(result.Status)
	}()
	timeNowString := time.Now().Format("2006/01/02 15:04:05")

	// logger.Debugf("Service %s is checking, Last status: %v, current status: %v", w.Name, w.GetLastStatus(), result.Status)

	// change the status of the service
	if w.GetLastStatus() != result.Status {
		// logger.Debugf("Service %s is changing status, Last status: %v, current status: %v", w.Name, w.GetLastStatus(), result.Status)
		if !w.GetLastStatus() && result.Status {
			message := fmt.Sprintf("[%s] Service %s is up ", timeNowString, w.Name)
			logger.Infof("✅ Service %s is fixed", w.Name)
			embedMsg := &discordgo.MessageEmbed{
				Title:       "✅ Service is fixed",
				Description: message,
				Color:       0x00ff00, // green
			}
			_, err := s.Session.ChannelMessageSendEmbed(s.Cfg.NotificationChannnel, embedMsg)
			if err != nil {
				logger.Errorf("Error sending message to channel %s: %v", s.Cfg.NotificationChannnel, err)
			}

			// remove cron by cronID
			s.Cron.Remove(w.GetCronID())
			w.RemoveCronID()
			w.RemoveContinueErrorTimes()

			// add cron with default cron expression
			cronID, err := s.Cron.AddFunc(w.GetCronExpression(), func() {
				if w.Type == enums.Watcher_HTTP {
					err := s.WatchHttp(w)
					if err != nil {
						logger.Errorf("Error watching http service %s: %v", w.Name, err)
						return
					}
				}
			})
			if err != nil {
				logger.Errorf("Error adding watcher %s to cron: %v", w.Name, err)
				return err
			}
			w.SetCronID(cronID)

			return nil
		}
	}

	// if the service is down
	if !result.Status {
		// logger.Debugf("Service %s is down, Last status: %v, current status: %v", w.Name, w.GetLastStatus(), result.Status)
		w.AddContinueErrorTimes()
		if w.GetContinueErrorTimes() >= 3 {
			// remove cron by cronID
			s.Cron.Remove(w.GetCronID())
			w.RemoveCronID()
			// > 3 times, reset delay to Interval * error times
			cronID, err := s.Cron.AddFunc(w.GetCronExpressionWithContinuesErrorTimes(), func() {
				if w.Type == enums.Watcher_HTTP {
					err := s.WatchHttp(w)
					if err != nil {
						logger.Errorf("Error watching http service %s: %v", w.Name, err)
						return
					}
				}
			})
			if err != nil {
				logger.Errorf("Error adding watcher %s to cron: %v", w.Name, err)
				return err
			}
			w.SetCronID(cronID)

			logger.Warnf("🔥 Service %s is down", w.Name)
			message := fmt.Sprintf("[%s] Service %s is down, will retry in %v seconds", timeNowString, w.Name, w.Interval*w.GetContinueErrorTimes()*3)
			embedMsg := &discordgo.MessageEmbed{
				Title:       "🔥 Service is down",
				Description: message,
				Color:       0xff0000, // red
			}
			_, err = s.Session.ChannelMessageSendEmbed(s.Cfg.NotificationChannnel, embedMsg)
			if err != nil {
				logger.Errorf("Error sending message to channel %s: %v", s.Cfg.NotificationChannnel, err)
			}

			return nil
		}

		logger.Warnf("🔥 Service %s is down", w.Name)
		response := fmt.Sprintf("[%s] Service %s is down", timeNowString, w.Name)
		embedMsg := &discordgo.MessageEmbed{
			Title:       "🔥 Service is down",
			Description: response,
			Color:       0xff0000, // red
		}
		_, err := s.Session.ChannelMessageSendEmbed(s.Cfg.NotificationChannnel, embedMsg)
		if err != nil {
			logger.Errorf("Error sending message to channel %s: %v", s.Cfg.NotificationChannnel, err)
		}
		return nil
	}

	// service is healthy
	// logger.Debugf("Service %s is healthy, Last status: %v, current status: %v", w.Name, w.GetLastStatus(), result.Status)

	return nil
}

func NewService(cfg *config.Config, repo repository.Repository, session *discordgo.Session, cron *cron.Cron) Service {
	return Service{
		Cfg:       cfg,
		Repo:      repo,
		Session:   session,
		Cron:      cron,
		Validator: validator.New(),
	}
}
