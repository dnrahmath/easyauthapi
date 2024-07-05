package middlewares

import (
	"easyauthapi/models/datastore"
	res "easyauthapi/models/response"
	"easyauthapi/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

//==================================================================================

func JWTMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Access-Token")

		var service services.TokenService
		tokenService := service.NewService()

		tokenModel, err := tokenService.VerifyToken(token, datastore.TokenTypeAccess)
		if err != nil {
			res.SendErrorResponse(c, http.StatusBadRequest, "Failed to verify access token, "+err.Error())
			return
		}

		c.Set("UuIdUser", tokenModel.UuIdUser)

		c.Next()
	}
}

//==================================================================================
