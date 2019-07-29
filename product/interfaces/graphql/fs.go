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

var _schemaGraphql = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xec\x56\xcd\x6e\xe3\x36\x10\xbe\xeb\x29\x14\xf4\x92\x02\x7d\x02\xdd\x12\xa7\x5b\x04\xdd\x2c\xd2\xda\xed\x65\x11\x14\x63\x6a\x6c\x13\xa6\x48\x81\x1c\x25\x16\x8a\xbc\x7b\x21\x92\x92\xf9\x23\xd9\x3e\xf6\xb0\x39\xc5\xf3\x7d\x9c\x19\x8e\x66\xbe\x21\x97\x84\x7a\x07\x0c\xcb\x95\x6a\x1a\xd4\x0c\xff\x79\xd5\xaa\xee\x18\x95\xff\x16\x65\x59\x96\x5b\x30\xf8\x04\x04\xd5\x99\xf0\x08\x86\x33\xcf\x1a\xa0\x3b\x4b\x24\x04\x83\x3a\xa1\x7a\xd6\x66\xc2\x1c\xd7\xb4\xc8\xf8\x8e\x33\x20\xae\xa4\xc9\xf9\xeb\x08\x77\x67\xb8\x59\x83\x40\xd8\x0a\xac\xca\x47\xa5\x04\x82\xf4\xce\xbc\x79\x3e\xf4\x78\xc8\x27\xd9\xb7\x58\x95\x6b\xd2\x5c\xee\x9d\x65\x8f\xf4\x5c\xa3\x24\xbe\xe3\xa8\x63\xe8\x00\xe6\x05\x6b\x0e\xf7\x7b\xad\xba\x76\xc2\x7e\x29\x3b\x03\xfb\xb3\x9b\x9f\x93\x7c\xf6\x48\x37\x1e\x4b\x73\xb5\xc7\xee\x8a\xcf\xa2\x18\xf2\x3c\xc3\x6b\xde\xb4\x02\xc7\xef\x62\x7f\x34\x28\xc9\xfc\xf8\x66\xff\xdf\x6f\x96\x16\xdc\x7f\x1a\xe2\x34\x14\x23\xf8\x0b\xb3\x07\x22\xcd\xb7\x1d\xa1\x19\x29\x69\xb8\x87\x89\xe1\xea\x78\x50\x9a\x9e\xd0\x30\xcd\xdb\xa1\xee\x71\x31\xea\x10\xc8\x82\x35\x43\xe6\x51\x2a\xdf\xe7\x2f\xf7\x56\x38\x3e\xe8\x23\xd2\xab\x00\x86\x2b\x55\x27\x9f\x44\x23\x01\x17\xa8\x1d\x92\x44\x1a\xc1\xf5\xb1\xab\xb2\x34\x46\xf0\x1b\x34\xf1\x49\x8b\x32\x8d\x40\x58\x3f\xd0\x00\x6d\x78\x83\xd6\xda\xb5\xf5\x8c\xf5\x9d\x1b\xbe\x15\xf8\x45\xab\xa6\xca\xac\x1b\x35\x71\xad\x79\x05\x84\x7b\xa5\xb9\x2b\xf5\xf9\xe6\xde\xde\xbb\xf6\xbf\x7b\xb3\xe4\x17\xe0\x72\x04\x82\x16\x48\xb8\x3e\xe3\xd1\xa8\x86\x5a\xbc\x40\xdb\x72\xb9\xaf\xca\xef\xfe\x5a\xbe\x98\x86\x14\x3b\x7e\xc5\x77\x14\x55\x7c\xe1\x23\xf6\x1f\x4a\xd7\x26\x3c\xe1\x46\xe9\x1b\x7e\xd8\xfa\x4c\x8d\x9b\xf5\x5c\x36\xb9\xbe\xe9\x6c\x97\x6c\x5c\xe7\x39\xa7\x17\x9b\x27\x90\x86\x57\xcd\x19\x96\x79\xd7\x5b\xfb\xb3\xdc\xa9\x94\xfb\x6c\x86\xf2\xdb\x7f\xa7\x11\xb3\x9c\x56\xe3\x1a\x05\x32\xc2\xfa\x6f\xd0\x1c\x24\xd9\x6e\x08\x22\xda\x7e\x2c\xab\xe5\x2e\xbc\xd0\x84\x41\x1a\x0f\xef\xc0\xc5\xa0\x1e\x36\x09\x33\xe7\x70\xca\xde\x3b\x75\x07\xbf\xaa\x1e\x04\xf5\x13\x98\xdf\x3a\x65\x2c\x7e\x81\x94\x38\x0e\x7f\xa6\x64\x35\xee\xa0\x13\x14\x85\xe2\x0c\x47\xf5\x7c\xe2\x86\xa9\x4e\x12\xd6\x89\x5e\xd5\x01\x30\x77\x74\xc4\x37\x78\xa2\x38\x62\xc3\xe5\xab\xe2\x92\xcc\x46\xad\x5b\x94\x54\x95\x5f\x84\x02\xf2\x20\x9c\x96\x41\xa6\x24\x59\x77\x71\xc0\x95\x33\xcf\xb6\xe3\x19\xf6\x15\x60\x9d\x21\xd5\xa0\xfe\x2d\x92\x59\x07\x1d\x40\x4a\x14\xb9\xb4\x08\xc5\x40\x04\xb6\xa5\xa2\xc7\x0b\xc8\x07\xb4\x7a\x6e\x66\x7a\x20\x62\xdb\x74\xee\xde\x6e\x73\x6d\xc9\xb1\x9c\x87\xe9\xa2\x24\xa7\x2a\x97\x43\xfe\x2a\x49\xf7\xb7\x86\xb4\x64\x1f\x52\xc0\x36\x14\x0d\x2b\x71\x20\x3a\x8c\x14\x63\xd1\xab\xdf\xac\xde\xd7\xd5\xfd\x9c\x68\xe9\x68\x1e\xc4\x74\x32\x02\x23\xfe\x8e\x7e\xe4\x2f\x0b\x05\xc4\xb3\x79\xc3\x68\x8a\x60\x92\xe6\xf8\xe9\xa4\x5d\xb8\xbb\x95\x91\x68\x12\xf3\x95\xc8\x1b\xdc\x38\x28\x34\xfb\x47\x40\xca\x0e\x76\x79\xbc\xd0\x76\xa8\x51\xb2\x1b\x5a\xf6\xbc\xcd\x7d\x5e\xd3\x03\xe0\x77\xec\xb3\x15\x10\xbe\x0e\xb2\x42\x4c\xae\x3c\xf9\x00\x66\x32\xdd\x1f\xb1\x9f\x79\xfa\x8c\x2f\x9f\x45\xde\x62\x8c\xeb\x17\x1a\xe7\x3d\x9f\xe6\xbc\x7d\x3b\xc9\x29\x1e\xfb\xcc\x7d\xbc\x69\x17\x9d\xb7\x40\x87\xd8\x22\xed\xab\x22\xe6\x68\xab\x6c\x0b\xbe\x17\xaf\x96\x8a\xf9\x82\x74\x5f\x91\xe7\x2b\xea\xec\x66\xe9\x11\x0c\x46\xd2\x7b\x36\x3f\x34\xc3\xc1\x05\xf0\x2f\xc9\x29\xda\x88\xb3\x1b\xc4\xbf\x52\x9a\x16\xf8\x5e\xfe\xd9\x09\xcc\x1a\xad\x46\xd9\xbf\x28\x8d\xe3\x61\x13\x9f\xfd\xe9\xfa\x36\x70\xf3\x01\xa7\x95\x00\x63\x82\x17\xce\x67\x51\xe0\x89\x50\xd6\x76\x02\xcb\x3f\x3a\x9c\x74\x2d\xad\xf7\xbd\x5b\xf5\x6d\xf6\xde\x9c\x69\xcc\xe2\xb3\xf8\x2f\x00\x00\xff\xff\x0f\xd1\x1e\xb6\xc9\x0e\x00\x00")

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

	info := bindataFileInfo{name: "schema.graphql", size: 3785, mode: os.FileMode(420), modTime: time.Unix(1564414810, 0)}
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
