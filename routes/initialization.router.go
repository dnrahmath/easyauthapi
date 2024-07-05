package routes

import (
	// "easyauthapi/configs"
	"easyauthapi/docs"
	"easyauthapi/middlewares"
	res "easyauthapi/models/response"
	"net/http"

	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func New() *gin.Engine {
	r := gin.New()
	initRoute(r)

	r.Use(gin.LoggerWithWriter(middlewares.LogWriter()))
	r.Use(gin.CustomRecovery(middlewares.AppRecovery()))
	r.Use(middlewares.CORSMiddleware())

	api := r.Group("/api")
	{
		AuthRoute(api)
		TokenRoute(api)
	}

	docs.SwaggerInfo.BasePath = api.BasePath() // adds /v1 to swagger base path

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	return r
}

func initRoute(r *gin.Engine) {
	_ = r.SetTrustedProxies(nil)
	r.RedirectTrailingSlash = false
	r.HandleMethodNotAllowed = true

	r.NoRoute(func(c *gin.Context) {
		res.SendErrorResponse(c, http.StatusNotFound, c.Request.RequestURI+" not found")
	})

	r.NoMethod(func(c *gin.Context) {
		res.SendErrorResponse(c, http.StatusMethodNotAllowed, c.Request.Method+" is not allowed here")
	})
}

func InitGin() {
	gin.DisableConsoleColor()
	// gin.SetMode(configs.UseConfig.Mode)
	// do some other initialization staff
}
