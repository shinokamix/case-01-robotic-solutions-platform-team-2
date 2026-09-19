package auth

import (
	"strings"
	"unicode"
	"unicode/utf8"

	"golang.org/x/net/idna"
)

const (
	minPasswordLength = 12
	maxPasswordLength = 128
	maxNameLength     = 100
	maxEmailLength    = 254
)

func CanonicalEmail(value string) (string, error) {
	value = strings.TrimSpace(value)
	if value == "" || strings.IndexFunc(value, unicode.IsSpace) >= 0 {
		return "", ErrInvalidEmail
	}

	local, domain, ok := strings.Cut(value, "@")
	if !ok || local == "" || domain == "" || strings.Contains(domain, "@") {
		return "", ErrInvalidEmail
	}
	if local[0] == '.' || local[len(local)-1] == '.' || strings.Contains(local, "..") {
		return "", ErrInvalidEmail
	}
	for _, char := range local {
		if !validLocalPartChar(char) {
			return "", ErrInvalidEmail
		}
	}

	asciiDomain, err := idna.Lookup.ToASCII(domain)
	if err != nil || asciiDomain == "" || strings.HasPrefix(asciiDomain, ".") || strings.HasSuffix(asciiDomain, ".") {
		return "", ErrInvalidEmail
	}
	for _, label := range strings.Split(asciiDomain, ".") {
		if label == "" || label[0] == '-' || label[len(label)-1] == '-' {
			return "", ErrInvalidEmail
		}
	}

	canonical := strings.ToLower(local + "@" + asciiDomain)
	if len(canonical) > maxEmailLength {
		return "", ErrInvalidEmail
	}
	return canonical, nil
}

func validLocalPartChar(char rune) bool {
	if char >= 'a' && char <= 'z' || char >= 'A' && char <= 'Z' || char >= '0' && char <= '9' {
		return true
	}
	return strings.ContainsRune(".!#$%&'*+/=?^_`{|}~-", char)
}

func ValidatePassword(password string) error {
	length := utf8.RuneCountInString(password)
	if length < minPasswordLength || length > maxPasswordLength {
		return ErrInvalidPassword
	}
	return nil
}

func NormalizeFirstName(name string) (string, error) {
	return normalizeName(name, ErrInvalidFirstName)
}

func NormalizeLastName(name string) (string, error) {
	return normalizeName(name, ErrInvalidLastName)
}

func normalizeName(name string, validationError error) (string, error) {
	name = strings.TrimSpace(name)
	length := utf8.RuneCountInString(name)
	if length == 0 || length > maxNameLength || strings.IndexFunc(name, unicode.IsControl) >= 0 {
		return "", validationError
	}
	return name, nil
}
