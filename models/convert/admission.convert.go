package convert

import (
	pl "easyauthapi/models/payload"
	"easyauthapi/utils"
	"github.com/google/uuid"

	"encoding/json"
	"errors"
	"time"
)

//===================================================================

type AdmissionConvert struct {
	Data *pl.Admission
}

//===================================================================

func (model *AdmissionConvert) ConvertMapOrByteToStruct(data interface{}) error {
	if dataMap, ok := data.(map[string]interface{}); ok {
		model.Data.UuId = dataMap["uuid"].(uuid.UUID)
		model.Data.UuIdUser = dataMap["uuiduser"].(uuid.UUID)
		model.Data.Username = dataMap["username"].(string)
		model.Data.Password = dataMap["password"].(string)
		model.Data.Email = dataMap["email"].(string)
		model.Data.EmailVerified = dataMap["emailverified"].(bool)
		model.Data.PhoneNumber = dataMap["phonenumber"].(string)
		model.Data.PhoneNumberVerified = dataMap["phonenumberverified"].(bool)
		model.Data.CreatedAt = dataMap["created_at"].(time.Time)
		model.Data.UpdatedAt = dataMap["updated_at"].(time.Time)
		return nil
	}

	if dataByte, ok := data.([]byte); ok {
		return json.Unmarshal(dataByte, model)
	}

	return errors.New("unsupported data type")
}

//==================================

func (model *AdmissionConvert) ConvertToMap(timeString bool) map[string]interface{} {
	return map[string]interface{}{
		"uuid":                model.Data.UuId,
		"uuiduser":            model.Data.UuIdUser,
		"username":            model.Data.Username,
		"password":            model.Data.Password,
		"email":               model.Data.Email,
		"emailverified":       model.Data.EmailVerified,
		"phonenumber":         model.Data.PhoneNumber,
		"phonenumberverified": model.Data.PhoneNumberVerified,
		"created_at": utils.Ternary(
			timeString,
			model.Data.CreatedAt.Format(time.RFC3339),
			model.Data.CreatedAt,
		),
		"updated_at": utils.Ternary(
			timeString,
			model.Data.UpdatedAt.Format(time.RFC3339),
			model.Data.UpdatedAt,
		),
	}
}

//==================================

func (model *AdmissionConvert) ConvertToByte(timeString bool) ([]byte, error) {
	dataMap := model.ConvertToMap(timeString)
	return json.Marshal(dataMap)
}

//===================================================================
