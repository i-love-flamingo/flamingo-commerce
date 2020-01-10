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

var _schemaGraphql = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xcc\x96\xdf\x6e\xeb\x36\x0c\xc6\xef\xf3\x14\x2c\x72\x73\xce\x80\xb3\x07\xf0\xe5\xb2\x9d\xb5\xc0\xba\x66\x49\xdb\x5d\x0c\x43\xc0\x4a\x4c\x2c\x4c\x96\x0c\x8a\x6e\x6a\x0c\x7d\xf7\x41\xb2\xea\x38\x4d\x13\xa4\x29\xda\x35\x17\x01\xac\x3f\xfc\x7e\xa4\x44\xfb\x1b\xc3\xc4\x57\x15\xb1\xa2\xc5\xa4\x24\xf5\x8f\x6f\x64\x31\xb5\xa8\xe8\x8a\x35\xf1\xc4\x3b\xa1\x07\x01\xa6\x9a\x29\x90\x93\x00\x52\x12\x30\x85\xc6\x0a\xf8\x65\x7a\x52\x0d\x33\x39\x81\x2f\xdc\x38\x67\xdc\xea\x2b\xd4\x31\x00\xf8\x18\x01\xaa\x46\x50\x8c\x77\x23\x69\x6b\x3a\x4a\xec\xdf\x11\x00\xc0\x18\xae\x4b\x82\x09\xb2\x80\x94\x28\x60\x02\xac\xbc\x71\x2b\x10\x0f\x77\xd4\x49\xe8\xb4\x52\x21\x4b\xb1\x89\xfc\x33\x29\xcf\x28\xa4\xe3\xde\x41\xa8\x6e\x47\xa6\x32\x0e\x14\x86\x27\x46\x13\x00\x2d\x13\xea\x76\x18\x37\xcd\x5d\xb8\xa5\x0f\xc5\x3e\x6e\x7d\xd5\xaf\xc9\x4a\x73\x41\x21\xd0\x54\x93\xd3\x91\xd6\xbb\x54\xa3\x90\x86\xfd\x12\x6a\x6c\xab\x58\x2c\x74\x7a\xab\x4c\xdf\xf2\x92\x0a\x5b\x50\xde\x09\x1a\x07\xa8\xb5\x89\xa5\x43\x0b\xa6\x97\x48\xcb\xf6\x02\x25\x9e\xc4\xb0\x48\xff\xa3\xc7\xd1\x68\xb4\xf7\x8c\x07\xf8\xf0\xad\xd3\x00\xbc\xf3\x8d\x24\xe6\x61\xc1\xe2\xbc\xb4\xb5\x51\x68\x6d\x0b\xa1\xf4\x6b\x17\x53\x43\x08\x8d\x22\x0a\x01\x6a\x5c\xd1\xc1\x23\x1e\x6a\x75\x27\x9c\x4b\x91\x2b\x9c\x7f\x7f\x1d\x4a\x6c\xba\xd9\x71\xf6\x77\x17\xe3\x59\xe8\xe2\x59\x0c\xe4\x1d\xf9\xbc\x93\x2a\x34\xb6\x97\xcd\xbf\xb9\xb0\x71\xab\xb3\x58\xb5\x94\xcb\xc1\x32\x0f\x68\x72\x46\x2b\x14\x5a\x63\x5b\xec\xc4\x1b\xa4\x3b\x65\x7f\x6f\x34\x71\xb1\x35\x59\x91\x94\x5e\x17\xbb\x24\xe9\x56\x29\x26\x6d\x64\x82\xac\x93\x18\xfc\x90\xd2\x4e\x07\xf3\xe3\x64\x6b\x2e\xad\xc7\xca\x37\x4e\x06\xc1\xfa\x2c\xa6\x6c\x14\x75\x41\xc5\x88\xa5\xe2\xe5\xdc\x8d\x13\xe2\x65\xbc\x9c\x47\x5e\xb3\x9c\xbf\x09\xdf\x8d\x43\x5b\xc0\x4f\xde\x5b\x42\x97\x06\x1d\x56\x54\xec\x54\xf6\xc8\xc0\x8b\xdf\x69\x0d\xa6\xaa\x2d\x55\xe9\xed\xf3\xbf\xf3\xfc\x89\x46\x3e\x15\xd0\xbc\x51\x2a\x36\xe0\x67\x62\xfa\x8e\x82\xf6\x17\x66\xcf\xef\x8a\x95\x7a\x83\xa2\xcc\xd3\x58\x24\x1d\xbf\xae\x7c\xa5\x5f\x5f\x2c\x19\x2b\x3a\x0d\x75\xfc\x22\xeb\x78\x17\x36\x0d\x8d\x6f\x66\xbf\xf5\xac\xe3\x93\x60\xcf\xa5\xb2\x1f\x82\x7a\x7e\x7d\xf9\x16\xd6\x19\x69\xc3\xa4\x4e\x6c\x96\x37\x96\xf5\xd5\x2f\xb0\xc5\x04\x9d\x22\x6b\x93\x55\x99\x11\x06\xef\xf2\x3d\xe4\xf4\x30\xbc\x60\xaf\x2a\xc3\x6e\xdc\x45\xfe\x6e\x9c\xda\x1f\x1f\x89\x7a\x8b\xd6\xe8\x34\xf0\xee\xb4\x71\xe8\xbe\x97\x9b\x25\x87\x39\xb4\x39\xf1\x4b\x7e\xfb\x6c\x3e\xb9\x1b\x7a\x10\x72\x1a\x52\xaa\x7f\x34\xc4\x6d\x6f\x20\x7f\xa5\xec\x57\x2d\x06\xd9\x72\x5b\x2a\x5b\xcd\xb5\x91\x12\x8c\x84\x6c\xbe\x36\x2e\xeb\x08\x9f\x7a\xd8\x83\xe5\x45\x67\x3b\x88\x97\xd9\x0f\x3f\x51\x5e\x39\xdb\x42\xed\x43\x30\x77\x96\xc0\x2c\x3b\x07\x59\x61\x50\xa5\x71\x04\xce\x0b\xa0\x12\x73\x1f\xd1\xa3\x6d\x45\x58\xc6\xa6\xe8\x88\xf7\xb0\xce\x05\x59\x36\x2c\x5f\x98\xa4\x61\x77\xc3\xb6\xef\x9b\xaf\x7d\x53\xe5\xd7\x28\x4c\x07\x08\xb9\x1a\x21\xa9\x27\xb9\x3d\x42\xdd\xe1\x6e\x94\x06\x51\x1f\x47\xff\x05\x00\x00\xff\xff\x50\xaa\x11\x7d\x5b\x0c\x00\x00")

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
