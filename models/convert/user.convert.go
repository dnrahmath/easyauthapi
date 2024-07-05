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

type UserConvert struct {
	Data *pl.User
}

//===================================================================

func (model *UserConvert) ConvertMapOrByteToStruct(data interface{}) error {
	if dataMap, ok := data.(map[string]interface{}); ok {
		model.Data.UuId = dataMap["uuid"].(uuid.UUID)
		model.Data.UuIdAddr = dataMap["uuidaddr"].(uuid.UUID)
		model.Data.UuIdAdmission = dataMap["uuidAdmission"].(uuid.UUID)
		model.Data.Gender = dataMap["gender"].(string)
		model.Data.Name = dataMap["name"].(string)
		model.Data.Noid = dataMap["noid"].(string)
		model.Data.Religion = dataMap["religion"].(string)
		model.Data.Roles = dataMap["role"].([]pl.Role)
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

func (model *UserConvert) ConvertToMap(timeString bool) map[string]interface{} {
	return map[string]interface{}{
		"uuid":          model.Data.UuId,
		"uuidaddr":      model.Data.UuIdAddr,
		"uuidAdmission": model.Data.UuIdAdmission,
		"gender":        model.Data.Gender,
		"name":          model.Data.Name,
		"noid":          model.Data.Noid,
		"religion":      model.Data.Religion,
		"role":          model.Data.Roles,
		"created_at":    utils.Ternary(timeString, model.Data.CreatedAt.Format(time.RFC3339), model.Data.CreatedAt),
		"updated_at":    utils.Ternary(timeString, model.Data.UpdatedAt.Format(time.RFC3339), model.Data.UpdatedAt),
	}
}

//==================================

func (model *UserConvert) ConvertToByte(timeString bool) ([]byte, error) {
	dataMap := model.ConvertToMap(timeString)
	return json.Marshal(dataMap)
}

//====================================================================
