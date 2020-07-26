package password

import (
	"ginweibo/controllers"
	"ginweibo/middleware/flash"
	passwordRequest "ginweibo/middleware/requests/password"
	passwordResetModel "ginweibo/models/password_reset"
	"ginweibo/routes/named"
	"ginweibo/utils/mail"

	"github.com/gin-gonic/gin"
)

func sendResetEmail(pwd *passwordResetModel.PasswordReset) error {
	subject := "重置密码！请确认你的邮箱。"
	tpl := "mail/reset_password.html"
	resetPasswordURL := named.G("password.reset", "token", pwd.Token)
	return mail.SendMail([]string{pwd.Email}, subject, tpl, gin.H{"resetPasswordURL": resetPasswordURL})
}

// 显示重置密码的邮箱发送页面
func ShowLinkRequestsForm(c *gin.Context) {
	controllers.Render(c, "password/email.html", gin.H{})
}

// 发送重设密码邮件
func SendResetLinkEmail(c *gin.Context) {
	email := c.PostForm("email")
	passwordForm := &passwordRequest.PasswordEmailForm{
		Email: email,
	}
	pwd, errors := passwordForm.ValidateAndGetToken()
	if len(errors) != 0 || pwd == nil {
		flash.SaveValidateMessage(c, errors)
		controllers.RedirectRouter(c, "password.request")
		return
	}
	if err := sendResetEmail(pwd); err != nil {
		flash.NewDangerFlash(c, "重置密码邮件发送失败: "+err.Error())
		// 删除 token
		passwordResetModel.DeleteByEmail(pwd.Email)
	} else {
		flash.NewSuccessFlash(c, "重置密码已发送到你的邮箱上，请注意查收。")
	}
	controllers.RedirectRouter(c, "password.request")
}

// 更新密码页面
func ShowResetForm(c *gin.Context) {
	token := c.Param("token")
	p, err := passwordResetModel.GetByToken(token)
	if err != nil {
		controllers.Render404(c)
		return
	}
	controllers.Render(c, "password/reset.html", gin.H{
		"token": token,
		"email": p.Email,
	})
}

// 更新密码
func Reset(c *gin.Context) {
	passwordForm := &passwordRequest.PassWordResetForm{
		Token:                c.PostForm("token"),
		Password:             c.PostForm("password"),
		PasswordConfirmation: c.PostForm("password_confirmation"),
	}
	user, errors := passwordForm.ValidateAndUpdateUser()
	if len(errors) != 0 || user == nil {
		flash.SaveValidateMessage(c, errors)
		controllers.RedirectRouter(c, "password.reset", "token", c.PostForm("token"))
		return
	}
	flash.NewSuccessFlash(c, "重置密码成功")
	controllers.RedirectRouter(c, "root")
}
