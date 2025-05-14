package utils

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	errors2 "github.com/pkg/errors"
	"microservices_course/week6/jwt/internal/model"
	"time"
)

func GenerateToken(info model.UserInfo, secretKey []byte, duration time.Duration) (string, error) {
	//Создаем данные, которые будут внутри токена
	claims := model.UserClaim{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(duration).Unix(),
		},
		Username: info.Username,
		Role:     info.Role,
	}

	// Создаем новый JWT-токен, используя алгоритм подписи ES256 и claims
	token := jwt.NewWithClaims(jwt.SigningMethodES256, claims)

	// Подписываем токен с помощью переданного секретного ключа и возвращаем его
	return token.SignedString(secretKey)
}

func VerifyToken(tokenStr string, secretKey []byte) (*model.UserClaim, error) {
	// Парсим токен и сразу указываем структуру claims, куда должны распаковаться данные.
	token, err := jwt.ParseWithClaims(tokenStr, &model.UserClaim{},
		func(token *jwt.Token) (interface{}, error) {
			// Проверяем, что используется ожидаемый метод подписи — HMAC.
			_, ok := token.Method.(*jwt.SigningMethodHMAC)
			if !ok {
				return nil, errors.New("unexpected token signing method")
			}
			// Возвращаем секретный ключ, который используется для проверки подписи токена.
			return secretKey, nil
		},
	)
	if err != nil {
		return nil, errors2.Errorf("invalid token: %s", err.Error())
	}
	// Приводим claims из токена к ожидаемому типу *model.UserClaim
	claims, ok := token.Claims.(*model.UserClaim)
	if !ok {
		return nil, errors.New("invalid token claims")
	}
	return claims, nil
}
