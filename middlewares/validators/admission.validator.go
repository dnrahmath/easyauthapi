package validators

import (
	req "easyauthapi/models/request"
	res "easyauthapi/models/response"
	"net/http"

	"github.com/gin-gonic/gin"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
)

func CreateAdmissionValidator() gin.HandlerFunc {
	return func(c *gin.Context) {
		var requestBody req.AdmissionPostReq
		// _ = c.ShouldBindBodyWith(&requestBody, binding.JSON)

		//mendapatkan data dari form-data
		err := c.ShouldBind(&requestBody)
		if err != nil {
			res.SendErrorResponse(c, http.StatusBadRequest, "failed bind form-data :"+err.Error())
		}

		if err = requestBody.Validate(); err != nil {
			res.SendErrorResponse(c, http.StatusBadRequest, "required data bind form-data , "+err.Error())
			return
		}

		c.Next()
	}
}

func GetAdmissionsValidator() gin.HandlerFunc {
	return func(c *gin.Context) {

		page := c.DefaultQuery("page", "0")
		err := validation.Validate(page, is.Int)
		if err != nil {
			res.SendErrorResponse(c, http.StatusBadRequest, "invalid page, "+page)
			return
		}

		c.Next()
	}
}

func UpdateAdmissionValidator() gin.HandlerFunc {
	return func(c *gin.Context) {
		var requestBody req.AdmissionPutReq
		// _ = c.ShouldBindBodyWith(&requestBody, binding.JSON)

		//mendapatkan data dari form-data
		err := c.ShouldBind(&requestBody)
		if err != nil {
			res.SendErrorResponse(c, http.StatusBadRequest, "failed bind form-data :"+err.Error())
		}

		if err = requestBody.Validate(); err != nil {
			res.SendErrorResponse(c, http.StatusBadRequest, "required data bind form-data , "+err.Error())
			return
		}

		c.Next()
	}
}
