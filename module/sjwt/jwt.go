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

// JwtSignBuffer is the extra days added to JWT ExpiresAt beyond the Redis TTL.
// This ensures Redis is the controlling factor for session expiry, not the JWT itself.
const JwtSignBuffer = 1

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
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24 * time.Duration(jwtExpire+JwtSignBuffer))),
		},
		UserID: userID,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(jwtKey))
}

func ParseToken(tokenString string) (*JwtClaims, error) {
	token, err := jwt.ParseWithClaims(
		tokenString,
		&JwtClaims{},
		func(token *jwt.Token) (any, error) {
			// 防御 alg confusion：拒绝任何非 HS256（包括 "none"）的签名。
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(jwtKey), nil
		},
	)
	if err != nil {
		return nil, err
	}
	if token == nil || !token.Valid {
		// 之前这里落到 "return nil, err" 会返回 (nil, nil)，上层靠 claims==nil 判断成功；
		// 改为显式错误，避免误判。
		return nil, fmt.Errorf("%s token is not valid", LogPrefix)
	}
	claims, ok := token.Claims.(*JwtClaims)
	if !ok {
		return nil, fmt.Errorf("%s unexpected claims type", LogPrefix)
	}
	return claims, nil
}

func MsgLog(msg string, params ...interface{}) string {
	if len(params) == 0 {
		return LogPrefix + msg
	}
	return LogPrefix + fmt.Sprintf(msg, params...)
}
