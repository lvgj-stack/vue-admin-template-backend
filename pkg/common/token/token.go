package token

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/Mr-LvGJ/gobase/pkg/common/errno"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

type Config struct {
	key         string
	identityKey string
}

var (
	config = Config{"Rtg8BPKNEf2mB4mgvKONGPZZQSaJWNLijxR42qRgq0iBb5", "identityKey"}
	once   sync.Once
)

func Sign(identityKey string) (tokenString string, err error) {

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		config.identityKey: identityKey,
		"nbf":              time.Now().Unix(),
		"iat":              time.Now().Unix(),
		"exp":              time.Now().Add(time.Hour).Unix(),
	})

	signedString, err := token.SignedString([]byte(config.key))

	return signedString, err

}

func ParseRequest(c *gin.Context) (string, error) {
	header := c.Request.Header.Get("Authorization")
	if len(header) == 0 {
		return "", errors.New("the length of the `Authorization` header is zero")
	}

	var t string

	fmt.Sscanf(header, "Bearer %s", &t)
	return Parse(t, config.key)
}

func secretFunc(secret string) jwt.Keyfunc {
	return func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return []byte(secret), nil
	}

}

func Parse(tokenString string, secret string) (string, error) {
	token, err := jwt.Parse(tokenString, secretFunc(secret))
	var identityKey string

	if err != nil {
		return "", err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		identityKey = claims[config.identityKey].(string)
	} else {
		return "", errno.ErrToken
	}
	return identityKey, nil

}

func Init(key string, identityKey string) {
	once.Do(func() {
		if key != "" {
			config.key = key
		}
		if identityKey != "" {
			config.identityKey = identityKey
		}
	})
}
