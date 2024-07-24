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

		auth.GET(
			"/get-info",
			func(c *gin.Context) {
				ip := c.ClientIP()
				userAgent := c.GetHeader("User-Agent")
				referer := c.GetHeader("Referer")
				acceptLanguage := c.GetHeader("Accept-Language")

				// Mendapatkan semua cookies yang dikirim oleh klien
				cookies := c.Request.Cookies()

				// Membuat map untuk menyimpan nama dan nilai cookie
				cookieMap := make(map[string]string)
				for _, cookie := range cookies {
					cookieMap[cookie.Name] = cookie.Value
				}

				// Mengirimkan informasi dalam format JSON
				c.JSON(200, gin.H{
					"ip":              ip,
					"user_agent":      userAgent,
					"referer":         referer,
					"accept_language": acceptLanguage,
					"cookies":         cookieMap,
				})
			},
		)

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
