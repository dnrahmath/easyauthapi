package services

import (
	"encoding/base64"
	"errors"
	"fmt"
	"strings"
	"time"

	"easyauthapi/configs"
	"easyauthapi/models/datastore"
	pl "easyauthapi/models/payload"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"

	"gorm.io/gorm"
)

//===================================================================

type TokenService struct {
	DB *gorm.DB
}

func (s *TokenService) NewService() *TokenService {
	return &TokenService{
		DB: configs.DB,
	}
}

//===================================================================

// CreateToken creates or updates a token record
func (s *TokenService) CreateOrUpdateToken(user *pl.User, tokenType string, expiresAt time.Time) (*pl.Token, error) {
	// Generate token string
	tokenString, err := generateToken(user, tokenType, expiresAt)
	if err != nil {
		return nil, err
	}

	// Filter criteria
	filterSrc := map[string]interface{}{
		"uuiduser": user.UuId,
		"type":     tokenType,
	}

	// Try to find existing token
	var existingToken pl.Token
	if err := s.DB.Where(filterSrc).First(&existingToken).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("error querying existing token: %v", err)
		}

		// Token not found, create new
		newToken := &pl.Token{
			UuId:        uuid.New(),
			UuIdUser:    user.UuId,
			Token:       tokenString,
			Type:        tokenType,
			ExpiresAt:   expiresAt,
			Blacklisted: false,
		}
		if err := s.DB.Create(newToken).Error; err != nil {
			return nil, fmt.Errorf("cannot save access token to db: %v", err)
		}
		return newToken, nil
	}

	// Token found, update existing
	existingToken.Token = tokenString
	existingToken.ExpiresAt = expiresAt
	if err := s.DB.Save(&existingToken).Error; err != nil {
		return nil, fmt.Errorf("cannot update access token, %v", err)
	}

	return &existingToken, nil
}

//==================================

// FindByUUID finds by UUID
func (s *TokenService) FindByUUIDUser(UuIdUser uuid.UUID, tokenType string) (*pl.Token, error) {
	var token pl.Token

	// Filter criteria
	filterSrc := map[string]interface{}{
		"uuiduser":    UuIdUser,
		"type":        tokenType,
		"blacklisted": false,
	}
	if err := s.DB.Where(filterSrc).First(&token).Error; err != nil {
		return nil, errors.New("cannot find token")
	}
	return &token, nil
}

// DeleteTokenByUuId deletes a token by its UuId
func (s *TokenService) DeleteTokenByUuIdUser(UuIdUser uuid.UUID) error {
	// Filter criteria
	filterSrc := map[string]interface{}{
		"uuiduser": UuIdUser,
	}
	//menghapus secara fisik Unscoped()
	if err := s.DB.Unscoped().Where(filterSrc).Delete(&pl.Token{}).Error; err != nil {
		return errors.New("cannot delete token")
	}
	return nil
}

//==================================

// ====
func (s *TokenService) GenerateAccessAndRefreshTokens(User *pl.User) (*pl.Token, *pl.Token, error) {
	accessExpiresAt := time.Now().Add(time.Duration(configs.JWTAccessExpirationMinutes) * time.Minute)
	refreshExpiresAt := time.Now().Add(time.Duration(configs.JWTRefreshExpirationDays) * 24 * time.Hour)

	accessToken, err := s.CreateOrUpdateToken(User, datastore.TokenTypeAccess, accessExpiresAt)
	if err != nil {
		return nil, nil, err
	}

	refreshToken, err := s.CreateOrUpdateToken(User, datastore.TokenTypeRefresh, refreshExpiresAt)
	if err != nil {
		return nil, nil, err
	}

	return accessToken, refreshToken, nil
}

// ====
func (s *TokenService) GetAccessAndRefreshTokens(UuIdUser uuid.UUID) (*pl.Token, *pl.Token, error) {
	accessToken, err := s.FindByUUIDUser(UuIdUser, datastore.TokenTypeAccess)
	if err != nil {
		return nil, nil, err
	}

	refreshToken, err := s.FindByUUIDUser(UuIdUser, datastore.TokenTypeRefresh)
	if err != nil {
		return nil, nil, err
	}

	return accessToken, refreshToken, nil
}

//==================================

// VerifyToken verifies the validity of a JWT token
func (s *TokenService) VerifyToken(tokenString string, tokenType string) (*pl.Token, error) {
	claims := &pl.UserClaims{}
	_, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(configs.JWTSecretKey), nil
	})

	if err != nil || claims.Type != tokenType {
		return nil, errors.New("token found to be invalid")
	}

	if time.Now().After(claims.ExpiresAt.Time) {
		return nil, errors.New("token is expired")
	}

	var token pl.Token
	filterSrc := map[string]interface{}{
		"token":       tokenString,
		"type":        tokenType,
		"blacklisted": false,
	}
	err = s.DB.Where(filterSrc).First(&token).Error
	if err != nil {
		return nil, fmt.Errorf("error finding token, %v", err)
	}

	return &token, nil
}

//==================================

// DecodeToken decodes a JWT token into its header, payload, and signature
func DecodeToken(tokenString string) (*pl.DecodeToken, error) {
	parts := strings.Split(tokenString, ".")
	if len(parts) != 3 {
		return nil, errors.New("invalid token format")
	}

	header, err := base64.RawURLEncoding.DecodeString(parts[0])
	if err != nil {
		return nil, fmt.Errorf("error decoding header: %v", err)
	}

	payload, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return nil, fmt.Errorf("error decoding payload: %v", err)
	}

	signature := parts[2]

	return &pl.DecodeToken{
		Header:    header,
		Payload:   payload,
		Signature: signature,
	}, nil
}

//==================================

// helper function to generate token
func generateToken(User *pl.User, tokenType string, expiresAt time.Time) (string, error) {
	claims := &pl.UserClaims{
		User: User,
		Type: tokenType,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			Subject:   fmt.Sprintf("%v", User.UuId),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(configs.JWTSecretKey))
}

//===================================================================

func (s *TokenService) CreateJWTVerify(AdmissionMap map[string]interface{}) string {
	// Set the expiration time for the token
	expirationTime := time.Now().Add(time.Duration(configs.JWTAccessExpirationMinutes) * time.Minute)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		// "user_data":        userMap,
		"Admission": AdmissionMap,
		"exp":       expirationTime.Unix(),
		// Add more claims as needed
	})

	// Sign and get the complete encoded token as a string
	tokenString, err := token.SignedString([]byte(configs.JWTSecretKey))
	if err != nil {
		// Handle error
		fmt.Println("Error creating JWT:", err)
		return ""
	}

	return tokenString
}

// =============================================================================================

// DecodeJWTVerifydecodes and verifies the JWT token
func (s *TokenService) DecodeJWTVerify(tokenString string) (map[string]interface{}, error) {
	// Parse the token with the provided secret key
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(configs.JWTSecretKey), nil
	})

	if err != nil {
		return nil, err
	}

	// Check if the token is valid
	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	// Extract the claims from the token
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("error extracting claims from token")
	}

	// Check the expiration time
	expirationTime := claims["exp"].(float64)
	if time.Now().Unix() > int64(expirationTime) {
		return nil, fmt.Errorf("token has expired")
	}

	// Extract the AdmissionMap from claims
	AdmissionMap, ok := claims["Admission"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("error extracting AdmissionMap from token")
	}

	return AdmissionMap, nil
}

// =============================================================================================
