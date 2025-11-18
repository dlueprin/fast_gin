package pwd

import (
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

// GenerateFromPassword 加密密码，需要明文密码，和加密花费值，用于密码注册
func GenerateFromPassword(password string) string {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		logrus.Errorf("密码加密错误 %s", err)
		return ""
	}
	return string(hashedPassword)
}

// CompareHashAndPassword 校验密码，需要明文密码和加密后的密码，用于登录验证
func CompareHashAndPassword(hashedPassword string, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}
