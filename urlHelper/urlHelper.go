package urlHelper

import (
	"strings"
)

// @TODO: implement it =)
func GenerateUrl(path string) string {
	return "/" + strings.Replace(path, "/", "", 1)
}
