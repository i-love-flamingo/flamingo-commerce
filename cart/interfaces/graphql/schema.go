// Code generated by go-bindata. (@generated) DO NOT EDIT.

// Package graphql generated by go-bindata.// sources:
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

var _schemaGraphql = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xec\x5b\x5f\x6f\xdc\x36\x12\x7f\xf7\xa7\xa0\xf7\x5e\xb6\x41\xee\x8a\xbb\x47\xbf\x6d\x76\x93\x60\xd1\xc6\x69\xec\x6d\xef\x21\x30\x02\x5a\x9a\xdd\x25\x42\x89\x32\x49\xd9\x16\x0e\xfd\xee\x07\xfe\x93\x48\x8a\x94\xe4\xb4\x57\xa0\x77\xd7\x87\xc6\x12\x87\xc3\x21\xf9\xe3\xcc\x6f\x86\x5a\xd9\x35\x80\xb6\xac\xaa\x80\x17\xf0\x65\x8b\xb9\xfc\xb2\x83\x82\x71\x2c\xa1\x54\x4f\xe8\x5f\x17\x08\x21\x54\x60\x2e\xaf\x22\x39\xf5\xbf\x4b\xdd\x5a\xba\x1e\x3b\xa0\xe4\x11\x38\x01\x71\x85\x3e\x67\xb4\x5a\x99\xee\xf2\x4e\xf7\x3d\x81\x1c\x35\xbd\xe9\xb6\xac\x84\x75\x69\x1f\xd5\xc3\x15\xba\x95\x9c\xd4\xa7\xcb\xef\x62\x33\x46\xbd\x9d\xda\x0d\xa5\x3f\xe1\xae\x82\x5a\xde\xc0\x43\x4b\x38\x94\x7b\x09\x95\x88\xfb\xff\xc4\x49\x61\x9b\x2e\xfb\xb9\xde\xb6\x55\x85\x79\x17\xcb\xda\xd7\x97\x17\xbf\x5e\x5c\x24\x56\xce\x36\xdb\x35\x2b\x89\x28\x58\x5b\xcb\xd1\x88\x9b\xa6\xa1\x04\xca\x9d\x6b\x37\xc3\x8a\xb6\x8a\x1b\xbc\x8e\xda\xca\x48\xee\x3d\x39\xca\x2d\xe6\x65\x56\xee\x3d\xc7\x75\x79\x60\x12\xd3\x7f\x12\x79\x9e\x15\xd7\x92\x6e\xf0\xa0\xc7\xa6\x52\xaf\x92\xfd\xce\x58\x8c\xcd\x7e\xc3\x18\x05\x5c\xf7\x13\x3b\xe0\x67\x18\x2d\x83\x7e\xe9\x24\xec\x4e\xdd\x02\x85\x42\x12\x56\x2b\x89\xdb\x86\x12\xf9\x0b\xa6\x2d\x98\xf1\xdf\x74\x1f\x40\x9e\x59\x29\xd6\x95\xf9\xf7\x0a\x7d\xb6\xa8\xb8\xfb\x6e\x64\x5c\x7a\x8b\x3c\x4c\x93\xf2\x0a\xed\x77\xc6\x46\xa8\x25\x91\xdd\x7e\xd7\xc3\x4c\xbf\xbd\x27\x94\x92\xfa\xb4\x29\x4b\x0e\x62\xbc\x8d\xe6\xb5\x96\x6c\x5a\x5e\x9c\xb1\x00\x3e\x42\x17\x70\xc1\x6a\x7b\x4a\x26\x0e\x47\x70\x26\x70\x59\x12\xb5\x08\x98\xee\xb0\xc4\x89\x71\xbd\x56\x63\x69\x13\x2d\xdf\xc8\x8c\xa8\xdd\x4c\x0f\x28\xab\x4f\xe2\xc0\x36\xad\x3c\xab\x15\x28\xd4\x31\xfa\x59\xcf\x22\xd8\x41\x1c\xb7\xc7\x0b\x85\x0d\x02\xb6\xac\x6d\x58\xad\x8e\xeb\x78\x8a\x43\x9b\x9d\x64\x09\x47\xdc\x52\xb9\x6d\x39\x87\xba\xe8\x42\x85\x52\x41\x91\x98\xe3\x1a\x29\x3a\xb8\x26\xab\x47\xfd\xb9\x35\xf0\xdc\xd7\xd6\x21\x35\x9c\x95\x6d\x21\xe3\xd7\x44\x04\xeb\x00\x65\x34\xcf\x53\x7f\x5e\x62\x34\x5d\x06\x67\xe4\x80\x9f\xd3\x27\xc2\x89\xdd\x6b\xb1\x6b\xc8\x08\xe0\xd1\xf9\xfd\x9c\xf4\x10\x4e\xc0\x77\x95\xdf\xe4\x21\x43\xc7\xb8\xf3\x7a\xf9\x67\xe8\xc2\x09\x7c\xc0\xa4\xbe\x3d\x93\xa6\x21\xf5\xe9\xed\x07\x4c\x68\xb8\x39\x44\xbc\xad\x1a\xd9\x45\x8b\x77\xc6\xc2\x29\x7e\xc7\xf8\xa4\x79\x7d\xbf\xf1\xb4\x94\x1f\xde\xef\xd6\x44\xff\x33\x3f\xa5\x4b\xa7\x61\x71\x4f\x25\xd6\xf7\xd2\xdb\xf4\x49\x76\xeb\x0a\xf3\xaf\x20\x7f\xa2\xb8\x80\xc0\xd8\xd7\xe8\x11\x73\x82\x6b\x19\x4f\x61\x5f\xcb\x61\xe8\xb7\xcf\x12\x78\x8d\xe9\x0d\x1c\x41\x81\x19\xd6\x1c\x8e\x73\x26\xb8\xee\xbf\xb0\xb6\x38\x03\xbf\xc5\x8f\xa4\x3e\x8d\x7c\x73\x6f\xaa\xea\x79\x80\x94\x8b\x31\x6f\xad\x42\xd1\x56\x6e\xe7\xb2\xf0\x0b\x65\x94\xa3\xcf\x46\x9c\x51\x87\xf7\x9c\x09\x31\xd3\xc5\xc1\xc1\xf5\xd9\x32\x31\x0a\x0a\x98\x52\xd7\x7c\x20\x92\x26\x70\xe8\x4e\x91\x1e\x71\xfa\xa0\x2d\x31\x2a\x3a\x98\xcb\x66\x1d\x44\xc4\xe9\x23\x5f\x5d\xb3\x5a\x6d\xec\x0d\x50\x4d\x46\x96\x75\x7a\x61\x8f\x21\xd8\x0e\xee\x34\x71\x9c\x7a\xda\x63\xe1\x18\x1e\x5f\x87\x7b\x4d\x79\xde\x74\x87\xae\x81\xb5\x8a\x94\x31\xc4\x67\xfc\xee\xe0\x2c\xb7\x67\xcc\x4f\x30\x5a\xc5\x2f\xf6\xfd\x00\x88\x31\x6f\x89\x3d\xc8\x0d\x54\x98\xd4\x0a\x67\x09\x99\x74\x44\xf7\xf8\x9b\x47\x56\x2d\xd5\x8b\x27\x61\x84\xb7\x56\xc0\xce\x44\x58\x20\x4e\xf6\xb9\xf5\x84\x6c\x3f\xd9\xaf\xe2\x68\xb1\x4c\x9f\x7e\x99\x2f\xef\x26\x8d\x77\xf6\x58\xfb\xf1\x04\x02\x22\xf7\x36\xa9\xd6\x37\x79\x81\x6a\xe7\xad\xf7\xf5\x91\x05\x58\x98\x1c\xa4\x9f\xe3\x82\x11\x8a\x05\x5a\x0f\xf8\x79\x81\x26\xd5\x31\x44\xb5\x4a\x01\xae\xd0\x3b\xca\xb0\xcc\x6b\x06\x07\x91\x0c\xb5\xc0\xcf\x77\x5e\x44\x31\x27\x03\x3f\x1f\xbc\xc1\x46\xce\xfc\x80\x9f\xb3\x53\xd1\x8e\xd9\x8e\x18\x52\x92\x3e\x7e\xec\x07\xf6\xa2\x1f\xed\xeb\x74\x88\xd6\x28\x22\xb5\x04\x7e\xc4\xc5\x68\x3b\x22\x8a\x67\xc7\x3d\x61\x09\x4f\x38\x45\xaf\x34\xb3\xce\x6c\x94\x63\xdf\x63\x60\x47\xa3\x7c\xd1\x62\x79\x7c\x27\xc5\xad\x69\x0f\x2d\xa6\xe4\x48\x12\xa4\x39\xd9\xeb\x93\x13\xb7\x36\x6a\xef\x92\x71\x3a\x59\xcc\x4e\x6b\xb6\x86\x8d\xd1\x65\x12\x8e\x08\x71\x63\xf7\x9a\x1e\x74\x67\x98\xee\x68\x83\x48\xd5\x50\x50\xaf\xc4\x9f\x61\x2b\x47\x79\xb6\x4b\x73\xed\xe3\x34\x45\xeb\x4b\x04\x69\x7f\xb9\xf3\x9b\x27\x4c\x48\x8e\xac\x1c\x56\x6e\x74\xd5\xd6\xaf\x43\xfa\xd4\x67\x22\xc1\x28\xf7\xf7\xda\xe6\x09\xca\x4c\x3e\xb1\x2c\x9d\x98\xcb\x26\x5e\x40\x53\xbe\x85\xa5\xbc\x98\xa4\xbc\x90\x94\x7d\x03\x27\x3b\x63\x61\x01\x34\xc3\x0a\xfc\xfd\x77\xb4\x20\x88\x3e\xea\xcd\x13\xe3\x5f\x8f\x94\x3d\xcd\x9f\xf5\x02\x73\xae\x1d\x95\xff\xd2\xe1\xef\x47\x56\xe0\x54\xd2\xbd\x8b\xda\x6d\x27\x41\x38\x94\x07\x52\xc1\x15\x52\xff\x5f\x90\xf1\x6f\x5b\x21\x59\xb5\x91\x92\x93\xfb\x56\x2a\x46\x35\x3d\x69\x37\x20\xca\xfa\x33\x3c\x5f\xd2\x68\x05\xbc\x89\xea\x1f\x01\x5b\x5b\x12\xcc\x13\x0c\x44\x2a\xb6\x1f\x1a\xd3\xa8\xfd\xcd\xc1\x5b\x4e\x1e\x0e\x9c\xaf\x98\x4d\x94\xda\xa4\x2b\x89\xb9\xf7\xfb\xba\x50\xa7\x30\xc3\x35\x82\x86\xb9\xa0\x1f\x8f\x38\x49\x38\x22\x61\xeb\x85\xee\xbb\x2d\xae\x1a\x4c\x4e\x9a\xdf\xaf\x0b\xef\xc1\xa3\x21\x8b\x66\x7a\x6f\x48\xcc\x91\x50\x09\x7c\x92\xc7\x8c\xfb\x2f\x9a\x5f\x4f\xb9\x7d\x1b\xc3\xa3\xe3\x65\x2a\x28\x6c\xa2\xf8\x1e\xa8\x21\x3e\x71\x93\xdd\x58\xd7\x98\xe7\x80\xc9\xde\x44\x78\x3e\x2b\x2e\x47\x32\x2e\x3f\xf2\x52\x1d\x66\xcb\xb8\xf2\xc7\xc9\x0b\x4a\xde\x3e\xa6\xf2\x79\x9f\xe3\x05\x30\xd2\x6f\x32\x03\xf8\x7a\xfd\x62\x64\x5c\x4e\x88\xfc\x93\xae\x55\x34\xa3\x5a\x85\x6e\xb4\xe5\x8a\x0f\x99\x7a\x86\x6f\xe6\x35\xae\xa2\x06\xc1\x5a\x5e\x40\x5c\xdc\x7b\x90\x9d\x57\x43\x0b\x3d\xd5\x0f\xd0\x05\x99\xfb\x58\x42\x53\x93\x91\xcc\x09\x64\x58\xc6\x5c\x7f\x85\xce\x4f\x3e\xcd\x5f\x7d\xe2\x18\x0f\x1a\x8b\xdb\x1d\x36\xb3\x20\xf5\x89\x82\x06\xca\x54\xed\x60\x90\xca\x16\x4a\x38\x7b\x9a\x53\xe3\x44\xe6\x6a\x7d\x2f\xf4\x50\x7f\xb1\xba\xe3\xc2\xb9\x7e\xce\x1e\x4d\xe3\xa7\x2d\xa2\x1e\xb1\xf4\x4e\x47\xfa\x9c\x1c\x09\x17\xb2\xd6\x38\xc8\xca\x50\x9c\x14\x09\x21\x49\xca\x92\xc2\xf5\x48\x2a\x20\xaa\xc6\xf1\x4f\xda\x23\x30\x6d\xa5\x8d\xa5\x59\x19\xc9\x01\x12\x53\x1b\xcb\x5c\xf3\x29\x9b\x07\x98\xda\x75\xfb\x91\xd4\x63\xa0\x16\xac\x6a\x70\xdd\x8d\x86\x0b\x3c\x1c\x91\x63\x81\x48\xa6\x61\x42\xf6\x3e\x30\x6b\xb5\xce\x61\x27\xf5\x70\x38\x11\xcf\x9b\xa6\xed\x51\x40\xe2\x33\x36\x1b\x99\x91\xa2\x60\xc7\x80\x42\x73\x66\xf5\x14\x3a\xa0\xd2\xe5\xe1\xac\xcd\x99\x34\x4c\xdf\x8a\xb8\x3c\x7f\xc1\xf5\x8a\x96\xc7\x74\x07\x12\x13\x3a\xbe\xc2\x0b\x9b\x9d\x17\x25\x42\x92\xfa\x64\xc8\x13\xf0\x14\xb1\x7a\x9b\x90\xc9\x98\x9c\x12\x8d\x7c\xf7\xd4\x5c\x7b\xe3\x5c\xee\x82\x25\x7c\x3c\xbe\x21\x5c\x9e\x23\xe7\x8c\x85\x68\x18\x37\x85\x02\xde\xa5\x1b\xaf\xdb\xea\x3e\x66\xa3\x35\x36\x68\xd6\x60\x9c\x5e\xfe\xd0\x9d\xba\x10\x1e\x91\xcc\x79\x12\x6a\xf0\x28\x80\x3f\x42\xa9\x23\xea\x7c\x61\x2a\x56\x62\xc7\x5e\xbd\x07\x89\x44\x03\x05\x39\x92\x02\x61\xd7\x8c\xee\x3b\xf4\x15\xba\x95\x0b\x19\xb1\xd7\x0f\x75\xff\x00\x9d\x0e\x36\x99\xa1\x5d\xb3\x1d\xd2\x57\x65\xfd\xa5\xce\xa1\x67\xaa\x53\xae\x08\x9a\xcd\x29\x72\xfc\x76\x49\x1d\x2b\xb3\x66\x3d\x81\x4a\x0e\x3a\x49\xd3\x5c\x19\x35\x6b\x6e\xcf\xb3\x92\xc1\xcd\x55\x63\xb3\xc9\xd8\xcd\x20\x31\x93\x90\xfd\x82\x29\x29\x35\x44\x6f\x40\xb4\xd4\x11\xc7\x33\x16\x4a\x8e\xd5\x6f\x39\x67\x83\xbf\x8e\x32\x8d\x5e\xe0\x03\x08\x81\x4f\xf0\x43\xbc\x79\x44\xd3\x3d\xa5\x57\xf8\xce\x28\x91\xe9\x0f\x76\x68\x85\xd9\x5a\x43\x42\xd6\xa3\x80\x0a\xe9\x69\x7f\x98\xb3\xf2\xd7\x8b\xe4\x30\x9f\x64\x77\x03\x42\x72\x52\x8c\x96\x86\x08\xd7\x32\xf0\xe0\x70\x61\x2a\xfc\xbc\xa1\x94\x3d\x79\xed\x68\x60\x6a\xfd\xee\xed\xc8\xb1\x27\x92\x5e\xab\xd1\xcd\xb8\x17\xb6\x67\xbc\xb6\xa2\x93\xf6\xa4\x0f\x69\x35\x53\xcf\xce\x21\x45\xab\x11\x5e\x00\xce\xe9\x0f\xd3\xcd\x77\x8c\xbb\x53\xb6\xb2\x2d\x2e\x54\xa0\xa3\x6a\x2b\xb1\xc4\xc6\x35\xa8\x47\xe3\xde\xcd\x7f\xc9\x40\xa2\xf4\x19\x6d\xc3\xb6\x22\x76\x44\xa2\x35\x47\xc0\xdd\xf6\xbb\x41\x5e\x23\xa8\x1a\xd9\x21\x72\xec\x87\x25\x42\x79\x09\x52\xae\x9c\xc3\xb0\x6a\x52\x05\x28\x35\x9c\x07\xfa\xbe\x0e\xb5\xba\x3d\xb3\x27\xa1\xb4\xca\x33\x20\x0e\x0f\x2d\x08\x89\x9e\xb0\x40\xa2\x2d\x0a\x10\xe2\xd8\x52\xda\x29\x8e\xae\x1e\xc0\x8e\xd5\x3f\x0e\x54\x37\x93\x71\xdb\x9b\xed\xfe\x0a\xc8\x03\xd4\xb7\x19\xbc\x78\xe8\x84\x02\xb7\x7f\xef\x08\xd0\x72\x70\xf0\x83\x21\xe6\xbc\x08\xbb\x8d\x4a\x4a\x9f\xb4\x71\xa6\xac\x95\xbf\xeb\x05\x2c\x3b\x5b\xbd\x87\x1a\x38\xa6\x39\x8d\x27\xd3\x3c\xa5\x73\xda\x0b\x0c\x22\x6e\x2a\x1b\x15\x39\x14\x6e\xd4\xf6\xe9\xb1\x50\x65\x8e\xfb\xdf\xd0\xc7\xa3\x84\x1a\xb5\x02\x4a\x05\x49\x24\x39\xae\x05\xd5\x56\xad\x6c\x65\x29\xed\xbd\x56\x1b\xb5\x36\xf8\xab\x42\x9f\x51\xa9\x33\xe3\x40\xa1\x64\x48\x9c\xd9\x93\xfa\x17\xea\x52\xbd\xe3\xe8\xaf\x88\xd4\xa8\xc0\x02\x50\xcd\xfc\xd1\x0c\xf7\xb1\x6b\x60\xbf\xb3\xf8\xd1\xe4\xda\xd3\x27\x30\x5a\xe5\xff\xb2\x39\xeb\x61\xf7\x25\xd4\xd2\x14\xfb\x95\xbd\xd8\xf8\x12\x0d\x3d\x0f\x85\x61\x6e\x3c\x99\x6e\x79\x7e\xea\xff\x29\xd7\x6c\xca\xe5\x12\xad\xbf\x5f\xcd\xcb\xfc\x23\x27\xf3\xbf\x9c\x94\xe9\x84\xcc\x0b\xb7\x29\x99\x05\x49\x19\xa9\x9b\x56\xe6\x01\xbd\xd7\xcd\x4b\x50\xfd\x07\x82\x7a\x01\xa6\x17\x40\x7a\x01\xa2\x17\x00\x7a\x01\x9e\x17\xc0\x79\x01\x9a\x17\x80\x79\x01\x96\x17\x40\x79\x01\x92\x17\x00\x79\x01\x8e\x17\xc0\xf8\xb7\xa0\xd8\x5d\x82\x58\x34\xfb\x48\x5e\xfd\x5c\x93\x87\x16\x7a\x5a\xaa\xf3\x21\x15\x5d\x88\x09\x0a\x9d\x0e\x70\xae\x75\x95\xa0\xb0\x41\x28\xe9\xaf\x23\x33\xb4\xb4\x0c\x2d\xc9\x14\x38\xfa\xe3\xd6\xd3\xc3\x96\x96\xc6\x10\x95\x9c\xdb\xa8\x1b\x91\x53\x74\x0f\x5e\xcc\x3d\x13\x11\x59\x3d\x77\x79\xb3\xfa\xd8\x98\x02\x00\x72\x57\x34\xc8\x7c\x5b\xeb\x82\xb6\x7f\x05\xb6\xa4\x47\x74\x41\x16\x75\xb1\x97\x5e\xc3\xc2\x97\x58\x02\xfa\x1e\x49\x52\x81\x5b\xab\xf8\x5a\x6c\xfa\x8a\x2b\x91\x24\xfc\xb1\x9b\xfb\xc2\x9c\x23\x60\xfd\x53\x1b\x2b\xcc\xfe\xff\xd6\xfd\x5d\xbc\xad\xbd\xe0\xd6\xec\xe0\x7f\x6a\x3b\x27\x53\xaf\x32\x5a\xec\x3f\x67\xee\xf5\x73\xa3\xd6\xc1\x01\xc7\x2d\xac\x59\x35\xf1\x25\x48\xc5\x7e\x4f\x4b\xa6\x9c\x60\x68\xc4\xcb\x7c\x21\xae\xd1\xab\x57\xae\x8a\xfa\xea\xd5\xf2\xa3\xb3\x00\x7b\xb1\xe4\x14\xf8\xb2\x9e\xde\x55\xf4\xfc\x69\x2d\x28\xeb\xcd\x04\x0d\xaf\x32\xea\x2b\xce\xcf\x37\xbe\xbb\xff\x3c\x61\xe5\xe5\x9d\x36\x01\x9e\xa5\x4a\x68\x34\x86\x3e\xb5\xc3\x17\x2d\x99\xfb\xc7\x6d\xe2\xd7\x37\x41\xe3\x65\xa2\xb7\xc5\x3f\x1b\x7d\x5d\x15\x57\xe2\xec\x5e\x4c\xd5\xa5\x10\x07\xd9\xf2\xba\x07\xac\xbd\x41\x54\x87\x92\xf7\x35\x2a\x95\x98\x49\xe0\x95\x70\xae\x0d\x9b\xda\x14\x7a\x68\xb1\xfe\x79\x85\x75\x66\x80\x0a\xfd\xcd\xbf\xd4\xdf\xe4\x20\x5c\x9b\x68\x77\x22\x8f\x50\x47\x20\x9b\xb2\x69\x9d\xbb\x07\x4d\x7e\xb3\xfd\x3a\xbd\x81\xa3\x22\x72\xaa\x1e\x37\xda\xb2\x0f\x96\x5e\x26\x77\x6d\x53\x96\x07\xa6\xfe\x1a\x1b\xb8\xdf\x5d\xbe\x1e\xae\x53\x97\x9a\x34\xbb\xd3\x3b\xa0\x20\x41\xfd\xe9\x60\xfc\xd2\x9f\x4e\x4d\x6b\xde\x4b\xa8\xfa\xef\xea\xf5\x1c\x7e\x27\xf5\xc6\x61\x2a\xf5\x9f\x64\xb7\x60\x04\x6f\xf1\x96\x0c\xb6\xda\x94\xa5\xf8\xde\x0c\x22\x34\xc6\x5c\x1d\xd1\xdd\x94\x5a\xa0\xfa\x78\x4c\x01\xcf\xa8\x08\x03\xee\x1a\x0f\x21\x7e\x8e\xd9\x8d\x8c\x1d\x57\x3a\xf3\xcb\x13\x55\xf4\xd6\xf1\x57\x80\xaf\x63\xd7\x3a\x1a\x2d\x59\x13\x4c\x0d\xb8\x69\x1a\xda\x0d\xd7\x0d\x1f\xb9\xbb\x3e\x58\x17\xcb\x37\x3a\xa1\xf7\x06\x2a\xf6\x08\xbd\xb2\x93\xfd\xe3\x05\xe8\xc9\x2a\x1d\xac\x5d\xfb\x1f\x9a\x2c\x57\x1a\x82\x84\xd5\xf0\x7d\xd5\x52\x49\x1a\x0a\xfd\x47\xb8\x0e\x2e\x20\xf2\xe0\x88\x78\x29\x88\xf5\xc0\xfe\xcd\x8b\xec\x8f\xbe\xfc\x44\x45\xff\x90\x6d\x5a\x4e\xa1\xe5\x2e\x05\xf0\xdf\x6c\x7b\xc4\x58\xd6\x22\x7c\xce\x1a\x96\x20\x19\xe1\x0f\xf2\x16\xb3\x23\x17\x8a\x14\xb5\x11\xb3\xc7\x52\x8b\xc5\xc4\x37\x58\x92\x21\x34\x1b\xce\xdd\x07\xa0\xc9\x83\x1e\x7d\x99\xf2\xc2\xf8\xfe\x72\xcf\x94\xb3\x32\x88\x88\x64\xc1\x0e\x12\x88\xbe\x94\x59\x97\x29\x8b\x27\xc8\xce\x02\xfb\x7f\xbd\xf8\x77\x00\x00\x00\xff\xff\xb2\xe1\xee\x39\x34\x3c\x00\x00")

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
	canonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[canonicalName]; ok {
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
	canonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[canonicalName]; ok {
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
// AssetDir("foo.txt") and AssetDir("nonexistent") would return an error
// AssetDir("") will return []string{"data"}.
func AssetDir(name string) ([]string, error) {
	node := _bintree
	if len(name) != 0 {
		canonicalName := strings.Replace(name, "\\", "/", -1)
		pathList := strings.Split(canonicalName, "/")
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
	canonicalName := strings.Replace(name, "\\", "/", -1)
	return filepath.Join(append([]string{dir}, strings.Split(canonicalName, "/")...)...)
}
