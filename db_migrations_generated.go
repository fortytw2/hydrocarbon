// Code generated by go-bindata.
// sources:
// schema/init.sql
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

var _schemaInitSql = []byte("\x1f\x8b\x08\x00\x00\x09\x6e\x88\x00\xff\xcc\x96\xcf\x6f\xea\x38\x10\xc7\xcf\xc9\x5f\x31\xb7\x82\x04\xab\x55\x77\xb5\x97\x3d\xd1\xd6\x95\xa2\xa5\xa1\x0b\x41\xa2\x7b\x89\xdc\x78\x08\x56\x1d\x3b\xeb\x1f\x50\xfe\xfb\x27\xe7\x47\x49\x80\x3e\x3d\x2a\xbd\xbe\x5e\x5a\x64\x8f\xc7\x33\xdf\xcf\xcc\x38\xb7\x73\x32\x49\x08\x90\x55\x42\xe2\x45\x34\x8b\x21\xba\x87\x78\x96\x00\x59\x45\x8b\x64\x01\x65\x9e\xe9\x7d\x69\xd5\xdf\x61\xd8\x58\x26\x93\x9b\x29\x01\x67\x50\x1b\x18\x84\x01\x67\xb0\x5c\x46\x77\xf0\x38\x8f\x1e\x26\xf3\x27\xf8\x87\x3c\xc1\x1d\xb9\x9f\x2c\xa7\x09\xe4\x28\x53\x4d\x25\x53\x45\xea\x1c\x67\x83\xe1\x28\x0c\x83\x4c\x23\xb5\xc8\x52\x6a\x21\x89\x1e\xc8\x22\x99\x3c\x3c\x26\xff\x55\x77\xc6\xcb\xe9\xf4\xed\xb0\x54\x3b\x7f\x20\x70\x25\xbb\xc4\x3e\x0c\xc6\x63\xc0\x2d\xea\x3d\x64\xce\x58\x55\xa0\x06\x6e\xa0\x44\x5d\x50\x89\xd2\x8a\x3d\x08\x2e\x5f\x90\x81\x55\x40\xc1\x58\xcd\x4b\xec\x98\xb2\xca\xc1\x5e\x39\xc8\xa8\x94\xca\x42\x1d\x2f\x50\x09\x34\xcb\x94\x93\x16\x76\xdc\x6e\x94\xb3\x80\xd2\xa2\xe6\x32\x87\x8c\x6a\x06\x0c\x2d\xe5\xc2\x84\x41\xed\x32\x6d\x5d\xa6\x9c\x41\x42\x56\xc9\x5b\xc4\xa3\x30\xc0\x82\x72\xd1\x5f\x0d\x87\x07\x89\x97\x71\xf4\xef\x92\x40\x14\xdf\x91\x55\xad\x74\x5a\x9d\x48\x9d\xe4\xff\xa7\x9c\xbd\xc2\x2c\x6e\x09\x08\xb5\x43\x3d\xa8\xb6\x87\xde\xc5\x78\x0c\x8a\x3a\xbb\x49\x69\x96\xa1\x31\xa9\x55\x2f\x28\x0d\x50\x8d\xfe\x44\x95\x75\x41\x25\xcd\x11\x6a\x03\xc8\x35\x95\xd6\xf8\x73\x6b\xad\x0a\x50\x4e\xd7\x0e\x40\xa8\x9c\x4b\x28\xb5\xda\x72\xe6\xaf\x1a\x43\xce\xed\xc6\x3d\x8f\xc0\xee\xb8\xb5\xa8\x47\x90\x2b\x95\x0b\xec\x57\xc6\xb9\xdb\x2f\xaf\x93\xc0\xa7\x97\xb6\x87\xe6\xe4\x9e\xcc\x49\x7c\x4b\x16\x4d\xda\x07\x2d\xc3\xa0\x8d\xf0\x58\xe5\x30\xe8\x06\x71\xc2\xa0\xbb\x99\xe2\x6b\xc9\xf5\xfe\x6c\x81\x35\x25\xb5\xab\xf4\x03\x8d\x6b\x8d\x66\x03\x4d\x62\x56\x41\x18\x34\x6b\xef\xdc\xd3\xdb\xfd\xde\x45\x61\x83\xaf\xd6\xbd\xc3\x4d\x49\x1c\x5b\x5e\x60\xbb\xd6\x72\xac\x0d\xf9\xba\xe1\x75\xfd\xdb\xef\xbe\xd2\x7d\xcd\x7a\x8b\x3e\x95\xca\xf6\x53\x70\x5c\xd8\xdf\x61\xd0\x91\xad\xdd\x43\x99\x29\x86\x83\x4e\x20\xcf\x7b\x8b\x66\xf0\xc7\xf5\x70\x04\x57\xcf\xd4\xe0\x5f\x7f\x5e\x35\x61\x31\xb8\x99\xcd\xa6\x64\x12\x9f\x5e\xb1\xa6\xc2\xe0\xbb\x8d\xd5\x95\xa4\xa1\xd3\xb4\x56\x5f\xac\xea\x7f\xc3\xc6\xa0\x31\x5c\xf5\xe8\xb4\x34\x3c\x03\x94\x96\x67\x7e\x5a\x34\x76\xa6\x0f\xa1\x5d\xfd\x52\x00\x6a\xcf\x34\x47\x69\x4f\x6a\x97\x97\xa7\x4d\xf5\xab\x70\xb5\xe2\xf5\x51\x1d\x24\xed\x62\x2a\xe9\xbe\x40\x69\x6b\x40\xbb\x0d\xfa\xbf\xbe\x77\x33\xa5\x19\x50\x21\x0e\x06\x05\x65\x08\x5b\xde\x3e\x04\x7d\x5e\x6f\x56\x5f\x8b\xd7\x25\x0f\x62\x3b\x54\xd6\x88\x2c\x5d\x2b\x51\x4d\x72\xaf\x4a\xfb\x5b\xad\xab\xbd\xa3\x4a\xed\x99\x7f\xb9\x57\xbe\x9b\x54\x9d\x0d\x97\x8c\x6f\x39\x73\x54\xbc\x97\xcd\x87\x20\x76\x64\x38\xc7\xb2\xa7\xd2\xa7\x21\xad\xde\x3c\xe1\xfc\xec\x3f\xee\x57\xa7\xc5\xc9\x9a\xe5\x56\xe0\x0f\x7e\x6b\x54\x09\xd5\xbe\x53\xa7\x45\xdb\x64\x95\x7e\x83\x7a\x7d\x04\x4e\x8b\xe1\xf1\xe7\x60\xa9\xcc\xc7\xda\x64\x43\xcd\xe6\x74\xc6\xfc\x74\xfd\xce\x88\x32\x0a\x83\x4c\x49\x7b\x32\x06\xbd\x58\xdf\x02\x00\x00\xff\xff\x90\x9e\x82\xf8\x25\x0b\x00\x00")

func schemaInitSqlBytes() ([]byte, error) {
	return bindataRead(
		_schemaInitSql,
		"schema/init.sql",
	)
}

func schemaInitSql() (*asset, error) {
	bytes, err := schemaInitSqlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "schema/init.sql", size: 2853, mode: os.FileMode(420), modTime: time.Unix(499137600, 0)}
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
	"schema/init.sql": schemaInitSql,
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
	"schema": &bintree{nil, map[string]*bintree{
		"init.sql": &bintree{schemaInitSql, map[string]*bintree{}},
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
