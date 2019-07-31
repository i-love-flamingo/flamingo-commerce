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

var _schemaGraphql = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x94\x57\xc1\x6e\xdc\x36\x10\x3d\xc7\x5f\x21\xc3\x97\x04\x08\xfa\x01\xba\xad\xa5\x22\x58\x34\xd9\xba\xf1\x16\x3d\x04\x46\x30\x16\x67\x77\x59\x53\x24\x4b\x8e\x62\x0b\x45\xff\xbd\xa0\xc4\xd5\x92\x14\xa5\x6d\x7d\x30\x0c\xce\xe3\xf0\xcd\xf0\xcd\x13\x4d\xbd\xc6\xa2\x52\x6d\x8b\xa6\xc1\xef\x35\x36\xca\x00\x21\xab\xc0\x50\xf1\xf7\x4d\x51\x14\x45\x03\x86\xca\x0b\xc4\x45\x6e\x87\x00\x3b\x83\xbf\x33\x14\xfc\x07\x1a\x8e\xb6\x2c\xbe\x45\xc8\x29\x61\x3d\x42\xfa\xdb\xa7\x9b\x7f\x6e\x6e\xe2\x53\x83\xc3\x38\x2b\x8b\x6d\x3d\xe6\x47\x49\x9c\xfa\x6d\x5d\x16\x8f\x64\xb8\x3c\x8e\xab\xcf\x5c\x08\x2e\x8f\x1b\xc6\x0c\x5a\x9b\x10\xf3\xab\x03\x50\x77\xa6\x39\x81\x45\x93\x60\x1e\xd0\x58\x25\x7d\x05\xcb\xbc\x27\xba\x0e\x78\xf7\x0e\x18\xe3\xc4\x95\x04\x51\x03\xc1\xfc\xd8\x20\x78\xeb\xb7\x68\xe8\x5b\x94\xf4\x88\x02\x1b\x17\x4d\x79\x24\x61\x5f\x1e\x0a\x25\x8f\x76\xaf\x36\x1d\x9d\x5c\x07\x1a\xd7\xbd\xdf\x87\x32\xee\x95\x12\x08\x1e\x08\x69\x3c\x6d\xd4\xdd\x3b\xd0\x5a\x70\x64\x95\xea\xb4\x92\x95\x62\xf3\x32\x2f\x21\x5f\x28\xc3\x03\x74\x82\xaa\xce\x18\x94\x4d\x7f\xc9\x78\xe7\xa2\xa4\x08\x04\x27\x6c\x67\x89\xf6\xe7\x88\xcf\xe3\xfe\xac\x54\x27\xa9\x2c\xb6\xd2\xeb\x45\x1b\xc5\xba\x86\xd2\x65\x6e\xa3\x4e\x20\x4b\x0a\x9d\xea\xf8\xc4\x0f\x54\x81\x61\xb3\xc3\x37\x71\x7c\x49\x62\x33\x2d\x7a\xcd\x79\x15\xf4\x65\x0a\xf7\x12\x88\xc5\xbe\xcd\x95\x5f\x87\xd1\xe5\xf3\xb3\xc7\x6e\xe5\x41\x2d\x1c\xed\x42\xd3\x0c\x66\xfb\xbe\xbd\xb4\xdc\x9e\xb8\xd6\x5c\x1e\xdd\x52\x92\xef\x31\x08\x8d\xd8\xee\x79\xb8\xb1\x4f\x46\x45\x33\xf4\x60\x78\x83\x03\xe2\x68\x40\xb2\x01\x93\x0d\xdb\xae\x1d\x82\x7b\x78\xdb\xb4\xe3\x85\x66\x51\xe3\x31\x3b\x5c\x8a\x8f\x59\x6a\x6e\x1b\x97\x64\x35\x55\xbb\x53\xd2\x15\xf0\x15\xc5\x70\x89\xff\x69\xcf\xff\xdd\x10\xb4\xe5\x0f\x4e\xa7\xf3\x9e\x7c\x8f\x82\xf2\xae\x83\x4f\x60\xbd\x74\x26\x75\xaf\xaa\xc4\xdd\xfd\xd9\x81\x15\xc3\x78\xb6\x5f\x95\x79\x39\x08\xf5\x1a\xaf\xb6\x48\x27\xc5\xe2\xb5\x06\x8c\xe1\xce\x3c\xc2\xc5\xb3\xf4\x3e\xab\x06\x32\xce\x54\x27\x61\xbf\xc7\x72\x83\x6c\xcf\x5b\x2c\x0b\xf7\x7b\x94\x09\x26\xe6\xf7\xfe\x05\x2f\xae\xf1\x21\x3e\x36\xf6\xd0\x5f\xb0\x77\x72\xf6\x80\xa7\xd1\x61\x02\x48\xd0\x07\x5b\x16\x2d\xe8\x6f\x76\x80\x3e\xfd\x69\x95\xfc\xe9\x2b\xbc\x7e\x41\x6b\xe1\x88\xeb\x6d\x3c\xd7\x50\xf8\x5e\x3a\xe4\x8c\xd5\x95\x2f\x49\x67\xf1\x3e\xf9\xea\x44\x0e\x15\xdf\x4f\x96\x4e\x38\x81\x67\x26\x9c\x44\x42\x45\x3b\xb5\x2c\x0d\x0b\xad\x8e\x1a\x5b\x57\xf7\xba\x21\x06\xac\xf8\xdc\x3d\x26\xd7\xf0\x06\x1e\x25\x1f\x56\xf2\xe9\xc3\xac\xc1\x67\xfd\x2f\xea\x07\xf7\xcf\x6f\xf2\x1d\xf6\xfb\x7e\x00\x95\x45\xe6\x27\x6c\xda\x81\x1b\x4b\x12\x9c\x2c\x17\x31\x02\xb2\x90\x78\x7a\x38\x63\x02\x77\x33\x54\x88\xf1\x77\xb6\xca\xc7\x82\xe8\xc8\x8f\xd5\x22\x86\x0c\x62\xa6\xb4\x39\x66\x67\xd6\x38\x5f\xc6\xc5\xf7\xed\x33\x97\x18\x4d\xd5\xa8\xcf\x56\x83\xec\x67\xc7\x45\x3e\xc1\x69\x0e\x48\x30\x5a\x59\xaa\x06\xad\xaf\xb1\x06\xba\xd6\x21\x83\x47\x3e\xbe\x3a\x96\x3b\x34\x88\xd9\x5c\xe1\x3c\x62\x66\x89\xa2\x1b\x43\x81\xfa\xa4\xe4\x9a\x3a\xb0\x05\x2e\x56\x38\x67\x85\x3a\x3e\x23\xbd\x4e\xaf\x7b\x88\x1e\xe0\xce\xd6\x08\xb8\x48\x91\x0f\x71\xd4\xd3\x7a\xe3\x96\xb8\x3c\x56\x9d\x25\xd5\xa2\xc9\x3c\x3c\x7f\xce\x40\xf2\x74\x73\xc8\x64\x38\x57\xca\x9c\x98\x9d\xdf\x2e\x40\xf8\xeb\xe1\x9e\x1b\x3a\x25\x0e\x06\xd6\x6a\x65\xc6\x47\x9e\xe9\xf3\xc1\x5d\xd7\x3e\xa7\x9f\x24\x09\xa3\x8e\x07\x19\x06\x8d\xc7\x37\x42\xc9\x06\xdb\x2e\x7e\xeb\x2e\xaf\xa7\x26\x24\x59\x2e\xfc\x07\x33\xcb\xf0\xc5\x4f\x66\x9a\x64\xc3\xd8\x5e\xb9\x1d\xef\x5b\x30\x2f\x48\x5a\x40\x83\xa3\xae\xb6\xf5\xed\xc7\xc9\xb3\x3e\x4e\x9f\xce\x2a\xf4\xfc\x0f\x6b\x04\xfe\x0d\x00\x00\xff\xff\x6e\xc8\xa9\xc9\x62\x0d\x00\x00")

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

	info := bindataFileInfo{name: "schema.graphql", size: 3426, mode: os.FileMode(420), modTime: time.Unix(1567424646, 0)}
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
