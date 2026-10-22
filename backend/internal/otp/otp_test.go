package otp

import "testing"

func TestGenerate(t *testing.T) {
	code, err := Generate()
	if err != nil {
		t.Fatal(err)
	}
	if len(code) != 6 {
		t.Errorf("len(code) = %d, want 6", len(code))
	}
	for _, c := range code {
		if c < '0' || c > '9' {
			t.Errorf("code %q contains non-digit %q", code, c)
		}
	}
}

func TestMatches(t *testing.T) {
	secret := []byte("test-secret")
	hash := Hash(secret, "123456")

	if !Matches(secret, "123456", hash) {
		t.Error("correct code did not match")
	}
	if Matches(secret, "654321", hash) {
		t.Error("wrong code matched")
	}
	if Matches([]byte("other-secret"), "123456", hash) {
		t.Error("code matched under a different secret")
	}
}
