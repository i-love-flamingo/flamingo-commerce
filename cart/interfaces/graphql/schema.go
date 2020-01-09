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

var _schemaGraphql = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xec\x5a\xcf\x6f\xdc\xb8\x15\xbe\xfb\xaf\xd0\x20\x97\x29\x10\x14\x68\x8f\x73\xb3\x67\x9c\x60\xb0\xb1\xe3\xd8\xb3\xe9\x21\x08\x0c\x5a\x7a\x33\xc3\x86\x22\xb5\x24\x65\x5b\x58\xec\xff\x5e\x90\x22\x25\xfe\x92\x44\xa7\xe8\x5e\xba\x3d\x74\x33\xe4\xf7\x1e\x1f\xc9\x8f\x8f\x1f\x9f\x2c\xbb\x06\x8a\x2d\xab\x6b\xe0\x25\x3c\xee\xa0\x64\x1c\x49\xa8\xb6\x88\xcb\xe2\xf7\x8b\xa2\x28\x8a\x12\x71\xb9\x19\x21\xaa\x67\xa5\x3b\x2a\x0b\xde\x01\xc1\xcf\xc0\x31\x88\x4d\xf1\xcd\x03\xee\x02\x48\xb7\xfa\xae\x4d\x4f\x10\x77\x5d\x75\x5b\x56\xc1\xba\x32\x3f\xd5\x8f\x4d\xf1\x20\x39\xa6\xa7\xd5\xdf\x82\x00\x22\x63\xeb\xf5\x92\x90\x3b\xd4\xd5\x40\xe5\x3d\xfc\xd6\x62\x0e\xd5\x5e\x42\x2d\x02\xf3\xc7\x3b\x8e\x4b\xd3\xb5\x1a\x26\xf9\xd0\xd6\x35\xe2\x5d\x88\x35\xcd\xab\x8b\x3f\x2e\x2e\xfc\xd5\x72\xbb\xcd\x62\xd5\xc0\x4f\xb0\xc3\xa2\x64\x2d\x95\xe1\xb0\x97\x4d\x43\x30\x54\x43\xb7\xb6\x10\x6d\x1d\xb6\x3b\x66\x3a\xd0\x00\xf7\x11\x1f\xe5\x16\xf1\x6a\x12\xf7\x91\x23\x5a\x1d\x98\x44\xe4\x5f\x58\x9e\x17\xe1\x1a\x69\x07\xf7\x2c\x2e\x6b\xd5\x94\xb4\x3b\x23\x11\x87\x7d\xc5\x18\x01\x44\x57\x83\x67\xf4\x0a\xd1\xda\xeb\xc6\xf4\x62\x9a\x45\xc4\xd5\xa6\xd8\xef\x7a\x2f\x40\x25\x96\xdd\x7e\x37\x50\x41\xb7\x3e\x61\x42\x30\x3d\x5d\x56\x15\x07\x11\x2d\x73\xdf\xaa\x81\x4d\xcb\xcb\x33\x12\xc0\x03\xcc\x1d\x70\xc1\xa8\xa1\xf1\x34\x7b\x3d\xd2\xa2\xaa\xc2\x12\x33\x8a\xc8\x0e\x49\x14\x0f\xea\x74\xf6\x51\x36\x3d\x13\x1f\x80\x40\xa9\xfa\x22\x16\x06\xfd\xfd\xd4\x80\x30\x7a\x12\x07\x76\xd9\xca\xb3\x9a\x7d\xa9\x78\xfe\xab\x9e\x82\xb7\xbe\x28\xec\x0f\x17\x09\xf5\xfb\xb3\x65\x6d\xc3\xa8\x3a\x4e\xd1\x04\xc7\x2e\x33\xc5\x0a\x8e\xa8\x25\x72\xdb\x72\x0e\xb4\xec\x7c\x7f\x52\xf1\x04\xf7\xc7\xc9\xf7\x73\xb0\x3d\xc6\x8d\xfa\xe7\xb6\xa7\xce\x9e\x9a\x6c\xd1\x70\x56\xb5\xa5\x0c\x9b\xb1\xf0\x56\x01\xaa\x60\x96\xa7\x81\xcb\x21\x0d\x57\x1e\x7f\x0f\xe8\x35\xcd\x56\x0b\x7b\xd2\xb0\x5b\x98\x00\xa0\xe8\x6c\x7d\x4b\x9d\x5d\xdb\xef\xe6\xb1\x9f\x49\x5f\x7e\xd6\xda\x39\x46\x6a\x64\x63\xf6\xfd\xc2\x02\x6e\x10\xa6\x0f\x67\xdc\x34\x98\x9e\xae\x6f\x10\x26\xfe\xce\x60\x71\x5d\x37\xb2\x0b\x96\xee\x8c\x84\x75\xfc\x81\xf1\xd9\xe8\x06\xbb\x78\x56\x2a\x49\xee\x77\x6b\xac\xff\xb3\x38\xa3\x95\x75\x90\x6b\xa8\x50\x83\x91\xde\xa2\x2f\xb2\x5b\xd7\x88\xff\x00\x79\x47\x50\x09\x5e\xa8\xef\x8b\x67\xc4\x31\xa2\x32\x9c\xc0\x9e\xca\x71\xe4\xeb\x57\x09\x9c\x22\x72\x0f\x47\x50\x3c\x86\x35\x87\xe3\x42\x04\xd6\xfa\x2b\x6b\xcb\x33\xf0\x07\xf4\x8c\xe9\x29\x4a\x99\x43\xa4\x9a\xf5\x90\x48\x2c\x8f\x7d\xab\x71\x28\xda\xda\x6e\xdb\x24\xf3\x7c\x8c\xca\xbf\x93\x17\xc1\xb0\xaf\xd6\x60\xcb\x44\x94\x77\x11\x21\xb6\xfb\x80\x25\x49\x10\xca\x1e\x86\x8f\x9c\x89\x89\x31\x3c\x48\x46\x4c\xce\xf9\xca\x42\xfb\x97\xce\xfc\xc9\xad\x6f\x19\x55\x9b\x74\x0f\x44\x5f\xf9\x79\x46\x6f\xb4\x18\xef\xb3\x31\x29\x26\xce\xc5\x20\x2e\x0c\xb3\xfc\x73\x68\x29\xac\x85\xc5\x55\x77\xe8\x1a\x58\xab\x5b\x2e\x64\xeb\x7c\xf6\x1c\x53\xde\xf6\x8c\xf8\x09\xa2\x45\x7c\x34\xed\x26\xac\x31\x74\x27\x7b\x85\x99\xe0\x1e\x6a\x84\x29\xa6\xa7\x14\x26\xad\x6c\x1c\x91\xe4\x48\x41\xa3\xa7\x82\x39\x18\xf0\x70\x9e\xfa\x99\x08\xc3\xc3\x59\x9b\x07\x07\x64\xec\xe4\xb0\x88\xe1\x5a\x19\x9b\x61\x95\x57\xdf\x67\x83\xb7\xf1\x98\xf8\xd1\x0c\x01\x82\x3c\x35\xeb\xd6\x0d\x39\xc3\xb5\xcd\xba\x7b\x7a\x64\x1e\x15\x66\x07\x19\xe6\x98\x31\x42\x99\xe1\xf5\x80\x5e\x33\x3c\x29\x43\x9f\xd4\x4a\x67\x6f\x8a\x0f\x84\x21\x39\xed\x19\x2c\x45\x92\xfa\x40\x21\xbe\x3b\x57\x43\x7f\x30\xd0\xeb\xc1\x19\x2c\x4c\xcb\xca\x66\x72\x2a\x3a\xc7\x9a\x11\x7d\x61\x31\xdc\x04\xfb\x51\x83\xe8\x9f\xa6\x39\x7d\xd5\x6a\x16\x61\x2a\x81\x1f\x51\x19\x6d\x47\x20\xd3\xcc\xb8\x27\x24\xe1\x05\xa5\x34\xd2\x57\x44\x5a\x88\x97\x37\x3d\x97\x5d\x2f\xb9\xa2\x41\x70\xdd\x10\x50\x4d\xe2\xcf\x0c\x27\x7a\x58\x19\xef\x96\xc2\xb3\xd7\xfe\xf0\x20\x4c\x1e\xdd\x9d\xdb\x9b\x78\x07\xda\xb3\x7a\xd5\xed\xab\xb5\x79\x02\x4c\xbe\xfb\x14\x70\x6a\x06\xc9\xc0\xd5\xd9\x9b\x08\x5e\x75\x0d\xe9\x2d\xc9\xdf\x89\x94\x16\xf8\x73\xb3\xc2\xf2\x3d\xbb\xa0\x6e\xf3\xc4\xed\x92\xb6\x7d\xc3\x6d\xfb\x33\x97\xed\x9b\xef\xda\x37\x6a\x8b\x9f\x90\x16\x67\x24\x0c\xfb\xe6\x6f\x37\x77\xf3\xed\xed\xe6\x25\x51\xd5\xf2\xc2\xf8\x8f\x23\x61\x2f\x7e\x6b\x0d\xf2\xcc\x2a\xbf\xad\x44\x9c\x63\x25\x06\xdd\x46\xcb\xbd\x4f\xac\x44\x89\xf7\xdf\x2e\xe8\x36\x36\x02\x73\xa8\x0e\xb8\x86\x4d\xa1\xfe\x7f\xa8\x6c\x78\x0f\xcc\xf5\x0f\xe8\x5c\x45\xe1\xbd\xfb\x3c\xe4\x2f\xd0\x79\x0a\x50\x21\xde\x05\x30\x67\x2d\xc4\xa6\xa8\x51\xf3\x4d\xf4\x69\xf1\xdf\x82\xd1\xbf\xdf\xa3\x97\x1b\x10\x02\x9d\x20\xc3\xf8\x06\x35\x23\xca\x0f\xdb\x01\x86\xe1\xdf\xa0\x26\x8a\xdd\x81\x87\x73\x98\xdd\x51\xbb\x9c\x85\xd9\xd6\xf8\x46\x43\x8b\x65\x83\x56\xc0\x55\x50\x62\xf0\x04\x55\xc6\x7d\x9b\xd0\x08\x52\xc9\x71\x3f\x94\x46\x31\x77\xea\xe0\xca\xd9\x63\x8f\xa6\xab\x46\xa9\x62\x93\x73\x21\xb8\xc7\x68\x4f\x4b\x95\x5e\x26\xc4\x80\xd7\xb1\x70\x2b\x87\x03\xce\x09\x82\x00\x6b\x68\xf9\xd4\x6d\x51\xdd\x20\x7c\xd2\xea\x7b\x5d\x3a\x3f\x1c\x95\x90\x33\xcd\xa7\x5e\x62\x1c\x31\x91\xc0\xe7\x54\x46\x6c\x9e\x33\xb7\x41\x0e\xbb\x01\xfa\xf9\xc0\x79\x44\x14\x7e\x17\x41\x4f\x40\x7a\x51\x12\x76\x99\x2d\xb5\x9d\xd3\xfa\x2c\x69\x8d\x85\x93\x87\xc3\x62\x1c\xe3\xf2\x33\xaf\x54\x86\x32\x6a\x68\x49\x00\x38\xbc\xc5\xf1\x5d\x37\xdc\x71\x46\x7d\x79\xfc\xd1\x2d\x69\xf7\xae\x57\xb7\xca\x17\xbe\xd8\x83\x8c\xab\xcb\x01\x4d\x54\x0e\xd0\x9d\xa6\x22\x70\x33\x51\x32\x70\xa3\xbc\x45\x75\xd0\x21\x58\xcb\x4b\x08\x2b\x67\xbf\xc9\xce\x29\x51\x2d\xe7\x53\x1f\xa1\xf5\x56\x84\xc9\x4c\xe1\xc3\x8b\x2e\x1c\x34\x84\x9b\xed\xed\x67\x81\xe9\x89\x80\x66\xc9\xdc\x9b\x7e\x44\x4d\x16\x23\x38\x7b\x59\x72\x63\x21\x4b\xa5\xb4\x37\x55\xc1\xdf\x19\xcf\x61\xc5\x58\xff\x9e\x3a\x94\x7d\x6a\x36\x74\x7a\x46\xd2\x39\x17\xe9\x13\x72\xc4\x5c\x48\xaa\x49\x30\x89\x21\x28\x09\xf1\xf9\x88\xab\x8a\xc0\x6d\x84\xf2\xa4\x77\x9f\xec\x67\xe3\x11\x88\xb4\xd2\x48\x83\x49\x8c\xe4\x00\x89\xa9\xc5\x98\x5b\x3e\x17\xf3\xc8\x51\xb3\x6e\x9f\x30\x8d\x59\x5a\xb2\xba\x41\xb4\x8b\x86\xf3\x72\x1b\x96\x31\x20\xc0\x34\x4c\xc8\x21\xfb\x4d\x46\xad\x5f\x96\xb3\x7e\x38\x9c\xb0\x93\x47\xd3\xf1\x28\x1a\xf1\x85\x98\x7b\x4c\xe4\xc8\xdb\x31\x20\xd0\x9c\x19\x9d\x63\x07\xd4\xba\xf8\x3a\x19\x73\x92\xa8\xfd\xc7\x06\xfb\xf8\x5e\xfe\x66\xa1\xe1\x4a\x01\x49\x84\x49\x88\xbc\xf3\x7b\x6d\xfe\xc4\x42\x62\x7a\xda\xb6\x42\xb2\x1a\x78\xe2\x03\xc5\x75\x02\x92\x0e\x37\x85\x0c\x72\xf6\xcc\x34\x87\xc8\xec\x03\x0c\x49\xf8\x7c\xbc\xc2\x5c\x9e\x83\x9c\x8c\x84\x68\x18\xef\x1f\xee\xbc\x4b\x77\xde\xb6\xf5\x53\x28\xab\x29\xea\x79\xac\x69\x38\xbb\xf0\x7e\x12\x35\x01\xe9\x54\x53\xea\xb9\x5d\x4a\xc9\xf1\x53\x2b\xc1\x11\xae\x1c\x04\xf0\x67\xa8\xf4\x6d\xb9\x58\x10\x1a\x6a\x77\x93\x6f\x88\x29\xd1\x97\x53\x7e\x49\x0e\x39\xd6\x27\x93\x63\xce\xe9\x17\x5b\xfb\x9b\x0c\x76\x10\x20\xc9\xc4\x6f\x4b\x88\x93\x2f\xaf\xfb\x11\xb1\x50\x5b\xfc\x8a\x08\xae\xf4\x3e\xde\x83\x68\x89\x55\x54\x67\x24\x14\x8e\xd1\x6b\xce\xd9\x98\xce\x02\xed\x3d\x00\xcc\xb3\xe4\x17\x08\xd8\xf3\x0e\x6b\x21\xa4\x1c\x0b\xe7\xb0\x06\x45\x29\xa5\x45\xc6\x38\xb4\xc3\xfe\x6d\xf1\x2e\x11\x70\x02\x5c\xfc\x7e\xa1\xb9\x64\x2b\x87\x51\x32\xd0\xbd\x30\x15\xe8\xbb\x89\xe2\x9f\x52\x30\x86\x7d\xe3\xdb\x94\xa9\xdf\xf6\x30\x04\x39\xc9\xff\xa8\xb3\x90\x8b\x1e\xfd\x77\xcd\x07\xc6\x2d\x75\x8f\x8c\xd7\x7d\xde\xe8\xff\xb7\x64\xa6\x13\x48\xa1\x6f\x5e\xbb\x2c\x89\x22\xcb\xa3\x82\x3a\xdb\x3d\x14\x5b\x1a\xce\x4a\x10\xc2\x51\xaa\x53\x1f\xd8\xcd\x87\xc1\xa1\xf4\xee\xf0\xe5\x7f\x3c\x74\xc2\x81\x5d\x2c\x0c\xa4\xd2\x2c\x88\x8b\x9d\xda\xea\xc3\x00\x18\xf4\x1f\x05\x8e\xc8\x9c\x8d\x43\xc1\xa9\x60\x2c\xf1\x0a\x5d\x87\x48\xb3\xdf\x7c\xc4\xfd\xd4\xbf\x33\xe6\xd9\x10\xc4\xfa\x53\x9e\x87\xf5\xf0\xe5\x75\x2e\xfd\x9c\xd4\xfc\x97\x86\x5b\xd4\x70\x56\xb9\xfd\x63\xb3\x8c\xf9\xe7\x14\xe6\xff\x59\xe5\x69\x85\xe7\x64\xd2\x14\x26\x43\xe5\x5d\x60\xda\xb4\x72\x24\x77\xcc\xeb\xbd\x06\xe4\x10\xfb\x4f\xe4\x75\x06\xad\x33\x58\x9d\x41\xea\x0c\x4e\x67\x50\x3a\x83\xd1\x19\x84\xce\xe0\x73\x06\x9d\x33\xd8\x9c\x41\xe6\x0c\x2e\x67\x50\x39\x83\xc9\xb9\x44\x86\x57\x09\xb4\xd2\xc2\xb3\xf8\xd2\x8e\x1f\x4e\xbc\xb4\xbd\x99\xf8\x13\xc1\xc8\xc3\x8d\x61\x4d\xe8\xe4\xb2\xaa\x0e\x4c\x59\xac\xa3\x5a\xce\x7e\xb7\x7a\x3f\x56\x5c\xde\x17\x8b\x7f\x41\x13\x44\xe0\x8d\xb3\x03\x02\x12\xdc\x72\xf0\xf2\x5f\xe4\x2c\xfb\x53\xe2\x6f\xf8\x43\x16\x1d\xef\x7f\xe5\xf4\xd7\x46\x3d\x8a\x94\xd3\x2f\xb2\xcb\xf0\xeb\x2c\x4f\xee\x10\xfa\xae\xed\xc7\xf1\x33\xd3\x1a\x8d\x19\x6a\xb3\x98\xbe\xa2\x8f\xb2\x31\x6e\x7a\xe0\x40\xb8\xad\xc3\x2f\x94\xef\xc3\x2f\x29\xd1\x68\x49\xe9\x97\x1a\x50\xbd\x6f\xba\xf1\x61\xf4\x99\xdb\x97\xce\xba\xcc\xda\xa0\x84\xcb\x7b\xa8\xd9\x33\x0c\x7e\x4e\xe6\x1f\x79\x1b\x3e\xe9\x6f\x8c\x71\xed\x96\x88\x17\xfd\xfd\x71\xf1\x9f\x00\x00\x00\xff\xff\xed\xd2\x2f\xbc\xa9\x2b\x00\x00")

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
