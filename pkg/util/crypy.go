package util

import "golang.org/x/crypto/bcrypt"

// HashPassword 函数将密码进行哈希处理
func HashPassword(password string) (string, error) {
	// 使用bcrypt生成密码的哈希值
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password),
		bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

// CheckPasswordHash 函数验证输入的密码是否正确
func CheckPasswordHash(password, hashedPassword string) bool {
	// 使用bcrypt比较密码和哈希值
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}
