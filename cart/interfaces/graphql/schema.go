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

var _schemaGraphql = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x84\x56\x41\x6f\xdb\x3c\x0c\x3d\x37\xbf\x42\x41\x2f\xdf\x07\xf4\x17\xf8\x96\xda\xc3\x10\x60\xcb\x32\x34\x3b\x0d\x43\xc1\x5a\x4c\x22\x40\x96\x3c\x89\x2e\x6a\x0c\xfb\xef\x83\x2d\x25\x91\x25\xd9\xc9\x51\xef\x89\x7c\xa4\x1e\x19\x53\xdf\x22\x2b\x75\xd3\xa0\xa9\xf1\xb5\xc2\x5a\x1b\x20\xe4\x25\x18\x62\x7f\x56\x8c\x31\x56\x83\xa1\xe2\x46\x19\x90\xf5\x08\xf0\x0b\xf9\x95\xa3\x14\xef\x68\x04\xda\x82\xfd\x9c\x30\xaf\x01\x2b\x47\xe9\xd7\xbf\x56\x7f\x57\xab\x69\xd6\x20\x99\xe0\x05\xdb\x56\x2e\x3e\x2a\x12\xd4\x6f\xab\x82\xbd\x90\x11\xea\xe4\x4e\xdf\x84\x94\x42\x9d\x36\xdc\xa0\xb5\x91\xae\x0d\x1f\x4f\x47\x5e\xdb\x99\xfa\x0c\x16\x4d\xc4\xd9\xa3\xb1\x5a\xf9\x02\xe6\x65\x5f\xd5\x0e\xc4\xc7\x07\xe0\x5c\x90\xd0\x0a\x64\x05\x04\x69\xda\x00\x5c\xfb\x2b\x2d\xf4\x0d\x2a\x7a\x41\x89\xf5\x80\xc6\x3a\x22\xd8\x57\x87\x52\xab\x93\x3d\xe8\x4d\x47\xe7\xa1\x01\xf5\xd0\xbc\x1f\x63\x19\xcf\x5a\x4b\x04\x4f\x84\x18\x8f\xfb\xf4\xf8\x00\x6d\x2b\x05\xf2\x52\x77\xad\x56\xa5\xe6\x69\x99\x37\xc8\x17\xca\xf1\x08\x9d\xa4\xb2\x33\x06\x55\xdd\xc7\x11\x49\x13\x48\x41\xd8\x24\x91\x0e\x17\xe4\xd6\x31\x97\xfd\xb3\x38\x52\x09\x86\x27\x37\x36\x53\x7c\xce\x17\x89\x81\xbc\x51\xfc\xdb\xf5\x45\x4c\xf7\x0f\x37\x75\xe8\x36\xa7\xb9\x0a\xd1\xf9\xfc\x93\xb4\xc3\x30\x64\x1b\xb0\x1c\x23\x48\x74\x31\x3a\x61\x13\x69\x1f\x50\x67\x5d\xa3\x79\x57\x87\x43\xb7\x77\x27\xf9\xf0\x61\xd4\x60\x7c\x7e\x53\x5f\xb0\xad\x9a\xb9\xe4\x47\xc5\xdf\x7b\x07\x2a\x58\xe6\x17\x3e\xff\x51\x18\x4b\x0a\x1a\x2c\xe6\x39\x12\xb2\x94\x09\xa7\x11\x9c\x4b\xdc\x25\xac\x90\x43\x82\x64\x12\x24\xe2\x58\x90\x1d\x81\x9b\xac\x59\x0e\x19\xc4\x4c\x69\x29\x67\x67\x96\x34\xdf\xc6\xdf\xf7\xed\x8b\x50\xe3\x38\x79\x92\x33\x7d\xad\x9b\x16\x54\x9f\xa4\x0b\x23\xd5\x82\x52\x42\xc4\x69\xb5\xa5\x61\x28\x97\x14\x59\x02\xba\xd7\x21\x83\x27\xe1\xc6\x7b\xbe\x43\xb5\xee\x14\x99\x3b\x9a\x1d\x27\x09\x34\x79\x31\x94\xd8\x9e\xb5\x5a\x72\x07\x36\x20\xe4\x82\xe6\xac\x51\xdd\xbe\xf6\x3e\x05\x7e\x77\xed\x8f\x74\x90\x15\x12\x08\x19\x33\xf7\x53\xd4\xcb\xfa\x10\x96\x84\x3a\x95\x9d\x25\xdd\xa0\xc9\x6c\xf8\x4f\x19\x4a\x5e\x6e\x8e\x19\x0d\xe7\x42\x99\x57\x65\x97\x2d\x07\x84\xdf\x8e\xcf\xc2\xd0\x79\xba\x8b\x5b\xb0\xb6\xd5\xe3\x06\x77\xcf\x97\x03\x77\x5d\xf3\x36\xfc\x71\x84\x98\x02\xe7\xe3\xd1\x86\x41\xe3\xf1\x83\x50\x71\x36\x0a\xfb\xde\x05\x0b\x2f\x14\x59\xcc\x7c\x29\x24\x11\xbe\xfa\xc9\x8c\x83\x6c\x38\x3f\xe8\xe1\xc6\x7f\xbe\x19\x4f\xd7\x35\xf5\x74\x5d\xe9\xce\x67\x5e\xdb\xff\x4b\x39\xff\x05\x00\x00\xff\xff\xf1\x72\x43\xea\xbd\x08\x00\x00")

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

	info := bindataFileInfo{name: "schema.graphql", size: 2237, mode: os.FileMode(420), modTime: time.Unix(1563984470, 0)}
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
