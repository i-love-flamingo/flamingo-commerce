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

var _schemaGraphql = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xec\x56\x41\x6f\xeb\x36\x0c\xbe\xfb\x57\xb8\xd8\xa5\x03\xf6\x0b\x7c\x6b\xd3\xbd\xa1\xd8\xeb\x43\xb7\x64\xbb\x3c\x14\x03\x23\x33\x8e\x10\x59\x32\x24\xba\x8d\x31\xf4\xbf\x0f\x96\x64\x5b\xb2\xec\x24\xc7\x1d\x5e\x4f\x0d\xbf\x4f\x24\x45\x93\x1f\xc5\x25\xa1\x3e\x00\xc3\x7c\xa3\xea\x1a\x35\xc3\x7f\x5e\xb5\x2a\x5b\x46\xf9\xbf\x59\x9e\xe7\xf9\x1e\x0c\x3e\x01\x41\x31\x11\x1e\xc1\x70\xe6\x59\x3d\x74\x67\x89\x84\x60\x50\xcf\xa8\x9e\xb5\x1b\x31\xc7\x35\x0d\x32\x7e\xe0\x0c\x88\x2b\x69\x52\xfe\x36\xc2\xdd\x19\x6e\xb6\x20\x10\xf6\x02\x8b\xfc\x51\x29\x81\x20\xbd\x33\x6f\x5e\x0e\x3d\x1c\xf2\x49\x76\x0d\x16\xf9\x96\x34\x97\x95\xb3\x54\x48\xcf\x25\x4a\xe2\x07\x8e\x3a\x86\x8e\x60\x5e\xb0\xe4\x70\x5f\x69\xd5\x36\x23\xf6\x4b\xde\x1a\xa8\x26\x37\x3f\xcf\xf2\xa9\x90\x6e\x3c\x36\xcf\xd5\x1e\xbb\xcb\x3e\xb3\xac\xcf\x73\x82\xb7\xbc\x6e\x04\x0e\xdf\xc5\xfe\xa8\x51\x92\xf9\xf1\xcd\xfe\xbf\xdf\x6c\x5e\x70\xff\x69\x88\x53\x5f\x8c\xe0\x2f\xcc\x1e\x88\x34\xdf\xb7\x84\x66\xa0\xcc\xc3\x3d\x8c\x0c\x57\xc7\xa3\xd2\xf4\x84\x86\x69\xde\xf4\x75\x8f\x8b\x51\x86\x40\x12\xac\xee\x33\x8f\x52\xf9\xbe\x7c\xb9\xb7\xcc\xf1\x41\x9f\x90\x5e\x05\x30\xdc\xa8\x72\xf6\x49\x34\x12\x70\x81\xda\x21\xb3\x48\x03\xb8\x3d\xb5\x45\x92\xc6\x00\x7e\x83\x3a\x3e\x69\x51\xa6\x11\x08\xcb\x07\xea\xa1\x1d\xaf\x7d\x4f\xb4\x4d\xb9\x64\x7e\xe7\x86\xef\x05\x7e\xd1\xaa\x2e\x52\xf3\x4e\x4d\x6c\x6b\xdf\x00\x61\xa5\x34\x77\xe5\x9e\x6e\xef\xed\x9d\x1b\x81\xbb\x37\x4b\x7e\x01\x2e\x07\x20\x68\x83\x19\xd7\x67\x3d\x18\x55\x5f\x8f\x17\x68\x1a\x2e\xab\x22\xff\xee\xaf\xe6\x0b\x6a\x48\xb1\xd3\x57\x7c\x47\x51\xc4\x97\x3e\x61\xf7\xa1\x74\x69\xc2\x13\x6e\x9c\xbe\xe1\x87\xad\xd1\xd8\xbc\x49\xdf\x25\xd3\xeb\x1b\xcf\x76\xca\xce\x75\x9f\x73\x7a\xb1\x81\x02\x79\x78\xd5\x9c\x61\x9e\x76\xbe\xb5\x3f\xcb\x83\x9a\x73\x9f\x4d\xff\x01\xec\xbf\xe3\x98\x59\x4e\xa3\x71\x8b\x02\x19\x61\xf9\x37\x68\x0e\x92\x6c\x47\x04\x11\x6d\x4f\xe6\xc5\x7a\x27\x5e\x68\xc4\x20\x8d\x87\x77\xe0\xa2\x57\x10\x9b\x84\x59\x72\x38\x66\xef\x9d\xba\x83\x5f\x55\x07\x82\xba\x11\x4c\x6f\x3d\x67\xac\x7e\x81\x39\x71\x10\x80\x44\xcd\x4a\x3c\x40\x2b\x28\x0a\xc5\x19\x0e\x0a\xfa\xc4\x0d\x53\xad\x24\x2c\x67\x9a\x55\x06\xc0\xd2\xd1\x01\xdf\xe1\x99\xe2\x88\x35\x97\xaf\x8a\x4b\x32\x3b\xb5\x6d\x50\x52\x91\x7f\x11\x0a\xc8\x83\x70\x5e\x07\x99\x92\x64\xdd\xc5\x01\x37\xce\xbc\xd8\x8e\x13\xec\x2b\xc0\x5a\x43\xaa\x46\xfd\x5b\x24\xb5\x0e\x3a\x82\x94\x28\x52\x79\x11\x8a\x81\x08\x6c\x6b\x45\x8f\x97\x90\x0f\x68\x35\xdd\x2c\xf4\x40\xc4\xb6\xe9\xdc\xbd\xdd\xe6\xda\x92\x63\x49\x0f\xd3\x45\x49\x4e\x55\x2e\x87\xfc\x55\x92\xee\x6e\x0d\x69\xc9\x3e\xa4\x80\x7d\x28\x1a\x56\xe3\x40\xb4\x18\x29\xc6\xaa\x57\xbf\x5d\xbd\xaf\xab\x3b\x7a\x52\xd3\xc8\xdc\xab\xe9\x68\x04\x46\xfc\x1d\xfd\xc8\x5f\x16\x0a\x88\x67\xf3\x86\xd1\x14\xc1\x24\x2d\xf1\xe7\x93\x76\xe1\xee\x56\x46\xa2\x49\x4c\xd7\x22\xaf\x71\xe7\xa0\xd0\xec\x1f\x02\x73\x76\xb0\xcf\xe3\xa5\x76\x40\x8d\x92\xdd\xd0\xb2\xd3\x46\xf7\x79\x8d\x8f\x80\xdf\xb1\x4b\x56\x40\xf8\x42\x48\x0a\x31\xba\xf2\xe4\x23\x98\xd1\x74\x7f\xc2\x6e\xe1\xf9\x33\xbc\x7e\x56\x79\xab\x31\xae\x5f\x68\x98\xf7\x74\x9a\xd3\xf6\x6d\x25\xa7\x78\xec\x13\xf7\xf1\xa6\x5d\x75\xde\x00\x1d\x63\x8b\xb4\x2f\x8b\x98\xa3\xad\xb2\xad\xf8\x5e\xbd\xda\x5c\xcc\x57\xa4\xfb\x8a\x3c\x5f\x51\x67\x37\x4b\x8f\x60\x30\x92\xde\xc9\xfc\x50\xf7\x07\x57\xc0\xbf\x24\xa7\x68\x23\x2e\x6e\x10\xff\x4a\xa9\x1b\xe0\x95\xfc\xb3\x15\x98\x34\x5a\x89\xb2\x7b\x51\x1a\x87\xc3\x26\x3e\xfb\xd3\xf5\x6d\xe0\xe6\x03\xce\x1b\x01\xc6\x04\x2f\x9c\xcf\x2c\xc3\x33\xa1\x2c\xed\x04\xe6\x7f\xb4\x38\xea\xda\xbc\xde\xf7\x6e\xd5\x37\xc9\x9b\x73\xa1\x31\xb3\xcf\xec\xbf\x00\x00\x00\xff\xff\x4d\x27\x2d\x5a\xcd\x0e\x00\x00")

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

	info := bindataFileInfo{name: "schema.graphql", size: 3789, mode: os.FileMode(438), modTime: time.Unix(1562166341, 0)}
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
