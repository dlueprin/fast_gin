package jwts

import (
	"errors"
	"fast_gin/global"
	"github.com/golang-jwt/jwt/v5"
	"github.com/sirupsen/logrus"
	"time"
)

// jwt缺点，不会主动过期，如果用户注销了，之前没过期的token还是能用的
// 应对方法是用黑名单策略，将注销的账号的token放入redis的一个黑名单中，用户请求时先判断token是否有效，然后判断是否在黑名单里面
// 还有双token机制，登录时生成两个token，一个存活时间久用于刷新：reflesh_token，一个存活时间短，用于验证：access_token,
// go get github.com/golang-jwt/jwt/v5

type Claims struct {
	UserID uint `json:"userID"`
	RoleID uint `json:"roleID"` //区分管理员和普通用户
}
type MyClaims struct { //千万不能存密码，因为jwt是编码不是加密
	Claims
	jwt.RegisteredClaims
}

func SetToken(data Claims) (string, error) {
	SetClaims := MyClaims{
		Claims: data,
		RegisteredClaims: jwt.RegisteredClaims{
			//newnumericdate:将go时间格式转换为jwt时间格式,使用全局变量的话要在前面加个time。duration
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(global.Config.Jwt.Expires) * time.Hour)), //有效时间
			Issuer:    global.Config.Jwt.Issuer,                                                                 //签发人

			//IssuedAt:  jwt.NewNumericDate(time.Now()),
			//NotBefore: jwt.NewNumericDate(time.Now()),
			//Issuer:    os.Getenv("JWT"),
		},
	}
	//使用指定的加密方式和签名创建的声明类型创建新令牌
	tokenStruct := jwt.NewWithClaims(jwt.SigningMethodHS256, SetClaims) //对称加密算法，适合单点服务，RS256非对称，适合分布式
	//获得完整的、签名的令牌
	token, err := tokenStruct.SignedString([]byte(global.Config.Jwt.Key)) //用签名密钥签名
	if err != nil {
		logrus.Errorf("颁发jwt失败%s", err)
		return "", err
	}
	return token, nil
}

// 验证token
func CheckToken(token string) (*MyClaims, error) {
	//用密钥解析token并获得
	tokenObj, err := jwt.ParseWithClaims(token, &MyClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(global.Config.Jwt.Key), nil //用签名密钥验证
	})
	if err != nil {
		logrus.Errorf("验证jwt失败:%s", err)
		return nil, err
	}
	logrus.Infoln("验证jwt成功，tokenObj:", tokenObj)
	//类型断言验证token是否有效（签名和过期时间）
	if claims, ok := tokenObj.Claims.(*MyClaims); ok && tokenObj.Valid {
		logrus.Infoln("token有效")
		return claims, nil
	} else {
		return nil, errors.New("token无效")
	}
}
