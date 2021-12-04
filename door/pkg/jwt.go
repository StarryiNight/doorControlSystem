package pkg

import (
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/pkg/errors"
)

const TokenExpireDuration = time.Hour * 2
const pass="door"

var mySecret = []byte("门禁系统")

type MyClaims struct {
	CardId   string  `json:"userId"`
	UserName string `json:"username"`
	jwt.StandardClaims
}

// GenToken 生成JWT
func GenToken(cardId string , userName string) (string, error) {
	// 创建一个我们自己的声明
	c := MyClaims{
		cardId,
		userName, // 自定义字段
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(TokenExpireDuration).Unix(), // 过期时间
			Issuer:    pass,                                       // 签发人
		},
	}
	// 使用指定的签名方法创建签名对象
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, c)
	// 使用指定的secret签名并获得完整的编码后的字符串token
	return token.SignedString(mySecret)
}

// ParseToken 解析JWT
func ParseToken(tokenString string) (*MyClaims, error) {
	// 解析token
	var mc = new(MyClaims)
	token, err := jwt.ParseWithClaims(tokenString, mc, func(token *jwt.Token) (i interface{}, err error) {
		return mySecret, nil
	})
	if err != nil {
		return nil, err
	}
	if token.Valid { // 校验token
		return mc, nil
	}
	return nil, errors.New("invalid token")
}
