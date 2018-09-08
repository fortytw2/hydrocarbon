package pg

import "testing"

func TestCompression(t *testing.T) {
	var text = `awiojposen&V9r800wenvuasnu cvopaS*N()ea-8dfv9asuy*(_DVN-`

	out, err := compressText(text)
	if err != nil {
		t.Fatal(err)
	}

	dec, err := decompressText(out)
	if err != nil {
		t.Fatal(err)
	}

	if dec != text {
		t.Fatal("did not get back the same thing after decompressing")
	}
}
