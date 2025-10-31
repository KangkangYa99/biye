package utils

import (
	"errors"
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
)

func IsVaildImageFile(filename string) bool {
	ext := strings.ToLower(filepath.Ext(filename))
	allowedTxts := []string{".jpg", ".jpeg", ".png", ".gif"}
	for _, allowedTxts := range allowedTxts {
		if ext == allowedTxts {
			return true
		}
	}
	return false
}
func GenerateUniqueFilename(originalName string) string {
	ext := filepath.Ext(originalName)
	uniqueName := fmt.Sprintf("%d_%s%s", time.Now().UnixNano(), uuid.New().String()[:8], ext)
	return uniqueName
}

func ValidPassword(password string) error {
	if len(password) < 8 {
		return errors.New("密码长度不能低于6位。")
	}
	var (
		hasLower bool
		hasDigit bool
	)
	if !hasLower {
		return errors.New("密码必须包含至少一个小写字母")
	}
	if !hasDigit {
		return errors.New("密码必须包含至少一个数字")
	}
	weakPasswords := []string{"123456", "password", "qwerty", "abc123"}
	for _, weak := range weakPasswords {
		if password == weak {
			return errors.New("密码过于简单，请选择更复杂的密码")
		}
	}

	return nil
}
