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

var _schemaGraphql = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x94\x57\xdd\x6e\xdb\x36\x14\xbe\x6e\x9e\x42\x41\x6e\x5a\xa0\xc8\x03\xe8\x2e\x91\x86\xc2\x58\xeb\x65\x8d\x87\x5d\x14\x41\x71\x22\x1e\xdb\x5c\x28\x92\x23\xa9\x26\xc2\xb0\x77\x1f\xf8\x63\x99\xa4\x28\x79\xcd\x45\x20\xf0\x7c\x3c\x7f\xfc\xce\x47\xda\x8c\x12\xab\x46\xf4\x3d\xaa\x0e\xbf\xb7\xd8\x09\x05\x06\x49\x03\xca\x54\xff\x5c\x55\x55\x55\x75\xa0\x4c\x7d\x86\x58\xcb\xb5\x33\x90\x13\xf8\x3b\x41\x46\x7f\xa0\xa2\xa8\xeb\xea\x5b\x82\x9c\x1c\xb6\x1e\x32\x5e\x3f\x5d\xfd\x7b\x75\x95\x46\x8d\x82\x51\x52\x57\x9b\xd6\xfb\x47\x6e\xa8\x19\x37\x6d\x5d\x3d\x1a\x45\xf9\xc1\xaf\x3e\x53\xc6\x28\x3f\xdc\x11\xa2\x50\xeb\x2c\xb1\xb0\xea\x80\x72\x50\xdd\x11\x34\xaa\x0c\xf3\x80\x4a\x0b\x1e\x2a\x58\xce\x7b\x4a\xd7\x02\x6f\xde\x01\x21\xd4\x50\xc1\x81\xb5\x60\x60\x1e\x36\x32\x5e\x87\x2d\x12\xc6\x1e\xb9\x79\x44\x86\x9d\xb5\xe6\x79\x64\xe6\x50\x1e\x32\xc1\x0f\x7a\x27\xee\x06\x73\xb4\x1d\xe8\x6c\xf7\xfe\x70\x65\xdc\x0b\xc1\x10\x02\x10\x72\x7b\xde\xa8\x9b\x77\x20\x25\xa3\x48\x1a\x31\x48\xc1\x1b\x41\xe6\x65\x9e\x4d\xa1\x50\x82\x7b\x18\x98\x69\x06\xa5\x90\x77\xe3\xd9\xe3\x8d\xb5\x1a\x61\x80\x51\x83\xfd\xcc\xd1\xee\x64\x09\x7e\xec\x67\x23\x06\x6e\xea\x6a\xc3\x03\x5f\xa4\x12\x64\xe8\x4c\xbe\x4c\x75\xd2\x09\x24\x59\xa1\x53\x1d\x9f\xe8\xde\x34\xa0\xc8\x2c\xf8\x5d\x6a\x5f\xa2\xd8\x8c\x8b\x81\x73\x81\x05\x63\x9d\xc3\x03\x05\x52\xb2\x6f\x4a\xe5\xb7\xb1\x75\x39\x7e\x31\xec\x86\xef\xc5\x42\x68\x6b\x9a\x66\xb0\xd8\xf7\xcd\xb9\xe5\xfa\x48\xa5\xa4\xfc\x60\x97\x32\x7f\x8f\x91\xc9\x63\x87\x67\x77\x62\x9f\x94\x48\x66\xe8\x41\xd1\x0e\x1d\xe2\xa0\x80\x13\x87\x29\x9a\xf5\xd0\x3b\xe3\x0e\xde\xee\x7a\x7f\xa0\x45\x94\x0f\xb3\xc5\x25\xbb\xf7\xd2\x52\xdd\x59\x27\xab\xae\xfa\xad\xe0\xb6\x80\xaf\xc8\xdc\x21\xfe\xaf\x3d\x3f\xbb\x21\x6a\xcb\x9f\xd4\x1c\x4f\x7b\xca\x3d\x8a\xca\xbb\x0c\x3e\x82\x0e\xd4\x99\xd8\xbd\xca\x12\x7b\xf6\x27\x05\x16\x04\xd3\xd9\x7e\x15\xea\x65\xcf\xc4\x6b\xba\xda\xa3\x39\x0a\x92\xae\x75\xa0\x14\xb5\xe2\x11\x2f\x9e\xa8\xf7\x59\x74\x50\x50\xa6\x36\x33\x87\x3d\x9a\x2a\x24\x3b\xda\x63\x5d\xd9\xff\x9e\x26\x98\x89\xdf\xfb\x17\x3c\xab\xc6\x87\x34\x6c\xaa\xa1\xbf\xe2\x68\xe9\x1c\x00\x4f\x5e\x61\x22\x48\xd4\x07\x5d\x57\x3d\xc8\x6f\xda\x41\x9f\xfe\xd2\x82\xdf\x7e\x85\xd7\x2f\xa8\x35\x1c\x70\xbd\x8d\xa7\x1a\xaa\xd0\x4b\x8b\x9c\x65\x75\xe1\x26\x19\x34\xde\x67\xb7\x4e\xa2\x50\xe9\xf9\x14\xd3\x89\x27\xf0\x94\x09\x35\x0c\x33\x81\x95\x96\x2e\xdc\x4e\x8b\xfb\x6a\x45\x0f\x94\xdf\x7a\x12\x79\x05\x86\x37\x70\xfc\xbd\x5d\x42\x90\xc0\x42\x08\x34\x5f\x72\x24\x88\x70\x5d\xad\x5e\x8f\xc8\x3d\xaa\xa2\xba\xa2\xbd\x64\x68\xd5\x18\xc9\x05\x15\x8d\x4a\xa1\x73\xc9\x99\xa4\x26\xa8\x7e\x32\x13\x6e\xa5\xec\x3e\xf6\x1a\xbd\x05\xfe\x36\xa3\xbb\x32\xca\x9b\xc2\xb1\x84\x7d\x3f\xc0\xd4\x55\xe1\x2f\x3e\xf4\x3d\x55\xda\x70\xb0\x5c\x5e\xc4\x30\x28\x42\xd2\x91\xa3\x84\x30\xdc\xce\x50\x31\x26\x1c\xf4\x6a\x3e\x1a\xd8\x60\xc2\x2c\x2e\x62\x8c\x42\x2c\x94\x36\xc7\x6c\xd5\x5a\xce\xe7\x19\x0b\x7d\xfb\x4c\x39\x26\xa3\xe8\x49\xdd\x4b\xe0\xe3\x2c\x5c\x22\x2e\xd4\xcc\x01\x19\x46\x0a\x6d\x1a\x37\x20\x6b\x59\x83\xb9\xd4\x21\x85\x07\xea\x9f\x2a\xcb\x1d\x72\xc4\x57\x17\x72\xf6\x98\x99\xa3\xe4\xc4\x90\xa1\x3c\x0a\xbe\xc6\x0e\xec\x81\xb2\x95\x9c\x8b\x44\xf5\x6f\xcf\xc0\xd3\xcb\xc2\x23\x1d\xdc\x6a\xa1\x01\xca\x72\xe4\x43\x6a\x0d\x69\xbd\x51\x6d\x28\x3f\x34\x83\x36\xa2\x47\x55\x78\xad\xfe\x52\x80\x94\xd3\x2d\x21\xb3\xe1\x5c\x29\x73\xca\xec\xf4\xe0\x01\x83\xbf\xed\xef\xa9\x32\xc7\x54\x81\x25\x68\x2d\x85\xf2\x2f\x43\x35\x96\x8d\xdb\xa1\x7f\xce\xef\x31\x0e\x9e\xc7\x8e\x86\x51\xe3\xf1\xcd\x20\x27\x4e\xeb\xab\xdf\x87\xf3\x93\xab\x8b\x93\xac\x17\x7e\xf6\xcc\x3c\x7c\x09\x93\x99\x3b\xb9\x23\x64\x27\xec\x8e\xf7\x3d\xa8\x17\x34\x92\x41\x87\x9e\x57\x9b\xf6\xfa\xe3\xa4\x59\x1f\xa7\xfb\xb6\x89\x2f\x8a\x0f\x6b\x09\xfc\x17\x00\x00\xff\xff\x09\x5f\xe1\x9f\x97\x0d\x00\x00")

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

	info := bindataFileInfo{name: "schema.graphql", size: 3479, mode: os.FileMode(420), modTime: time.Unix(1567424635, 0)}
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
