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

var _schemaGraphql = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x84\x56\xcd\x6e\xe4\x36\x0c\x3e\x6f\x9e\x42\x41\x2e\x2d\xb0\xc8\x03\xf8\x96\xd8\x45\x31\xc0\x26\x4d\x9b\xf4\xb4\x08\x16\x5c\x8b\x33\xc3\xae\x2c\xa9\x12\x9d\x89\x51\xf4\xdd\x17\xb2\xe5\x19\x4b\xfe\x49\x0e\x19\x43\xfc\x4c\x7e\xa4\x3e\x92\xe6\xce\xa2\x28\x4d\xd3\xa0\xab\xf1\x5b\x85\xb5\x71\xc0\x28\x4b\x70\x2c\xfe\xbb\x12\x42\x88\x1a\x1c\x17\x17\x48\xb0\x5c\xf7\x06\x39\x82\xbf\x49\x54\xf4\x86\x8e\xd0\x17\xe2\x6b\x82\x3c\x3b\xac\x06\x48\x77\xfd\x7a\xf5\xff\xd5\x55\x1a\x75\x12\x8c\x64\x21\x76\xd5\xe0\x1f\x35\x13\x77\xbb\xaa\x10\xcf\xec\x48\x1f\x86\xd3\xef\xa4\x14\xe9\xc3\x9d\x94\x0e\xbd\xcf\x88\xc5\xd3\x1e\x68\x5b\x57\x1f\xc1\xa3\xcb\x30\x4f\xe8\xbc\xd1\x31\x83\x75\xde\x67\xba\x01\x78\xf3\x09\xa4\x24\x26\xa3\x41\x55\xc0\x30\x0f\x3b\x31\x5e\xc7\x57\x2c\x74\x0d\x6a\x7e\x46\x85\x75\xb0\xe6\x3c\x32\x73\x4c\x0f\x95\xd1\x07\xff\x62\xee\x5a\x3e\x86\x0a\xd4\xa1\x7a\x7f\xf7\x69\xdc\x1b\xa3\x10\x22\x10\x72\x7b\x5e\xa8\x9b\x4f\x60\xad\x22\x94\xa5\x69\xad\xd1\xa5\x91\xf3\x34\x2f\xa6\x98\xa8\xc4\x3d\xb4\x8a\xcb\xd6\x39\xd4\x75\x77\xf1\x78\x13\xac\x6c\x18\x14\x31\x36\x33\x47\x2f\xa3\x25\xfa\x09\x8f\xa5\x69\x35\x17\x62\xa7\xa3\x5e\xac\x33\xb2\xad\x39\x3f\x26\x9f\x54\x02\x65\x96\xe8\x39\x8f\xdf\x69\xcf\x25\x38\x39\x0b\x7e\x97\xda\xd7\x24\x36\xd3\x62\xd4\x5c\x54\x41\x57\xe4\xf0\x28\x81\x54\xec\xbb\xa5\xf4\xab\xa9\x75\x3d\xfe\x62\xd8\x9d\xde\x9b\x95\xd0\xc1\x74\xee\xc1\xc5\xba\xef\x2e\x25\xf7\x47\xb2\x96\xf4\x21\x1c\x65\xfe\x9e\x27\xa6\x6d\x6e\x21\xe2\xd8\xf7\x46\x62\xaa\xa8\x93\x71\x3f\xf6\xca\x9c\xd2\xd3\x06\xf9\x68\x64\x7a\x56\x83\x73\x14\x24\x3b\x3d\x1c\x13\xfe\x62\x6a\x58\xe8\x87\x2a\x33\xc7\x77\x3c\x39\x94\x2f\xd4\x60\x21\xc2\xff\x41\x88\x79\x37\x36\x60\xbf\xfa\x3e\xd4\xeb\xf0\x33\x83\x4d\xf2\xf3\x09\xfe\x1f\x6f\xf4\xed\x5f\x70\x7a\x40\xef\xe1\x80\xa3\xd0\xa5\x29\xc4\xc9\x11\xa3\x38\x60\xd6\xdf\xe2\xa1\xcf\x78\xbb\x90\x63\x16\x22\x56\x33\x20\xd3\x72\xc0\x87\x13\xac\xf5\x78\x9f\x4d\xbb\xa4\x33\xd2\x1b\x5a\xa4\x33\xbd\xf9\x91\x09\xb1\xc2\xac\xb1\xad\xa3\x1a\x35\x72\x31\x3c\x55\xa6\x01\xd2\xb7\x4f\xe1\x39\x16\x04\xde\xa1\x09\x8d\x7b\xbb\x86\x90\xe4\xeb\x00\x18\x60\xeb\x8e\xfa\xca\x36\x60\xc5\xe9\x88\x7a\x40\x09\xf2\x82\x1a\xab\x30\x4c\x01\x5c\x2d\xec\xa4\xc5\xc6\x6d\x31\x97\x7a\x2f\xf1\xc9\xb4\x99\x98\x9f\x86\x93\x65\xf7\x53\xaf\x93\x1d\xf4\x2f\x77\xfd\xa8\x5a\x7e\x29\x5e\x4b\x7c\xef\x0d\xb8\x10\x0b\x7f\xd3\x4b\xdf\x93\xf3\xac\x21\xa8\x79\x15\xa3\x60\x11\x92\x36\x1d\x49\xa9\xf0\x71\x86\x9a\x62\xe2\x45\x6f\xf2\xf1\xa0\x5a\x8e\xdd\xb8\x8a\x61\x87\xb8\x90\xda\x1c\xf3\xe8\xb6\x38\x5f\xba\x31\xd6\xed\x0b\xe9\x7e\x25\x45\xd0\x6b\x14\x75\x63\x41\x77\xb3\x70\xc9\x78\x21\x9e\x03\x32\x8c\x35\x9e\xcb\xbe\x41\xb6\x58\x03\x7f\x54\x21\x87\x07\x1a\x56\xe4\x7a\x85\x7a\xe1\xbb\x0f\x38\x0f\x98\x99\xa3\xe4\xc6\x50\xa1\x3d\x1a\xbd\xa5\x0e\x6c\x80\xd4\x06\xe7\x45\xa1\x0e\xdf\x3c\x51\xa7\x1f\x0f\x1e\xdb\xc3\xc3\xd4\x64\x20\x95\x23\x9f\x52\x6b\xa4\xf5\x4e\x9e\x49\x1f\xca\xd6\xb3\x69\xd0\x2d\x7c\x25\xfd\xb6\x00\x59\xa6\xbb\x84\xcc\x9a\x73\x23\xcd\x33\xb3\x71\xd1\x02\xe3\x1f\xfb\x7b\x72\x7c\x4c\x27\xb0\x05\xef\xad\x71\xc3\x17\x89\xeb\x96\x8d\x8f\x6d\xf3\x3d\xdf\x64\x1a\x06\x1d\xf7\x32\x9c\x14\x1e\xdf\x19\xb5\xec\x67\xbd\xf8\xb3\xbd\xac\xfa\x7a\x4a\xb2\x58\xf9\xdc\x9e\x79\x78\x88\x9d\x99\x3b\xb9\x93\xf2\xc5\x84\x37\x7e\x69\xc0\xfd\x40\xb6\x0a\x6a\x1c\x74\xb5\xab\xae\x3f\x9f\x67\xd6\xe7\xf3\xc6\x2d\xa7\x8b\xe2\xd7\x2d\x02\x3f\x03\x00\x00\xff\xff\x41\x61\xa3\xb3\x0f\x0c\x00\x00")

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

	info := bindataFileInfo{name: "schema.graphql", size: 3087, mode: os.FileMode(420), modTime: time.Unix(1567424540, 0)}
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
