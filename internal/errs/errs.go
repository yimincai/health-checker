package errs

import (
	"errors"

	"github.com/yimincai/health-checker/pkg/logger"
)

var (
	ErrInternalError    = errors.New("internal error, please contact admin 👨‍💻")
	ErrForbidden        = errors.New("forbidden, please contact admin 👨‍💻")
	ErrSendingMessage   = errors.New("error while sending message, 👨‍💻")
	ErrUserNotFound     = errors.New("user not found 👨‍💻")
	ErrUserNotEnabled   = errors.New("user not enabled 🤕, please contact admin 👨‍💻")
	ErrInvalidDate      = errors.New("invalid date, please check the date format 📅")
	ErrWatcherNotFound  = errors.New("service not found, please add it first 👨‍💻")
	ErrDuplicateWatcher = errors.New("watcher already exists, please use another name or location 👨‍💻")
)

func LogError(err error) {
	if err != nil {
		logger.Error(err)
	}
}
