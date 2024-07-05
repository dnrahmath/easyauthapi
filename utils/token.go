package utils

import (
	"encoding/base64"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
)

// CreateToken using RSA private key
func CreateToken(ttl time.Duration, payload interface{}, privateKey string) (string, error) {
	decodedPrivateKey, err := base64.StdEncoding.DecodeString(privateKey)
	if err != nil {
		return "", fmt.Errorf("could not decode key: %w", err)
	}
	key, err := jwt.ParseRSAPrivateKeyFromPEM(decodedPrivateKey)
	if err != nil {
		return "", fmt.Errorf("create: parse key: %w", err)
	}

	now := time.Now().UTC()

	claims := make(jwt.MapClaims)
	claims["sub"] = payload
	claims["exp"] = now.Add(ttl).Unix()
	claims["iat"] = now.Unix()
	claims["nbf"] = now.Unix()

	token, err := jwt.NewWithClaims(jwt.SigningMethodRS256, claims).SignedString(key)
	if err != nil {
		return "", fmt.Errorf("create: sign token, %w", err)
	}

	return token, nil
}

// ValidateToken using RSA public key
func ValidateToken(token string, publicKey string) (interface{}, error) {
	decodedPublicKey, err := base64.StdEncoding.DecodeString(publicKey)
	if err != nil {
		return nil, fmt.Errorf("could not decode: %w", err)
	}

	key, err := jwt.ParseRSAPublicKeyFromPEM(decodedPublicKey)
	if err != nil {
		return nil, fmt.Errorf("validate: parse key: %w", err)
	}

	parsedToken, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected method: %s", t.Header["alg"])
		}
		return key, nil
	})
	if err != nil {
		return nil, fmt.Errorf("validate: %w", err)
	}

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok || !parsedToken.Valid {
		return nil, fmt.Errorf("validate: invalid token")
	}

	return claims["sub"], nil
}

// // GenerateToken using HMAC secret key
// func GenerateToken(ttl time.Duration, payload interface{}, secretJWTKey string) (string, error) {
// 	token := jwt.New(jwt.SigningMethodHS256)

// 	now := time.Now().UTC()
// 	claims := token.Claims.(jwt.MapClaims)

// 	claims["sub"] = payload
// 	claims["exp"] = now.Add(ttl).Unix()
// 	claims["iat"] = now.Unix()
// 	claims["nbf"] = now.Unix()

// 	tokenString, err := token.SignedString([]byte(secretJWTKey))
// 	if err != nil {
// 		return "", fmt.Errorf("generating JWT Token failed: %w", err)
// 	}

// 	return tokenString, nil
// }

// // ValidateTokenHMAC using HMAC secret key
// func ValidateTokenHMAC(token string, signedJWTKey string) (interface{}, error) {
// 	tok, err := jwt.Parse(token, func(jwtToken *jwt.Token) (interface{}, error) {
// 		if _, ok := jwtToken.Method.(*jwt.SigningMethodHMAC); !ok {
// 			return nil, fmt.Errorf("unexpected method: %s", jwtToken.Header["alg"])
// 		}

// 		return []byte(signedJWTKey), nil
// 	})
// 	if err != nil {
// 		return nil, fmt.Errorf("invalidate token, %w", err)
// 	}

// 	claims, ok := tok.Claims.(jwt.MapClaims)
// 	if !ok || !tok.Valid {
// 		return nil, fmt.Errorf("invalid token claim")
// 	}

// 	return claims["sub"], nil
// }
