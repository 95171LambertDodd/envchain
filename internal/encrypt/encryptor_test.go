package encrypt_test

import (
	"strings"
	"testing"

	"github.com/yourorg/envchain/internal/encrypt"
)

func key16() []byte { return []byte("0123456789abcdef") }
func key32() []byte { return []byte("0123456789abcdef0123456789abcdef") }

func TestNewEncryptorInvalidKey(t *testing.T) {
	_, err := encrypt.NewEncryptor([]byte("short"))
	if err != encrypt.ErrInvalidKey {
		t.Fatalf("expected ErrInvalidKey, got %v", err)
	}
}

func TestNewEncryptorValidKey(t *testing.T) {
	for _, k := range [][]byte{key16(), key32()} {
		_, err := encrypt.NewEncryptor(k)
		if err != nil {
			t.Fatalf("unexpected error for key len %d: %v", len(k), err)
		}
	}
}

func TestEncryptDecryptRoundtrip(t *testing.T) {
	enc, _ := encrypt.NewEncryptor(key16())
	plaintext := "super-secret-value"
	cipher, err := enc.Encrypt(plaintext)
	if err != nil {
		t.Fatalf("Encrypt: %v", err)
	}
	if cipher == plaintext {
		t.Fatal("ciphertext should differ from plaintext")
	}
	got, err := enc.Decrypt(cipher)
	if err != nil {
		t.Fatalf("Decrypt: %v", err)
	}
	if got != plaintext {
		t.Fatalf("expected %q, got %q", plaintext, got)
	}
}

func TestEncryptProducesUniqueOutputs(t *testing.T) {
	enc, _ := encrypt.NewEncryptor(key32())
	a, _ := enc.Encrypt("value")
	b, _ := enc.Encrypt("value")
	if a == b {
		t.Fatal("two encryptions of the same value should differ due to random nonce")
	}
}

func TestDecryptInvalidBase64(t *testing.T) {
	enc, _ := encrypt.NewEncryptor(key16())
	_, err := enc.Decrypt("!!!not-base64!!!")
	if err != encrypt.ErrInvalidCiphertext {
		t.Fatalf("expected ErrInvalidCiphertext, got %v", err)
	}
}

func TestDecryptTamperedCiphertext(t *testing.T) {
	enc, _ := encrypt.NewEncryptor(key16())
	cipher, _ := enc.Encrypt("hello")
	// Flip last character to corrupt the MAC.
	tampered := cipher[:len(cipher)-1] + strings.Map(func(r rune) rune {
		if r == 'A' {
			return 'B'
		}
		return 'A'
	}, string(cipher[len(cipher)-1:]))
	_, err := enc.Decrypt(tampered)
	if err != encrypt.ErrInvalidCiphertext {
		t.Fatalf("expected ErrInvalidCiphertext for tampered input, got %v", err)
	}
}

func TestDecryptWrongKey(t *testing.T) {
	enc1, _ := encrypt.NewEncryptor(key16())
	enc2, _ := encrypt.NewEncryptor([]byte("fedcba9876543210"))
	cipher, _ := enc1.Encrypt("secret")
	_, err := enc2.Decrypt(cipher)
	if err != encrypt.ErrInvalidCiphertext {
		t.Fatalf("expected ErrInvalidCiphertext with wrong key, got %v", err)
	}
}
