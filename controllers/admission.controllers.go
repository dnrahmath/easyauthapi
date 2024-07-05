package controllers

import (
	"easyauthapi/models/convert"
	pl "easyauthapi/models/payload"
	req "easyauthapi/models/request"
	res "easyauthapi/models/response"
	"easyauthapi/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

//===============================================================================================================

type AdmissionController struct {
	AdmissionService *services.AdmissionService
	AdmissionConvert *convert.AdmissionConvert
}

func (ctrl *AdmissionController) NewController() *AdmissionController {
	AdmissionService := services.AdmissionService{}
	return &AdmissionController{
		AdmissionService: AdmissionService.NewService(),
	}
}

//===============================================================================================================

func (ctrl *AdmissionController) Create(c *gin.Context) {
	response := &res.Response{
		StatusCode: http.StatusBadRequest,
		Success:    false,
	}

	var requestBody req.AdmissionPutReq

	//mendapatkan data dari form-data
	err := c.ShouldBind(&requestBody)
	if err != nil {
		res.SendErrorResponse(c, http.StatusBadRequest, "failed bind form-data :"+err.Error())
	}

	//-------------------------
	// UuIdUser, exists := c.Get("UuIdUser")
	// if !exists {
	// 	response.Message = "cannot get user"
	// 	response.SendResponse(c)
	// 	return response
	// }

	//-------------------------
	requestBodyAdmission := &pl.Admission{
		UuId:                uuid.New(),
		UuIdUser:            requestBody.UuIdUser,
		Username:            requestBody.Username,
		Password:            requestBody.Password,
		Email:               requestBody.Email,
		EmailVerified:       false,
		PhoneNumber:         requestBody.PhoneNumber,
		PhoneNumberVerified: false,
	}

	Admission, err := ctrl.AdmissionService.Create(
		requestBodyAdmission,
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

	response.StatusCode = http.StatusCreated
	response.Success = true
	response.Message = "successfully created the data"
	response.Data = gin.H{"Admission": AdmissionMap}
	response.SendResponse(c)
}

// ================================================================

func (ctrl *AdmissionController) Get(c *gin.Context, optionparam string, optionsrcparam string) {
	response := &res.Response{
		StatusCode: http.StatusBadRequest,
		Success:    false,
	}

	filterSrc := map[string]interface{}{
		"gender":   "male",
		"religion": "islam",
	}

	//===============
	pageQuery := c.DefaultQuery("page", "0")
	page, _ := strconv.Atoi(pageQuery)
	limit := 5

	//===============
	var Admissions []pl.Admission

	//-----
	Admissions, err := ctrl.AdmissionService.GetByQuery(filterSrc, page, limit)
	if err != nil {
		response.Message = "error fetching Admissions by query, " + err.Error()
		return
	}
	//-----

	//===============
	hasPrev := page > 0
	hasNext := len(Admissions) > limit

	if hasNext {
		Admissions = Admissions[:limit]
	}

	//============================
	// SETUP MapInterface retun
	//============================
	var AdmissionsMap []interface{}
	for _, Admission := range Admissions {
		ctrl.AdmissionConvert.Data = &Admission
		AdmissionsMap = append(AdmissionsMap, ctrl.AdmissionConvert.ConvertToMap(true))
	}
	//============================

	response.StatusCode = http.StatusOK
	response.Success = true
	response.Data = gin.H{"Admission": AdmissionsMap, "prev": hasPrev, "next": hasNext}
	response.SendResponse(c)
}

// ================================================================

func (ctrl *AdmissionController) GetOne(c *gin.Context, optionparam string, optionsrcparam string) {
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

	//-----
	Admissions, err := ctrl.AdmissionService.GetByQuery(filterSrc, page, limit)
	if err != nil {
		response.Message = "failed get one data " + err.Error()
		response.SendResponse(c)
		return
	}
	//-----

	//============================
	// SETUP MapInterface retun
	//============================
	var AdmissionsMap []interface{}
	for _, Admission := range Admissions {
		ctrl.AdmissionConvert.Data = &Admission
		AdmissionsMap = append(AdmissionsMap, ctrl.AdmissionConvert.ConvertToMap(true))
	}
	//============================

	response.StatusCode = http.StatusOK
	response.Success = true
	response.Data = gin.H{"Admission": AdmissionsMap}
	response.SendResponse(c)
}

// ================================================================
func (ctrl *AdmissionController) UpdateOne(c *gin.Context, optionparam string, optionsrcparam string, allowEmptyQuery bool) {
	response := &res.Response{
		StatusCode: http.StatusBadRequest,
		Success:    false,
	}

	// filterSrc map[string]interface{}

	// filterSrc := bson.D{{Key: "metadata.uuidcol", Value: "GANTI"}}
	filterSrc := map[string]interface{}{
		"gender":   "male",
		"religion": "islam",
	}

	var requestBody req.AdmissionPutReq

	//mendapatkan data dari form-data
	err := c.ShouldBind(&requestBody)
	if err != nil {
		res.SendErrorResponse(c, http.StatusBadRequest, "failed bind form-data :"+err.Error())
	}

	//-------------------------
	//UPDATE DATA
	//-------------------------
	requestBodyAdmission := &pl.Admission{
		UuIdUser:    requestBody.UuIdUser,
		Username:    requestBody.Username,
		Password:    requestBody.Password,
		Email:       requestBody.Email,
		PhoneNumber: requestBody.PhoneNumber,
	}

	//----------------------------------------------
	//jika diperbarui maka akan mengembalikan nilai "false"
	if requestBodyAdmission.Email != "" {
		requestBodyAdmission.EmailVerified = false
	}
	if requestBodyAdmission.PhoneNumber != "" {
		requestBodyAdmission.PhoneNumberVerified = false
	}
	//----------------------------------------------

	Admission, err := ctrl.AdmissionService.Update(
		filterSrc,
		requestBodyAdmission,
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

/*
- lebih baik [-deleteOne-] dan [-updateOne-]
- karena [-updateMultiple-] dan [-deleteMultiple-] memiliki resiko
- jika ingin tetap [-Multiple-] maka di setup looping request pada [-Frontend-]
*/
// ================================================================
func (ctrl *AdmissionController) DeleteOne(c *gin.Context, optionparam string, optionsrcparam string, allowEmptyQuery bool) {
	response := &res.Response{
		StatusCode: http.StatusBadRequest,
		Success:    false,
	}

	// filterSrc := bson.D{{Key: "metadata.uuidcol", Value: "GANTI"}}
	// filterSrc := make(map[string]interface{})
	filterSrc := map[string]interface{}{
		"gender":   "male",
		"religion": "islam",
	}

	//-------------------------
	//DELETE DATA
	//-------------------------
	// deleteResult, err := ctrl.AdmissionService.DeleteOne(filterSrc)
	err := ctrl.AdmissionService.DeleteOne(filterSrc)
	if err != nil {
		response.Message = err.Error()
		return
	}

	//-------------------------
	response.StatusCode = http.StatusOK
	response.Success = true
	response.Message = "successfully deleted data"
	response.SendResponse(c)
}

//==============================================================================================================
//==============================================================================================================
