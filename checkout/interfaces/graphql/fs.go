// Code generated by go-bindata. (@generated) DO NOT EDIT.

// Package graphql generated by go-bindata.// sources:
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

var _schemaGraphql = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xd4\x97\x4d\x6f\xdb\x38\x13\xc7\xef\xfe\x14\x63\xe4\xd2\x02\xe9\x17\xd0\xad\xf1\xf3\x74\x37\x40\x8b\x78\x93\xb6\x7b\x28\x0a\x63\x22\x8e\x6d\x22\x14\xc7\x3b\x1c\xc5\x15\x16\xfb\xdd\x17\xa4\x68\x45\xb6\x6c\xaf\xf3\xd2\x76\xd7\x87\xa2\x11\xc9\x99\xdf\xbc\x70\xa4\xbf\x36\x2b\x82\x09\x57\x15\x49\x49\xb3\xc9\x92\xca\x3b\xae\x75\x76\xa3\x28\x3a\x75\x58\xd2\x95\x18\x92\xd9\x35\x85\xda\x29\xfc\x39\x02\x00\xa8\x6b\x6b\x0a\xb8\x51\xb1\x7e\x31\x1e\xfd\x35\x3a\xdb\x63\xe0\xe1\xec\x84\xbd\xd2\x37\x05\xa1\x95\x50\x20\xaf\x01\x74\x49\x20\xad\x45\x9e\xa7\xbf\xca\x5a\x84\xbc\xc2\x2b\xa9\xbd\xb7\x7e\xf1\x1a\x56\xd1\x00\x70\xb4\x00\x55\xad\xa8\x96\xfd\xe8\x00\xed\xd0\x59\x0b\x7a\x06\x1f\x97\x04\x13\x14\x05\x5d\xa2\x82\x0d\xb0\x60\xeb\x17\xa0\x0c\xb7\xd4\xba\x30\x69\x67\x89\xa2\xc5\x83\xe5\xff\x51\xc9\x82\x4a\x26\x9e\xed\x99\x6a\x4f\x64\x2a\xeb\xa1\xc4\xb0\x61\xb4\x01\xd0\x09\xa1\x69\xfa\x76\xd3\xda\xa5\x9f\x73\x28\x0e\x71\x9b\xab\x6e\x4f\xf6\x74\xa3\xa8\x04\x86\x56\xe4\x4d\xa4\x65\x9f\x72\x14\xd2\x63\x9e\xc3\x0a\x9b\x2a\x26\x0b\xbd\xd9\x4a\xd3\x9b\xbc\xa5\xc2\x06\x4a\xf6\x8a\xd6\x03\x1a\x63\x63\xea\xd0\x81\xed\x5c\xa4\x6d\x07\x81\x12\x4f\x62\x98\xa5\x7f\xc7\x19\xeb\x2d\xd4\xde\xfe\x51\x13\x58\x03\x73\x96\xc4\xb4\x12\x2e\x29\x84\xbd\x6d\x31\x3a\xdc\x18\xbd\x98\xe1\x4d\x0b\x06\x78\xcb\xb5\xb6\x46\x7b\x59\x8e\xeb\xda\xac\x6c\x89\xce\x35\x10\x96\xbc\xf6\x31\x1f\x08\xa1\x2e\x89\x42\x80\x15\x2e\xe8\x68\x5f\xf4\x7d\xb5\x6d\x91\xf3\x97\xcb\x92\x7f\x5f\x8e\x65\x63\xfa\x70\x62\xfc\xb5\xb5\xb1\x63\xba\xd8\xb1\x81\x32\x70\x9f\x4f\x52\x85\xd6\x75\x6e\xf3\xaf\x97\xb5\x14\xcb\xd1\xda\xf4\x68\x72\x44\x0b\x54\x5a\x63\x53\x0c\xec\xf5\xc2\x9d\x0a\xdf\x5b\x43\x52\x6c\x2d\x56\xa4\x4b\x36\xc5\x90\x24\xfe\x1f\x2b\xae\xbd\xf6\x16\x3b\xaa\xa9\xd8\x32\x37\x86\x5a\x75\x54\xec\x8f\x65\x64\xbd\x92\xcc\x63\x8b\x9e\xd8\x6c\x39\x20\x8f\x15\x15\x83\xac\x9c\x68\x63\xf6\x3b\x5a\x05\x5b\xad\x1c\x55\x69\xde\xfc\x68\xdf\xef\x58\x26\x75\x50\xae\xe2\x5c\xf8\x39\x18\x37\x75\x19\x2f\xe6\xcf\x72\xff\x0e\xad\x23\xf3\x52\xde\xe3\x13\x21\x0c\xec\x4f\x1e\x5a\x99\xe0\x3a\x9d\x7a\x42\xfa\x96\xbc\xbe\x9c\x0b\x56\xf4\x92\x31\x7c\xba\x7e\xff\x8c\x8a\x2e\x79\xfd\xeb\xc7\x0f\xef\x5f\x12\x28\xda\x7b\x3a\xd1\x35\x19\x2b\x54\xbe\xd8\x55\x7b\x76\x8a\xa6\x1c\xf4\xbb\x43\xc5\x07\x53\x8c\xad\xa1\x24\xa1\x38\xfe\xe2\xc8\xbd\xc8\x52\xcd\xba\x33\xe3\xaf\x31\xb0\x47\x4f\xc6\xad\x8e\xce\xa8\x9b\x4b\xd1\xb2\xa5\x89\xfb\x84\x6b\xda\x9a\x9c\xfd\x5f\x84\x9f\x32\xaf\x4e\x03\x7b\x3a\x57\x7e\xdd\xfd\x5b\xf1\x26\xe8\x4b\x72\x64\x2e\x9a\x67\xcc\xfc\x1f\x94\xc3\xff\x04\x6b\xfc\x76\xfa\x8c\xce\x9a\xf4\xcd\xff\x1d\xcb\x1e\x1f\xdd\x77\x8e\x5a\x89\xd3\x7f\xc1\xc4\x6f\xb8\xcf\x3b\xeb\x8f\x19\x4a\xdb\xd7\x3e\x33\xdc\x51\xb3\x3d\x4b\xee\xd1\xd5\x54\xc0\x97\xfc\x2c\x0d\x07\xfa\xa6\xe4\x0d\x24\x3f\xbf\xd5\x24\x4d\x27\x68\x2e\x93\x7a\x12\x02\x04\x2c\xd5\xde\xd3\x96\x02\xe8\x7f\x8e\x0f\xf9\xde\xa6\x03\x0f\x94\x05\x5c\x30\x3b\x42\x3f\x3e\x70\x60\xd2\xaa\xb2\x2c\xaa\x8e\xbf\x7b\xf3\xa6\xf1\x2e\xfe\x87\xac\xdd\xba\x08\xae\xbc\x6b\x60\xc5\x21\xd8\x5b\x47\x60\xe7\x9d\x68\x29\x97\xd6\x13\x78\xd6\x4d\x64\x9c\x34\x16\xc2\xdc\x46\xf9\x92\xb6\x9d\x03\xc7\xf0\xd7\x36\x44\x05\xa9\xb5\xf8\xb0\xa5\x1f\xb3\x7c\xfc\x87\x44\xec\xa8\xdb\x57\xad\xa5\x4f\xe2\xba\xca\xbc\xde\x17\xec\x7e\x51\xbc\x91\x47\xed\xdd\x0a\x51\x5b\x0e\x68\x86\x25\x3a\xdf\x93\x02\x1b\x52\xf4\x29\xdc\x43\x15\x49\x4e\x0e\x96\xf0\x0c\x26\x8e\x50\xda\x9c\x38\x0c\x0a\x41\x59\xc8\x3c\x82\xe0\xa8\xf7\x68\xfc\x88\xf3\x5f\x48\x8f\xbb\x6e\xbd\x44\xdd\x4a\x3e\xd4\x42\xa1\xd5\xe5\x0f\xea\x76\xd3\x05\x89\x91\x4c\x38\x07\xcf\x1e\x6e\x1d\x97\x77\x9b\x1b\x3b\xc4\xba\xa6\xb9\x50\x58\xf6\xc1\x4e\xea\xd4\x1d\xe8\x8a\x83\x82\x50\x19\x2b\x37\x84\xbe\x6d\x60\x8d\x56\x63\x39\x37\xe2\x77\x1b\x59\x79\x43\x7d\xde\x5e\xd1\x39\x0b\x3d\x1a\xfd\x22\xef\x3f\xf9\xb2\xfd\x1d\x00\x00\xff\xff\x0e\x2a\x60\x2d\xc1\x11\x00\x00")

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
	canonicalName := strings.Replace(name, "\\", "/", -1)
	return filepath.Join(append([]string{dir}, strings.Split(canonicalName, "/")...)...)
}
