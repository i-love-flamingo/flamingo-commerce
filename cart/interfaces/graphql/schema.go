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

var _schemaGraphql = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xec\x5b\x4d\x6f\xdc\x38\xd2\xbe\xfb\x57\xa8\xfd\x5e\x7a\x82\xbc\x33\xd8\x3d\xfa\x66\x77\x27\x41\x63\xc6\x4e\x62\x3b\xd9\x43\x60\x18\xb4\x54\xdd\xcd\x0d\x45\x2a\x24\x65\x5b\x58\xcc\x7f\x5f\xf0\x4b\x22\x29\x4a\x62\x12\xcc\x2c\x66\x37\x73\x98\xa4\xc9\x62\xb1\x8a\x7c\x58\x7c\xaa\xa8\xc8\xae\x81\x62\xc3\xea\x1a\x78\x09\xf7\x5b\x28\x19\x47\x12\xaa\x0d\xe2\xb2\xf8\xd7\x49\x51\x14\x45\x89\xb8\x3c\x1b\x44\x54\xcf\x4a\x77\x54\x4e\x78\x0b\x04\x3f\x02\xc7\x20\xce\x8a\x4f\x81\xe0\x36\x12\xe9\x56\x77\x7a\xe8\x01\xc6\x5d\x17\xdd\x86\x55\xb0\xae\xec\x4f\xf5\xe3\xac\xb8\x91\x1c\xd3\xc3\xea\xa7\xc8\x80\xd1\x60\xa7\xf5\x9c\x90\x77\xa8\xab\x81\xca\x6b\xf8\xd2\x62\x0e\xd5\x4e\x42\x2d\xa2\xe1\xf7\xef\x38\x2e\x6d\xd7\xaa\x77\xf2\xa6\xad\x6b\xc4\xbb\x58\xd6\x36\xaf\x4e\x7e\x3f\x39\x09\x57\xcb\xef\xb6\x8b\x55\x61\x51\xb2\x96\xca\x78\xc6\xf3\xa6\x21\x18\xaa\xad\xeb\xd6\xc2\xa2\xad\xe3\x76\x6f\x98\xb6\x31\x92\x7b\x83\xf7\x72\x83\x78\x35\x29\xf7\x86\x23\x5a\xdd\x32\x89\xc8\x3f\xb0\x3c\x2e\x8a\x6b\x49\x37\x79\x30\xe2\xbc\x56\x4d\xc9\x71\x47\x24\xc6\x66\x5f\x30\x46\x00\xd1\x55\xaf\x19\x3d\xc3\x68\xd9\x75\x63\x7a\x1d\xed\xfa\xe1\xea\xac\xd8\x6d\x8d\x16\xa0\x12\xcb\x6e\xb7\xed\x51\xa0\x5b\x1f\x30\x21\x98\x1e\xce\xab\x8a\x83\x18\x2d\xb3\x69\xd5\x82\x4d\xcb\xcb\x23\x12\xc0\x23\x99\x77\xc0\x05\xa3\x16\xc1\xd3\xc0\x0d\xf0\x8a\xaa\x0a\x4b\xcc\x28\x22\x5b\x24\xd1\x78\x52\xaf\xd3\x58\xd9\x18\x10\xde\x00\x81\x52\xf5\x8d\x00\x18\xf5\x1b\xd7\x80\x30\x7a\x10\xb7\xec\xbc\x95\x47\xe5\x7d\xa9\x20\xfe\x41\xbb\x10\xac\x2f\x8a\xfb\xe3\x45\x42\x66\x7f\x36\xac\x6d\x18\x55\x27\x69\xe4\xe0\xd0\x65\x5d\xac\x60\x8f\x5a\x22\x37\x2d\xe7\x40\xcb\x2e\xd4\x27\x15\x4e\xb0\x39\x49\xa1\x9e\x5b\xd7\x63\xd5\xa8\xbf\x6e\x0c\x74\x76\xd4\x06\x8a\x86\xb3\xaa\x2d\x65\xdc\x8c\x45\xb0\x0a\x50\x45\x5e\x1e\x7a\x2c\xc7\x30\x5c\x05\xf8\xbd\x45\xcf\x69\xb4\x3a\xb1\x07\x2d\x76\x05\x13\x02\x68\x74\xb6\x3e\xa5\xce\xae\xeb\xf7\x43\xd8\xb7\x44\xae\x30\x60\x6d\xbd\x41\x6a\x66\x3b\xec\xee\xc4\x09\x5c\x22\x4c\x6f\x8e\xb8\x69\x30\x3d\xbc\xba\x44\x98\x84\x3b\x83\xc5\xab\xba\x91\x5d\xb4\x74\x47\x24\x9c\xe2\xd7\x8c\xcf\x5a\xd7\x8f\x1b\x7b\xa5\xe2\xe3\x6e\xbb\xc6\xfa\x8f\x45\x8f\x56\x4e\x41\xee\x40\x25\xd5\x0f\xd2\x5b\xf4\x5e\x76\xeb\x1a\xf1\xcf\x20\xdf\x11\x54\x42\x60\xea\xcb\xe2\x11\x71\x8c\xa8\x8c\x1d\xd8\x51\x39\xcc\xfc\xea\x59\x02\xa7\x88\x5c\xc3\x1e\x14\x8e\x61\xcd\x61\xbf\x60\x81\x1b\xfd\x91\xb5\xe5\x11\xf8\x0d\x7a\xc4\xf4\x30\x0a\x99\xbd\xa5\x1a\xf5\x90\x08\x2c\xf7\xa6\xd5\x2a\x14\x6d\xed\xb6\x6d\x12\x79\xa1\x8c\x8a\xbf\x93\x17\x41\xbf\xaf\x6e\xc0\x86\x89\x51\xdc\x45\x84\xb8\xee\x5b\x2c\x49\x02\x50\xee\x30\xbc\xe1\x4c\x4c\xcc\x11\x88\x64\xd8\xe4\x9d\xaf\x2c\xe9\xf0\xd2\x99\x3f\xb9\xf5\x15\xa3\x6a\x93\xae\x81\xe8\xdb\x3e\x6f\xd0\x57\x8e\x18\xee\xb3\x21\x28\x26\xce\x45\xcf\x2b\x2c\xb2\xc2\x73\xe8\x20\xac\x39\xc5\x45\x77\xdb\x35\xb0\x56\xb7\x5c\x8c\xd6\xf9\xe8\x39\x84\xbc\xcd\x11\xf1\x03\x8c\x16\xf1\xde\xb6\x5b\xb3\x06\xd3\xbd\xe8\x15\x47\x82\x6b\xa8\x11\xa6\x98\x1e\x52\x32\x69\x52\xe3\xf1\x23\x8f\x05\x5a\x2a\x15\xf9\x60\x85\xfb\xf3\x64\x3c\x11\x16\x87\xb3\x63\x6e\x3c\x21\x3b\x4e\xf6\x8b\x18\xaf\x95\x1d\xd3\xaf\xf2\xea\x6e\xd6\x78\x67\x8f\xb5\x1f\xcd\x00\x20\x8a\x53\xb3\x6a\x7d\x93\x33\x54\xbb\xa8\xbb\xa3\x7b\x16\x40\x61\x76\x92\xde\xc7\x8c\x19\xca\x0c\xad\xb7\xe8\x39\x43\x93\x1a\x18\x82\x5a\x51\xec\xb3\xe2\x35\x61\x48\x4e\x6b\x06\x07\x91\x24\x3f\x50\x12\x77\xde\xd5\x60\x0e\x06\x7a\xbe\xf5\x26\x8b\xc3\xb2\x1a\x33\xe9\x8a\x8e\xb1\x76\xc6\x90\x58\xf4\x37\xc1\x6e\xe0\x20\xfa\xa7\x6d\x4e\x5f\xb5\x1a\x45\x98\x4a\xe0\x7b\x54\x8e\xb6\x23\xa2\x69\x76\xde\x03\x92\xf0\x84\x52\x1c\xe9\x23\x22\x2d\x8c\x97\x37\xed\xcb\xd6\x50\xae\xd1\x24\xb8\x6e\x08\xa8\x26\xf1\x67\x9a\x33\xca\xa9\x5c\x4a\x63\x7f\xce\x5e\xfb\x7d\x2e\x98\x3c\xba\x5b\xbf\x37\x91\x02\xba\xb3\x7a\xd1\xed\xaa\xb5\x4d\x01\x26\x53\x3e\x25\x38\xe5\x41\xd2\x70\x75\xf6\x26\x8c\x57\x5d\x7d\x78\x4b\xe2\x77\x22\xa4\x45\xfa\xfc\xa8\xb0\x7c\xcf\x2e\xb0\xdb\x3c\x72\xbb\xc4\x6d\xbf\xe2\xb6\xfd\x96\xcb\xf6\xab\xef\xda\xaf\xe4\x16\xdf\x40\x2d\x8e\x48\x58\xf4\xcd\xdf\x6e\xfe\xe6\xbb\xdb\x2d\x08\xa2\xaa\xe5\x89\xf1\xcf\x7b\xc2\x9e\xc2\xd6\x1a\xe4\x91\x55\x61\x5b\x89\x38\xc7\x8a\x0c\xfa\x8d\x0e\x7b\xbf\xb1\x12\x25\xf2\xbf\x6d\xd4\x6d\xc7\x08\xcc\xa1\xba\xc5\x35\x9c\x15\xea\xff\x7d\x51\x23\x48\x30\xd7\x9f\xa1\xf3\x19\x45\x90\xf7\x05\x92\xbf\x42\x17\x30\x40\x25\xf1\x7f\x91\x98\xb7\x16\xe2\xac\xa8\x51\xf3\x49\x98\xb0\xf8\x4f\xc1\xe8\xcf\xd7\xe8\xe9\x12\x84\x40\x07\xc8\x18\x7c\x89\x9a\x41\x2a\x34\xdb\x13\x8c\xcd\xbf\x44\xcd\xc8\x76\x4f\x3c\xf6\x61\x76\x47\xdd\x72\x16\x76\x5b\xc7\x37\x1a\x5a\x2c\x1b\xb4\x02\x2e\xa2\x12\x43\x40\xa8\x32\xee\xdb\x04\x47\x90\x8a\x8e\x87\xa6\x34\x0a\xb9\x53\x07\x57\xce\x1e\x7b\x34\x5d\x35\x4a\x15\x9b\xbc\x0b\xc1\x3f\x46\x3b\x5a\xaa\xf0\x32\x41\x06\x82\x8e\x85\x5b\x39\x9e\x70\x8e\x10\x44\xb2\x16\x96\x0f\xdd\x06\xd5\x0d\xc2\x07\xcd\xbe\xd7\xa5\xf7\xc3\x63\x09\x39\x6e\x3e\x18\x8a\xb1\xc7\x44\x02\x9f\x63\x19\xe3\xe1\x39\xbe\xf5\x74\xd8\x37\x30\x8c\x07\x5e\x12\x51\x84\x5d\x04\x3d\x00\x31\xa4\x24\xee\xb2\x5b\xea\x3a\xa7\xf9\x59\x72\x34\x16\x5e\x1c\x8e\x8b\x71\x8c\xcb\xb7\xbc\x52\x11\xca\xb2\xa1\x25\x02\xe0\xe1\x16\x8f\xef\xba\xfe\x8e\xb3\xec\x2b\xc0\x8f\x6e\x49\xab\xf7\xb5\xfa\x55\xbe\x38\x63\x8f\x22\xae\x2e\x07\x34\xa3\x72\x80\xee\xb4\x15\x81\xcb\x89\x92\x81\x6f\xe5\x15\xaa\xa3\x0e\xc1\x5a\x5e\x42\x5c\x39\xfb\x22\x3b\xaf\x44\xb5\x1c\x4f\x43\x09\xcd\xb7\x46\x32\x99\x21\xbc\xcf\xe8\xe2\x49\x63\x71\xbb\xbd\xc6\x0b\x4c\x0f\x04\x34\x4a\xe6\x72\xfa\x41\x6a\xb2\x18\xc1\xd9\xd3\x92\x1a\x27\xb2\x54\x4a\xcb\x0d\x4c\xc3\x75\xc1\xd9\x53\x5c\x31\xd6\xbf\xa7\x0e\xa5\x09\xcd\x16\x4e\x8f\x48\x7a\xe7\x22\x7d\x42\xf6\x98\x0b\x49\x35\x08\x26\x65\x08\x4a\x8a\x84\x78\xc4\x55\x45\xe0\x6a\x24\x15\x50\x6f\x13\xec\x67\xed\x11\x88\xb4\xd2\x52\x83\x49\x19\xc9\x01\x12\xae\x8d\x65\xae\xf8\x9c\xcd\x03\x46\xed\xba\xfd\x86\xe9\x18\xa5\x25\xab\x1b\x44\xbb\xd1\x74\x41\x6c\xc3\x72\x2c\x10\xc9\x34\x4c\xc8\x3e\xfa\x4d\x5a\xad\x33\xcb\x59\x3d\x1c\x0e\xd8\x8b\xa3\x69\x7b\x14\x8c\xf8\x82\xcd\x46\x66\xa4\x28\xd8\x31\x20\xd0\x1c\x19\x9d\x43\x07\xd4\xba\xf8\x3a\x69\x73\x12\xa8\xe6\xb1\xc1\x25\xdf\xcb\x6f\x16\x5a\x5c\x31\x20\x89\x30\x89\x25\xdf\x85\xbd\x2e\x7e\x62\x21\x31\x3d\x6c\x5a\x21\x59\x0d\x3c\xf1\x40\xf1\x2a\x21\x92\x36\x37\x25\x19\xc5\xec\x19\x37\x7b\xcb\x5c\x02\x86\x24\xbc\xdd\x5f\x60\x2e\x8f\x51\x4c\x46\x42\x34\x8c\x9b\xc4\x9d\x77\xe9\xce\xab\xb6\x7e\x88\x69\x35\x45\x06\xc7\x1a\x86\xb3\x0b\x1f\x06\x51\x6b\x90\x0e\x35\xa5\xf6\xed\x5c\x4a\x8e\x1f\x5a\x09\x1e\x71\xe5\x20\x80\x3f\x42\xa5\x6f\xcb\xc5\x82\x50\x5f\xbb\x9b\xcc\x21\xa6\x48\x5f\x4e\xf9\x25\x39\xe5\x50\x9f\x4c\xce\x39\xc7\x5f\x5c\xed\x6f\xd2\xd8\x9e\x80\x24\x03\xbf\x2b\x21\x4e\x66\x5e\xd7\x83\xc4\x42\x6d\xf1\x23\x22\xb8\xd2\xfb\x78\x0d\xa2\x25\x8e\x51\x1d\x91\x50\x72\x8c\xbe\xe2\x9c\x0d\xe1\x2c\xe2\xde\xbd\x80\x4d\x4b\x7e\x85\x08\x3d\x58\xf3\x20\xa5\x57\xf8\x67\x35\x2a\x4a\x29\x2e\x32\xd8\xa1\x15\x4e\x96\x13\x13\xb2\x1e\x39\x52\x30\x49\x87\x8b\x29\x2b\x27\x0a\x7f\x8a\xbd\x58\xe4\x0d\x79\x29\x53\xbf\xdd\x41\x88\x26\x08\x1f\x74\x16\xe2\xd0\x7d\x98\xd3\xbc\x66\xdc\xc1\x76\xcf\x78\x6d\x62\x86\xf9\x6f\x69\x98\x0e\x1e\x85\xbe\x75\xdd\x9a\x24\x0a\x2c\xf7\x4a\xd4\xdb\xea\xbe\xd0\xd2\x70\x56\x82\x10\x1e\x4b\x9d\x7a\x57\xb7\x8f\x82\x7d\xd9\xdd\xc3\xca\x1f\x3c\x75\x42\x81\x9d\xf8\xf4\x35\x06\x52\x15\xa2\x81\x12\xef\x71\xe9\x19\x62\xf6\x5b\x9c\x5a\xb2\x01\xa4\xd2\x48\x19\xd7\x43\xb5\xf2\xd7\xbd\x80\xbd\x7c\x4f\xdf\x00\x05\x8e\xc8\x94\xc6\x83\xe9\x9e\xd3\x39\x8f\xe2\x41\xc4\xb9\x72\x5e\x7c\x86\xae\x60\xfb\x42\x1e\xc1\xcc\x55\xd4\x06\xae\x3f\x17\x6f\xf7\x12\xa8\x4a\x85\x2b\x85\x8f\x42\x72\x44\x05\xd1\x56\x9d\xda\x3a\x48\xfa\xf4\x9d\x9e\xab\xb5\x41\x9f\x31\x3d\x58\x95\x3a\xe5\x09\x14\x4a\x56\x88\x23\x7b\x52\x7f\x02\xad\x54\x1b\x2f\xfe\xbf\xc0\xb4\x28\x91\x80\x82\x32\x7f\x36\x73\xbb\xd9\x35\xb0\x2f\xd4\xbf\x99\x24\x6a\x1e\xee\xd1\x2a\xff\x97\xf9\xac\xa7\xdd\x55\x40\x25\xde\x63\xe0\xda\x5e\xa4\x4f\xb2\x81\x9e\x87\xc2\x30\xef\xc9\x8d\x0d\xde\x9d\xf9\x83\x5c\x2f\x92\x6b\x47\xa9\xff\x76\xb6\x2c\xf3\xf7\x29\x99\xff\x65\xfa\xad\xa9\xb7\x77\xcd\xa5\x64\x32\xe8\xf7\x09\xa6\x4d\x2b\x07\x70\x8f\x71\xbd\xd3\x02\x39\xc0\xfe\x13\x71\x9d\x01\xeb\x0c\x54\x67\x80\x3a\x03\xd3\x19\x90\xce\x40\x74\x06\xa0\x33\xf0\x9c\x01\xe7\x0c\x34\x67\x80\x39\x03\xcb\x19\x50\xce\x40\x72\x2e\x90\x23\x1c\xbb\xf2\xf6\x0f\x20\xff\x00\xf2\x5f\x1a\xc8\xf6\xd5\x3b\x40\xb3\x8f\xe4\xd3\x0f\x14\x7f\x69\xa1\x4f\x6d\x74\x9a\xaa\x48\x13\x36\x5c\xa7\xd3\xbc\xcd\xf5\x9e\x26\xd2\xa0\x80\x21\xf5\xef\xc1\x76\x93\x0d\x43\xaa\x90\x44\xe1\xd0\xf1\xb7\xa4\x53\x27\xce\xe8\xbd\x39\xb2\x96\x54\xc6\x16\x45\x94\x2c\x9f\xb4\x5f\xa6\xf6\xb3\x3d\x80\xc7\x26\x8f\x58\x44\x86\x2f\x3d\x34\x9d\xbe\x6d\x4c\x01\xa3\x70\xef\x49\xc5\xa5\x7e\x7e\x74\x74\xd4\x7f\x8a\xcc\x19\x11\x3d\x54\x4e\x7e\x94\x30\x72\xfd\x3f\xb3\x3b\xf3\x89\x69\xc2\xcc\x3e\x33\x55\x1b\xf4\x24\x0a\xbc\x5f\xdc\x22\x61\x76\xf2\x7b\x77\x2a\x7b\x83\x7a\xc1\x8d\xd9\x8b\xf4\xc6\x68\xd1\x21\xff\x54\xc6\x8b\xd6\x94\x67\x86\xc5\xb7\x2e\xbc\x2c\xa0\x6e\x64\xa7\x9c\x75\x4e\x61\x61\xb2\xc8\xd3\xef\xc8\x95\xc3\x25\xe4\xf0\xa5\x05\x21\x8b\x27\x24\x0a\xd1\x96\x25\x88\x7d\x4b\x48\x37\x64\xd4\xa7\x5f\x97\x60\x4f\x6c\xde\x8f\x8c\xe3\x47\xc6\xf1\xd7\xca\x38\xe6\xee\x37\x77\xd6\x4d\x50\xce\x8e\xa1\x88\x16\x2f\x5e\xb8\xa2\xfe\x8b\x17\xf9\xf1\x34\x23\x08\xc5\x92\x73\x51\x48\xfb\x07\xcf\x12\x68\xa5\x4b\xd3\xc5\xfb\x76\xf8\xb4\x2a\xf0\xf8\x6c\xe2\xdf\x0f\xad\xc6\xa2\x2e\xd4\xb0\xd1\xb7\xd3\x71\x55\x78\x34\xfd\xa5\x3d\x5f\xb1\x05\xe7\x55\x75\xcb\x94\x8a\xf5\xe8\xa9\x78\xb7\x5d\xbd\x1c\x1e\x74\x5f\xa6\x17\xef\xa7\x4c\xf3\xb7\x40\x40\x82\xff\xb5\xc9\xf2\x07\xff\xcb\xfa\x76\x12\xea\xfe\x3b\x79\x6d\xef\x77\x29\xfd\xd0\x54\xc8\x28\x7d\x2f\xbb\x0c\xbd\xde\xf2\xe4\x4e\xa1\x77\xcb\xcc\x13\x5e\x8a\x6b\x34\x44\xf3\xb3\xc5\x24\x7c\xf4\xcd\xe7\x58\x6e\x7a\xe2\xa8\x36\xbc\x8e\x3f\x80\x7c\x19\xe3\x7e\x34\x5b\xb2\xba\x9c\x9a\xf0\xbc\x69\x48\x37\xbc\xbb\xbc\xe5\xee\x21\x65\x5d\x66\x6d\x50\x42\xe5\x35\xd4\xec\x11\x7a\x3d\x07\xfb\x97\xbc\x0d\x9f\xd4\x37\xd8\xb8\xf6\xbf\x40\xc9\xd2\x77\x7a\x5e\x55\xe2\x17\xb3\xb6\xa2\x60\x14\x7e\xa9\x5b\x22\x71\x43\xa0\xff\x68\xb6\xb0\x1b\x03\xb6\x30\x99\xda\x96\xe8\x5a\x07\xb1\x1e\x18\xb6\x69\x18\x55\xae\x53\xc9\xc0\xea\x6e\xf4\x9d\x7c\x8a\x31\xdc\xad\xfe\x08\xdb\xc3\xa8\x2d\xd6\x22\xfc\x3d\x69\x58\x38\x2e\xdf\x85\xdf\x4f\xfe\x1d\x00\x00\xff\xff\xe4\xb3\x13\x6d\x82\x39\x00\x00")

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
