package validators

import (
	"easyauthapi/models/datastore"
	pl "easyauthapi/models/payload"
	res "easyauthapi/models/response"
	"easyauthapi/services"

	"encoding/json"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

type PayloadToken pl.PayloadToken

// checkRoleToken memeriksa apakah dua slice string sama.
func (pyloadToken PayloadToken) checkRoleToken(token string, roleCheck string) ([]byte, error) {
	// claims := &db.UserClaims{}
	if token == "" {
		return nil, errors.New("auth header is missing")
	}

	decodeToken, err := services.DecodeToken(token)
	if err != nil {
		return nil, err
	}

	/*=====================================================*/

	// Mendapatkan payload dari decodeToken
	statusUnm := json.Unmarshal(decodeToken.Payload, &pyloadToken)
	if statusUnm != nil { //jika statusUnm == error
		return nil, statusUnm
	}

	/*=====================================================*/

	// Iterate through roles
	for _, role := range pyloadToken.User.Roles {
		if role.Name == roleCheck {
			return json.Marshal(pyloadToken)
		}
	}

	// If no match is found, return an error
	return nil, errors.New("role is not defined")

	/*=====================================================*/
}

/*===============*/

func checkRoleAndContinue(c *gin.Context, token string, roleCheck string) {
	var PayloadToken PayloadToken

	_, err := PayloadToken.checkRoleToken(token, roleCheck)
	if err != nil {
		res.SendErrorResponse(c, http.StatusBadRequest, "access not allowed, role is not "+roleCheck+". "+err.Error())
		return
	}

	c.Next()
}

/*===============*/

func IsGuest() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Access-Token")
		checkRoleAndContinue(c, token, datastore.RoleGuest)
	}
}

/*===============*/

func IsUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Access-Token")
		checkRoleAndContinue(c, token, datastore.RoleUser)
	}
}

/*===============*/

func IsMerchant() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Access-Token")
		checkRoleAndContinue(c, token, datastore.RoleMerchant)
	}
}

/*===============*/

func IsAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Access-Token")
		checkRoleAndContinue(c, token, datastore.RoleAdmin)
	}
}

/*===============*/
