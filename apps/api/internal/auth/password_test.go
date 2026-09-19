package auth

import (
	"strings"
	"testing"
)

func TestPasswordHasher(t *testing.T) {
	hasher := PasswordHasher{}
	encoded, err := hasher.Hash("correct horse battery staple")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.HasPrefix(encoded, "$argon2id$v=19$m=65536,t=3,p=1$") {
		t.Fatalf("unexpected encoded hash: %s", encoded)
	}
	if !hasher.Compare(encoded, "correct horse battery staple") {
		t.Fatal("correct password did not match")
	}
	if hasher.Compare(encoded, "incorrect password") {
		t.Fatal("incorrect password matched")
	}
}
