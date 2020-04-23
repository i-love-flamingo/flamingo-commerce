// Package graphql Code generated by go-bindata. (@generated) DO NOT EDIT.
// sources:
// schema.graphql
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

var _schemaGraphql = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x8c\x94\xc1\x6e\xf2\x30\x0c\xc7\xef\x7d\x8a\xf0\x1a\xbd\x7d\xc0\x77\x40\x42\x68\x1b\xbb\x4d\x13\xca\x88\x29\x96\x12\xa7\x72\x5c\x89\x6a\xe2\xdd\xa7\xb6\x61\x6b\xda\x54\x82\x13\xf8\xff\xb3\xdd\xbf\x6d\x2a\x6d\x0d\x6a\xe3\x9d\x03\x3e\xc3\x69\xd3\x04\xf1\x0e\xf8\x74\x14\x2d\x4d\x78\x83\xd0\x58\x51\xdf\x85\x52\x4a\x61\xd8\xfb\xaa\x02\xb3\xa3\x52\xad\xbd\xb7\xa0\x69\xd5\x0b\x4d\x00\xde\x6d\x4b\x75\x14\x46\xaa\x56\xc5\xbd\x28\xf2\x55\xd3\x72\xe6\x2f\xa3\xfb\x5d\x03\x07\x4f\xda\x6e\xb5\xe8\x32\xf3\x44\x2f\xbd\xde\xa9\x03\xaf\x8d\x61\x08\x01\x42\xa9\x3e\xe6\xf4\xbf\x41\x5d\x7d\xf6\xac\x81\x8b\x6e\xac\x1c\xaf\x58\xd7\x48\x55\x14\x73\x5d\x1e\x79\xe3\xb4\x35\x5a\xfb\x64\xd6\xa2\xf5\xd1\xe3\x47\xff\x15\x90\x01\x2e\xbb\xaf\xc9\x1c\x2e\xc8\x41\x0e\xda\x41\x99\xc6\xad\xfe\x0d\x27\x71\x87\xc6\x58\x18\x94\x24\xae\x91\xfe\x3b\x8d\x76\x52\xa7\x66\xb8\xe0\x6d\xe8\x9b\x08\x5f\xc8\x72\x35\xba\xed\xa5\x77\x74\x30\x84\x49\x0b\x76\x7b\x41\x69\x9f\x58\xf1\x63\x14\xd1\x24\x43\x85\x9e\x36\xde\x40\x6c\xa8\xe6\x7d\xcf\xbe\x21\xe1\x76\x06\xa5\x8c\xab\x35\xb5\x69\x91\x09\x13\x84\x01\x64\x86\x64\x98\x03\x4f\xa9\x31\xa3\x8d\xc1\xc1\x72\xf4\xb2\x47\xea\xaf\x2c\x42\xc3\x49\x09\x58\xa8\xaf\x9e\x26\xc6\xd2\x51\xfb\x20\x73\xef\x13\x67\xfd\x5c\xe7\x9f\xd9\x49\x50\xdc\xfd\x02\xd3\x9d\x47\x06\x49\x18\x88\xe7\xb0\xd4\xeb\x5e\x14\x70\x13\x20\xa3\xfa\xe5\xbe\x36\xc0\x6d\x5c\xe4\xd2\x1b\x22\xf7\x6f\x18\xbf\x3b\xf2\xc9\x99\xac\x88\xdf\x8b\x9f\x00\x00\x00\xff\xff\xd9\xe6\xfb\x5e\x91\x04\x00\x00")

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

	info := bindataFileInfo{name: "schema.graphql", size: 0, mode: os.FileMode(0), modTime: time.Unix(0, 0)}
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