package models

import (
	"fmt"
	"sync"

	"github.com/robfig/cron/v3"
)

var CronScheduledMap map[string]cron.EntryID
var scmLocker = &sync.RWMutex{}

var WatcherLastStatusMap map[string]bool
var wlsmLocker = &sync.RWMutex{}

var WatcherContinueErrorTimesMap map[string]int
var wcemLocker = &sync.RWMutex{}

func init() {
	CronScheduledMap = make(map[string]cron.EntryID)
	WatcherLastStatusMap = make(map[string]bool)
	WatcherContinueErrorTimesMap = make(map[string]int)
}

func (w *Watcher) GetCronExpression() string {
	return fmt.Sprintf("@every %ds", w.Interval)
}

func (w *Watcher) GetCronExpressionWithContinuesErrorTimes() string {
	wcemLocker.RLock()
	defer wcemLocker.RUnlock()
	return fmt.Sprintf("@every %ds", w.Interval*WatcherContinueErrorTimesMap[w.ID]*3)
}

func (w *Watcher) GetCronID() cron.EntryID {
	scmLocker.RLock()
	defer scmLocker.RUnlock()
	return CronScheduledMap[w.ID]
}

func (w *Watcher) SetCronID(id cron.EntryID) {
	scmLocker.Lock()
	defer scmLocker.Unlock()
	CronScheduledMap[w.ID] = id
}

func (w *Watcher) RemoveCronID() {
	scmLocker.Lock()
	defer scmLocker.Unlock()
	delete(CronScheduledMap, w.ID)
}

func (w *Watcher) SetLastStatus(b bool) {
	wlsmLocker.Lock()
	defer wlsmLocker.Unlock()
	WatcherLastStatusMap[w.ID] = b
}

func (w *Watcher) GetLastStatus() bool {
	wlsmLocker.RLock()
	defer wlsmLocker.RUnlock()
	if _, ok := WatcherLastStatusMap[w.ID]; !ok {
		return true
	}

	return WatcherLastStatusMap[w.ID]
}

func (w *Watcher) RemoveLastStatus() {
	wlsmLocker.Lock()
	defer wlsmLocker.Unlock()
	delete(WatcherLastStatusMap, w.ID)
}

func (w *Watcher) AddContinueErrorTimes() {
	wcemLocker.Lock()
	defer wcemLocker.Unlock()
	if _, ok := WatcherContinueErrorTimesMap[w.ID]; !ok {
		WatcherContinueErrorTimesMap[w.ID] = 0
	}

	WatcherContinueErrorTimesMap[w.ID]++
}

func (w *Watcher) GetContinueErrorTimes() int {
	wcemLocker.RLock()
	defer wcemLocker.RUnlock()
	if _, ok := WatcherContinueErrorTimesMap[w.ID]; !ok {
		return 0
	}

	return WatcherContinueErrorTimesMap[w.ID]
}

func (w *Watcher) RemoveContinueErrorTimes() {
	wcemLocker.Lock()
	defer wcemLocker.Unlock()
	delete(WatcherContinueErrorTimesMap, w.ID)
}
