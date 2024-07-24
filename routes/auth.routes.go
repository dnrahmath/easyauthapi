package routes

import (
	"crypto/md5"
	"easyauthapi/controllers"
	"easyauthapi/middlewares"
	"easyauthapi/middlewares/validators"
	"encoding/hex"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func AuthRoute(router *gin.RouterGroup, handlers ...gin.HandlerFunc) {
	auth := router.Group("/auth", handlers...)
	{
		//==================================
		call := controllers.AuthController{}
		ctrl := call.NewController()
		//==================================

		auth.GET("/set-cookie", func(c *gin.Context) {
			cookie := &http.Cookie{
				Name:     "unique-id",
				Value:    "1234567890", // Harus digenerate secara unik untuk setiap pengguna
				Expires:  time.Now().Add(365 * 24 * time.Hour),
				HttpOnly: true,
			}
			http.SetCookie(c.Writer, cookie)
			c.JSON(http.StatusOK, gin.H{"status": "cookie set"})
		})

		auth.GET(
			"/get-info",
			func(c *gin.Context) {
				var ip string

				// Mendapatkan IP dari header "X-Forwarded-For"
				XForwardedFor := c.GetHeader("X-Forwarded-For")
				if XForwardedFor != "" {
					ip = XForwardedFor
					fmt.Println("X-Forwarded-For : " + ip)
				}

				// Jika "X-Forwarded-For" tidak tersedia, gunakan "X-Real-IP"
				XRealIP := c.GetHeader("X-Real-IP")
				if ip == "" && XRealIP != "" {
					ip = XRealIP
					fmt.Println("X-Real-IP : " + ip)
				}

				// Jika header tidak tersedia, gunakan metode bawaan ClientIP [Carrier-Grade NAT (CGNAT)]
				ClientIP := c.ClientIP()
				if ip == "" && ClientIP != "" {
					ip = ClientIP
					fmt.Println("Client IP : " + ip)
				}

				userAgent := c.GetHeader("User-Agent")
				referer := c.GetHeader("Referer")
				acceptLanguage := c.GetHeader("Accept-Language")

				// Generate fingerprint
				fingerprintSource := ip + userAgent + referer + acceptLanguage
				hash := md5.Sum([]byte(fingerprintSource))
				fingerprint := hex.EncodeToString(hash[:])

				cookies := c.Request.Cookies()
				cookieMap := make(map[string]string)
				for _, cookie := range cookies {
					cookieMap[cookie.Name] = cookie.Value
				}

				headers := c.Request.Header

				c.JSON(http.StatusOK, gin.H{
					"XForwardedFor":   XForwardedFor,
					"XRealIP":         XRealIP,
					"ClientIP":        ClientIP,
					"SelectedIP":      ip, // IP yang dipilih berdasarkan prioritas
					"user_agent":      userAgent,
					"referer":         referer,
					"accept_language": acceptLanguage,
					"fingerprint":     fingerprint,
					"cookies":         cookieMap,
					"headers":         headers,
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
