package lib

import (
	"bytes"
	"testing"
)

func TestEncryptAndDecrypt(t *testing.T) {
	cleartext := []byte("The quick brown fox jumps over the lazy dog.")
	password := []byte("mysuperdupersecurepassword")
	salt := NewNonce()
	key := NewKey(salt, password)
	nonce := NewNonce()
	ciphertext, err := Encrypt(cleartext, nonce, key)
	if err != nil {
		t.Fatal(err)
	}

	if bytes.Equal(ciphertext, cleartext) {
		t.Fatal("Ciphertext and cleartext are the same.")
	}

	decrypted, err := Decrypt(ciphertext, nonce, key)

	if err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal(cleartext, decrypted) {
		t.Fatal("Original cleartext and decrypted version are different.")
	}
}
