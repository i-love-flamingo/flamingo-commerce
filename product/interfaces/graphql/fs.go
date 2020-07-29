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

var _schemaGraphql = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xec\x57\xcd\x6e\xdb\x38\x10\xbe\xe7\x29\x18\xef\x25\x05\x8a\x3e\x80\x6f\x89\xdb\x2c\x82\xc6\x45\xb6\xf6\xee\x65\x11\x04\x63\x6a\x2c\x13\xa1\x48\x2d\x49\x39\x11\x16\xfb\xee\x0b\xfe\x48\x26\x45\x49\x76\x6f\x3d\x34\xa7\x78\xe6\x9b\x1f\xce\xbf\x98\x30\xa8\xf6\x40\x91\xac\x64\x55\xa1\xa2\xf8\xf2\xa4\x64\xd1\x50\x43\xfe\xbd\x22\x84\x90\x1d\x68\xfc\x0c\x06\x96\x27\xc0\x1d\x68\x46\x03\xca\xb2\xae\x1d\xd0\x20\x68\x54\x03\x68\x40\x6d\x7b\x9e\xc7\xea\x1a\x29\xdb\x33\x0a\x86\x49\xa1\x73\xfc\x26\xe1\x7b\x19\xa6\x37\xc0\x11\x76\x1c\x97\xe4\x4e\x4a\x8e\x20\x82\xb2\x40\x1e\x37\xdd\x09\x05\x27\xdb\x1a\x97\x64\x63\x14\x13\xa5\xa7\x94\x68\x1e\x0a\x14\x86\xed\x19\xaa\x94\x75\x00\xbd\xc6\x82\xc1\x4d\xa9\x64\x53\xf7\xbc\x8f\xa4\xd1\x50\x9e\xd4\x7c\x18\xf8\x53\xa2\xb9\x50\x6c\xe8\xab\x13\xbb\xbe\xfa\xef\xea\xca\xfa\x79\x62\x6f\x58\x55\x73\xec\xf2\xe2\x7e\x54\x28\x8c\xfe\x95\xb3\x9f\x29\x67\x83\xa4\xad\xa4\xd8\xb3\xb2\x51\xf6\x21\xbf\x52\xf7\x13\xa7\xce\x2a\x39\x82\x62\x20\x8c\x5e\x92\xbf\x87\xa8\x97\xbf\x3c\xef\xfa\x39\x46\x3a\xa2\x0d\xd7\xad\x31\x8a\xed\x1a\x83\x56\x36\xd8\x7a\xce\x5b\x78\xa0\xec\x47\xb3\x7d\x69\xd0\x33\xbb\x43\x75\xc1\xb0\x61\xc6\x26\x36\xfa\x8b\x33\x01\xd1\x9b\xfc\xdf\xd0\xe2\xe9\xd5\xde\xbd\x83\x54\xe6\x33\x6a\xaa\x58\x6d\x83\x92\x26\xb6\x88\x19\x99\xb1\xca\x66\x21\x71\x25\x4b\x81\x4f\xd4\xf3\x95\xc7\x83\x7a\x45\xf3\xc4\x81\xe2\x4a\x16\x83\xf2\x52\x68\x80\x71\x54\x9e\x33\xb0\xd4\x31\x37\xaf\xcd\x32\x73\xa3\x63\x7e\x83\x2a\x95\x74\x5c\xaa\x10\x0c\x16\xb7\xc6\xb2\xb6\xac\x42\x47\x6d\xea\x62\x84\x7a\x64\x9a\xed\x38\xde\x2b\x59\x2d\x33\xea\x56\xf6\x58\x47\x5e\x81\xc1\x52\x2a\xe6\x43\x7d\x7a\x79\xa0\xb7\xbe\x95\x43\xed\xad\x81\x89\x8e\x11\x55\xc1\x00\x1b\x3c\xee\x88\xd2\xc6\x62\x0d\x75\xcd\x44\x19\x97\xa8\x4f\x9c\x91\xf4\xf5\x11\x8f\xc8\x97\xe9\x83\x5f\xb1\x7d\x93\xaa\x48\x8a\xda\x8f\x85\x6f\xf8\xe6\xe2\xd3\x37\xe1\x54\xad\x9f\xa6\x50\x28\x3a\x57\x25\x5b\x5f\x79\x5e\xe9\x6c\xf1\x44\x63\xee\x49\x31\x8a\x24\x2f\x7c\x47\x7f\x10\x7b\x39\xc4\x3e\x68\x1b\x7e\xf7\x6f\x3f\x2e\x1c\xa6\x56\xb8\x41\x8e\xd4\x60\x11\x5a\xd1\x55\x43\x64\xd1\xd5\x23\x19\x19\x04\x5d\x15\xce\x14\x61\xe4\xc6\xed\x11\x18\xf7\xb3\x9f\x51\xd4\x63\x0a\x7b\xef\x83\x52\x2f\xf8\x28\x5b\xe0\xa6\xed\x99\xf9\xab\x87\x08\x27\xbc\x48\xa4\xbf\x80\x12\x4c\x94\x96\x4b\xf6\x52\x11\x73\x40\x42\x1b\xa5\x50\x18\x52\x7b\x2d\x8b\xdc\x66\x24\x35\x69\x35\xc2\x4c\x66\x7e\xe8\x60\x37\x74\xb2\x6d\x50\xe0\x1e\x1a\x6e\x12\x63\x8c\x62\xb7\x81\x3e\x33\x4d\x65\x23\x0c\x16\x83\x99\x5f\x44\x8c\x31\xd1\x8e\xbf\xc5\x77\x93\x5a\xac\x98\x78\x92\x4c\x18\xbd\x95\x9b\x1a\x85\x59\x92\x7b\x2e\xc1\x04\x26\xbc\x4f\x33\xa9\x14\xc6\xa9\x4b\x0d\xae\x3c\xd9\xb5\xc1\x62\x24\xfa\xfa\x20\xdf\xb4\x8b\xbf\x0b\x15\x88\xc2\xfd\xa8\x9d\x1d\x82\xa0\x04\x16\x8b\xd9\x30\xc6\xda\x7c\x20\x17\xdb\x4e\x9d\xdc\x3b\x6d\x39\xf2\x23\xc1\x4f\xe5\x27\xb2\x66\x1c\xf5\xad\x28\xd6\x52\xe1\x62\x22\x09\x4e\xdb\x11\x78\x33\xab\xce\x57\x0f\x6d\x09\x05\x41\x76\xe8\xd5\x87\x57\x48\x45\x2a\x6b\x68\x31\x9f\xd4\x91\x72\x39\x05\x30\x3c\x8d\x36\xda\xc8\x0a\xd5\xef\xc9\x32\xf7\xac\x03\x08\x81\x3c\x1f\xfa\x5c\x52\xe0\x11\x6d\xaa\x2c\xd3\x33\x27\x18\x74\x57\xc3\xd8\xde\x4f\xd0\xce\x9d\x99\xbd\x9e\x83\xd3\x45\x1b\xbb\x8b\xc2\xf8\x79\x3f\x6f\xf2\x8b\x30\xaa\xbd\xd4\xa4\x03\x07\x93\x1c\x76\xf1\x38\xf7\x47\x0b\x6f\x2e\x3b\x50\xba\x53\x22\xe8\x3a\x7b\x05\x0e\xb6\x5c\x47\xb6\x6b\xae\x27\x02\x35\xec\x88\x61\x18\xcf\x8f\x70\x48\xa7\xe6\x05\x43\x93\x47\xb3\x66\x0c\x3f\x9c\x45\x41\x6c\xc1\x7f\x70\x50\xa6\xf8\x19\x4b\x91\xc6\x99\x38\xbb\x65\x92\xcc\xc5\xfc\x30\x62\x15\x6e\x3d\x2b\x26\x87\xb3\x76\x88\x8e\x2e\xba\xf4\xac\xd9\xa3\xed\xdb\x0b\xda\xe3\x74\xd3\x05\xbf\xfa\x33\xf0\x2b\xb6\xd9\x21\x10\xdf\x88\x59\x28\x7a\x55\x01\x7c\x00\xdd\x93\x6e\x5e\xb1\x1d\x39\xe6\xbb\x5b\x7e\x12\x37\x69\x23\x93\xd4\x77\xed\x57\x6c\xad\x7c\xec\xf5\x87\x33\x7e\x9e\x0d\x4b\x37\xa1\xb2\xf9\x63\x29\x8f\x79\xd3\x8d\xf4\x61\x23\x98\xc9\xe7\xd7\x05\xdd\x99\x9e\x78\x93\x9e\xd4\x60\x0e\x29\x45\xb8\x73\x36\xc5\x28\xb7\xda\x26\x74\x4f\xc6\x61\xb8\xcd\x27\xc6\xfc\x99\xfd\x7c\x66\x3d\xfb\x51\x71\x07\x1a\x93\xdd\x7b\x22\xdf\x56\x56\x70\x82\xf9\xa7\x60\x26\x39\xc5\x46\x4f\x88\x70\x1e\x57\x35\xb0\x52\x7c\x6f\x38\x66\xb5\x5d\xa0\x68\xed\xc6\xec\x84\x75\x2a\xfb\xdb\xf9\x73\xc0\xb7\x24\xbc\xaf\x38\x68\x1d\x9d\xd6\x93\xdf\x85\x1b\x04\x45\x0f\xdf\x51\x37\xbc\x5b\x85\x61\xfc\x8c\xf5\x57\xf0\x73\x0f\x14\x53\xbe\x57\xf3\x72\x6f\x19\xd7\xcf\x61\x54\x37\x65\x89\x3a\x7c\xfa\x67\xd0\x4d\xcf\x0d\x4a\xb5\xa3\xaf\x31\xf9\xd0\x0c\x60\x4b\xed\x3f\xd6\xbb\x3b\xda\x19\x5b\x26\xdf\x03\xf8\x6e\xd0\x1e\x3a\xf6\xa9\x7f\x34\xd8\xaf\xa6\xe1\x3b\x6e\xfc\x1d\x5d\x67\x1f\x73\x23\xfd\x3e\xaa\x20\x38\x76\xa3\x43\xf8\xfe\x69\x50\x9b\xdc\xef\xc0\x18\xd1\x9a\x04\xde\xfa\xfe\x7f\x00\x00\x00\xff\xff\x9e\x2b\x07\xe4\x11\x15\x00\x00")

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
