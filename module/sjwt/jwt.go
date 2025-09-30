package sjwt

import (
	"context"
	"fmt"
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/golang-jwt/jwt/v5"
)

type JwtClaims struct {
	jwt.RegisteredClaims
	UserID int `json:"user_id"`
}

const LogPrefix = "[SUASK-JWT]"

var jwtKey = g.Cfg().MustGet(context.TODO(), "jwt.signKey").String()
var jwtExpire = g.Cfg().MustGet(context.TODO(), "jwt.expire").Int64()

func GetKey() []byte {
	return []byte(jwtKey)
}

func GetExpireDay() int64 {
	return jwtExpire
}

func GetExpireSecond() int64 {
	return jwtExpire * 60 * 60 * 24
}

func GenerateToken(userID int) (string, error) {
	claims := JwtClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "SuAsk",
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24 * time.Duration(jwtExpire))),
		},
		UserID: userID,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(jwtKey))
}

func ParseToken(tokenString string) (*JwtClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JwtClaims{}, func(token *jwt.Token) (any, error) {
		return []byte(jwtKey), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*JwtClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, err
}

func MsgLog(msg string, params ...interface{}) string {
	if len(params) == 0 {
		return LogPrefix + msg
	}
	return LogPrefix + fmt.Sprintf(msg, params...)
}
