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

var _schemaGraphql = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xec\x57\xc1\x6e\xe3\x36\x10\xbd\xfb\x2b\x14\xec\xc5\x0b\x14\xfd\x00\xdf\x12\x6f\x53\x04\x8d\x17\xe9\xda\xed\x65\x11\x04\x63\x6a\x2c\x0f\x4c\x91\x2a\x39\xca\x46\x28\xfa\xef\x85\x48\x4a\x16\x45\xc9\xf6\xb1\x87\xee\x69\x33\xef\x71\x66\x38\x9a\x79\x43\x93\x62\x34\x07\x10\x98\xad\x75\x59\xa2\x11\xf8\xf6\x62\x74\x5e\x0b\xce\xfe\x5e\x64\x59\x96\xed\xc1\xe2\x17\x60\x58\x9d\x09\x0f\x60\x49\x04\x56\x0b\xdd\x39\x22\x23\x58\x34\x23\x6a\x60\xed\x7a\xcc\x73\x6d\x85\x82\x0e\x24\x80\x49\x2b\x9b\xf2\xb7\x11\xee\xcf\x90\xdd\x82\x44\xd8\x4b\x5c\x65\x0f\x5a\x4b\x04\x15\x9c\x05\xf3\x74\xe8\xee\x50\x48\xb2\xa9\x70\x95\x6d\xd9\x90\x2a\xbc\xa5\x40\x7e\xca\x51\x31\x1d\x08\x4d\x0c\x1d\xc1\x6e\x30\x27\x58\x16\x46\xd7\x55\x8f\xfd\x94\xd5\x16\x8a\xb3\x9b\xcf\xa3\x7c\x0a\xe4\x1b\x8f\x8d\x73\x75\xc7\xee\x16\xff\x2c\x16\x6d\x9e\x67\x78\x4b\x65\x25\xb1\xfb\x2e\xee\x8f\x12\x15\xdb\xff\xbf\xd9\x7f\xf7\x9b\x8d\x0b\x1e\x3e\x0d\x13\xb7\xc5\x18\xfc\x1b\x66\x0f\xcc\x86\xf6\x35\xa3\xed\x28\xe3\x70\xf7\x3d\xc3\xd7\xf1\xa8\x0d\x7f\x41\x2b\x0c\x55\x6d\xdd\xe3\x62\xe4\x43\x20\x09\x56\xb6\x99\x47\xa9\x7c\x9f\xbe\xdc\xeb\xc2\xf3\xc1\x9c\x90\x5f\x24\x08\x5c\xeb\x7c\xf4\x49\x0c\x32\x90\x44\xe3\x91\x51\xa4\x0e\xdc\x9e\xea\x55\x92\x46\x07\x7e\x85\x32\x3e\xe9\x50\x61\x10\x18\xf3\x7b\x6e\xa1\x1d\x95\xe8\xac\x75\x95\x4f\x58\xdf\xc9\xd2\x5e\xe2\xa3\xd1\xe5\x2a\xb1\xee\x74\xcf\x75\xe6\x35\x30\x16\xda\x90\x2f\xf5\xf9\xe6\xc1\xde\xf8\xf6\xbf\x7b\x75\xe4\x0d\x90\xea\x80\x41\x0b\x8c\xb8\x21\xe3\xce\xa8\xdb\x5a\x6c\xa0\xaa\x48\x15\xab\xec\x7b\xb8\x56\x28\xa6\x65\x2d\x4e\xcf\xf8\x8e\x72\x15\x5f\xf8\x84\xcd\x0f\x6d\x72\x3b\x3c\xe1\x47\xe9\x2b\xfe\x70\xf5\xe9\x1b\x37\xe9\xb9\x64\x72\x43\xd3\xb9\x2e\xd9\xf9\xce\xf3\x4e\x2f\x36\xcf\x40\x1a\x5e\x0c\x09\xcc\xd2\xae\x77\xf6\x27\x75\xd0\x63\xee\x93\x6d\xcb\xef\xfe\xdb\x8f\x98\xe3\x54\x06\xb7\x28\x51\x30\xe6\x7f\x82\x21\x50\xec\xba\x61\x10\xd1\xf5\x63\xb6\x9a\xef\xc2\x0b\x4d\x38\x48\xe3\xfe\x1d\x48\xb6\xea\xe1\x92\xb0\x53\x0e\xfb\xec\x83\x53\x7f\xf0\x59\x37\x20\xb9\xe9\xc1\xf4\xd6\x63\xc6\xec\x17\x18\x13\xbb\xe1\x4f\x94\x2c\xc7\x03\xd4\x92\xa3\x50\x24\xb0\x53\xcf\x2f\x64\x85\xae\x15\x63\x3e\xd2\xab\x7c\x00\x4c\x1d\xed\xf0\x1d\x7e\x70\x1c\xb1\x24\xf5\xa2\x49\xb1\xdd\xe9\x6d\x85\x8a\x57\xd9\xa3\xd4\xc0\x01\x84\x8f\x79\x50\x68\xc5\xce\x5d\x1c\x70\xed\xcd\x93\xed\x78\x86\x43\x05\x44\x6d\x59\x97\x68\x7e\x8d\x64\xd6\x43\x47\x50\x0a\x65\x2a\x2d\x52\x0b\x90\x03\xdb\x5c\xd1\xe3\x05\x14\x02\x3a\x3d\xb7\x13\x3d\x10\xb1\x5d\x3a\x77\xaf\xb7\xb9\x76\xe4\x58\xce\x87\xe9\xa2\x62\xaf\x2a\x97\x43\xfe\xa2\xd8\x34\xb7\x86\x74\xe4\x10\x52\xc2\x7e\x28\x1a\x4e\xe2\x40\xd6\x18\x29\xc6\xac\xd7\xb0\x59\x83\xaf\xab\xfb\x79\xa4\xa5\x9d\xb9\x15\xd3\xde\x08\x82\xe9\x1d\xc3\xc8\x5f\x16\x0a\x88\x67\xf3\x86\xd1\x94\x83\x49\x9a\xe2\x8f\x27\xed\xc2\xdd\x9d\x8c\x44\x93\x98\xae\x44\x2a\x71\xe7\xa1\xa1\x39\x3c\x02\xc6\xec\xc1\x2e\x8f\x17\xda\x01\x0d\x2a\x71\x43\xcb\x9e\xb7\x79\xc8\xab\x7f\x00\xfc\x86\x4d\xb2\x02\x86\xaf\x83\xa4\x10\xbd\xab\x40\x3e\x82\xed\x4d\xcb\x13\x36\x13\x4f\x9f\xee\xe5\x33\xcb\x9b\x8d\x91\x9c\xb4\x0f\x4d\x3b\xb8\x4b\xa1\xf3\xa8\x11\x3f\x5f\xc9\xf4\x6a\x61\x3a\xdd\x48\x54\xa1\xb5\x3c\xa7\xa3\x30\x31\x1d\xb5\x22\x8e\x55\x25\x89\x1a\x2f\xf2\xd9\x98\x15\xf0\x31\xb6\x28\xf7\x68\x89\x39\xc6\x09\xe7\x8c\xef\xd9\x1b\x8f\x77\xc5\xcc\x66\xb8\xa2\xfe\x57\xc4\xdf\x8f\xea\x03\x58\x8c\x94\xfd\x6c\xbe\x2f\xdb\x83\x33\xe0\x1f\x8a\x38\x5a\xb8\x93\x0b\x2a\x3c\x82\xca\x0a\xa8\x50\xdf\x6a\x89\x49\x1f\xe7\xa8\x9a\x8d\x36\xd8\x1d\xb6\xf1\xd9\x4f\xd7\x97\x8d\x1f\x3f\xf8\x58\x4b\xb0\x76\xf0\x80\x9a\x2b\xee\xdb\x16\xc1\x88\xe3\x37\xb4\xb5\xec\x56\x51\xe5\xa1\xa9\x59\x0a\x79\x7e\xca\x1e\x41\x20\xdb\xcc\xba\xd3\xb9\x2e\x81\xd4\xcf\xce\xb6\xd6\xb2\x7d\xcb\x90\xf6\x29\xdb\xba\x28\xd0\x86\x1f\x46\x67\x77\x3e\xea\xdb\xb6\x47\x83\x63\xef\x6f\x83\xd1\x6f\x9f\x40\x6e\xad\xae\x43\xf1\x83\x51\xe5\x4e\xa9\xb2\xdf\x6b\xec\xf5\x7f\x9c\xec\xd2\x3f\x89\xaa\xe4\x5d\x3e\x31\xc0\x93\x0e\x42\xe4\xa5\x0d\x35\xfa\xab\x46\xcb\x69\x62\x01\x98\xf0\x1a\x55\xb7\xcd\xfd\xdf\x00\x00\x00\xff\xff\x40\x46\xc5\x34\x55\x10\x00\x00")

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
