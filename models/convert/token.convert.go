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

type TokenConvert struct {
	Data *pl.Token
}

//===================================================================

func (model *TokenConvert) ConvertMapOrByteToStruct(data interface{}) error {
	if dataMap, ok := data.(map[string]interface{}); ok {
		model.Data.UuId = dataMap["uuid"].(uuid.UUID)
		model.Data.UuIdUser = dataMap["uuiduser"].(uuid.UUID)
		model.Data.Token = dataMap["token"].(string)
		model.Data.Type = dataMap["type"].(string)
		model.Data.ExpiresAt = dataMap["expires_at"].(time.Time)
		model.Data.Blacklisted = dataMap["blacklisted"].(bool)
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

func (model *TokenConvert) ConvertToMap(timeString bool) map[string]interface{} {
	return map[string]interface{}{
		"uuid":        model.Data.UuId,
		"uuiduser":    model.Data.UuIdUser,
		"token":       model.Data.Token,
		"type":        model.Data.Type,
		"expires_at":  model.Data.ExpiresAt,
		"blacklisted": model.Data.Blacklisted,
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

func (model *TokenConvert) ConvertToByte(timeString bool) ([]byte, error) {
	dataMap := model.ConvertToMap(timeString)
	return json.Marshal(dataMap)
}

// ==================================

func (model *TokenConvert) GetResponseJson() map[string]string {
	return map[string]string{
		"token":   model.Data.Token,
		"expires": model.Data.ExpiresAt.Format(time.RFC3339),
	}
}

//===================================================================
