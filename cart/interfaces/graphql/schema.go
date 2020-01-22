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

var _schemaGraphql = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xec\x5a\x4d\x6f\xdc\x38\x12\xbd\xfb\x57\xa8\x91\x4b\x2f\x10\x2c\xb0\x7b\xec\x9b\xdd\xed\x04\x8d\x89\x1d\xc7\xee\xc9\x1e\x82\xc0\xa0\xa5\xea\x6e\x6e\x28\x52\x43\x52\xb6\x85\xc1\xfc\xf7\x05\x29\x52\xe2\x97\x24\xc6\x8b\x9d\xcb\xce\x1c\x26\x96\xf8\x58\x2c\x92\x8f\xc5\x57\xa5\x96\x5d\x03\xc5\x96\xd5\x35\xf0\x12\x1e\x77\x50\x32\x8e\x24\x54\x5b\xc4\x65\xf1\xfb\x45\x51\x14\x45\x89\xb8\xdc\x8c\x10\xd5\xb2\xd2\x0d\x95\x05\xef\x80\xe0\x67\xe0\x18\xc4\xa6\xf8\xe6\x01\x77\x01\xa4\x5b\x7d\xd7\x5d\x4f\x10\x37\x5d\x75\x5b\x56\xc1\xba\x32\x8f\xea\x61\x53\x3c\x48\x8e\xe9\x69\xf5\xb7\xc0\x81\xa8\xb3\xb5\x7a\x49\xc8\x1d\xea\x6a\xa0\xf2\x1e\x7e\x6b\x31\x87\x6a\x2f\xa1\x16\x41\xf7\xc7\x3b\x8e\x4b\xd3\xb4\x1a\x26\xf9\xd0\xd6\x35\xe2\x5d\x88\x35\xaf\x57\x17\x7f\x5c\x5c\xf8\xab\xe5\x36\x9b\xc5\xaa\xb0\x28\x59\x4b\x65\x38\xe2\x65\xd3\x10\x0c\xd5\xce\x36\x6b\xb0\x68\xeb\xf0\xbd\xd3\x4d\xfb\x18\xe0\x3e\xe2\xa3\xdc\x22\x5e\x4d\xe2\x3e\x72\x44\xab\x03\x93\x88\xfc\x0b\xcb\xf3\x22\x5c\x23\xed\xe0\x5e\x8f\xcb\x5a\xbd\x4a\xf6\x3b\x23\x11\xbb\x7d\xc5\x18\x01\x44\x57\x83\x65\xf4\x0a\xd1\xb2\xeb\x97\xe9\x75\x34\xeb\x87\xab\x4d\xb1\xdf\xf5\x56\x80\x4a\x2c\xbb\xfd\x6e\x60\x81\x7e\xfb\x84\x09\xc1\xf4\x74\x59\x55\x1c\x44\xb4\xcc\xfd\x5b\x0d\x6c\x5a\x5e\x9e\x91\x00\x1e\x60\xee\x80\x0b\x46\x0d\x83\xa7\x89\xeb\xf1\x15\x55\x15\x96\x98\x51\x44\x76\x48\xa2\x78\x50\xa7\xb1\xf7\xb2\xe9\x49\xf8\x00\x04\x4a\xd5\x16\x11\x30\x68\xef\xa7\x06\x84\xd1\x93\x38\xb0\xcb\x56\x9e\xd5\xec\x4b\x45\xf1\x5f\xf5\x14\xbc\xf5\x45\x61\x7b\xb8\x48\x15\x1c\x51\x4b\xe4\xb6\xe5\x1c\x68\xd9\xf9\x8d\x52\x6d\x3a\xee\x8f\x85\x3f\xeb\x83\x6d\x31\xd3\x56\x7f\x6e\x7b\x1e\xec\xa9\x39\xf5\x0d\x67\x55\x5b\xca\xf0\x35\x16\xde\x94\xa0\x0a\x5c\x3e\x0d\xc4\x0c\x39\xb5\xf2\xc8\x78\x40\xaf\x69\xea\x59\xd8\x93\x86\xdd\xc2\x04\x00\x45\x07\xe5\x5b\xea\x20\xda\x76\x37\x1e\xbd\x25\x0c\xf9\xd1\x67\xe7\x74\x52\x23\x9b\x6e\xdf\x2f\x2c\xe0\x06\x61\xfa\x70\xc6\x4d\x83\xe9\xe9\xfa\x06\x61\xe2\xef\x0c\x16\xd7\x75\x23\xbb\x60\xe9\xce\x48\x58\xc3\x1f\x18\x9f\xf5\x6e\xe8\x17\xcf\x4a\x05\xbb\xfd\x6e\x8d\xf5\x3f\x8b\x33\x5a\x59\x03\xb9\x1d\x15\x6a\xe8\xa4\xb7\xe8\x8b\xec\xd6\x35\xe2\x3f\x40\xde\x11\x54\x82\xe7\xea\xfb\xe2\x19\x71\x8c\xa8\x0c\x27\xb0\xa7\x72\x1c\xf9\xfa\x55\x02\xa7\x88\xdc\xc3\x11\x14\x8f\x61\xcd\xe1\xb8\xe0\x81\xed\xfd\x95\xb5\xe5\x19\xf8\x03\x7a\xc6\xf4\x14\xc5\xbf\xc1\x53\xcd\x7a\x48\x44\x89\xc7\xfe\xad\x31\x28\xda\xda\x6e\xdb\x24\xf3\x7c\x8c\x0a\xa6\x93\x51\x7d\xd8\x57\xdb\x61\xcb\x44\x14\x44\x11\x21\xb6\xf9\x80\x25\x49\x10\xca\x1e\x86\x8f\x9c\x89\x89\x31\x3c\x48\x86\x4f\xce\xf9\xca\x42\xfb\x37\xc8\xfc\xc9\xad\x6f\x19\x55\x9b\x74\x0f\x44\x5f\xdd\x79\x9d\x7e\xb2\xc7\x78\x39\x6d\x59\xdb\x30\xda\x13\x2c\x3a\x17\x83\x48\x30\xcc\xf2\xcf\xa1\xa5\xb0\x16\x08\x57\xdd\xa1\x6b\x60\xad\xae\xac\x90\xad\xf3\xd1\x73\x0c\x79\xdb\x33\xe2\x27\x88\x16\xf1\xd1\xbc\x37\x6e\x8d\xae\x3b\xd1\x2b\x8c\x04\xf7\x50\x23\x4c\x31\x3d\xa5\x30\x69\x85\xe2\x88\x1d\x47\xd2\x19\x5d\x14\xcc\xc1\x80\x87\xf3\xd4\xcf\x44\x18\x1e\xce\xf6\x79\x70\x40\xa6\x9f\x1c\x16\x31\x5c\x2b\xd3\x67\x58\xe5\xd5\xf7\x59\xe7\xad\x3f\xc6\x7f\x34\x43\x80\x20\x4e\xcd\x9a\x75\x5d\xce\x30\x6d\xa3\xee\x9e\x1e\x99\x47\x85\xd9\x41\x86\x39\x66\x8c\x50\x66\x58\x3d\xa0\xd7\x0c\x4b\xaa\xa3\x4f\x6a\xa5\x97\x37\xc5\x07\xc2\x90\x9c\xb6\x0c\x96\x22\x49\x7d\xa0\x10\xdf\x9d\xab\xa1\x3f\x18\xe8\xf5\xe0\x0c\x16\x86\x65\xd5\x67\x72\x2a\x3a\xc6\x9a\x11\x7d\x61\x31\xdc\x04\xfb\x51\x83\xe8\x47\xf3\x3a\x7d\xd5\x6a\x16\x61\x2a\x81\x1f\x51\x19\x6d\x47\xa0\xb9\xcc\xb8\x27\x24\xe1\x05\xa5\x34\xd2\x57\x44\x5a\x88\x97\x37\x3d\x97\x5d\x2f\xb9\xa2\x41\x70\xdd\x10\x50\xaf\xc4\x9f\xe9\x4e\x94\x20\xd9\xfc\xc4\x3c\xce\x5e\xfb\x43\x62\x97\x3c\xba\x3b\xb7\x35\x91\xcf\xd9\xb3\x7a\xd5\xed\xab\xb5\xd1\xf3\x93\xf9\x9b\x02\x4e\xcd\x20\xe9\xb8\x3a\x7b\x13\xce\xab\xa6\x21\xbc\x25\xf9\x3b\x11\xd2\x02\x7b\x6e\x54\x58\xbe\x67\x17\xd4\x6d\x9e\xb8\x5d\xd2\xb6\x3f\x71\xdb\xbe\xe5\xb2\xfd\xe9\xbb\xf6\x27\xb5\xc5\x1b\xa4\xc5\x19\x09\xc3\xbe\xf9\xdb\xcd\xdd\x7c\x7b\xbb\x79\x41\x54\xbd\x79\x61\xfc\xc7\x91\xb0\x17\xff\x6d\x0d\xf2\xcc\x2a\xff\x5d\x89\x38\xc7\x4a\x0c\xfa\xd9\x54\x3f\xc6\x27\x56\xa2\x44\x32\xb7\x0b\x9a\x4d\x1f\x81\x39\x54\x07\x5c\xc3\xa6\x50\xff\x1f\x2a\x14\x5e\xb6\xb8\xfe\x01\x9d\xab\x28\xdc\x61\xfd\xa4\xf3\x17\xe8\x3c\x05\xa8\x10\xef\x02\x98\xb3\x16\x62\x53\xd4\xa8\xf9\x26\xfa\xb0\xf8\x6f\xc1\xe8\xdf\xef\xd1\xcb\x0d\x08\x81\x4e\x90\xd1\xf9\x06\x35\x23\xca\x77\xdb\x01\x86\xee\xdf\xa0\x26\xf2\xdd\x81\x87\x73\x98\xdd\x51\xbb\x9c\x85\xd9\xd6\xf8\x46\x43\x8b\x35\x80\x56\xc0\x55\x50\x2f\xf0\x04\x55\xc6\x7d\x9b\xd0\x08\x52\xc9\x71\xdf\x95\x46\x31\x77\xea\xe0\xca\xd9\x63\x8f\xa6\x4b\x40\xa9\xca\x91\x73\x21\xb8\xc7\x68\x4f\x4b\x15\x5e\x26\xc4\x80\xd7\xb0\x70\x2b\x87\x03\xce\x09\x82\x00\x6b\x68\xf9\xd4\x6d\x51\xdd\x20\x7c\xd2\xea\x7b\x5d\x3a\x0f\x8e\x4a\xc8\x99\xe6\x53\x2f\x31\x8e\x98\x48\xe0\x73\x2a\x23\xee\x9e\x33\xb7\x41\x0e\xbb\x0e\xfa\xf1\xc0\x49\x22\x0a\xbf\x89\xa0\x27\x20\xbd\x28\x09\x9b\xcc\x96\xda\xc6\x69\x7d\x96\xec\x8d\x85\x13\x87\xc3\xca\x1a\xe3\xf2\x33\xaf\x54\x84\x32\x6a\x68\x49\x00\x38\xbc\xc5\xf1\x5d\x37\xdc\x71\x46\x7d\x79\xfc\xd1\x6f\xd2\xe6\x5d\xab\x6e\xc9\x2e\xcc\xd8\x83\x88\xab\xcb\x01\x4d\x54\x0e\xd0\x8d\xa6\x22\x70\x33\x51\x32\x70\xbd\xbc\x45\x75\xd0\x20\x58\xcb\x4b\x08\xcb\x60\xbf\xc9\xce\x29\x51\x2d\xc7\x53\x1f\xa1\xf5\x56\x84\xc9\x0c\xe1\x43\x46\x17\x0e\x1a\xc2\xcd\xf6\xf6\xb3\xc0\xf4\x44\x40\xb3\x64\x2e\xa7\x1f\x51\x93\xc5\x08\xce\x5e\x96\xcc\x58\xc8\x52\x29\x2d\x37\x30\x8d\xd7\x05\x67\x2f\x61\xf9\x57\x3f\x4f\x1d\xca\x3e\x34\x1b\x3a\x3d\x23\xe9\x9c\x8b\xf4\x09\x39\x62\x2e\x24\xd5\x24\x98\xc4\x10\x94\x84\xf8\x7c\xc4\x55\x45\xe0\x36\x42\x79\xd2\xbb\x0f\xf6\xb3\xfe\x08\x44\x5a\x69\xa4\xc1\x24\x46\x72\x80\xc4\xd4\x62\xcc\x2d\x9f\xf3\x79\xe4\xa8\x59\xb7\x4f\x98\xc6\x2c\x2d\x59\xdd\x20\xda\x45\xc3\x79\xb1\x0d\xcb\x18\x10\x60\x1a\x26\xe4\x10\xfd\x26\xbd\xd6\x99\xe5\xac\x1d\x0e\x27\xec\xc4\xd1\xb4\x3f\x8a\x46\x7c\xc1\xe7\x1e\x13\x19\xf2\x76\x0c\x08\x34\x67\x46\xe7\xd8\x01\xb5\x2e\xbe\x4e\xfa\x9c\x24\x6a\xff\xe5\xc0\x26\xdf\xcb\x1f\x20\x34\x5c\x29\x20\x89\x30\x09\x91\x77\x7e\xab\x8d\x9f\x58\x48\x4c\x4f\xdb\x56\x48\x56\x03\x4f\x7c\x6d\xb8\x4e\x40\xd2\xee\xa6\x90\x41\xcc\x9e\x99\xe6\xe0\x99\x4d\xc0\x90\x84\xcf\xc7\x2b\xcc\xe5\x39\x88\xc9\x48\x88\x86\xf1\x3e\x71\xe7\x5d\xba\xf1\xb6\xad\x9f\x42\x59\x4d\x51\xcf\x63\x4d\xc3\xd9\x85\xf7\x83\xa8\x71\x48\x87\x9a\x52\xcf\xed\x52\x4a\x8e\x9f\x5a\x09\x8e\x70\xe5\x20\x80\x3f\x43\xa5\x6f\xcb\xc5\x82\xd0\x50\xbb\x9b\xcc\x21\xa6\x44\x5f\x4e\xf9\x25\x39\xe4\x58\x9f\x4c\x8e\x39\xa7\x5f\x6c\xed\x6f\xd2\xd9\x41\x80\x24\x03\xbf\x2d\x21\x4e\x66\x5e\xf7\x23\x62\xa1\xb6\xf8\x15\x11\x5c\xe9\x7d\xbc\x07\xd1\x12\xab\xa8\xce\x48\x28\x1c\xa3\xd7\x9c\xb3\x31\x9c\x05\xda\x7b\x00\x98\xb4\xe4\x17\x08\xd8\xf3\x0e\x6b\x21\xa4\x0c\x0b\xe7\xb0\x06\x45\x29\xa5\x45\x46\x3f\xb4\xc1\x3e\xb7\x78\x97\x70\x38\x01\x2e\x7e\xbf\xd0\x5c\xb2\x95\xc3\x28\x18\xe8\x56\x98\x72\xf4\xdd\x44\xf1\x4f\x29\x18\xc3\xbe\x31\x37\x65\xea\xd9\x1e\x86\x20\x26\xf9\x1f\x75\x16\x62\xd1\xa3\x9f\xd7\x7c\x60\xdc\x52\xf7\xc8\x78\xdd\xc7\x8d\xfe\xbf\xa5\x6e\x3a\x80\x14\xfa\xe6\xb5\xcb\x92\x28\xb2\x3c\x2a\xa8\xb3\xdd\x43\xb1\xa5\xe1\xac\x04\x21\x1c\xa5\x3a\xf5\xa1\xdc\x7c\x18\x1c\x4a\xef\x0e\x5f\xfe\xc7\x43\x27\x0c\xd8\xc5\xc2\x40\x2a\xcd\x82\xb8\xd8\xa9\x7b\x7d\x18\x00\x83\xfe\xa3\xc0\x11\x99\xeb\xe3\x50\x70\xca\x19\x4b\xbc\x42\xd7\x21\xd2\xec\x37\x1f\x71\x3f\xf5\x79\xc6\x3c\x1b\x02\x5f\xdf\x64\x79\x58\x0f\x5f\x5e\xe7\xd2\xcf\x09\xcd\x7f\x69\xb8\x45\x0d\x67\x95\xdb\x3f\x36\xcb\x98\x7f\x4e\x61\xfe\x9f\x55\x9e\x56\x78\x4e\x24\x4d\x61\x32\x54\xde\x05\xa6\x4d\x2b\x47\x72\xc7\xbc\xde\x6b\x40\x0e\xb1\xff\x44\x5e\x67\xd0\x3a\x83\xd5\x19\xa4\xce\xe0\x74\x06\xa5\x33\x18\x9d\x41\xe8\x0c\x3e\x67\xd0\x39\x83\xcd\x19\x64\xce\xe0\x72\x06\x95\x33\x98\x9c\x4b\x64\x78\x95\x40\x2b\x2d\x3c\x8b\x2f\xed\xf8\xe1\xc4\x0b\xdb\x9b\x89\x9f\xfa\x45\x16\x6e\x0c\x6b\x42\x23\x97\x55\x75\x60\xaa\xc7\x3a\xaa\xe5\xec\x77\xab\xf7\x63\xc5\xe5\x7d\xb1\xf8\x0b\x9a\xc0\x03\x6f\x9c\x1d\x10\x90\xe0\x96\x83\x97\x7f\x91\xb3\x6c\x4f\x89\xbf\xe1\x87\x2c\xda\xdf\xff\xca\xe8\xaf\x8d\x4a\x8a\x94\xd1\x2f\xb2\xcb\xb0\xeb\x2c\x4f\xee\x10\xfa\xae\xed\xc7\xf1\x23\xd3\x1a\x8d\x11\x6a\xb3\x18\xbe\xa2\x8f\xb2\x31\x6e\x7a\xe0\x40\xb8\xad\xc3\x2f\x94\xef\xc3\x2f\x29\xd1\x68\x49\xe9\x97\x1a\x50\xe5\x37\xdd\x98\x18\x7d\xe6\x36\xd3\x59\x97\x59\x1b\x94\x30\x79\x0f\x35\x7b\x86\xc1\xce\xc9\xfc\x91\xb7\xe1\x93\xf6\x46\x1f\xd7\x6e\x89\x78\xd1\xde\x1f\x17\xff\x09\x00\x00\xff\xff\x1f\x97\x4b\x5f\x71\x2b\x00\x00")

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
