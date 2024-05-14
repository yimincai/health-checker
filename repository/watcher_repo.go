package repository

import (
	"github.com/yimincai/health-checker/internal/database"
	"github.com/yimincai/health-checker/internal/errs"
	"github.com/yimincai/health-checker/models"
)

func (r *Repo) CreateWatcher(watcher *models.Watcher) (*models.Watcher, error) {
	tx := r.Db.Begin()

	result := tx.Create(watcher)
	if result.Error != nil {
		tx.Rollback()
		return nil, result.Error
	}

	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	return watcher, nil
}

func (r *Repo) DeleteWatcher(watcherID string) error {
	tx := r.Db.Begin()
	result := tx.Delete(&models.Watcher{}, "id = ?", watcherID)
	if result.Error != nil {
		tx.Rollback()
		return result.Error
	}

	if result.RowsAffected == 0 {
		tx.Rollback()
		return errs.ErrWatcherNotFound
	}

	if err := tx.Commit().Error; err != nil {
		return err
	}

	return nil
}

func (r *Repo) FindWatcherByID(watcherID string) (*models.Watcher, error) {
	var w *models.Watcher

	result := r.Db.Find(&w, "id = ?", watcherID)
	if result.Error != nil {
		return nil, result.Error
	}

	return w, nil
}

func (r *Repo) FindWatcherByName(watcherName string) (*models.Watcher, error) {
	var w *models.Watcher

	result := r.Db.Find(&w, "name = ?", watcherName)
	if result.Error != nil {
		return nil, result.Error
	}

	return w, nil
}

func (r *Repo) FindWatchers() ([]*models.Watcher, error) {
	var ss []*models.Watcher

	result := r.Db.Find(&ss)
	if result.Error != nil {
		return nil, result.Error
	}

	return ss, nil
}

func (r *Repo) UpdateWatcher(watcher *models.Watcher) (*models.Watcher, error) {
	tx := r.Db.Begin()

	var s *models.Watcher
	updateData := database.ParseRDBUpdateData(s)

	result := tx.Model(&models.Watcher{}).Where("id = ?", watcher.ID).Updates(updateData)
	if result.Error != nil {
		tx.Rollback()
		return nil, result.Error
	}

	if result.RowsAffected == 0 {
		tx.Rollback()
		return nil, errs.ErrWatcherNotFound
	}

	var nw *models.Watcher

	result = tx.Find(&nw, "id = ?", watcher.ID)
	if result.Error != nil {
		tx.Rollback()
		return nil, result.Error
	}

	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	return nw, nil
}
