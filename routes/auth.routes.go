package routes

import (
	"easyauthapi/controllers"
	"easyauthapi/middlewares"
	"easyauthapi/middlewares/validators"

	"github.com/gin-gonic/gin"
)

func AuthRoute(router *gin.RouterGroup, handlers ...gin.HandlerFunc) {
	auth := router.Group("/auth", handlers...)
	{
		//==================================
		call := controllers.AuthController{}
		ctrl := call.NewController()
		//==================================

		auth.POST(
			"/register",
			validators.RegisterValidator(),
			func(c *gin.Context) {
				ctrl.Register(c)
			},
		)

		auth.POST(
			"/login",
			validators.LoginValidator(),
			func(c *gin.Context) {
				ctrl.Login(c)
			},
		)

		auth.GET(
			"/refresh",
			validators.RefreshValidator(),
			func(c *gin.Context) {
				ctrl.Refresh(c)
			},
		)

		auth.GET(
			"/logout",
			middlewares.JWTMiddleware(),
			func(c *gin.Context) {
				ctrl.LogoutUser(c)
			},
		)

		//========================================

		auth.GET(
			"/me",
			validators.IsUser(),
			validators.GetUsersValidator(),
			func(c *gin.Context) {
				ctrl.MeGet(c)
			},
		)

		// auth.PUT(
		// 	"/me",
		// 	validators.IsUser(),
		// 	validators.UpdateUserValidator(),
		// 	func(c *gin.Context) {
		// 		ctrl.MePut(c)
		// 	},
		// )

		// auth.DELETE(
		// 	"/me",
		// 	validators.IsUser(),
		// 	func(c *gin.Context) {
		// 		ctrl.MeDel(c)
		// 	},
		// )

		// //========================================
	}
}
