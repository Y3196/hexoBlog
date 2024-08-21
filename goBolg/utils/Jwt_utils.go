package utils

import (
	"time"

	"github.com/dgrijalva/jwt-go"
)

// JWTSecret 密钥
const JWTSecret = "FelixBlog"

// TokenDuration 过期时间
const TokenDuration = time.Hour * 24 * 3

// Claims 自定义声明
type Claims struct {
	UserID   int    `json:"user_id"`
	Username string `json:"username"`
	jwt.StandardClaims
}

// GenerateJWT 生成包含用户信息和过期时间的 JWT
func GenerateJWT(userID int, username string) (string, error) {
	claims := &Claims{
		UserID:   userID,
		Username: username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(TokenDuration).Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(JWTSecret))
}

// ValidateJWT 验证 JWT 并返回其声明
func ValidateJWT(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(JWTSecret), nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}
	return nil, jwt.ErrSignatureInvalid // 返回具体的错误类型
}
