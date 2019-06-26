// Code generated by go-bindata.
// sources:
// schema.graphql
// DO NOT EDIT!

package graphql

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

var _schemaGraphql = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xec\x56\x4b\x6f\xdb\x38\x10\x3e\xc7\x80\xff\x03\x83\x5c\xb2\x8b\xfd\x05\xba\xe5\xb1\xbb\x08\x9a\x57\x61\xa3\x97\xa0\x28\xc6\xd4\xd8\x26\x2c\x91\x02\x39\x72\x2a\x14\xfd\xef\x05\x49\x51\x0f\x8a\x51\x9c\x63\x81\xe6\x14\xcf\x7c\xc3\x79\x7d\x33\x23\x21\x09\xf5\x16\x38\xb2\x1b\x55\x96\xa8\x39\x7e\x7b\xd6\x2a\xaf\x39\xb1\x1f\xcb\x05\x63\x8c\x6d\xc0\xe0\x2d\x10\x64\x3d\xe2\x1a\x8c\xe0\x2d\xcc\xaa\xce\x3d\x92\x10\x0c\xea\x08\xdb\xc2\xd6\x9d\xae\x05\x9b\x0a\xb9\xd8\x0a\x0e\x24\x94\x34\x53\x83\xd5\x48\xdf\x1a\x09\xb3\x82\x02\x61\x53\x60\xc6\xae\x95\x2a\x10\x64\x78\xae\x95\xa7\xbd\x07\xab\x10\x68\x53\x61\xc6\x56\xa4\x85\xdc\xb5\xa2\x1d\xd2\x5d\x8e\x92\xc4\x56\xa0\x8e\x74\x7b\x30\x0f\x98\x0b\xb8\xdc\x69\x55\x57\x9d\xf2\x1f\x56\x1b\xd8\xf5\x0f\xfd\x15\xc7\xb4\x43\x3a\xd1\x2e\x8e\xd7\x99\x9d\x2f\x17\x3f\x97\x8b\xe5\xc2\x46\xdb\x23\x56\xa2\xac\x0a\x0c\x3d\x72\x3f\x4a\x94\x64\xfe\xf4\xef\xf7\xe8\x5f\x5c\xfb\xd0\x26\x12\x64\xab\x32\xf8\x1b\xe5\x00\x44\x5a\x6c\x6a\x42\x13\x30\xb1\xcf\xab\x0e\xd1\x56\x74\xaf\x34\xdd\xa2\xe1\x5a\x54\xb6\x07\x51\x51\xf2\xa1\x66\xea\xaf\xb4\x19\x8c\xc2\x79\x49\x27\xf9\xd5\xa6\xe8\x2c\x40\x1f\x90\x9e\x0b\xe0\x78\xa3\xf2\xb8\x3f\x1a\x09\x44\x81\xda\xab\x62\x6f\x41\xbb\x3a\xd4\xd9\x34\x96\xa0\x7d\x84\x32\xb2\xf5\x7a\xae\x11\x08\xf3\x2b\xb2\xca\xb5\x28\x03\x4d\xea\x2a\x4f\xca\x8f\xc2\x88\x4d\x81\xff\x69\x55\x66\x09\xf9\x5a\x0d\xf0\x5e\x73\x03\x84\x3b\xa5\x85\xaf\x7e\x5f\x89\x56\xde\xf8\xe9\xb0\xb5\xb0\xe8\x07\x10\x32\x68\x06\xdc\x88\xc0\x5d\xf4\x41\xac\x6c\x6d\x1e\xa0\xaa\x84\xdc\x65\xec\xa5\x4d\xb2\xab\xaf\x21\xc5\x0f\xf7\x78\xc4\x22\x8b\x0b\x70\xc0\xe6\x55\xe9\xdc\x8c\xac\xfc\xbc\x3d\xe2\xab\x2b\x59\xcf\xec\x04\x27\x27\x33\x1e\x48\xe9\x38\xb4\xf6\xcc\xf4\x2f\xcf\x73\xcb\x6b\x2f\xd8\xfa\xe9\xf6\xa9\xdf\x4d\x4c\x48\x56\x69\xc1\x91\x11\x98\x43\x16\x40\x7e\xdd\x3c\x3b\xc5\x74\x84\x9c\xfc\x4e\x6e\xd5\x70\x37\x79\xa1\xb1\x9d\x73\xff\x76\x13\xeb\x41\x95\xc6\x15\x16\xc8\x09\xf3\x2f\xa0\x05\x48\x72\x84\x1a\xc6\xe6\x78\xcd\xb2\x19\x36\xcf\x70\xf9\x23\xe9\xf9\x80\xaf\x8e\x20\x0a\xbb\xb9\x5c\xb8\xc6\x3a\xee\x12\x0b\xce\x5a\xe8\xbd\x6a\xa0\xa0\xa6\x53\xb3\xbf\x63\xc9\x5c\xef\xc6\xeb\x36\xf4\xcf\x2d\x2e\x93\xc8\x76\x04\xff\xdf\xa2\x6c\x34\x27\x3e\xef\xf0\xd1\xde\x1a\x0d\x2c\x4a\xf2\xb3\x32\xef\xf7\x5f\x49\xba\xf9\x80\x5f\x87\x0f\x7e\x0b\xd8\x8c\x26\xc1\x0d\x30\x14\x35\x8e\xa7\x60\xe6\xe9\xf6\xa6\x84\x07\xdf\xbf\x4d\xfd\xc2\x18\xcb\xed\xc2\xe8\xa5\xa7\xd1\xe3\x0c\x38\x89\x23\xb6\x2c\x7e\x87\xfc\x17\x67\x30\xe6\x51\xa2\xb2\x53\x5a\x9d\x15\x03\xfe\xa4\x4c\x62\x7e\xcd\xd7\xcb\x0d\x48\xd7\x75\x77\x81\x13\x67\x43\x94\xb8\xf6\xba\x91\xbc\xbd\x98\x13\xfc\xe0\xea\x45\x2b\x7f\x8b\x1a\x25\x1f\x32\xeb\xed\xd0\xfa\xd3\xd7\x11\x1f\x7b\xe1\x27\x6c\xa6\x9b\x71\x08\x48\xd5\xa6\x53\x06\x83\x3d\x98\x4e\x76\x79\xc0\x26\xf1\xd9\x30\x7d\x39\x06\xbe\xe9\xe5\xa4\xec\x42\x72\x7c\x7a\x5e\x53\xd3\x50\x4b\x41\xd1\x25\x4e\x38\x19\x5f\xa5\x19\x17\x15\xd0\x3e\x12\x49\x77\x90\x23\x94\x46\x49\x6f\x5e\xbd\x36\x04\xfc\x4e\x28\x73\xc7\x22\xf6\xb9\xc6\x7e\xaa\xe3\xdc\x2f\xfd\x26\xae\x26\x5f\x15\x89\x62\xba\xb7\x7f\x05\x00\x00\xff\xff\xe4\xe3\xb1\xe4\xd0\x0c\x00\x00")

func schemaGraphqlBytes() ([]byte, error) {
	return bindataRead(
		_schemaGraphql,
		"schema.graphql",
	)
}

func schemaGraphql() (*asset, error) {
	bytes, err := schemaGraphqlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "schema.graphql", size: 3280, mode: os.FileMode(438), modTime: time.Unix(1561558095, 0)}
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
	"schema.graphql": schemaGraphql,
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
	"schema.graphql": &bintree{schemaGraphql, map[string]*bintree{}},
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
