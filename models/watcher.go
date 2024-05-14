package models

import (
	"github.com/yimincai/health-checker/internal/enums"
	"github.com/yimincai/health-checker/pkg/snowflake"
	"gorm.io/gorm"
)

type Watcher struct {
	ID       string            `gorm:"primaryKey"`
	Type     enums.WatcherType `json:"type"` // Type of the watcher
	Name     string            `json:"name" validate:"required"`
	Location string            `json:"location" validate:"required,url"`    // URL or IP
	Interval int               `json:"interval" validate:"required,min=10"` // Interval to check the service
	IsEnable bool              `json:"is_enable"`                           // Enable or disable the watcher
}

type CheckResult struct {
	Name         string
	Status       bool
	ResponseTime int64 // ms
}

// BeforeCreate will set snowflake id rather than numeric id.
func (w *Watcher) BeforeCreate(_ *gorm.DB) (err error) {
	w.ID = snowflake.GetID()
	return nil
}
