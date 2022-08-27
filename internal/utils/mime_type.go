package utils

import (
	"strings"

	"github.com/pkg/errors"
)

var (
	errNotFoundImageType = errors.New("mimeType: not found mime type")
)

var signatures map[string]string

func init() {
	signatures = map[string]string{
		"JVBERi0":     "image/gif",
		"R0lGODdh":    "image/gif",
		"R0lGODlh":    "image/gif",
		"iVBORw0KGgo": "image/png",
		"/9j/":        "image/jpg",
	}
}

func DetectMimeTypeFromBase64(b64 string) (string, error) {
	for k, s := range signatures {
		if strings.Contains(b64, k) {
			return s, nil
		}
	}

	return "", errNotFoundImageType
}

func IsValidMimeType(mimeType string) bool {
	switch mimeType {
	case "image/jpeg", "image/jpg", "image/png", "image/gif":
		return true
	}

	return false
}
