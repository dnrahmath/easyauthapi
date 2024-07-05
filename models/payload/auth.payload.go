package payload

import "github.com/golang-jwt/jwt/v4"

//===================================================================

type LoginAndRegis struct {
	Value    string `form:"value" json:"value" bson:"value"`
	Password string `form:"password" json:"password" bson:"password"`
}

type UserPostPut struct {
	Username    string   `form:"username" json:"username" bson:"username"`
	Password    string   `form:"password" json:"password" bson:"password"`
	Email       string   `form:"email" json:"email" bson:"email"`
	PhoneNumber string   `form:"phonenumber" json:"phonenumber" bson:"phonenumber"`
	Gender      string   `form:"gender" json:"gender" bson:"gender"`
	Name        string   `form:"name" json:"name" bson:"name"`
	Noid        string   `form:"noid" json:"noid" bson:"noid"`
	Religion    string   `form:"religion" json:"religion" bson:"religion"`
	Role        []string `form:"role" json:"role" bson:"role"`
}

//===================================================================

// jwt
type UserClaims struct {
	jwt.RegisteredClaims
	User *User  `form:"user" json:"user"`
	Type string `form:"type" json:"type"`
}

//===================================================================
