package tools

import (
	"regexp"
	"strings"
)

const (
	regUrlPattern = `^(https?|ftp)://[^\s/$.?#].[^\s]*$`
	webPrefix     = "www."
)

var regUrl = regexp.MustCompile(regUrlPattern)

func MatchUrl(url string) bool {
	return regUrl.MatchString(url)
}

func StripWebPrefix(url string) string {
	return strings.Replace(url, webPrefix, "", -1)
}
