package utils

import (
	"strings"
)

func GenerateSlug(name string) string {
	return strings.ToLower(strings.Join(strings.Fields(name), "-"))
}
