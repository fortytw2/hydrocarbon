package pg

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"io/ioutil"
	"strings"
)

const compressionPrefix = "gzip_"

// compressText gzips text
func compressText(in string) (string, error) {
	var buf bytes.Buffer

	zw, err := gzip.NewWriterLevel(&buf, 5)
	if err != nil {
		return "", err
	}

	_, err = zw.Write([]byte(in))
	if err != nil {
		return "", err
	}

	err = zw.Flush()
	if err != nil {
		return "", err
	}

	err = zw.Close()
	if err != nil {
		return "", err
	}

	return compressionPrefix + base64.StdEncoding.EncodeToString(buf.Bytes()), nil
}

func decompressText(in string) (string, error) {
	if !strings.HasPrefix(in, compressionPrefix) {
		return in, nil
	}

	decoded, err := base64.StdEncoding.DecodeString(strings.TrimPrefix(in, compressionPrefix))
	if err != nil {
		return "", err
	}

	zr, err := gzip.NewReader(bytes.NewReader(decoded))
	if err != nil {
		return "", err
	}

	decomp, err := ioutil.ReadAll(zr)
	if err != nil {
		return "", err
	}

	err = zr.Close()
	if err != nil {
		return "", err
	}

	return string(decomp), nil
}
