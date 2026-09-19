package auth

import (
	"errors"
	"strings"
	"testing"
)

func TestCanonicalEmail(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name  string
		input string
		want  string
		err   error
	}{
		{name: "trim and lowercase", input: " USER@EXAMPLE.COM ", want: "user@example.com"},
		{name: "keep plus suffix", input: "User+tag@Example.com", want: "user+tag@example.com"},
		{name: "IDNA domain", input: "user@пример.рф", want: "user@xn--e1afmkfd.xn--p1ai"},
		{name: "consecutive dots", input: "user..name@example.com", err: ErrInvalidEmail},
		{name: "leading dot", input: ".user@example.com", err: ErrInvalidEmail},
		{name: "internal space", input: "user name@example.com", err: ErrInvalidEmail},
		{name: "unicode local part", input: "юзер@example.com", err: ErrInvalidEmail},
		{name: "label starts with hyphen", input: "user@-example.com", err: ErrInvalidEmail},
		{name: "too long", input: strings.Repeat("a", 250) + "@x.dev", err: ErrInvalidEmail},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			got, err := CanonicalEmail(test.input)
			if !errors.Is(err, test.err) {
				t.Fatalf("CanonicalEmail() error = %v, want %v", err, test.err)
			}
			if got != test.want {
				t.Fatalf("CanonicalEmail() = %q, want %q", got, test.want)
			}
		})
	}
}

func TestNormalizeNames(t *testing.T) {
	t.Parallel()
	if got, err := NormalizeFirstName("  Иван  "); err != nil || got != "Иван" {
		t.Fatalf("NormalizeFirstName() = %q, %v", got, err)
	}
	if _, err := NormalizeFirstName("\n"); !errors.Is(err, ErrInvalidFirstName) {
		t.Fatalf("empty first name error = %v", err)
	}
	if _, err := NormalizeLastName(strings.Repeat("я", 101)); !errors.Is(err, ErrInvalidLastName) {
		t.Fatalf("long last name error = %v", err)
	}
}

func TestValidatePasswordCountsCharacters(t *testing.T) {
	t.Parallel()
	if err := ValidatePassword(strings.Repeat("я", 12)); err != nil {
		t.Fatalf("12 character password: %v", err)
	}
	if !errors.Is(ValidatePassword(strings.Repeat("a", 11)), ErrInvalidPassword) {
		t.Fatal("11 character password must be invalid")
	}
	if !errors.Is(ValidatePassword(strings.Repeat("a", 129)), ErrInvalidPassword) {
		t.Fatal("129 character password must be invalid")
	}
}
