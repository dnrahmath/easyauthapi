package validators

import (
	res "easyauthapi/models/response"
	"net/http"

	"github.com/gin-gonic/gin"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
)

func PathIdValidator() gin.HandlerFunc {
	return func(c *gin.Context) {

		id := c.Param("id")
		err := validation.Validate(id, is.MongoID)
		if err != nil {
			res.SendErrorResponse(c, http.StatusBadRequest, "invalid id, "+id)
			return
		}

		c.Next()
	}
}
