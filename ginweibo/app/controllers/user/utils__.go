package user

import (
	"ginweibo/app/helpers"
	userModel "ginweibo/app/models/user"
	"ginweibo/routes/named"

	"github.com/gin-gonic/gin"
)

func sendConfirmEmail(u *userModel.User) error {
	subject := "感谢注册 Weibo 应用！请确认你的邮箱。"
	tpl := "mail/confirm.html"
	confirmURL := named.G("signup.confirm", "token", u.ActivationToken)
	return helpers.SendMail([]string{u.Email}, subject, tpl, gin.H{"confirmURL": confirmURL})
}
