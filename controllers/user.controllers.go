package controllers

import (
	"net/http"
	"strconv"

	"easyauthapi/models/convert"
	pl "easyauthapi/models/payload"
	req "easyauthapi/models/request"
	res "easyauthapi/models/response"
	"easyauthapi/services"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

//===============================================================================================================

type UserController struct {
	UserService      *services.UserService
	UserConvert      *convert.UserConvert
	AdmissionService *services.AdmissionService
	AdmissionConvert *convert.AdmissionConvert
	TokenService     *services.TokenService
	TokenConvert     *convert.TokenConvert
	TypeVal          pl.TypeVal
}

func (ctrl *UserController) NewController() *UserController {
	UserService := services.UserService{}
	AdmissionService := services.AdmissionService{}
	TokenService := services.TokenService{}
	return &UserController{
		UserService:      UserService.NewService(),
		AdmissionService: AdmissionService.NewService(),
		TokenService:     TokenService.NewService(),
	}
}

//===============================================================================================================

func (ctrl *UserController) Create(c *gin.Context) {
	response := &res.Response{
		StatusCode: http.StatusBadRequest,
		Success:    false,
	}

	var requestBody req.UserPutReq
	err := c.ShouldBind(&requestBody)
	if err != nil {
		response.Message = "failed to bind form-data: " + err.Error()
		response.SendResponse(c)
		return
	}

	/* =================================================================== */
	//check
	if requestBody.Email != "" {
		err = ctrl.AdmissionService.CheckByValue(ctrl.TypeVal.DefaultValue().Email, requestBody.Email)
		if err != nil {
			response.Message = err.Error()
			return
		}
	}
	if requestBody.PhoneNumber != "" {
		err = ctrl.AdmissionService.CheckByValue(ctrl.TypeVal.DefaultValue().PhoneNumber, requestBody.PhoneNumber)
		if err != nil {
			response.Message = err.Error()
			return
		}
	}
	if requestBody.Username != "" {
		err = ctrl.AdmissionService.CheckByValue(ctrl.TypeVal.DefaultValue().Username, requestBody.Username)
		if err != nil {
			response.Message = err.Error()
			return
		}
	}
	/* =================================================================== */

	// create User record
	UuIdAdmission := uuid.New()
	requestBodyUser := &pl.User{
		UuId:          uuid.New(),
		UuIdAddr:      uuid.Nil,
		UuIdAdmission: UuIdAdmission,
		Gender:        requestBody.Gender,
		Name:          requestBody.Name,
		Noid:          requestBody.Noid,
		Religion:      requestBody.Religion,
		Roles:         nil,
	}
	// create Admission record
	requestBodyAdmission := &pl.Admission{
		UuId:                UuIdAdmission,
		EmailVerified:       false,
		PhoneNumberVerified: false,
	}

	/* =================================================================== */

	User, err := ctrl.UserService.Create(requestBodyUser)
	if err != nil {
		response.Message = "failed to create user: " + err.Error()
		response.SendResponse(c)
		return
	}

	Admission, err := ctrl.AdmissionService.Create(requestBodyAdmission)
	if err != nil {
		response.Message = "failed to create user type sign: " + err.Error()
		response.SendResponse(c)
		return
	}
	ctrl.UserConvert.Data = User
	UserMap := ctrl.UserConvert.ConvertToMap(true)

	ctrl.AdmissionConvert.Data = Admission
	AdmissionMap := ctrl.AdmissionConvert.ConvertToMap(true)
	//====================

	accessToken, refreshToken, err := ctrl.TokenService.GenerateAccessAndRefreshTokens(requestBodyUser)
	if err != nil {
		response.Message = "failed to create token, " + err.Error()
		response.SendResponse(c)
		return
	}

	ctrl.TokenConvert.Data = accessToken
	accessTokenRes := ctrl.TokenConvert.GetResponseJson()

	ctrl.TokenConvert.Data = refreshToken
	refreshTokenRes := ctrl.TokenConvert.GetResponseJson()

	//====================

	response.StatusCode = http.StatusCreated
	response.Success = true
	response.Message = "successfully created user"
	response.Data = gin.H{
		"User":      UserMap,
		"Admission": AdmissionMap,
		"token": gin.H{
			"access":  accessTokenRes,
			"refresh": refreshTokenRes,
		},
	}
	response.SendResponse(c)
}

//===============================================================================================================

func (ctrl *UserController) Get(c *gin.Context) {
	response := &res.Response{
		StatusCode: http.StatusBadRequest,
		Success:    false,
	}

	//=======================================
	pageQuery := c.DefaultQuery("page", "0")
	page, err := strconv.Atoi(pageQuery)
	if err != nil {
		response.Message = "Invalid page number"
		response.SendResponse(c)
		return
	}

	limit := 5
	filterSrc := map[string]interface{}{
		"gender":   "male",
		"religion": "islam",
	}
	//=======================================

	Users, err := ctrl.UserService.GetByQuery(filterSrc, page, limit)
	if err != nil {
		response.Message = "error fetching Users by query, " + err.Error()
		return
	}
	//-----

	//===============
	hasPrev := page > 0
	hasNext := len(Users) > limit

	if hasNext {
		Users = Users[:limit]
	}

	//============================
	// SETUP MapInterface retun
	//============================
	var UsersMap []interface{}
	for _, User := range Users {
		ctrl.UserConvert.Data = &User
		UsersMap = append(UsersMap, ctrl.UserConvert.ConvertToMap(true))
	}
	//============================
	response.StatusCode = http.StatusOK
	response.Success = true
	response.Data = gin.H{
		"User": UsersMap,
		"prev": hasPrev,
		"next": hasNext,
	}
	response.SendResponse(c)
}

//===============================================================================================================

func (ctrl *UserController) GetOne(c *gin.Context) {
	response := &res.Response{
		StatusCode: http.StatusBadRequest,
		Success:    false,
	}

	uuidStr := c.Param("uuid")
	userId, err := uuid.Parse(uuidStr)
	if err != nil {
		response.Message = "Failed to parse uuid: " + err.Error()
		response.SendResponse(c)
		return
	}
	User, err := ctrl.UserService.FindByUUID(userId)
	if err != nil {
		response.Message = "Failed to get user: " + err.Error()
		response.SendResponse(c)
		return
	}

	// Filter criteria
	filterSrc := map[string]interface{}{
		"uuid": User.UuId,
	}

	Admission, err := ctrl.AdmissionService.FindByValue(filterSrc)
	if err != nil {
		response.Message = "Failed to get user: " + err.Error()
		response.SendResponse(c)
		return
	}

	//====================
	ctrl.UserConvert.Data = User
	UserMap := ctrl.UserConvert.ConvertToMap(true)

	ctrl.AdmissionConvert.Data = Admission
	AdmissionMap := ctrl.AdmissionConvert.ConvertToMap(true)
	//====================

	response.StatusCode = http.StatusOK
	response.Success = true
	response.Data = gin.H{
		"User":      UserMap,
		"Admission": AdmissionMap,
	}
	response.SendResponse(c)
}

//===============================================================================================================

func (ctrl *UserController) UpdateOne(c *gin.Context) {
	response := &res.Response{
		StatusCode: http.StatusBadRequest,
		Success:    false,
	}

	var requestBody req.UserPutReqByAdm
	if err := c.ShouldBind(&requestBody); err != nil {
		response.Message = "Failed to bind form-data: " + err.Error()
		response.SendResponse(c)
		return
	}

	/* =================================================================== */

	// create User record
	UuIdAdmission := uuid.New()
	requestBodyUser := &pl.User{
		UuId:          uuid.New(),
		UuIdAddr:      uuid.Nil,
		UuIdAdmission: UuIdAdmission,
		Gender:        requestBody.Gender,
		Name:          requestBody.Name,
		Noid:          requestBody.Noid,
		Religion:      requestBody.Religion,
		Roles:         nil,
	}

	/* =================================================================== */

	// userId := c.Param("id")
	filterSrc := map[string]interface{}{
		"gender":   "male",
		"religion": "islam",
	}
	User, err := ctrl.UserService.Update(filterSrc, requestBodyUser)
	if err != nil {
		response.Message = "Failed to update user: " + err.Error()
		response.SendResponse(c)
		return
	}

	//====================
	ctrl.UserConvert.Data = User
	UserMap := ctrl.UserConvert.ConvertToMap(true)
	//====================
	// AdmissionMap := Admission.ConvertToMap(true)

	response.StatusCode = http.StatusOK
	response.Success = true
	response.Message = "Successfully updated user"
	response.Data = gin.H{
		"User": UserMap,
	}
	response.SendResponse(c)
}

//===============================================================================================================

func (ctrl *UserController) DeleteOne(c *gin.Context) {
	response := &res.Response{
		StatusCode: http.StatusBadRequest,
		Success:    false,
	}

	// userId := c.Param("id")
	filterSrc := map[string]interface{}{
		"gender":   "male",
		"religion": "islam",
	}
	err := ctrl.UserService.Delete(filterSrc)
	if err != nil {
		response.Message = "Failed to delete user: " + err.Error()
		response.SendResponse(c)
		return
	}

	response.StatusCode = http.StatusOK
	response.Success = true
	response.Message = "Successfully deleted user"
	response.SendResponse(c)
}

//===============================================================================================================
