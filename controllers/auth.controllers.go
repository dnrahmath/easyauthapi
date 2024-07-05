package controllers

import (
	"net/http"

	"golang.org/x/crypto/bcrypt"

	"easyauthapi/middlewares/validators"
	"easyauthapi/models/convert"
	"easyauthapi/models/datastore"
	pl "easyauthapi/models/payload"
	req "easyauthapi/models/request"
	res "easyauthapi/models/response"
	"easyauthapi/services"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// =============================================================================================

type AuthController struct {
	UserService      *services.UserService
	UserConvert      *convert.UserConvert
	AdmissionService *services.AdmissionService
	AdmissionConvert *convert.AdmissionConvert
	TokenController  TokenController
	TokenService     *services.TokenService
	TokenConvert     *convert.TokenConvert
	TypeVal          pl.TypeVal
}

func (ctrl *AuthController) NewController() *AuthController {
	UserService := services.UserService{}
	UserConvert := &convert.UserConvert{}
	AdmissionService := services.AdmissionService{}
	AdmissionConvert := &convert.AdmissionConvert{}
	TokenController := TokenController{}
	TokenService := services.TokenService{}
	TokenConvert := &convert.TokenConvert{}
	return &AuthController{
		UserService:      UserService.NewService(),
		UserConvert:      UserConvert,
		AdmissionService: AdmissionService.NewService(),
		AdmissionConvert: AdmissionConvert,
		TokenController:  *TokenController.NewController(),
		TokenService:     TokenService.NewService(),
		TokenConvert:     TokenConvert,
	}
}

// =============================================================================================

// Register handles user registration
func (ctrl *AuthController) Register(c *gin.Context) {
	var requestBody req.LoginOrRegisReq
	err := c.ShouldBind(&requestBody)
	if err != nil {
		res.SendErrorResponse(c, http.StatusBadRequest, "Failed to bind form-data, "+err.Error())
		return
	}

	UuIdAdmission := uuid.New()
	requestBodyUser := &pl.User{
		UuId:          uuid.New(),
		UuIdAdmission: UuIdAdmission,
	}
	requestBodyAdmission := &pl.Admission{
		UuId: UuIdAdmission,
	}

	if validators.IsValidEmail(requestBody.Value) {
		err = ctrl.AdmissionService.CheckByValue(ctrl.TypeVal.DefaultValue().Email, requestBody.Value)
		requestBodyAdmission.Email = requestBody.Value
	} else if validators.IsValidPhoneNumber(requestBody.Value) {
		err = ctrl.AdmissionService.CheckByValue(ctrl.TypeVal.DefaultValue().PhoneNumber, requestBody.Value)
		requestBodyAdmission.PhoneNumber = requestBody.Value
	} else {
		err = ctrl.AdmissionService.CheckByValue(ctrl.TypeVal.DefaultValue().Username, requestBody.Value)
		requestBodyAdmission.Username = requestBody.Value
	}
	if err != nil {
		res.SendErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	//====================

	User, err := ctrl.UserService.Create(requestBodyUser)
	if err != nil {
		res.SendErrorResponse(c, http.StatusBadRequest, "Failed to create user, "+err.Error())
		return
	}

	requestBodyAdmission.UuIdUser = User.UuId
	requestBodyAdmission.Password = requestBody.Password
	Admission, err := ctrl.AdmissionService.Create(requestBodyAdmission)
	if err != nil {
		res.SendErrorResponse(c, http.StatusBadRequest, "Failed to create admission, "+err.Error())
		return
	}

	//====================

	ctrl.UserConvert.Data = User
	UserMap := ctrl.UserConvert.ConvertToMap(true)

	ctrl.AdmissionConvert.Data = Admission
	AdmissionMap := ctrl.AdmissionConvert.ConvertToMap(true)

	//====================

	accessToken, refreshToken, err := ctrl.TokenService.GenerateAccessAndRefreshTokens(User)
	if err != nil {
		res.SendErrorResponse(c, http.StatusBadRequest, "Failed to generate tokens, "+err.Error())
		return
	}

	ctrl.TokenConvert.Data = accessToken
	accessTokenRes := ctrl.TokenConvert.GetResponseJson()
	c.Header("Access-Token", accessTokenRes["token"])

	ctrl.TokenConvert.Data = refreshToken
	refreshTokenRes := ctrl.TokenConvert.GetResponseJson()
	c.Header("Refresh-Token", refreshTokenRes["token"])

	//====================

	response := &res.Response{
		StatusCode: http.StatusCreated,
		Success:    true,
		Data: gin.H{
			"User":      UserMap,
			"Admission": AdmissionMap,
			// "token": gin.H{
			// 	"access":  accessTokenRes,
			// 	"refresh": refreshTokenRes,
			// },
		},
	}
	response.SendResponse(c)
}

// Login handles user login
func (ctrl *AuthController) Login(c *gin.Context) {
	var requestBody req.LoginOrRegisReq
	err := c.ShouldBind(&requestBody)
	if err != nil {
		res.SendErrorResponse(c, http.StatusBadRequest, "Failed to bind form-data, "+err.Error())
		return
	}

	val := ctrl.TypeVal.DefaultValue()
	var Admission *pl.Admission
	if validators.IsValidEmail(requestBody.Value) {
		filterSrc := map[string]interface{}{val.Email: requestBody.Value}
		Admission, err = ctrl.AdmissionService.FindByValue(filterSrc)
	} else if validators.IsValidPhoneNumber(requestBody.Value) {
		filterSrc := map[string]interface{}{val.PhoneNumber: requestBody.Value}
		Admission, err = ctrl.AdmissionService.FindByValue(filterSrc)
	} else {
		filterSrc := map[string]interface{}{val.Username: requestBody.Value}
		Admission, err = ctrl.AdmissionService.FindByValue(filterSrc)
	}
	if err != nil {
		res.SendErrorResponse(c, http.StatusBadRequest, "Login not found")
		return
	}

	//====================

	User, err := ctrl.UserService.FindByUUID(Admission.UuIdUser)
	if err != nil {
		res.SendErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(Admission.Password), []byte(requestBody.Password))
	if err != nil {
		res.SendErrorResponse(c, http.StatusBadRequest, "Invalid credentials")
		return
	}

	ctrl.UserConvert.Data = User
	UserMap := ctrl.UserConvert.ConvertToMap(true)

	ctrl.AdmissionConvert.Data = Admission
	AdmissionMap := ctrl.AdmissionConvert.ConvertToMap(true)
	//====================

	accessToken, refreshToken, err := ctrl.TokenService.GenerateAccessAndRefreshTokens(User)
	if err != nil {
		res.SendErrorResponse(c, http.StatusBadRequest, "Failed to generate tokens, "+err.Error())
		return
	}

	ctrl.TokenConvert.Data = accessToken
	accessTokenRes := ctrl.TokenConvert.GetResponseJson()
	c.Header("Access-Token", accessTokenRes["token"])

	ctrl.TokenConvert.Data = refreshToken
	refreshTokenRes := ctrl.TokenConvert.GetResponseJson()
	c.Header("Refresh-Token", refreshTokenRes["token"])

	//====================

	response := &res.Response{
		StatusCode: http.StatusOK,
		Success:    true,
		Data: gin.H{
			"User":      UserMap,
			"Admission": AdmissionMap,
			// "token": gin.H{
			// 	"access":  accessTokenRes,
			// 	"refresh": refreshTokenRes,
			// },
		},
	}
	response.SendResponse(c)
}

// Refresh handles token refresh
func (ctrl *AuthController) Refresh(c *gin.Context) {
	tokenGet := c.GetHeader("Refresh-Token")
	token, err := ctrl.TokenService.VerifyToken(tokenGet, datastore.TokenTypeRefresh)
	if err != nil {
		res.SendErrorResponse(c, http.StatusBadRequest, "Failed to verify refresh token, "+err.Error())
		return
	}

	//====================

	User, err := ctrl.UserService.FindByUUID(token.UuIdUser)
	if err != nil {
		res.SendErrorResponse(c, http.StatusBadRequest, "Failed to find user, "+err.Error())
		return
	}

	Admission, err := ctrl.AdmissionService.FindByUUID(User.UuIdAdmission)
	if err != nil {
		res.SendErrorResponse(c, http.StatusBadRequest, "Failed to find admission, "+err.Error())
		return
	}

	ctrl.UserConvert.Data = User
	UserMap := ctrl.UserConvert.ConvertToMap(true)

	ctrl.AdmissionConvert.Data = Admission
	AdmissionMap := ctrl.AdmissionConvert.ConvertToMap(true)
	//====================

	err = ctrl.TokenService.DeleteTokenByUuIdUser(token.UuIdUser)
	if err != nil {
		res.SendErrorResponse(c, http.StatusBadRequest, "Failed to delete old token, "+err.Error())
		return
	}

	accessToken, refreshToken, err := ctrl.TokenService.GenerateAccessAndRefreshTokens(User)
	if err != nil {
		res.SendErrorResponse(c, http.StatusBadRequest, "Failed to generate new tokens, "+err.Error())
		return
	}

	ctrl.TokenConvert.Data = accessToken
	accessTokenRes := ctrl.TokenConvert.GetResponseJson()
	c.Header("Access-Token", accessTokenRes["token"])

	ctrl.TokenConvert.Data = refreshToken
	refreshTokenRes := ctrl.TokenConvert.GetResponseJson()
	c.Header("Refresh-Token", refreshTokenRes["token"])

	//====================

	response := &res.Response{
		StatusCode: http.StatusOK,
		Success:    true,
		Data: gin.H{
			"User":      UserMap,
			"Admission": AdmissionMap,
			// "token": gin.H{
			// 	"access":  accessTokenRes,
			// 	"refresh": refreshTokenRes,
			// },
		},
	}
	response.SendResponse(c)
}

// Logout user
func (ctrl *AuthController) LogoutUser(c *gin.Context) {
	response := &res.Response{
		StatusCode: http.StatusBadRequest,
		Success:    false,
	}

	//-----
	decodeToken, err := ctrl.TokenController.GetDecodeToken(c)
	if err != nil {
		response.Message = err.Error()
		response.SendResponse(c)
		return
	}

	//-----
	err = ctrl.TokenService.DeleteTokenByUuIdUser(decodeToken.User.UuId)
	if err != nil {
		res.SendErrorResponse(c, http.StatusBadRequest, "Failed to delete old token, "+err.Error())
		return
	}

	c.Header("Access-Token", "")
	c.Header("Refresh-Token", "")

	//====================
	response.StatusCode = http.StatusOK
	response.Success = true
	response.Message = "Logout Successfully"
	response.SendResponse(c)
}

/* ==================================================================== */

// yang di gunakan dalam route get, bukan berasal dari database, melainkan hasil decode dari tokem , jadi jika kawtu dari decode token belum habis masih bisa digunakan
func (ctrl *AuthController) MeGet(c *gin.Context) {
	response := &res.Response{
		StatusCode: http.StatusBadRequest,
		Success:    false,
	}

	//-----
	decodeToken, err := ctrl.TokenController.GetDecodeToken(c)
	if err != nil {
		response.Message = err.Error()
		response.SendResponse(c)
		return
	}

	//-----
	User, err := ctrl.UserService.FindByUUID(decodeToken.User.UuId)
	if err != nil {
		response.Message = "Failed to get user, " + err.Error()
		response.SendResponse(c)
		return
	}

	// Filter criteria
	filterSrc := map[string]interface{}{
		"uuiduser": User.UuId,
	}

	Admission, err := ctrl.AdmissionService.FindByValue(filterSrc)
	if err != nil {
		response.Message = "Failed to get user, " + err.Error()
		response.SendResponse(c)
		return
	}

	//====================
	ctrl.UserConvert.Data = User
	UserMap := ctrl.UserConvert.ConvertToMap(true)

	ctrl.AdmissionConvert.Data = Admission
	AdmissionMap := ctrl.AdmissionConvert.ConvertToMap(true)

	//====================
	accessToken, refreshToken, err := ctrl.TokenService.GetAccessAndRefreshTokens(decodeToken.User.UuId)
	if err != nil {
		res.SendErrorResponse(c, http.StatusBadRequest, "Please login again, "+err.Error())
		return
	}

	ctrl.TokenConvert.Data = accessToken
	accessTokenRes := ctrl.TokenConvert.GetResponseJson()
	c.Header("Access-Token", accessTokenRes["token"])

	ctrl.TokenConvert.Data = refreshToken
	refreshTokenRes := ctrl.TokenConvert.GetResponseJson()
	c.Header("Refresh-Token", refreshTokenRes["token"])

	//====================
	response.StatusCode = http.StatusOK
	response.Success = true
	response.Data = gin.H{
		"User":      UserMap,
		"Admission": AdmissionMap,
	}
	response.SendResponse(c)
}

// /* ----------------------------------------------- */

// func (ctrl *AuthController) MePut(c *gin.Context) {
// 	var requestBody req.UserPutReq

// 	//mendapatkan data dari form-data
// 	err := c.ShouldBind(&requestBody)
// 	if err != nil {
// 		res.SendErrorResponse(c, http.StatusBadRequest, "failed bind form-data :"+err.Error())
// 	}

// 	response := &res.Response{
// 		StatusCode: http.StatusBadRequest,
// 		Success:    false,
// 	}

// 	//-----
// 	decodeToken, err := ctrl.TokenController.GetDecodeToken(c)
// 	if err != nil {
// 		response.Message = err.Error()
// 		response.SendResponse(c)
// 		return
// 	}
// 	filterSrc := map[string]interface{}{
// 		"uuid": decodeToken.User.UuId,
// 	}
// 	//-----

// 	response = FuncUpdateOneUser(c, colMap, response, filterSrc)
// 	response.SendResponse(c)
// }

// /* ----------------------------------------------- */

// func (ctrl *AuthController) MeDel(c *gin.Context) {
// 	response := &res.Response{
// 		StatusCode: http.StatusBadRequest,
// 		Success:    false,
// 	}

// 	//-----
// 	decodeToken, err := ctrl.TokenController.GetDecodeToken(c)
// 	if err != nil {
// 		response.Message = err.Error()
// 		response.SendResponse(c)
// 		return
// 	}
// 	filterSrc := map[string]interface{}{
// 		"uuid": decodeToken.User.UuId,
// 	}
// 	//-----

// 	response = FuncDeleteOneUser(c, NameCol, response, filterSrc)
// 	response.SendResponse(c)
// }

// /* ==================================================================== */
