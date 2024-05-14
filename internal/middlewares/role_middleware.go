package middlewares

import (
	"github.com/yimincai/health-checker/internal/bot"
	"github.com/yimincai/health-checker/repository"
)

type RequiredAdminPermission struct {
	Repo repository.Repository
}

func (m *RequiredAdminPermission) Exec(ctx *bot.Context, cmd bot.Command) (next bool, err error) {
	// if !cmd.IsAdminRequired() {
	// 	return true, nil
	// }
	//
	// user, err := m.Repo.FindUserByDiscordUserID(ctx.Message.Author.ID)
	// if err != nil {
	// 	return false, errs.ErrUserNotFound
	// }
	//
	// if !user.IsEnable {
	// 	return false, errs.ErrUserNotEnabled
	// }
	//
	// if user.Role != enums.RoleType_Admin {
	// 	return false, errs.ErrForbidden
	// }

	return true, nil
}

func NewRequiredAdminPermission(repo repository.Repository) *RequiredAdminPermission {
	return &RequiredAdminPermission{
		Repo: repo,
	}
}
