// Package encrypt provides AES-GCM symmetric encryption and decryption
// for sensitive environment variable values managed by envchain.
//
// # Overview
//
// Encryptor wraps AES-GCM to offer authenticated encryption. Each call to
// Encrypt generates a fresh random nonce, so repeated encryption of the same
// plaintext yields different ciphertexts. The nonce is prepended to the
// ciphertext and the combined blob is base64-encoded for safe storage in
// JSON snapshots or YAML layer files.
//
// # Key Sizes
//
// AES-128 (16-byte key), AES-192 (24-byte key), and AES-256 (32-byte key)
// are all supported. Keys of any other length are rejected with ErrInvalidKey.
//
// # Usage
//
//	enc, err := encrypt.NewEncryptor([]byte(os.Getenv("ENVCHAIN_KEY")))
//	if err != nil { /* handle */ }
//
//	cipher, err := enc.Encrypt("my-secret")
//	plain,  err := enc.Decrypt(cipher)
package encrypt
