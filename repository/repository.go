package repository

import (
	"github.com/yimincai/health-checker/models"
	"gorm.io/gorm"
)

type Repository interface {
	// Watcher

	CreateWatcher(watcher *models.Watcher) (*models.Watcher, error)
	FindWatchers() ([]*models.Watcher, error)
	FindWatcherByID(watcherID string) (*models.Watcher, error)
	FindWatcherByName(watcherName string) (*models.Watcher, error)
	UpdateWatcher(watcher *models.Watcher) (*models.Watcher, error)
	DeleteWatcher(watcherID string) error
}

type Repo struct {
	Db *gorm.DB
}

func New(db *gorm.DB) Repository {
	return &Repo{Db: db}
}
