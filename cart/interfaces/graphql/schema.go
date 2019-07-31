// Code generated for package graphql by go-bindata DO NOT EDIT. (@generated)
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

// Mode return file modify time
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

var _schemaGraphql = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x84\x56\xc1\x6e\xdb\x30\x0c\x3d\x37\x5f\xe1\xa0\x97\x0d\xe8\x17\xf8\x96\xda\xc3\x10\x60\xeb\x32\x34\x3b\x0d\x43\xc1\x5a\x4c\x22\x4c\x96\x34\x89\x2e\x6a\x0c\xfb\xf7\xc1\x96\x92\xc8\x92\xec\xe4\xa8\xf7\x44\x3e\x52\x8f\x8c\xa9\xd7\x58\x54\xaa\x6d\xd1\x34\xf8\x52\x63\xa3\x0c\x10\xb2\x0a\x0c\x15\x7f\x57\x45\x51\x14\x0d\x18\x2a\xaf\x94\x01\x59\x8f\x00\x3b\x93\x5f\x18\x0a\xfe\x86\x86\xa3\x2d\x8b\x9f\x13\xe6\x25\x60\xed\x28\xfd\xfa\xd7\xea\xdf\x6a\x35\xcd\x1a\x24\xe3\xac\x2c\xb6\xb5\x8b\x8f\x92\x38\xf5\xdb\xba\x2c\x9e\xc9\x70\x79\x74\xa7\xaf\x5c\x08\x2e\x8f\x1b\xc6\x0c\x5a\x1b\x09\xf3\xa7\x23\x51\x77\xa6\x39\x81\x45\x13\x71\x76\x68\xac\x92\xbe\x82\x79\xdd\x17\xb9\x03\xf1\xfe\x0e\x18\xe3\xc4\x95\x04\x51\x03\x41\x9a\x36\x00\xd7\xfe\x8a\x86\xbe\x45\x49\xcf\x28\xb0\x19\xd0\x58\x47\x04\xfb\xf2\x50\x28\x79\xb4\x7b\xb5\xe9\xe8\x34\x74\xa0\x19\xba\xf7\x63\x2c\xe3\x51\x29\x81\xe0\x89\x10\xe3\x71\xa3\xee\xef\x40\x6b\xc1\x91\x55\xaa\xd3\x4a\x56\x8a\xa5\x65\x5e\x21\x5f\x28\xc3\x03\x74\x82\xaa\xce\x18\x94\x4d\x1f\x47\x24\x45\x20\x38\x61\x9b\x44\xda\x9f\x91\x6b\xc7\x5c\xf6\xcf\xfc\x40\x15\x18\x96\xdc\xd8\x4c\xf1\x39\x63\x24\x0e\xf2\x4e\xf1\x6f\xd7\x97\x31\xdd\x3f\xdc\xd4\xa2\xdb\x9c\xe6\x3a\x44\xe7\xf3\x4f\xd2\x0e\xd3\x90\x6d\xc0\x72\x8c\x20\xd1\xd9\xe9\x84\x6d\xa4\x7d\x40\x9d\x75\x8d\x62\x5d\x13\x4e\xdd\xce\x9d\xe4\xc3\x87\x51\x83\xf9\xf9\x43\x7d\x59\x6c\xe5\xcc\x25\x3f\x2a\xfe\xde\x1b\x50\x59\x64\x7e\xe1\xf3\x1f\xb8\xb1\x24\xa1\xc5\x72\x9e\x23\x20\x4b\x99\x70\x5a\xce\x98\xc0\xa7\x84\x15\x72\x88\x93\x48\x82\x44\x1c\x0b\xa2\x23\x70\x93\x35\xcb\x21\x83\x98\x29\x2d\xe5\x3c\x99\x25\xcd\xd7\xf1\xf7\x7d\xfb\xc2\xe5\x38\x4e\x9e\xe4\x4c\xdf\xa8\x56\x83\xec\x93\x74\x61\xa4\x86\x53\x4a\x88\x38\x5a\x59\x1a\x86\x72\x49\x91\x25\xa0\x5b\x1d\x32\x78\xe4\x6e\xbc\xe7\x3b\xd4\xa8\x4e\x92\xb9\xa1\xd9\x71\x92\x40\x93\x17\x43\x81\xfa\xa4\xe4\x92\x3b\xb0\x05\x2e\x16\x34\x67\x8d\xea\xf6\xb5\xf7\x29\xdc\x5e\xfb\x23\x1d\x44\x8d\x04\x5c\xc4\xcc\xdd\x14\xf5\xb2\xde\xb9\x25\x2e\x8f\x55\x67\x49\xb5\x68\x32\x1b\xfe\x53\x86\x92\x97\x9b\x63\x46\xc3\xb9\x50\xe6\x45\xd9\x79\xcb\x01\xe1\xb7\xc3\x23\x37\x74\x9a\xee\x62\x0d\xd6\x6a\x35\x6e\x70\xf7\x7c\x39\xf0\xa9\x6b\x5f\x87\x3f\x8e\x10\x93\xe0\x7c\x3c\xda\x30\x68\x3c\xbe\x13\x4a\x56\x8c\xc2\xbe\x77\xc1\xc2\x0b\x45\x96\x33\x9f\x0a\x49\x84\xaf\x7e\x32\xe3\x20\x1b\xc6\xf6\x6a\xb8\xf1\xa1\x05\xf3\x1b\x49\x0b\x68\xd0\xf9\x6a\x5b\xaf\x1f\x2e\x3b\xeb\xe1\xb2\xdf\x1d\xe8\x85\x7e\x5c\x12\xf0\x3f\x00\x00\xff\xff\xd8\x12\x29\xb2\xcb\x08\x00\x00")

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

	info := bindataFileInfo{name: "schema.graphql", size: 2251, mode: os.FileMode(420), modTime: time.Unix(1564550596, 0)}
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
