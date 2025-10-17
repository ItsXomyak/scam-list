package utils

import (
	"errors"
	"net/url"
	"regexp"
	"strings"
)

var schemeRe = regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9+.-]*://`)

// ExtractDomain принимает URL или голый домен и возвращает только host (без порта).
// Примеры:
//   - "https://chatgpt.com/g/..." -> "chatgpt.com"
//   - "chatgpt.com"               -> "chatgpt.com"
//   - "//example.com/path"        -> "example.com"
func ExtractDomain(raw string) (string, error) {
	s := strings.TrimSpace(raw)
	if s == "" {
		return "", errors.New("empty input")
	}

	// Поддержка protocol-relative URL: //example.com/path
	if strings.HasPrefix(s, "//") {
		s = "http:" + s
	} else if !schemeRe.MatchString(s) {
		// Если схемы нет — подставим http://
		s = "http://" + s
	}

	u, err := url.Parse(s)
	if err != nil {
		return "", err
	}

	host := u.Hostname()
	if host == "" {
		return "", errors.New("invalid URL or host")
	}

	// Нормализация: нижний регистр и убрать финальную точку.
	host = strings.ToLower(strings.TrimSuffix(host, "."))

	return host, nil
}
