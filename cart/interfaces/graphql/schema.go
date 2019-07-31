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

var _schemaGraphql = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x84\x56\xc1\x6e\xdb\x30\x0c\x3d\x37\x5f\xe1\xa0\x97\x0d\xe8\x17\xf8\x96\xda\xc3\x10\x60\xeb\x3a\xb4\x3b\x0d\x43\xc1\x5a\x4c\x22\x4c\x96\x34\x89\x2e\x6a\x0c\xfb\xf7\x41\x96\xe2\xc8\xb2\xec\xf4\x54\xe8\x3d\x91\x8f\xf4\x23\x15\xea\x35\x16\x95\x6a\x5b\x34\x0d\xbe\xd4\xd8\x28\x03\x84\xac\x02\x43\xc5\xdf\x4d\x51\x14\x45\x03\x86\xca\x0b\xc5\x21\xdb\x01\x60\x67\xf2\x0b\x43\xc1\xdf\xd0\x70\xb4\x65\xf1\x73\xc2\x1c\x03\xd6\x9e\xd2\x6f\x7f\x6d\xfe\x6d\x36\xd3\xac\x51\x32\xce\xca\x62\x5f\xfb\xf8\x28\x89\x53\xbf\xaf\xcb\xe2\x89\x0c\x97\x47\x7f\xfa\xca\x85\xe0\xf2\xb8\x63\x06\xad\x4d\x74\xed\xd8\x70\x3a\xf0\x74\x67\x9a\x13\x58\x34\x09\xe7\x11\x8d\x55\x32\x14\xb0\x2c\x7b\x54\xeb\x88\xb7\x37\xc0\x18\x27\xae\x24\x88\x1a\x08\xe6\x69\x23\x70\x1b\xae\x68\xe8\x5b\x94\xf4\x84\x02\x1b\x87\xa6\x3a\x12\x38\x54\x87\x42\xc9\xa3\x7d\x56\xbb\x8e\x4e\xae\x01\x8d\x6b\xde\x8f\xa1\x8c\x7b\xa5\x04\x42\x20\x42\x8a\xa7\x7d\xba\xbd\x01\xad\x05\x47\x56\xa9\x4e\x2b\x59\x29\x36\x2f\xf3\x02\x85\x42\x19\x1e\xa0\x13\x54\x75\xc6\xa0\x6c\xfa\x4b\xc4\x5b\x87\x92\x22\x10\x9c\xb0\x9d\x05\x7a\x3e\x23\x21\x8e\xfb\xb7\x52\x9d\xa4\xb2\xd8\xcb\x60\x17\x6d\x14\xeb\x1a\x4a\x8f\xb9\x9d\x74\x02\x59\x52\xe8\x58\xc7\x67\x7e\xa0\x0a\x0c\x9b\x25\xdf\x4d\xf1\x25\x87\xcd\xac\x18\x2c\x17\x5c\xd0\x97\x29\x3d\x58\x60\xea\xf5\x7d\xae\xfc\x3a\x46\x97\xf3\x4f\xd2\xba\xb1\xca\xf6\x72\x3d\x46\x94\xe8\x3c\x32\x84\x6d\xa2\xdd\xa1\x71\xcf\x23\xf8\xd1\x9f\xe4\xc3\xc7\x51\xa3\x41\xfc\x43\xfd\xf0\xc1\xf2\x97\xc2\xd0\x85\x7b\x6f\x40\x65\x91\xf9\x8b\xad\x79\xe0\xc6\x92\x84\x16\xcb\x65\x8e\x80\x2c\x65\xc2\x69\x39\x63\x02\x1f\x66\xac\x98\x43\x9c\xc4\x2c\x48\xc2\xb1\x20\x3a\x02\x3f\xa3\x8b\x1c\x32\x88\x99\xd2\xe6\x9c\x07\xb3\xa6\xf9\xb2\x48\x42\xdf\xbe\x70\x39\x0c\x66\x20\xf9\xf9\x69\x54\xab\x41\xf6\xb3\x74\x71\xa4\x86\xd3\x9c\x90\x70\xb4\xb2\xe4\xc6\x7b\x4d\x91\x25\xa0\x6b\x1d\x32\x78\xe4\x7e\x51\x2c\x77\xa8\x71\x83\x6d\xae\x68\xf6\x9c\x59\xa0\xc9\x17\x43\x81\xfa\xa4\xe4\x9a\x3b\xb0\x05\x2e\x56\x34\x67\x8d\xea\x37\x7f\xf0\x29\xb0\xab\x0f\xc8\x40\x07\x51\x23\x01\x17\x29\xf3\x71\x8a\x06\x59\xef\xdc\x12\x97\xc7\xaa\xb3\xa4\x5a\x34\x99\xb7\xe2\x53\x86\x92\x97\x9b\x63\x26\xc3\xb9\x52\xe6\xa8\xec\xbc\xe5\x80\xf0\xdb\xe1\x9e\x1b\x3a\x4d\xdf\x09\x0d\xd6\x6a\x65\xfc\x5e\x36\x7d\x1e\x7c\xe8\xda\x57\xf7\x04\xc5\x98\x04\xef\xe3\xc1\x86\x51\xe3\xf1\x9d\x50\xb2\x62\x10\xf6\xbd\x8b\x16\x5e\x2c\xb2\x5c\xf8\xcd\x31\x8b\xf0\x35\x4c\x66\x1a\x64\xc7\xd8\xb3\x72\x37\x3e\xb4\x60\x7e\x23\x69\x01\x0d\x7a\x5f\xed\xeb\xed\xdd\xb8\xb3\xee\xc6\xfd\xee\xc1\x20\xf4\xe3\x9a\x80\xff\x01\x00\x00\xff\xff\x82\xd0\x9c\x43\x14\x09\x00\x00")

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

	info := bindataFileInfo{name: "schema.graphql", size: 2324, mode: os.FileMode(438), modTime: time.Unix(1564562065, 0)}
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
