package utils

import (
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
)

func Join(elem ...string) string {
	return strings.ReplaceAll(filepath.Join(elem...), "\\", "/")
}

func RandFilename(ext string) string {
	filename := fmt.Sprintf(
		"%s_%v",
		strings.ReplaceAll(uuid.NewString()[:6], "-", ""),
		time.Now().UnixMilli(),
	)
	if ext != "" {
		filename += fmt.Sprintf(".%s", ext)
	}
	return filename
}
