// Code generated by go-bindata.
// sources:
// schema/01_init.sql
// DO NOT EDIT!

package hydrocarbon

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func bindataRead(data []byte, name string) ([]byte, error) {
	gz, err := gzip.NewReader(bytes.NewBuffer(data))
	if err != nil {
		return nil, fmt.Errorf("Read %q: %v", name, err)
	}

	var buf bytes.Buffer
	_, err = io.Copy(&buf, gz)
	clErr := gz.Close()

	if err != nil {
		return nil, fmt.Errorf("Read %q: %v", name, err)
	}
	if clErr != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

type asset struct {
	bytes []byte
	info  os.FileInfo
}

type bindataFileInfo struct {
	name    string
	size    int64
	mode    os.FileMode
	modTime time.Time
}

func (fi bindataFileInfo) Name() string {
	return fi.name
}
func (fi bindataFileInfo) Size() int64 {
	return fi.size
}
func (fi bindataFileInfo) Mode() os.FileMode {
	return fi.mode
}
func (fi bindataFileInfo) ModTime() time.Time {
	return fi.modTime
}
func (fi bindataFileInfo) IsDir() bool {
	return false
}
func (fi bindataFileInfo) Sys() interface{} {
	return nil
}

var _schema01_initSQL = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xcc\x96\x4f\x73\xda\x3a\x14\xc5\xd7\xf6\xa7\xb8\x3b\x60\x1e\x79\x93\x97\x79\xf3\x36\x59\x11\x50\xe6\xb9\x25\x26\xf5\x9f\x0e\xe9\xc6\x23\xac\x0b\xa8\x18\xc9\xb5\xe4\x04\xbe\x7d\x47\xc6\x22\x36\x90\x36\xb4\x9d\x4c\x56\xcc\x88\x73\xad\xab\xdf\x39\xbe\xf2\x30\x20\x83\x88\x00\x99\x46\xc4\x0f\xbd\x89\x0f\xde\x2d\xf8\x93\x08\xc8\xd4\x0b\xa3\x10\xf2\x45\x5a\x6c\x73\x2d\xaf\xdd\x9f\x08\x53\xae\x71\xa3\xaf\x5d\xab\x8b\x06\x37\x63\x02\xa5\xc2\x42\x41\xd7\x75\x38\x83\x38\xf6\x46\x70\x1f\x78\x77\x83\xe0\x01\x3e\x92\x07\x18\x91\xdb\x41\x3c\x8e\x60\x81\x22\x29\xa8\x60\x72\x9d\x94\x25\x67\xdd\x5e\xdf\x75\x9d\xb4\x40\xaa\x91\x25\x54\x43\xe4\xdd\x91\x30\x1a\xdc\xdd\x47\x5f\xaa\x1d\xfd\x78\x3c\xde\x17\x0b\xf9\x64\x0a\x9c\x32\x67\xe7\xe8\x5d\x07\xd7\x94\x67\x30\xf4\x22\x32\x8d\xf6\x32\xf3\x47\xec\x7b\x9f\x62\x02\xdd\x4a\xd0\x73\x7b\xd7\xae\x7b\x71\x01\x99\x5c\x70\x01\x5a\xae\x50\x28\xa0\x05\x82\x14\x78\xa1\xf9\x1a\xed\x5a\xa9\x90\x81\x96\xb5\x90\xcf\x41\xd2\x52\x2f\xe1\xea\xef\x4b\xe0\x0a\x84\xd4\x95\xa2\x8d\xa7\xd2\x26\xf5\x03\xce\xa7\xe4\x18\xbc\x89\x2d\x0a\xc8\x2d\x09\x88\x3f\x24\x61\x8d\xbd\x79\xa8\x73\x71\xe2\x26\xe7\x05\xaa\xd7\xe9\xe1\x2f\xf0\xfc\x88\x04\x9f\x07\x63\xe8\x5c\xfd\x0b\xff\x4f\xe2\x20\xec\x98\x6d\xab\x06\xe9\x02\x85\x86\x03\xce\x0e\xcf\x61\xe8\x8d\x82\x56\x97\x15\x89\x9d\xd2\x6e\x80\x22\x95\x0c\xbb\x8d\xd3\xcf\xb6\x1a\x55\xf7\x9f\xff\x7a\x7d\xe8\x2c\x71\xd3\xa9\x41\x30\xb8\x99\x4c\xc6\x64\xe0\x1f\x37\x39\xa7\x99\xc2\xca\xc8\x9a\x7e\x6d\xb1\xe7\x8f\xc8\xb4\x65\xc2\xee\x27\xe1\x6c\x03\x13\xff\xc0\x9e\xea\xb7\x4e\x83\x42\xa5\xb8\x6c\xe5\xc1\xfa\x6f\x5c\x47\xa1\x79\x4a\x35\x5a\x9d\x6a\xdb\x6e\x57\xdf\x97\xe5\x67\x7a\xb5\xc2\xad\x7d\x7b\xce\xf3\x8a\xa6\x9a\x3f\xe2\xcb\x6e\xe9\xa2\x7c\xd9\x2c\x8b\x2e\x59\xe1\xd6\xda\xf4\x8c\x73\x85\xdb\xda\xa0\xb9\xcc\x98\x01\xd2\x74\x66\x4d\xb9\xd0\x94\x0b\x48\x65\x96\x61\xaa\xab\x1a\x39\x87\x39\x22\x3b\x70\xc8\x96\xff\x41\x83\xe2\x90\x04\x61\x13\xea\x1b\xcc\x38\x41\xd7\xd8\x36\x73\xaf\xe9\x30\x9c\xd3\x32\xd3\x9d\xe6\xc8\xab\x7b\xef\x83\x29\xdc\x8f\xbe\x8a\x4f\x45\x92\x0b\xc6\x1f\x39\x2b\x69\x76\x12\x5a\xa5\x7b\x87\xc3\xde\xc9\xa8\xd2\x09\x8a\x6f\x25\x96\x47\x55\x66\xff\x3c\x2b\xcd\xd0\x3e\x8c\x7d\x59\x64\x47\x6b\x9a\xeb\x0c\x0f\x57\x9f\x09\xee\x9e\xd4\x87\xb2\xc8\x5a\xfc\x12\x9b\x28\xae\x80\xc2\x57\x69\xee\x12\x3a\xcb\x10\x66\xa8\x9f\x10\x85\x85\x2c\x98\xcd\xde\x31\xdb\xa4\x91\xca\xd7\x8e\x01\x67\x57\x73\x4a\x6a\x9f\xd6\x14\x9b\x5d\x4e\x49\xab\xe6\x7e\x6b\xb8\x9c\x1b\xdd\xbc\xe0\xb2\xe0\x7a\x6b\x6e\x96\x63\xd9\xa5\x91\x34\xe3\xf5\x1c\xdd\xfd\x81\xfb\x50\x1f\xa7\xd7\x9c\x26\x3b\x9c\xb9\x54\xfa\x97\xde\xee\x1f\x12\x7a\x13\x30\xa9\x14\x1a\x85\x4e\x96\x54\x2d\x8f\x3f\x5f\x4e\xe7\xd3\x31\x97\x92\x2c\x5e\x1a\x05\x9d\xbe\xeb\xcc\x24\xdb\xbe\xe6\x05\x30\x9f\x06\xba\xa0\xf0\x21\x9c\xf8\x37\xcd\xe4\x37\xfb\xaa\x88\x7f\x0f\x00\x00\xff\xff\xbd\x55\x0e\x68\x5a\x0a\x00\x00")

func schema01_initSQLBytes() ([]byte, error) {
	return bindataRead(
		_schema01_initSQL,
		"schema/01_init.sql",
	)
}

func schema01_initSQL() (*asset, error) {
	bytes, err := schema01_initSQLBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "schema/01_init.sql", size: 2650, mode: os.FileMode(420), modTime: time.Unix(499137600, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

// Asset loads and returns the asset for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func Asset(name string) ([]byte, error) {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[cannonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("Asset %s can't read by error: %v", name, err)
		}
		return a.bytes, nil
	}
	return nil, fmt.Errorf("Asset %s not found", name)
}

// MustAsset is like Asset but panics when Asset would return an error.
// It simplifies safe initialization of global variables.
func MustAsset(name string) []byte {
	a, err := Asset(name)
	if err != nil {
		panic("asset: Asset(" + name + "): " + err.Error())
	}

	return a
}

// AssetInfo loads and returns the asset info for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func AssetInfo(name string) (os.FileInfo, error) {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[cannonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("AssetInfo %s can't read by error: %v", name, err)
		}
		return a.info, nil
	}
	return nil, fmt.Errorf("AssetInfo %s not found", name)
}

// AssetNames returns the names of the assets.
func AssetNames() []string {
	names := make([]string, 0, len(_bindata))
	for name := range _bindata {
		names = append(names, name)
	}
	return names
}

// _bindata is a table, holding each asset generator, mapped to its name.
var _bindata = map[string]func() (*asset, error){
	"schema/01_init.sql": schema01_initSQL,
}

// AssetDir returns the file names below a certain
// directory embedded in the file by go-bindata.
// For example if you run go-bindata on data/... and data contains the
// following hierarchy:
//     data/
//       foo.txt
//       img/
//         a.png
//         b.png
// then AssetDir("data") would return []string{"foo.txt", "img"}
// AssetDir("data/img") would return []string{"a.png", "b.png"}
// AssetDir("foo.txt") and AssetDir("notexist") would return an error
// AssetDir("") will return []string{"data"}.
func AssetDir(name string) ([]string, error) {
	node := _bintree
	if len(name) != 0 {
		cannonicalName := strings.Replace(name, "\\", "/", -1)
		pathList := strings.Split(cannonicalName, "/")
		for _, p := range pathList {
			node = node.Children[p]
			if node == nil {
				return nil, fmt.Errorf("Asset %s not found", name)
			}
		}
	}
	if node.Func != nil {
		return nil, fmt.Errorf("Asset %s not found", name)
	}
	rv := make([]string, 0, len(node.Children))
	for childName := range node.Children {
		rv = append(rv, childName)
	}
	return rv, nil
}

type bintree struct {
	Func     func() (*asset, error)
	Children map[string]*bintree
}

var _bintree = &bintree{nil, map[string]*bintree{
	"schema": {nil, map[string]*bintree{
		"01_init.sql": {schema01_initSQL, map[string]*bintree{}},
	}},
}}

// RestoreAsset restores an asset under the given directory
func RestoreAsset(dir, name string) error {
	data, err := Asset(name)
	if err != nil {
		return err
	}
	info, err := AssetInfo(name)
	if err != nil {
		return err
	}
	err = os.MkdirAll(_filePath(dir, filepath.Dir(name)), os.FileMode(0755))
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(_filePath(dir, name), data, info.Mode())
	if err != nil {
		return err
	}
	err = os.Chtimes(_filePath(dir, name), info.ModTime(), info.ModTime())
	if err != nil {
		return err
	}
	return nil
}

// RestoreAssets restores an asset under the given directory recursively
func RestoreAssets(dir, name string) error {
	children, err := AssetDir(name)
	// File
	if err != nil {
		return RestoreAsset(dir, name)
	}
	// Dir
	for _, child := range children {
		err = RestoreAssets(dir, filepath.Join(name, child))
		if err != nil {
			return err
		}
	}
	return nil
}

func _filePath(dir, name string) string {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	return filepath.Join(append([]string{dir}, strings.Split(cannonicalName, "/")...)...)
}
