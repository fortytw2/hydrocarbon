package hydrocarbon

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"
)

// A KeySigner is used to verify the integrity of SessionKeys at the system borders
type KeySigner struct {
	key []byte
}

// NewKeySigner creates a new KeySigner with the given key
func NewKeySigner(key string) *KeySigner {
	return &KeySigner{
		key: []byte(key),
	}
}

// Sign appends an HMAC to a value
func (ks *KeySigner) Sign(val string) (string, error) {
	h := hmac.New(sha256.New, ks.key)
	_, err := h.Write([]byte(val))
	if err != nil {
		return "", nil
	}

	return fmt.Sprintf("%s:%s", val, base64.StdEncoding.EncodeToString(h.Sum(nil))), nil
}

// Verify checks a value signed with Sign
func (ks *KeySigner) Verify(pubVal string) (string, error) {
	h := hmac.New(sha256.New, ks.key)

	spl := strings.Split(pubVal, ":")
	if len(spl) != 2 {
		return "", errors.New("invalid token")
	}

	_, err := h.Write([]byte(spl[0]))
	if err != nil {
		return "", nil
	}

	hmacBytes, err := base64.StdEncoding.DecodeString(spl[1])
	if err != nil {
		return "", err
	}

	if !hmac.Equal(hmacBytes, h.Sum(nil)) {
		return "", errors.New("invalid signature")
	}

	return spl[0], nil
}
