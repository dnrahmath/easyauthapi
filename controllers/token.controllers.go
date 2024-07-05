package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"easyauthapi/models/convert"
	pl "easyauthapi/models/payload"
	res "easyauthapi/models/response"
	"easyauthapi/services"
	"easyauthapi/utils"
	"github.com/google/uuid"

	"github.com/gin-gonic/gin"
)

// =============================================================================================

type TokenController struct {
	UserService      *services.UserService
	AdmissionService *services.AdmissionService
	AdmissionConvert *convert.AdmissionConvert
	TokenService     *services.TokenService
	TypeVal          pl.TypeVal
}

func (ctrl *TokenController) NewController() *TokenController {
	UserService := services.UserService{}
	AdmissionService := services.AdmissionService{}
	AdmissionConvert := &convert.AdmissionConvert{}
	TokenService := services.TokenService{}
	return &TokenController{
		UserService:      UserService.NewService(),
		AdmissionService: AdmissionService.NewService(),
		AdmissionConvert: AdmissionConvert,
		TokenService:     TokenService.NewService(),
	}
}

// =============================================================================================

func (ctrl *TokenController) SendCode(c *gin.Context, typeparam string) {
	response := &res.Response{
		StatusCode: http.StatusBadRequest,
		Success:    false,
	}
	host := c.Request.Host

	//=======================================
	decodeToken, err := ctrl.GetDecodeToken(c)
	if err != nil {
		response.Message = err.Error()
		response.SendResponse(c)
		return
	}

	//=======================================
	Admission, err := ctrl.AdmissionService.FindByUUID(decodeToken.User.UuIdAdmission)
	if err != nil {
		res.SendErrorResponse(c, http.StatusBadRequest, "Failed to find user, "+err.Error())
		return
	}
	if Admission == nil {
		res.SendErrorResponse(c, http.StatusBadRequest, "User not found")
		return
	}

	//----------------------------------------------
	//jika diperbarui maka akan mengembalikan nilai "true"

	val := ctrl.TypeVal.DefaultValue()
	if typeparam == val.Email {
		Admission.EmailVerified = true
	} else if typeparam == val.PhoneNumber {
		Admission.PhoneNumberVerified = true
	} else {
		response.Message = "param type is not allowed"
		response.SendResponse(c)
		return
	}
	//----------------------------------------------

	//============================
	// SETUP MapInterface retun
	//============================
	ctrl.AdmissionConvert.Data = Admission
	AdmissionMap := ctrl.AdmissionConvert.ConvertToMap(true)
	//============================

	//============================
	tokenVerify := ctrl.TokenService.CreateJWTVerify(
		AdmissionMap,
	)
	//============================

	response.StatusCode = http.StatusOK
	response.Success = true
	response.Data = gin.H{
		"Url": fmt.Sprintf("%s/api/token/%s", host, tokenVerify),
		// "Admission": AdmissionMap, //nanti dihilangin
	}
	response.SendResponse(c)
}

func (ctrl *TokenController) UpdateOneVerify(c *gin.Context, tokenVerify string) {
	response := &res.Response{
		StatusCode: http.StatusBadRequest,
		Success:    false,
	}

	//============================
	tokenDecode, err := ctrl.TokenService.DecodeJWTVerify(
		tokenVerify,
	)
	if err != nil {
		response.Message = "failed decode token, " + err.Error()
		response.SendResponse(c)
	}
	//============================

	//-------------------------
	//UPDATE DATA
	//-------------------------
	tokenDecode["uuid"], _ = uuid.Parse(tokenDecode["uuid"].(string))
	tokenDecode["uuiduser"], _ = uuid.Parse(tokenDecode["uuiduser"].(string))
	tokenDecode["created_at"], _ = utils.ConvertStringToTime(tokenDecode["created_at"].(string))
	tokenDecode["updated_at"], _ = utils.ConvertStringToTime(tokenDecode["updated_at"].(string))
	err = ctrl.AdmissionConvert.ConvertMapOrByteToStruct(tokenDecode)
	if err != nil {
		response.Message = "failed convert to struct, " + err.Error()
		response.SendResponse(c)
	}
	//===============

	// filterSrc := bson.D{{Key: "uuiduser", Value: tokenDecode["uuiduser"]}}
	filterSrc := map[string]interface{}{
		"uuiduser": tokenDecode["uuiduser"],
	}
	Admission, err := ctrl.AdmissionService.Update(
		filterSrc,
		ctrl.AdmissionConvert.Data,
	)

	if err != nil {
		response.Message = err.Error()
		return
	}

	//============================
	// SETUP MapInterface retun
	//============================
	ctrl.AdmissionConvert.Data = Admission
	AdmissionMap := ctrl.AdmissionConvert.ConvertToMap(true)
	//============================

	response.StatusCode = http.StatusOK
	response.Success = true
	response.Message = "successfully updated the data"
	response.Data = gin.H{"Admission": AdmissionMap}
	response.SendResponse(c)
}

//==============================================================================================================

func (ctrl *TokenController) GetDecodeToken(c *gin.Context) (*pl.PayloadToken, error) {
	dataToken := &pl.PayloadToken{}

	token := c.GetHeader("Access-Token")
	if token == "" {
		return nil, errors.New("auth header is missing")
	}

	decodeToken, err := services.DecodeToken(token)
	if err != nil {
		return nil, err
	}

	// Mendapatkan payload dari decodeToken
	statusUnm := json.Unmarshal(decodeToken.Payload, &dataToken)
	if statusUnm != nil { //jika statusUnm == error
		return nil, statusUnm
	}

	return dataToken, nil
}

//==============================================================================================================
