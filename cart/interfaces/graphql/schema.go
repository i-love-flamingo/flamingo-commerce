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

var _schemaGraphql = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x84\x96\xc1\x6e\xdb\x3c\x0c\xc7\xcf\x0d\x90\x77\x50\xd1\xcb\xf7\x01\x7d\x02\xdf\x52\x7b\x18\x02\x6c\x5d\x86\x66\xa7\x61\x28\x58\x8b\x49\x04\xc8\x92\x27\xd1\x45\x8d\x61\xef\x3e\xd8\x92\x1d\x4b\xb2\x9d\x1c\xcd\x9f\xc8\x3f\x29\x92\x0a\xb5\x35\xb2\x5c\x57\x15\x9a\x12\x5f\x0b\x2c\xb5\x01\x42\x9e\x83\x21\xf6\x67\xbb\x61\x8c\xb1\x12\x0c\x65\x57\xa6\x33\xdd\x3b\x0b\x1f\xf0\x57\x8e\x52\xbc\xa3\x11\x68\x33\xf6\x33\x40\x47\x97\x85\x43\xda\xfb\x5f\xdb\xcd\xdf\xed\x66\xbb\x09\x43\x4f\x23\x0a\x9e\xb1\x7d\xe1\x83\xa0\x22\x41\xed\xbe\xc8\xd8\x0b\x19\xa1\xce\xfe\xf3\x9b\x90\x52\xa8\xf3\x8e\x1b\xb4\x36\x92\xb7\xe3\xfd\x57\x07\xd6\x8d\x29\x2f\x60\xd1\x44\xd0\x01\x8d\xd5\x6a\x48\x64\x59\xfe\x55\x75\x47\x3e\xdc\x01\xe7\x82\x84\x56\x20\x0b\x20\x48\x23\x4f\x8c\xf7\xc3\x99\x1a\xda\x0a\x15\xbd\xa0\xc4\xb2\x33\xc7\x52\x22\xf3\x90\x22\x4a\xad\xce\xf6\xa8\x77\x0d\x5d\xba\x32\x94\x5d\x1d\x7f\xf4\xa9\x3c\x69\x2d\x11\x06\x12\x62\x20\x29\xd7\xc3\x1d\xd4\xb5\x14\xc8\x73\xdd\xd4\x5a\xe5\x9a\xa7\xb9\x5e\x4d\x43\xb6\x1c\x4f\xd0\x48\xca\x1b\x63\x50\x95\x6d\xe2\x93\x34\x81\x14\x84\x55\xe2\xeb\x38\x58\x26\x85\x73\x02\x3e\x8b\x13\xe5\x60\x78\x72\x64\x17\xda\x57\xfa\x24\xe9\xa9\xa1\x71\xfc\x3d\xb6\x59\xcc\xfb\x4b\x8c\xda\x76\x3f\xa7\xbc\x98\x5a\x57\x45\x84\xb1\xbb\x31\x99\x2d\xc5\x4d\x3f\x93\x78\xe3\x04\x10\x56\x51\x12\x9d\xd9\x77\xb4\xd1\xbc\x29\xa7\x33\x79\x70\x5f\x16\x63\x04\xae\xa7\xc3\xf5\x9b\xda\x8c\xed\xd5\xf2\x49\x3f\x4a\xc3\xe1\x77\xa0\x8c\xcd\xfc\x82\xbe\x38\x09\x63\x49\x41\x85\xd9\x0a\x24\x61\x96\x09\xa1\x4a\x70\x2e\xf1\x39\xc1\x02\x88\x04\xc9\xc4\x4d\x0c\x59\x90\x0d\x81\x9b\xbd\x65\x88\x0c\xe2\x4c\x82\x33\xd0\xb3\x59\x15\x7e\x5d\x12\xbe\x80\x5f\x84\xea\x07\xce\x53\x7e\x28\x4a\x5d\xd5\xa0\xda\x24\x62\xe0\xab\x14\x94\x12\x31\x54\x6b\x4b\xdd\xe4\xae\xaa\xb2\x04\x74\xb3\x52\x06\xcf\xc2\x6d\x81\x95\x4a\x95\xba\x51\x64\x6e\x09\x77\x50\xe2\x2a\xbc\x3d\x94\x58\x5f\xb4\x5a\x6d\x16\xac\x40\xc8\x35\xe1\x0b\xdd\xeb\x76\xfc\xd0\xbc\xc0\x6f\x3f\x16\xfd\x01\x90\x05\x12\x08\x19\xa3\x87\xd0\x3a\x88\xfb\x10\x96\x84\x3a\xe7\x8d\x25\x5d\xa1\x99\x79\x16\x3e\xcd\x20\x8b\xa2\xe7\xe0\x64\x78\x57\x13\x1e\x15\x8e\x7b\x11\x08\xbf\x9d\x9e\x84\xa1\x4b\xb4\xc3\x6b\xb0\xb6\xd6\xfd\xee\x77\x17\x3a\x6b\x7d\x6e\xaa\xb7\xee\xd5\x09\x8c\x0a\x5c\x87\xf7\xed\x19\xde\x04\x7e\x10\x2a\xce\x7a\x7d\xdf\x9b\xe9\x92\x9c\x8a\xcd\x16\xfe\x78\xcc\x79\xf9\xea\xa7\x37\x71\xb4\xe3\xfc\xa8\xbb\x53\xff\xf9\xda\x3c\x8e\x4b\xed\x71\x7c\x0f\x5c\x0b\x7a\x8d\xff\xaf\xc7\xfd\x17\x00\x00\xff\xff\x8d\xa7\x0b\x8a\x12\x09\x00\x00")

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

	info := bindataFileInfo{name: "schema.graphql", size: 2322, mode: os.FileMode(438), modTime: time.Unix(1563887042, 0)}
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
