package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io"

	"golang.org/x/crypto/hkdf"

	"github.com/kcarretto/paragon/pkg/script"
)

var (
	NONCE_SIZE = 12
	KEY_ENTROPY_SIZE = 32
	KEY_SIZE = 16
)

// CreateKey creates a new Key object to be passed around given the value. Will also error if
// passed value is not valid base64/the key is not KEY_SIZE bytes.
func CreateKey(value string) (Key, error) {
	b, err := base64.StdEncoding.DecodeString(value)
	if err != nil {
		return Key{""}, err
	}
	if len(b) != KEY_SIZE {
		return Key{""}, fmt.Errorf("passed value for key is not %d bytes long", KEY_SIZE)
	}
	return Key{value}, nil
}

// GenerateKey creates a new Key object to be passed around.
//
//go:generate go run ../gendoc.go -lib crypto -func generateKey -retval key@Key -retval err@Error -doc "GenerateKey creates a new Key object to be passed around."
//
// @callable:	crypto.generateKey
// @retval:		key 	@Key
// @retval:		err 	@Error
//
// @usage:		k, err = crypto.generateKey()
func GenerateKey() (Key, error) {
	// KEY_ENTROPY_SIZE bytes of entropy should be gucci gang
	seed := make([]byte, KEY_ENTROPY_SIZE)
	_, err := rand.Read(seed)
	if err != nil {
		return Key{""}, err
	}
	// Generate 48 bytes of key material. (nil is for salt and info which we don't need)
	hkdf := hkdf.New(sha256.New, seed, nil, nil)

	// take first KEY_SIZE bytes
	key := make([]byte, KEY_SIZE)
    if _, err := io.ReadFull(hkdf, key); err != nil {
        return Key{""}, err
    }

	b64Key := Key{base64.StdEncoding.EncodeToString(key)}
	return b64Key, nil
}

func generateKey(parser script.ArgParser) (script.Retval, error) {
	return script.WithError(GenerateKey()), nil
}

// Encrypt takes a Key and some data and returns the AESGCM encrypted IV+ciphertext.
//
//go:generate go run ../gendoc.go -lib crypto -func encrypt -param key@Key -param data@String -retval ciphertext@String -retval err@Error -doc "Encrypt takes a Key and some data and returns the AESGCM encrypted IV+ciphertext."
//
// @callable:	crypto.encrypt
// @param:		key 		@Key
// @param:		data 		@String
// @retval:		ciphertext 	@String
// @retval:		err 		@Error
//
// @usage:		c, err = crypto.encrypt(key, data)
func Encrypt(key Key, data string) (string, error) {
	rawKey, err := base64.StdEncoding.DecodeString(key.String())
	if err != nil {
		return "", err
	}
	rawData := []byte(data)

	aesCipher, err := aes.NewCipher(rawKey)
	if err != nil {
		return "", err
	}

	// ALWAYS do 12 byte nonces with GCM o/w you will cry
	nonce := make([]byte, NONCE_SIZE)
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	aesgcm, err := cipher.NewGCM(aesCipher)
	if err != nil {
		return "", err
	}

	ciphertext := aesgcm.Seal(nil, nonce, rawData, nil)

	return string(append(nonce, ciphertext...)), nil

}

func encrypt(parser script.ArgParser) (script.Retval, error) {
	key, err := ParseParam(parser, 0)
	if err != nil {
		return nil, err
	}
	data, err := parser.GetString(1)
	if err != nil {
		return nil, err
	}
	return script.WithError(Encrypt(key, data)), nil
}


// Decrypt takes a Key and AESGCM encrypted IV+ciphertext data and returns the plaintext.
//
//go:generate go run ../gendoc.go -lib crypto -func decrypt -param key@Key -param data@String -retval plaintext@String -retval err@Error -doc "Decrypt takes a Key and AESGCM encrypted IV+ciphertext data and returns the plaintext."
//
// @callable:	crypto.decrypt
// @param:		key 		@Key
// @param:		data 		@String
// @retval:		plaintext 	@String
// @retval:		err 		@Error
//
// @usage:		p, err = crypto.decrypt(key, data)
func Decrypt(key Key, data string) (string, error) {
	rawKey, err := base64.StdEncoding.DecodeString(key.String())
	if err != nil {
		return "", err
	}
	rawData := []byte(data)
	nonce := rawData[:NONCE_SIZE]
	ciphertext := rawData[NONCE_SIZE:]

	aesCipher, err := aes.NewCipher(rawKey)
	if err != nil {
		return "", err
	}
	aesgcm, err := cipher.NewGCM(aesCipher)
	if err != nil {
		return "", err
	}
	plaintext, err := aesgcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", err
	}
	return string(plaintext), nil
}


func decrypt(parser script.ArgParser) (script.Retval, error) {
	key, err := ParseParam(parser, 0)
	if err != nil {
		return nil, err
	}
	data, err := parser.GetString(1)
	if err != nil {
		return nil, err
	}
	return script.WithError(Decrypt(key, data)), nil
}