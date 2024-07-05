package datastore

import (
	"time"

	pl "easyauthapi/models/payload"
	"easyauthapi/utils"
	"github.com/google/uuid"
)

//===================================================================

type TokenDatastore pl.Token

//===================================================================

const (
	TokenTypeAccess  = "access"
	TokenTypeRefresh = "refresh"
)

//==================================

func (DataNew *TokenDatastore) Post() *pl.Token {
	data := &pl.Token{
		UuId:        DataNew.UuId,
		UuIdUser:    DataNew.UuIdUser,
		Token:       DataNew.Token,
		Type:        DataNew.Type,
		ExpiresAt:   DataNew.ExpiresAt,
		Blacklisted: DataNew.Blacklisted,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	return data
}

//==================================

func (DataPrev *TokenDatastore) Put(DataNew *pl.Token) *pl.Token {
	// Using separate conditions for each field
	UuIdUser := utils.Ternary(DataNew.UuIdUser != uuid.Nil, DataNew.UuIdUser, DataPrev.UuIdUser).(uuid.UUID)
	Token := utils.Ternary(DataNew.Token != "", DataNew.Token, DataPrev.Token).(string)
	Type := utils.Ternary(DataNew.Type != "", DataNew.Type, DataPrev.Type).(string)
	ExpiresAt := utils.Ternary(!DataNew.ExpiresAt.IsZero(), DataNew.ExpiresAt, DataPrev.ExpiresAt).(time.Time)
	Blacklisted := utils.Ternary(DataNew.Blacklisted, DataNew.Blacklisted, DataPrev.Blacklisted).(bool)

	// Creating a new object with the set values
	data := &pl.Token{
		UuIdUser:    UuIdUser,
		Token:       Token,
		Type:        Type,
		ExpiresAt:   ExpiresAt,
		Blacklisted: Blacklisted,
		CreatedAt:   DataPrev.CreatedAt,
		UpdatedAt:   time.Now(),
	}

	return data
}

//===================================================================
