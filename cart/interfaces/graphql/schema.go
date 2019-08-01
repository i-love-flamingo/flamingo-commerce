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

var _schemaGraphql = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x94\x57\xcb\x6e\xe3\x36\x14\x5d\x4f\xbe\x42\x41\x36\x33\xc0\xa0\x1f\xa0\x9d\x63\x17\x03\xa3\x13\x37\x4d\xdc\x76\x31\x08\x06\x37\xe2\xb5\xcd\x86\x22\x55\xf2\x2a\x19\xa1\xe8\xbf\x17\x94\x28\x99\x2f\xc9\x69\x16\x41\xc0\x73\x74\x5f\x3c\x3c\x64\xa8\x6b\xb0\x58\xab\xba\x46\x5d\xe1\xf7\x0d\x56\x4a\x03\x21\x5b\x83\xa6\xe2\x9f\xab\xa2\x28\x8a\x0a\x34\x95\x67\x8a\x45\xae\x7b\x80\x8d\xe4\xef\x0c\x05\x7f\x45\xcd\xd1\x94\xc5\xb7\x80\x39\x05\xdc\x0c\x94\xee\xfa\xe9\xea\xdf\xab\xab\x30\xab\x97\x8c\xb3\xb2\xd8\x6e\x86\xf8\x28\x89\x53\xb7\xdd\x94\xc5\x23\x69\x2e\x8f\xc3\xea\x33\x17\x82\xcb\xe3\x8a\x31\x8d\xc6\x44\x85\xb9\xd5\x9e\xd8\xb4\xba\x3a\x81\x41\x1d\x71\xee\x51\x1b\x25\x5d\x07\xf3\x75\x4f\xe5\x5a\xe2\xcd\x07\x60\x8c\x13\x57\x12\xc4\x06\x08\xd2\xb4\x1e\x78\xed\x3e\x69\xa0\xab\x51\xd2\x23\x0a\xac\x2c\x1a\xd7\x11\xc1\xae\x3d\x14\x4a\x1e\xcd\x5e\xad\x5a\x3a\xd9\x09\x54\x76\x7a\xbf\xf7\x6d\xdc\x2a\x25\x10\x1c\x11\x62\x3c\x1e\xd4\xcd\x07\x68\x1a\xc1\x91\xad\x55\xdb\x28\xb9\x56\x2c\x6d\xf3\x0c\xb9\x46\x19\x1e\xa0\x15\xb4\x6e\xb5\x46\x59\x75\xe7\x88\x37\x16\x25\x45\x20\x38\x61\x9d\x04\xda\x8f\x88\x8b\x63\xff\x5c\xab\x56\x52\x59\x6c\xa5\xd3\x4b\xa3\x15\x6b\x2b\x8a\x97\xb9\x09\x26\x81\x2c\x6a\x74\xea\xe3\x0b\x3f\xd0\x1a\x34\x4b\x92\xaf\x42\x7c\x4e\x62\x89\x16\x9d\xe6\x9c\x0a\xba\x32\xa6\x3b\x09\x84\x62\xdf\xe6\xda\xdf\xf8\xe8\x7c\xfe\x6c\xda\xad\x3c\xa8\x99\xd4\x16\x9a\xce\x60\x76\xee\xdb\xf3\xc8\xcd\x89\x37\x0d\x97\x47\xbb\x14\xc5\x7b\xf4\xa0\x81\xdb\x3e\xf7\x3b\xf6\x45\xab\xe0\x0c\xdd\x6b\x5e\x61\xcf\x38\x6a\x90\xac\xe7\x64\x61\xd3\xd6\x3d\xb8\x87\x1f\xab\x7a\xd8\xd0\x2c\x6b\x48\xb3\xc3\x39\x7c\x88\xb2\xe1\xa6\xb2\x41\x16\x43\xd5\x3b\x25\x6d\x03\x0f\x28\xfa\x4d\x7c\xd7\x37\xff\xf7\x03\x6f\x2c\x7f\x72\x3a\x8d\xdf\xe4\x67\xe4\xb5\x77\x99\x7c\x02\xe3\xa4\x33\xa9\x7b\x51\x25\x76\xef\x47\x07\x56\x0c\xc3\xb3\xfd\xa6\xf4\xcb\x41\xa8\xb7\x70\xb5\x46\x3a\x29\x16\xae\x55\xa0\x35\xb7\xe6\xe1\x2f\x8e\xd2\xfb\xaa\x2a\xc8\x38\xd3\x26\x82\xdd\x37\x86\x6b\x64\x7b\x5e\x63\x59\xd8\xdf\x83\x4c\x30\x32\xbf\x8f\x2f\x78\x76\x8d\x4f\x61\xda\xd0\x43\x7f\xc1\xce\xca\xd9\x11\x9e\x06\x87\xf1\x28\xde\x1c\x4c\x59\xd4\xd0\x7c\x33\x3d\xf5\xe9\x2f\xa3\xe4\x4f\x0f\xf0\x76\x87\xc6\xc0\x11\x97\xc7\x38\xf6\x50\xb8\x59\x5a\x66\x52\xd5\x85\x9b\xa4\x35\x78\x1b\xdd\x3a\x81\x43\x85\xfb\x93\x2d\xc7\x3f\x81\x63\x25\x9c\x44\x54\x4a\x63\xd5\x32\x77\x58\x68\xf1\xa8\xb1\x65\x75\x2f\x1b\xa2\x57\x15\x4f\xdd\x63\x72\x0d\x67\xe0\x41\xf0\x7e\x25\x1f\xde\x8f\xea\x5f\xeb\x3f\x08\xb5\x04\xf1\x80\x07\xb4\x97\x4c\x34\x83\x1a\xf4\x0b\x52\x23\xa0\xc2\x75\xa2\xfb\x57\xd0\x1c\x24\xdd\xf5\x9c\xfb\x3c\xc7\x55\xb9\x83\x3a\x02\x8c\x6a\x75\x85\xf1\x35\xf9\x37\x75\xde\x65\xb4\x2c\xd1\x94\xf1\x07\x88\x16\x13\xce\x3b\x4f\xc5\xe8\x0b\xab\x38\x69\x4c\x77\x6a\xcb\x8f\xd9\x69\xd2\x4d\xfa\x15\xa8\x2c\x32\x3f\x7e\xcb\x07\xae\x0d\xc9\x7e\x3e\xb3\x1c\x01\x59\x4a\xb8\x55\x9c\x31\x81\xbb\x84\xe5\x73\x9c\xca\x17\xeb\x31\x20\x5a\x72\x46\x34\xcb\x21\x8d\x98\x69\x2d\xe5\xec\xf4\x52\xcd\xe7\xed\x73\x73\xfb\xca\x65\xba\x81\x95\xaa\x1b\x90\x5d\x92\x2e\x70\x56\x4e\x29\x21\xe2\x34\xca\xd0\xa0\xd0\xa5\xaa\x81\x2e\x4d\x48\xe3\x91\x0f\xef\xb4\xf9\x09\xf5\xc7\x5f\x5f\xa8\x79\xe0\x24\x81\x82\x1d\x43\x81\xcd\x49\xc9\x25\x75\x60\x0d\x5c\x2c\xd4\x9c\x15\xea\xf0\xf0\x76\x3a\xbd\xec\xba\x4d\x4f\xb7\x17\x01\x01\x17\x31\xf3\x3e\x44\x47\x6b\xe1\x86\xb8\x3c\xae\x5b\x43\xaa\x46\x9d\x79\xaa\xff\x9c\xa1\xe4\xcb\xcd\x31\x23\x3b\x5b\x68\x73\xaa\x6c\x7c\xed\x01\xe1\xaf\x87\x5b\xae\xe9\x14\xd9\x15\x18\xd3\x28\x3d\x3c\x8b\x75\x97\x07\x77\x6d\xfd\x1c\x5f\xe2\x12\x06\x1d\xf7\x32\xf4\x06\x6f\xfd\x55\xb2\xfe\xa2\x2b\x7e\x6b\xcf\xef\xcd\xca\x2f\xb2\x9c\xf9\x9f\x2f\x89\x70\xe7\x4e\x66\x1c\x64\xc5\xd8\x5e\xd9\x2f\x3e\x26\x86\xbd\xdd\x5c\x7f\x9e\x6c\xf5\xf3\xf4\xd8\x08\x9c\xfa\xd3\x52\x01\xff\x05\x00\x00\xff\xff\x54\x26\xe2\xb3\x94\x0e\x00\x00")

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

	info := bindataFileInfo{name: "schema.graphql", size: 3732, mode: os.FileMode(420), modTime: time.Unix(1567424654, 0)}
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
