// Code generated by go-bindata. (@generated) DO NOT EDIT.

// Package fake generated by go-bindata.// sources:
// mock/categoryTree.json
// mock/clothing.json
// mock/electronics.json
// mock/flat-screen_tvs.json
// mock/headphone_accessories.json
// mock/headphones.json
// mock/jumpsuits.json
// mock/tablets.json
package fake

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
		return nil, fmt.Errorf("read %q: %v", name, err)
	}

	var buf bytes.Buffer
	_, err = io.Copy(&buf, gz)
	clErr := gz.Close()

	if err != nil {
		return nil, fmt.Errorf("read %q: %v", name, err)
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

// Name return file name
func (fi bindataFileInfo) Name() string {
	return fi.name
}

// Size return file size
func (fi bindataFileInfo) Size() int64 {
	return fi.size
}

// Mode return file mode
func (fi bindataFileInfo) Mode() os.FileMode {
	return fi.mode
}

// ModTime return file modify time
func (fi bindataFileInfo) ModTime() time.Time {
	return fi.modTime
}

// IsDir return file whether a directory
func (fi bindataFileInfo) IsDir() bool {
	return fi.mode&os.ModeDir != 0
}

// Sys return file is sys mode
func (fi bindataFileInfo) Sys() interface{} {
	return nil
}

var _categorytreeJson = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x9c\x93\x3f\x6b\xf3\x30\x10\x87\x77\x7d\x8a\x43\xc3\x3b\xe5\xc5\x7b\xb6\xe2\x36\x94\x0e\xa5\x10\xd3\x25\x84\x70\x51\xae\x91\x8a\x22\x05\xe9\x5c\x28\x21\xdf\xbd\x38\x7f\x5a\x9b\x44\x8a\x5a\x2f\xf6\xf0\xfc\xce\x8f\xee\x74\x33\x01\xb0\x13\x00\x00\xb2\x46\xa6\xb5\x0f\x9f\xb5\x5f\x91\x1c\xc3\xf1\x91\x64\x49\x71\xf0\xce\xa8\x28\x47\x43\xf0\x19\x37\x3d\xf0\x21\x0d\xbe\x20\xeb\x2c\x38\x6d\x97\x4d\x20\x8a\xf7\xc8\x28\xc7\x30\x13\x47\x74\x77\x7a\x27\xe5\xde\x2c\xf2\xff\xa8\x02\x91\x5b\xf0\xc7\xb9\x5c\x46\x72\x62\x91\x61\x7a\x08\x44\xf8\x07\xcd\xeb\xb5\x48\x52\xb7\xba\x8c\x9f\xd2\xfb\x51\xa9\xb2\x26\x5c\x6d\xb5\x77\x54\x60\xfb\x98\x65\xd3\x9a\xd7\x73\x89\x26\x0f\xad\x4b\xcc\x17\xa8\x14\xc5\xe8\x83\x19\xfc\x20\x73\x90\xbb\xdb\x81\x92\xd3\x54\xfd\x3a\xbd\x32\xfb\xef\xef\xf9\xaf\xe7\xc1\xb8\xb4\xc4\x05\xc3\x68\xd2\x60\xda\xfd\x1c\x12\x7d\xd1\x4e\xf2\x20\x98\x5d\x3c\x65\x3d\x6b\xe3\xd6\xf9\xad\xab\x3b\x8a\xba\xcb\x38\xc1\xa8\x8d\x77\xf9\xdd\x4b\xe2\x7f\xdd\xc0\xf7\x76\xb3\x8d\xad\x29\x69\xe0\x53\x0e\xbd\xa1\x59\xfd\x84\x2f\x5b\x29\xe6\x5f\x01\x00\x00\xff\xff\xbb\xff\x03\x12\xc6\x04\x00\x00")

func categorytreeJsonBytes() ([]byte, error) {
	return bindataRead(
		_categorytreeJson,
		"categoryTree.json",
	)
}

func categorytreeJson() (*asset, error) {
	bytes, err := categorytreeJsonBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "categoryTree.json", size: 0, mode: os.FileMode(0), modTime: time.Unix(0, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _clothingJson = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xaa\xe6\x52\x50\x50\x72\x4e\x2c\x49\x4d\xcf\x2f\xaa\x74\xce\x4f\x49\x55\xb2\x52\x50\x4a\xce\xc9\x2f\xc9\xc8\xcc\x4b\x57\xd2\x41\x96\xf5\x4b\xcc\x05\xcb\x3a\x83\x64\x53\x8b\x15\xd4\x14\xdc\x12\x8b\x33\x32\xf3\xf3\x50\x95\x05\x24\x96\x64\x60\x57\xc6\x55\x0b\x08\x00\x00\xff\xff\xa9\x92\xc7\x2f\x6e\x00\x00\x00")

func clothingJsonBytes() ([]byte, error) {
	return bindataRead(
		_clothingJson,
		"clothing.json",
	)
}

func clothingJson() (*asset, error) {
	bytes, err := clothingJsonBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "clothing.json", size: 0, mode: os.FileMode(0), modTime: time.Unix(0, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _electronicsJson = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xaa\xe6\x52\x50\x50\x72\x4e\x2c\x49\x4d\xcf\x2f\xaa\x74\xce\x4f\x49\x55\xb2\x52\x50\x4a\xcd\x49\x4d\x2e\x29\xca\xcf\xcb\x4c\x2e\x56\xd2\x41\x56\xe0\x97\x98\x0b\x56\xe0\x8a\x4b\x41\x40\x62\x49\x06\xba\x02\xae\x5a\x40\x00\x00\x00\xff\xff\x05\xf9\xd5\x8b\x65\x00\x00\x00")

func electronicsJsonBytes() ([]byte, error) {
	return bindataRead(
		_electronicsJson,
		"electronics.json",
	)
}

func electronicsJson() (*asset, error) {
	bytes, err := electronicsJsonBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "electronics.json", size: 0, mode: os.FileMode(0), modTime: time.Unix(0, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _flatScreen_tvsJson = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xaa\xe6\x52\x50\x50\x72\x4e\x2c\x49\x4d\xcf\x2f\xaa\x74\xce\x4f\x49\x55\xb2\x52\x50\x4a\xcb\x49\x2c\xd1\x2d\x4e\x2e\x4a\x4d\xcd\x8b\x2f\x29\x2b\x56\xd2\x41\x56\xe4\x97\x98\x0b\x56\xe4\x96\x93\x58\xa2\x10\x0c\x56\x54\xac\xa0\xa6\x10\x12\x86\xaa\x2c\x20\xb1\x24\x03\xa4\xcc\x35\x27\x35\xb9\xa4\x28\x3f\x2f\x33\xb9\x58\x1f\x53\x0b\x57\x2d\x20\x00\x00\xff\xff\x3c\x7e\xb0\x02\x81\x00\x00\x00")

func flatScreen_tvsJsonBytes() ([]byte, error) {
	return bindataRead(
		_flatScreen_tvsJson,
		"flat-screen_tvs.json",
	)
}

func flatScreen_tvsJson() (*asset, error) {
	bytes, err := flatScreen_tvsJsonBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "flat-screen_tvs.json", size: 0, mode: os.FileMode(0), modTime: time.Unix(0, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _headphone_accessoriesJson = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xaa\xe6\x52\x50\x50\x72\x4e\x2c\x49\x4d\xcf\x2f\xaa\x74\xce\x4f\x49\x55\xb2\x52\x50\xca\x48\x4d\x4c\x29\xc8\xc8\xcf\x4b\x8d\x4f\x4c\x4e\x4e\x2d\x2e\xce\x2f\xca\x4c\x2d\x56\xd2\x41\x56\xea\x97\x98\x0b\x56\xea\x88\x4b\x41\x40\x62\x49\x06\x48\x81\x6b\x4e\x6a\x72\x49\x51\x7e\x5e\x66\x72\xb1\xbe\x07\xcc\xdc\x62\x7d\x64\x7d\x5c\xb5\x80\x00\x00\x00\xff\xff\x89\xd8\x8b\x00\x86\x00\x00\x00")

func headphone_accessoriesJsonBytes() ([]byte, error) {
	return bindataRead(
		_headphone_accessoriesJson,
		"headphone_accessories.json",
	)
}

func headphone_accessoriesJson() (*asset, error) {
	bytes, err := headphone_accessoriesJsonBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "headphone_accessories.json", size: 0, mode: os.FileMode(0), modTime: time.Unix(0, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _headphonesJson = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xaa\xe6\x52\x50\x50\x72\x4e\x2c\x49\x4d\xcf\x2f\xaa\x74\xce\x4f\x49\x55\xb2\x52\x50\xca\x48\x4d\x4c\x29\xc8\xc8\xcf\x4b\x2d\x56\xd2\x41\x96\xf7\x4b\xcc\x05\xcb\x7b\xe0\x90\x0f\x48\x2c\xc9\x00\xc9\xbb\xe6\xa4\x26\x97\x14\xe5\xe7\x65\x26\x17\xeb\x23\xa9\xe5\xaa\x05\x04\x00\x00\xff\xff\x11\x32\xb4\x48\x6e\x00\x00\x00")

func headphonesJsonBytes() ([]byte, error) {
	return bindataRead(
		_headphonesJson,
		"headphones.json",
	)
}

func headphonesJson() (*asset, error) {
	bytes, err := headphonesJsonBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "headphones.json", size: 0, mode: os.FileMode(0), modTime: time.Unix(0, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _jumpsuitsJson = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xaa\xe6\x52\x50\x50\x72\x4e\x2c\x49\x4d\xcf\x2f\xaa\x74\xce\x4f\x49\x55\xb2\x52\x50\xca\x2a\xcd\x2d\x28\x2e\xcd\x2c\x29\x56\xd2\x41\x96\xf6\x4b\xcc\x05\x4b\x7b\x61\x97\x0e\x48\x2c\xc9\x00\x49\x3b\xe7\xe4\x97\x64\xa4\x16\x2b\xa8\x29\xb8\x25\x16\x67\x64\xe6\xe7\xe9\x23\x34\x70\xd5\x02\x02\x00\x00\xff\xff\x9f\x0a\xa6\x82\x71\x00\x00\x00")

func jumpsuitsJsonBytes() ([]byte, error) {
	return bindataRead(
		_jumpsuitsJson,
		"jumpsuits.json",
	)
}

func jumpsuitsJson() (*asset, error) {
	bytes, err := jumpsuitsJsonBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "jumpsuits.json", size: 0, mode: os.FileMode(0), modTime: time.Unix(0, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _tabletsJson = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xaa\xe6\x52\x50\x50\x72\x4e\x2c\x49\x4d\xcf\x2f\xaa\x74\xce\x4f\x49\x55\xb2\x52\x50\x2a\x49\x4c\xca\x49\x2d\x29\x56\xd2\x41\x96\xf4\x4b\xcc\x05\x4b\x86\x60\x93\x0c\x48\x2c\xc9\x00\x49\xba\xe6\xa4\x26\x97\x14\xe5\xe7\x65\x26\x17\xeb\xc3\x14\x72\xd5\x02\x02\x00\x00\xff\xff\xd1\x64\x9b\x19\x65\x00\x00\x00")

func tabletsJsonBytes() ([]byte, error) {
	return bindataRead(
		_tabletsJson,
		"tablets.json",
	)
}

func tabletsJson() (*asset, error) {
	bytes, err := tabletsJsonBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "tablets.json", size: 0, mode: os.FileMode(0), modTime: time.Unix(0, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

// Asset loads and returns the asset for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func Asset(name string) ([]byte, error) {
	canonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[canonicalName]; ok {
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
	canonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[canonicalName]; ok {
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
	"categoryTree.json":          categorytreeJson,
	"clothing.json":              clothingJson,
	"electronics.json":           electronicsJson,
	"flat-screen_tvs.json":       flatScreen_tvsJson,
	"headphone_accessories.json": headphone_accessoriesJson,
	"headphones.json":            headphonesJson,
	"jumpsuits.json":             jumpsuitsJson,
	"tablets.json":               tabletsJson,
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
// AssetDir("foo.txt") and AssetDir("nonexistent") would return an error
// AssetDir("") will return []string{"data"}.
func AssetDir(name string) ([]string, error) {
	node := _bintree
	if len(name) != 0 {
		canonicalName := strings.Replace(name, "\\", "/", -1)
		pathList := strings.Split(canonicalName, "/")
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
	"categoryTree.json":          &bintree{categorytreeJson, map[string]*bintree{}},
	"clothing.json":              &bintree{clothingJson, map[string]*bintree{}},
	"electronics.json":           &bintree{electronicsJson, map[string]*bintree{}},
	"flat-screen_tvs.json":       &bintree{flatScreen_tvsJson, map[string]*bintree{}},
	"headphone_accessories.json": &bintree{headphone_accessoriesJson, map[string]*bintree{}},
	"headphones.json":            &bintree{headphonesJson, map[string]*bintree{}},
	"jumpsuits.json":             &bintree{jumpsuitsJson, map[string]*bintree{}},
	"tablets.json":               &bintree{tabletsJson, map[string]*bintree{}},
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
	canonicalName := strings.Replace(name, "\\", "/", -1)
	return filepath.Join(append([]string{dir}, strings.Split(canonicalName, "/")...)...)
}
