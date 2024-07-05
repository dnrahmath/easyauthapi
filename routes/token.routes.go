package routes

import (
	"easyauthapi/controllers"
	"easyauthapi/middlewares"
	"easyauthapi/middlewares/validators"

	"github.com/gin-gonic/gin"
)

func TokenRoute(router *gin.RouterGroup, handlers ...gin.HandlerFunc) {
	Token := router.Group("/token", handlers...)
	{
		//==================================
		call := controllers.TokenController{}
		ctrl := call.NewController()
		//==================================

		//membuat dan mengirim token verify , berasal dari Access-Token uuid duser
		Token.POST(
			"/sendcode",
			middlewares.JWTMiddleware(),
			validators.IsUser(),
			func(c *gin.Context) {
				typeparam := c.Query("type")
				ctrl.SendCode(c, typeparam)
			},
		)

		// mengembalikan kode verify
		Token.GET(
			"/:token",
			func(c *gin.Context) {
				tokenparam := c.Param("token")
				ctrl.UpdateOneVerify(c, tokenparam)
			},
		)

		//==================================
	}
}
