package status

import (
	"ginweibo/controllers"
	userModel "ginweibo/models/user"

	"github.com/gin-gonic/gin"
)

func backTo(c *gin.Context, currentUser *userModel.User) {
	back := c.DefaultPostForm("back", "")
	if back != "" {
		controllers.Redirect(c, back, true)
		return
	}
	controllers.RedirectRouter(c, "users.show", currentUser.ID)
}
