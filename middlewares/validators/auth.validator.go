package validators

import (
	req "easyauthapi/models/request"
	res "easyauthapi/models/response"
	"net/http"
	"regexp"

	"github.com/gin-gonic/gin"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
)

func RegisterValidator() gin.HandlerFunc {
	return func(c *gin.Context) {
		var requestBody req.LoginOrRegisReq
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

func RegisterByAdminValidator() gin.HandlerFunc {
	return func(c *gin.Context) {
		var requestBody req.LoginOrRegisReq
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

func LoginValidator() gin.HandlerFunc {
	return func(c *gin.Context) {
		var requestBody req.LoginOrRegisReq
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

//===================================================================

// isValidEmail memeriksa apakah string adalah alamat email yang valid
func IsValidEmail(email string) bool {
	regex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return regex.MatchString(email)
}

// isValidPhoneNumber memeriksa apakah string adalah nomor telepon yang valid
func IsValidPhoneNumber(phoneNumber string) bool {
	// Sesuaikan dengan aturan validasi nomor telepon yang sesuai dengan kebutuhan Anda
	// Dalam contoh ini, hanya memeriksa apakah itu terdiri dari angka
	regex := regexp.MustCompile(`^\d+$`)
	return regex.MatchString(phoneNumber)
}

//===================================================================

func RefreshValidator() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenGet := c.GetHeader("Refresh-Token")
		if tokenGet == "" {
			res.SendErrorResponse(c, http.StatusBadRequest, "failed header Refresh-Token token not found. ")
		}
		c.Next()
	}
}

func GetUsersValidator() gin.HandlerFunc {
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

func UpdateUserValidator() gin.HandlerFunc {
	return func(c *gin.Context) {
		var requestBody req.UserPutReq
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

func CreateUserValidatorByAdm() gin.HandlerFunc {
	return func(c *gin.Context) {
		var requestBody req.UserPostReqByAdm
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

func UpdateUserValidatorByAdm() gin.HandlerFunc {
	return func(c *gin.Context) {
		var requestBody req.UserPutReqByAdm
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

//===================================================================
