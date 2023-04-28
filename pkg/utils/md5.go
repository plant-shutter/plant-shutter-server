package utils

import (
	"crypto/md5"
	"fmt"

	"github.com/gabriel-vasile/mimetype"
)

func GetFileName(data []byte) string {
	return fmt.Sprintf("%x%s", md5.Sum(data), mimetype.Detect(data).Extension())
}
