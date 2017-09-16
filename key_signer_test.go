package hydrocarbon

import (
	"testing"
)

func TestKeySigner(t *testing.T) {
	ks := NewKeySigner("test-signing-key")

	ogKey := "12345678"
	val, err := ks.Sign(ogKey)
	if err != nil {
		t.Fatal(err)
	}

	key2, err := ks.Verify(val)
	if err != nil {
		t.Fatal(err)
	}

	if key2 != ogKey {
		t.Fatal("did not get back original key")
	}
}
