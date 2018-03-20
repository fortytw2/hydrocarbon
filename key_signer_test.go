package hydrocarbon

import (
	"errors"
	"testing"
)

func TestKeySigner(t *testing.T) {
	var cases = []struct {
		name string
		run  func(k *KeySigner) error
	}{
		{
			"basic",
			func(ks *KeySigner) error {
				ogKey := "12345678"
				val, err := ks.Sign(ogKey)
				if err != nil {
					return err
				}

				key2, err := ks.Verify(val)
				if err != nil {
					return err
				}

				if key2 != ogKey {
					return errors.New("keys did not match")
				}

				return nil
			},
		},
	}

	ks := NewKeySigner("test-signing-key")
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.run(ks)
			if err != nil {
				t.Fatal(err)
			}

		})
	}
}
