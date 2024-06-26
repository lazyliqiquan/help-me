package middlewares

import (
	"context"
	"errors"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"github.com/lazyliqiquan/help-me/models"
)

type AuthType int

var (
	TokenPrivateKey []byte //token加密私钥
)

// UserClaims 用来生成token的结构，主要用来鉴权
type UserClaims struct {
	Id int //只有用户id是不会变的,而权限可能发生改变，为获取到最新的权限，可以调用一次数据库查询
	jwt.StandardClaims
}

func init() {
	// 生成一个随机字符串，用来当作私钥，这样每次重启后的私钥都是不一样的，
	// 这样也会导致之前的token全部失效
	TokenPrivateKey = []byte(uuid.New().String())
}

// GenerateToken 好像每次程序重新启动，之前的token全部失效了
// 生成 token
func GenerateToken(id int) (string, error) {
	tokenDuration, err := models.RDB.Get(context.Background(), "tokenDuration").Int()
	if err != nil {
		return "", err
	}
	UserClaim := &UserClaims{
		Id: id,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Duration(tokenDuration * int(time.Hour))).Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, UserClaim)
	tokenString, err := token.SignedString(TokenPrivateKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

// AnalyseToken 解析 token
func AnalyseToken(tokenString string) (*UserClaims, error) {
	userClaim := new(UserClaims)
	claims, err := jwt.ParseWithClaims(tokenString, userClaim,
		func(token *jwt.Token) (interface{}, error) {
			return TokenPrivateKey, nil
		})
	if err != nil {
		return nil, err
	}
	if !claims.Valid {
		return nil, errors.New("analyse token fail")
	}
	return userClaim, nil
}
