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

var _schemaGraphql = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x84\x56\xcd\x6e\x1c\x37\x0c\x3e\xc7\x4f\x21\xc3\x97\x04\x08\xfc\x00\x73\xb3\x67\x8b\x62\xd1\xc4\x75\x6b\xf7\x14\x18\x01\x33\xe2\xee\xb2\xd6\x48\xaa\xc4\x89\x3d\x28\xfa\xee\x85\x66\xb4\xbb\x92\xe6\x67\x7d\x30\x16\xe4\x27\xf2\x23\xf5\x91\x1a\xee\x2d\x8a\xda\xb4\x2d\xba\x06\xbf\x6f\xb0\x31\x0e\x18\x65\x0d\x8e\xc5\xbf\x57\x42\x08\xd1\x80\xe3\xea\x0c\x09\x9e\xeb\xc1\x21\x8f\xe0\xef\x12\x15\xfd\x44\x47\xe8\x2b\xf1\x2d\x43\x9e\x02\x6e\x46\x48\x7f\xfd\x72\xf5\xdf\xd5\x55\x9e\x35\x49\x46\xb2\x12\xdb\xcd\x18\x1f\x35\x13\xf7\xdb\x4d\x25\x9e\xd8\x91\xde\x8f\xd6\x1f\xa4\x14\xe9\xfd\x9d\x94\x0e\xbd\x2f\x88\x45\xeb\x00\xb4\x9d\x6b\x0e\xe0\xd1\x15\x98\x47\x74\xde\xe8\x58\xc1\x32\xef\x13\xdd\x00\xbc\xf9\x00\x52\x12\x93\xd1\xa0\x36\xc0\x30\x4d\x9b\x38\xaf\xe3\x11\x0b\x7d\x8b\x9a\x9f\x50\x61\x13\xbc\x25\x8f\xc2\x1d\xcb\x43\x65\xf4\xde\x3f\x9b\xbb\x8e\x0f\xa1\x03\x4d\xe8\xde\x5f\x43\x19\xf7\xc6\x28\x84\x08\x84\xd2\x5f\x36\xea\xe6\x03\x58\xab\x08\x65\x6d\x3a\x6b\x74\x6d\xe4\xb4\xcc\xb3\x2b\x16\x2a\x71\x07\x9d\xe2\xba\x73\x0e\x75\xd3\x9f\x23\xde\x04\x2f\x1b\x06\x45\x8c\xed\x24\xd0\xf3\xd1\x13\xe3\x84\x9f\xb5\xe9\x34\x57\x62\xab\xa3\x5e\xac\x33\xb2\x6b\xb8\x34\x93\xcf\x3a\x81\xb2\x28\xf4\x54\xc7\xaf\xb4\xe3\x1a\x9c\x9c\x24\xbf\xcb\xfd\x4b\x12\x9b\x68\x31\x6a\x2e\xaa\xa0\xaf\x4a\x78\x94\x40\x2e\xf6\xed\x5c\xf9\x9b\xd4\xbb\x9c\x7f\x36\xed\x56\xef\xcc\x42\xea\xe0\x3a\xcd\xe0\x6c\xdf\xb7\xe7\x96\xfb\x03\x59\x4b\x7a\x1f\x4c\x45\xbc\xa7\xc4\xb5\xce\x2d\x64\x3c\xce\xbd\x91\x98\x2b\xea\xcd\xb8\xd7\x9d\x32\x6f\xb9\xb5\x45\x3e\x18\x99\xdb\x1a\x70\x8e\x82\x64\x53\xe3\xb1\xe0\x2f\xa6\x81\x99\x79\xd8\x14\xee\x78\xc6\x93\x43\xf9\x4c\x2d\x56\x22\xfc\x1f\xac\x7b\x2c\x46\xee\xe3\x2b\x9e\xb5\xfa\x29\x4f\x9b\x4f\xee\x6f\xd8\x87\x26\x46\xc0\xcb\xa8\xeb\x04\x92\xf4\xc1\x57\xa2\x05\xfb\xcd\x0f\xd0\x97\xbf\xbd\xd1\xb7\x7f\xc2\xdb\x57\xf4\x1e\xf6\xb8\xde\xc6\x63\x0d\x22\xf6\x32\x20\x27\xac\x2e\xec\xaf\xce\xe3\x7d\xb1\xeb\xb2\xb9\xc8\xef\x67\x96\x4e\x7a\xef\x47\x26\xc4\x0a\x8b\xb1\xb6\x8e\x1a\xd4\xc8\xd5\xf8\x6b\x63\x5a\x20\x7d\xfb\x18\x7e\xc7\xb9\x87\x77\x68\xc3\xd8\xde\x2e\x21\x24\xf9\x26\x00\x46\xd8\x72\x20\x23\xcd\xd0\x55\xf1\x76\x40\x3d\xa2\x04\x79\x41\xad\x55\x18\x76\x00\xca\x0b\xb3\x9b\x94\x42\x53\xa1\x0f\x02\x4f\x76\x4d\xe2\x7e\x1c\x2d\xf3\xe1\xd3\xa8\xc9\x0b\xf4\x0f\xf7\xc3\xa2\x9a\x3f\x14\xaf\x25\x9e\xfb\x09\x5c\x89\x99\xbf\xf4\xd2\x77\xe4\x3c\x6b\x08\x5a\x5e\xc4\x28\x98\x85\xe4\x23\x47\x52\x2a\x7c\x98\xa0\x52\x4c\xbc\xe8\x55\x3e\x1e\x54\xc7\x71\x16\x17\x31\xec\x10\x67\x4a\x9b\x62\x1e\xdc\x1a\xe7\xf3\x8c\xc5\xbe\x7d\x21\x8d\xd9\x28\x8e\xa2\x6e\x2d\xe8\x7e\x92\x2e\x5b\x2e\xc4\x53\x40\x81\xb1\xc6\x73\x3d\x0c\xc8\x1a\x6b\xe0\x4b\x1d\x72\xb8\xa7\xf1\x81\x5c\xee\xd0\x20\x7c\x77\x81\xf3\x88\x99\x04\xca\x6e\x0c\x15\xda\x83\xd1\x6b\xea\xc0\x16\x48\xad\x70\x9e\x15\xea\xf8\xc5\x13\x75\x7a\x79\xf1\xd8\x01\x1e\x76\x21\x03\xa9\x12\xf9\x98\x7b\x23\xad\x77\xf2\x4c\x7a\x5f\x77\x9e\x4d\x8b\x6e\xe6\x1b\xe9\x97\x19\xc8\x3c\xdd\x39\x64\x31\x9c\x2b\x65\x9e\x98\x1d\x9f\x59\x60\xfc\x7d\x77\x4f\x8e\x0f\xf9\x06\xb6\xe0\xbd\x35\x6e\xfc\x1e\x71\xfd\xbc\xf3\xa1\x6b\x7f\x94\xef\x98\x86\x51\xc7\x83\x0c\x93\xc6\xe3\x3b\xa3\x96\xc3\xae\x17\x7f\x74\xe7\x87\xbe\x49\x49\x56\x0b\x1f\xdb\x93\x08\x5f\xe3\x64\x96\x41\xee\xa4\x7c\x36\xe1\xc4\xc7\x16\xdc\x2b\xb2\x55\xd0\xe0\xa8\xab\xed\xe6\xfa\xf3\x69\x67\x7d\x3e\xbd\xb7\x75\xfa\x50\x7c\x5a\x23\xf0\x7f\x00\x00\x00\xff\xff\x2e\x3d\x27\x6f\x0d\x0c\x00\x00")

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

	info := bindataFileInfo{name: "schema.graphql", size: 3085, mode: os.FileMode(420), modTime: time.Unix(1567424601, 0)}
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
