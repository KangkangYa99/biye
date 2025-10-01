package utils

import (
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
