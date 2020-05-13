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

var _schemaGraphql = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xec\x5a\x4d\x6f\xdc\x38\xd2\xbe\xfb\x57\xa8\xfb\xbd\xf4\x1b\x64\x67\xb0\x7b\xf4\xad\xdd\x9d\x04\x8d\x19\x3b\x19\xdb\x33\x7b\x18\x18\x01\x2d\x55\x77\x73\x43\x91\x0a\x49\xd9\x16\x16\xf9\xef\x0b\x7e\x49\x24\x45\x49\x4c\x02\x04\xd8\x8f\x39\x4c\xdc\x64\xb1\x58\x24\x1f\x56\x3d\x55\x94\xec\x1a\x28\x76\xac\xae\x81\x97\xf0\x71\x0f\x25\xe3\x48\x42\xb5\x43\x5c\x16\xff\xbc\x28\x8a\xa2\x28\x11\x97\x97\x83\x88\xea\x59\xe9\x8e\xca\x09\xef\x81\xe0\x27\xe0\x18\xc4\x65\xf1\x67\x20\xb8\x8f\x44\xba\xd5\x83\x1e\x7a\x82\x71\xd7\x55\xb7\x63\x15\x6c\x2a\xfb\x53\xfd\xb8\x2c\xee\x24\xc7\xf4\xb4\xfa\xff\xc8\x80\xd1\x60\xa7\x75\x4b\xc8\x07\xd4\xd5\x40\xe5\x2d\x7c\x6e\x31\x87\xea\x20\xa1\x16\xd1\xf0\x8f\x1f\x38\x2e\x6d\xd7\xaa\x5f\xe4\x5d\x5b\xd7\x88\x77\xb1\xac\x6d\x5e\x5d\x7c\xb9\xb8\x08\x77\xcb\xef\xb6\x9b\x55\x61\x51\xb2\x96\xca\x78\xc6\x6d\xd3\x10\x0c\xd5\xde\x75\x6b\x61\xd1\xd6\x71\xbb\x37\x4c\xdb\x18\xc9\xbd\xc3\x47\xb9\x43\xbc\x9a\x94\x7b\xc7\x11\xad\xee\x99\x44\xe4\xef\x58\x9e\x17\xc5\xb5\xa4\x9b\x3c\x18\xb1\xad\x55\x53\x72\xdc\x19\x89\xb1\xd9\x57\x8c\x11\x40\x74\xd5\x6b\x46\x2f\x30\xda\x76\xdd\x98\xde\x47\xbb\x7f\xb8\xba\x2c\x0e\x7b\xa3\x05\xa8\xc4\xb2\x3b\xec\x7b\x14\xe8\xd6\x47\x4c\x08\xa6\xa7\x6d\x55\x71\x10\xa3\x6d\x36\xad\x5a\xb0\x69\x79\x79\x46\x02\x78\x24\xf3\x01\xb8\x60\xd4\x22\x78\x1a\xb8\x01\x5e\x51\x55\x61\x89\x19\x45\x64\x8f\x24\x1a\x4f\xea\x75\x1a\x2b\x1b\x03\xc2\x3b\x20\x50\xaa\xbe\x11\x00\xa3\x7e\xb3\x34\x20\x8c\x9e\xc4\x3d\xdb\xb6\xf2\xac\x56\x5f\x2a\x88\xff\xae\x97\x10\xec\x2f\x8a\xfb\xe3\x4d\x42\xe6\x7c\x76\xac\x6d\x18\x55\x37\x69\xb4\xc0\xa1\xcb\x2e\xb1\x82\x23\x6a\x89\xdc\xb5\x9c\x03\x2d\xbb\x50\x9f\x54\x38\xc1\xe6\x26\x85\x7a\xee\x5d\x8f\x55\xa3\xfe\xdc\x19\xe8\x1c\xa8\x75\x14\x0d\x67\x55\x5b\xca\xb8\x19\x8b\x60\x17\xa0\x8a\x56\x79\xea\xb1\x1c\xc3\x70\x15\xe0\xf7\x1e\xbd\xa4\xd1\xea\xc4\x1e\xb5\xd8\x0d\x4c\x08\xa0\xd1\xdd\xfa\x33\x75\x77\x5d\xbf\xef\xc2\xbe\xc5\x73\x85\x0e\x6b\xef\x0d\x52\x33\xdb\x61\x0f\x17\x4e\xe0\x1a\x61\x7a\x77\xc6\x4d\x83\xe9\xe9\xcd\x35\xc2\x24\x3c\x19\x2c\xde\xd4\x8d\xec\xa2\xad\x3b\x23\xe1\x14\xbf\x65\x7c\xd6\xba\x7e\xdc\x78\x55\xca\x3f\x1e\xf6\x1b\xac\xff\x59\x5c\xd1\xca\x29\xc8\x1d\xa8\xa4\xfa\x41\xfa\x88\x7e\x93\xdd\xa6\x46\xfc\x13\xc8\x0f\x04\x95\x10\x98\xfa\xba\x78\x42\x1c\x23\x2a\xe3\x05\x1c\xa8\x1c\x66\x7e\xf3\x22\x81\x53\x44\x6e\xe1\x08\x0a\xc7\xb0\xe1\x70\x5c\xb0\xc0\x8d\xfe\x83\xb5\xe5\x19\xf8\x1d\x7a\xc2\xf4\x34\x72\x99\xbd\xa5\x1a\xf5\x90\x70\x2c\x1f\x4d\xab\x55\x28\xda\xda\x1d\xdb\x24\xf2\x42\x19\xe5\x7f\x27\x03\x41\x7f\xae\x6e\xc0\x8e\x89\x91\xdf\x45\x84\xb8\xee\x7b\x2c\x49\x02\x50\xee\x32\xbc\xe3\x4c\x4c\xcc\x11\x88\x64\xd8\xe4\xdd\xaf\x2c\xe9\x30\xe8\xcc\xdf\xdc\xfa\x86\x51\x75\x48\xb7\x40\x74\xb4\xcf\x1b\xf4\x95\x23\x86\x78\x36\x38\xc5\xc4\xbd\xe8\x79\x85\x45\x56\x78\x0f\x1d\x84\x35\xa7\xb8\xea\xee\xbb\x06\x36\x2a\xca\xc5\x68\x9d\xf7\x9e\x83\xcb\xdb\x9d\x11\x3f\xc1\x68\x13\x3f\xda\x76\x6b\xd6\x60\xba\xe7\xbd\x62\x4f\x70\x0b\x35\xc2\x14\xd3\x53\x4a\x26\x4d\x6a\x3c\x7e\xe4\xb1\x40\x4b\xa5\xa2\x35\x58\xe1\xfe\x3e\x99\x95\x08\x8b\xc3\xd9\x31\x77\x9e\x90\x1d\x27\xfb\x4d\x8c\xf7\xca\x8e\xe9\x77\x79\xf5\x30\x6b\xbc\xb3\xc7\xda\x8f\x66\x00\x10\xf9\xa9\x59\xb5\xbe\xc9\x19\xaa\x9d\xd7\x3d\xd0\x23\x0b\xa0\x30\x3b\x49\xbf\xc6\x8c\x19\xca\x0c\xad\xf7\xe8\x25\x43\x93\x1a\x18\x82\x5a\x51\xec\xcb\xe2\x2d\x61\x48\x4e\x6b\x06\x07\x91\x24\x3f\x50\x12\x0f\x5e\x68\x30\x17\x03\xbd\xdc\x7b\x93\xc5\x6e\x59\x8d\x99\x5c\x8a\xf6\xb1\x76\xc6\x90\x58\xf4\x91\xe0\x30\x70\x10\xfd\xd3\x36\xa7\x43\xad\x46\x11\xa6\x12\xf8\x11\x95\xa3\xe3\x88\x68\x9a\x9d\xf7\x84\x24\x3c\xa3\x14\x47\xfa\x03\x91\x16\xc6\xdb\x9b\x5e\xcb\xde\x50\xae\xd1\x24\xb8\x6e\x08\xa8\x26\xf1\x23\xcd\x19\xe5\x54\x2e\xa5\xb1\x3f\x67\xc3\x7e\x9f\x0b\x26\xaf\xee\xde\xef\x4d\xa4\x80\xee\xae\x5e\x75\x87\x6a\x63\x53\x80\xc9\x94\x4f\x09\x4e\xad\x20\x69\xb8\xba\x7b\x13\xc6\xab\xae\xde\xbd\x25\xf1\x3b\xe1\xd2\x22\x7d\xbe\x57\x58\x8e\xb3\x0b\xec\x36\x8f\xdc\x2e\x71\xdb\xaf\x88\xb6\xdf\x12\x6c\xbf\x3a\xd6\x7e\x25\xb7\xf8\x06\x6a\x71\x46\xc2\xa2\x6f\x3e\xba\xf9\x87\xef\xa2\x5b\xe0\x44\x55\xcb\x33\xe3\x9f\x8e\x84\x3d\x87\xad\x35\xc8\x33\xab\xc2\xb6\x12\x71\x8e\x15\x19\xf4\x1b\x1d\xf6\x7e\x65\x25\x4a\xe4\x7f\xfb\xa8\xdb\x8e\x11\x98\x43\x75\x8f\x6b\xb8\x2c\xd4\xff\xfb\xa2\x46\x90\x60\x6e\x3e\x41\xe7\x33\x8a\x20\xef\x0b\x24\x7f\x81\x2e\x60\x80\x4a\xe2\xff\x22\x31\x6f\x2f\xc4\x65\x51\xa3\xe6\x4f\x61\xdc\xe2\x3f\x04\xa3\x3f\xdd\xa2\xe7\x6b\x10\x02\x9d\x20\x63\xf0\x35\x6a\x06\xa9\xd0\x6c\x4f\x30\x36\xff\x1a\x35\x23\xdb\x3d\xf1\x78\x0d\xb3\x27\xea\xb6\xb3\xb0\xc7\x3a\x8e\x68\x68\xb1\x6c\xd0\x0a\xb8\x8a\x4a\x0c\x01\xa1\xca\x88\xb7\x09\x8e\x20\x15\x1d\x0f\x4d\x69\x14\x72\xa7\x2e\xae\x9c\xbd\xf6\x68\xba\x6a\x94\x2a\x36\x79\x01\xc1\xbf\x46\x07\x5a\x2a\xf7\x32\x41\x06\x82\x8e\x85\xa8\x1c\x4f\x38\x47\x08\x22\x59\x0b\xcb\xc7\x6e\x87\xea\x06\xe1\x93\x66\xdf\x9b\xd2\xfb\xe1\xb1\x84\x9c\x65\x3e\x1a\x8a\x71\xc4\x44\x02\x9f\x63\x19\xe3\xe1\x39\x6b\xeb\xe9\xb0\x6f\x60\xe8\x0f\xbc\x24\xa2\x08\xbb\x08\x7a\x04\x62\x48\x49\xdc\x65\x8f\xd4\x75\x4e\xf3\xb3\xe4\x68\x2c\x3c\x3f\x1c\x17\xe3\x18\x97\xef\x79\xa5\x3c\x94\x65\x43\x4b\x04\xc0\xc3\x2d\x1e\xc7\xba\x3e\xc6\x59\xf6\x15\xe0\x47\xb7\xa4\xd5\xfb\x5a\xfd\x2a\x5f\x9c\xb1\x47\x1e\x57\x97\x03\x9a\x51\x39\x40\x77\xda\x8a\xc0\xf5\x44\xc9\xc0\xb7\xf2\x06\xd5\x51\x87\x60\x2d\x2f\x21\xae\x9c\x7d\x96\x9d\x57\xa2\x5a\xf6\xa7\xa1\x84\xe6\x5b\x23\x99\x4c\x17\xde\x67\x74\xf1\xa4\xb1\xb8\x3d\x5e\xb3\x0a\x4c\x4f\x04\x34\x4a\xe6\x72\xfa\x41\x6a\xb2\x18\xc1\xd9\xf3\x92\x1a\x27\xb2\x54\x4a\xcb\x75\x4c\x43\xb8\xe0\xec\x39\xae\x18\xeb\xdf\x53\x97\xd2\xb8\x66\x0b\xa7\x27\x24\xbd\x7b\x91\xbe\x21\x47\xcc\x85\xa4\x1a\x04\x93\x32\x04\x25\x45\x42\x3c\xe2\xaa\x22\x70\x33\x92\x0a\xa8\xb7\x71\xf6\xb3\xf6\x08\x44\x5a\x69\xa9\xc1\xa4\x8c\xe4\x00\x89\xa5\x8d\x65\x6e\xf8\x9c\xcd\x03\x46\xed\xbe\xfd\x8a\xe9\x18\xa5\x25\xab\x1b\x44\xbb\xd1\x74\x81\x6f\xc3\x72\x2c\x10\xc9\x34\x4c\xc8\xde\xfb\x4d\x5a\xad\x33\xcb\x59\x3d\x1c\x4e\xd8\xf3\xa3\x69\x7b\x14\x8c\xf8\x82\xcd\x46\x66\xa4\x28\x38\x31\x20\xd0\x9c\x19\x9d\x43\x07\xd4\xba\xf8\x3a\x69\x73\x12\xa8\xe6\xb1\xc1\x25\xdf\xcb\x6f\x16\x5a\x5c\x31\x20\x89\x30\x89\x25\x3f\x84\xbd\xce\x7f\x62\x21\x31\x3d\xed\x5a\x21\x59\x0d\x3c\xf1\x40\xf1\x26\x21\x92\x36\x37\x25\x19\xf9\xec\x99\x65\xf6\x96\xb9\x04\x0c\x49\x78\x7f\xbc\xc2\x5c\x9e\x23\x9f\x8c\x84\x68\x18\x37\x89\x3b\xef\xd2\x9d\x37\x6d\xfd\x18\xd3\x6a\x8a\x0c\x8e\x35\x0c\x67\x37\x3e\x74\xa2\xd6\x20\xed\x6a\x4a\xbd\xb6\xad\x94\x1c\x3f\xb6\x12\x3c\xe2\xca\x41\x00\x7f\x82\x4a\x47\xcb\xc5\x82\x50\x5f\xbb\x9b\xcc\x21\xa6\x48\x5f\x4e\xf9\x25\x39\xe5\x50\x9f\x4c\xce\x39\xc7\x5f\x5c\xed\x6f\xd2\xd8\x9e\x80\x24\x1d\xbf\x2b\x21\x4e\x66\x5e\xb7\x83\xc4\x42\x6d\xf1\x0f\x44\x70\xa5\xcf\xf1\x16\x44\x4b\x1c\xa3\x3a\x23\xa1\xe4\x18\x7d\xc3\x39\x1b\xdc\x59\xc4\xbd\x7b\x01\x9b\x96\xfc\x02\x11\x7a\xb0\xe6\x41\x4a\xaf\xf0\xef\x6a\x54\x94\x52\x5c\x64\xb0\x43\x2b\x9c\x2c\x27\x26\x64\x3d\x72\xa4\x60\x92\x76\x17\x53\x56\x7e\xb9\x48\x4e\x73\x0b\x2a\xf1\x2a\x47\xfb\x82\x85\xeb\x19\xd8\x61\xb8\x2b\x35\x7a\xd9\x12\xc2\x9e\xbd\xfe\x62\xa0\x30\xfd\xd1\xed\xf1\xb1\x67\x58\x5e\xaf\xd1\xcd\xb8\x17\xd2\x16\x4a\x94\x8a\x67\xd9\x3b\x32\x64\xd0\x4c\xfd\x76\x57\x36\xda\x8a\xf0\xe9\x69\x49\x7f\x98\x7d\xbd\x65\xdc\x5d\xb0\xb5\xed\x71\x7e\xb4\x38\xaa\xbe\x0a\x49\xb4\x36\x21\x9e\xf1\xda\x78\x3f\xf3\x5f\xa8\xd6\xd3\x67\xb4\x0d\x67\x5a\xb0\x63\x21\x5a\x83\x7f\xf7\xbe\xec\x26\x79\x5d\x40\xdd\xc8\xae\xc0\xc7\x7e\x5a\x2c\x8a\x27\x35\x76\x6d\xc9\x87\x53\x93\xa8\x33\x7d\x54\xd3\x79\x88\xef\xeb\x4d\xeb\xbb\x33\x7b\x16\x4a\xab\x3c\x43\xc1\xe1\x73\x0b\x42\x16\xcf\x48\x14\xa2\x2d\x4b\x10\xe2\xd8\x12\xd2\x29\xf2\xaa\x7e\x80\x9d\xab\xff\x39\x70\xc0\x89\x8f\x12\xec\x8b\x6a\xff\x66\xe1\x01\xea\xdb\x0c\xce\x9e\x3a\xa1\xc0\x9d\xdf\x5b\x0c\xa4\x2a\x44\x03\x25\x3e\xe2\xd2\x33\xc4\x5c\x16\x61\x8f\x51\x49\xe9\x6b\x36\x2e\x26\x6b\xe5\x6f\x7b\x01\xcb\x5c\xd6\xef\x80\x02\x47\x64\x4a\xe3\xc9\x74\xcf\xe9\x9c\x77\x01\x83\x88\x5b\xca\xb6\xf8\x04\x9d\xc2\x8d\x3a\x3e\x3d\x57\x51\x9b\xbb\xfe\x53\xf1\xfe\x28\x81\x16\xad\x80\x4a\x41\xb2\x90\x1c\x51\x41\xb4\x55\x6b\x5b\x44\x4a\xbb\xae\xf5\x56\xed\x0d\xfa\xa4\xd0\x67\x54\xea\x7c\x31\x50\x28\x59\x21\xce\xec\x59\xfd\x0b\xb4\x52\x6d\xbc\xf8\x4b\x81\x69\x51\x22\x01\x05\x65\xfe\x6c\x86\x1a\xd8\x3d\xb0\xcf\xfb\xbf\x9a\x0c\x74\xfe\x06\x46\xbb\xfc\x1f\xb6\x66\x3d\xed\xa1\x02\x2a\xf1\x11\x03\xd7\xf6\x22\xe3\x4b\x34\xf4\x3c\x14\x86\x49\x63\x7a\xb3\xc6\x7e\xea\x7f\xe9\xc8\x62\x3a\xe2\x92\x90\xbf\x5e\x2e\xcb\xfc\x6d\x4a\xe6\xbf\x39\x61\xd1\xc9\x8a\x17\x6e\x53\x32\x19\x09\x0b\xa6\x4d\x2b\xa7\x01\x7d\xd0\xdd\x39\xa8\xfe\x81\xa0\xce\xc0\x74\x06\xa4\x33\x10\x9d\x01\xe8\x0c\x3c\x67\xc0\x39\x03\xcd\x19\x60\xce\xc0\x72\x06\x94\x33\x90\x9c\x01\xe4\x0c\x1c\x67\xc0\xf8\x7b\x50\xec\x9e\x04\x2c\x9a\x7d\x24\xaf\x7f\xa7\xf8\x73\x0b\x3d\x2d\xd5\xc9\x90\x8a\x2e\xd8\x04\x85\x4e\x07\x38\xd7\xbb\x4e\x50\xd8\x20\x94\xf4\xaf\x8e\x13\xb4\xb4\x0a\x2d\x89\x09\x57\x7c\xdd\x7a\x7a\xd8\x92\xca\x18\xa2\xf2\x57\x1b\x75\x23\x72\x5a\x3c\x82\x17\x73\xcf\x58\x44\x56\x2f\xbd\x65\xac\xdf\x37\x26\x47\x2e\xdc\x93\x45\x71\xad\x5f\xb8\x5c\xd0\xf6\x5f\xbb\x72\x46\x44\x6f\x61\x93\xef\xde\xc1\x8e\xf8\x14\xff\xc7\x1e\xcd\x57\x66\x0c\x01\x67\x9f\x3b\x16\x61\x4e\xef\x7b\x4f\x27\xfb\x50\x7a\xc1\x9d\xd9\xff\xf4\x61\x14\x33\x59\x4f\x15\xed\xd4\xbf\x41\xda\x33\x77\xeb\xdd\x86\x18\xb4\x66\x83\x0b\xd1\xe2\xd5\x2b\x57\x50\x7b\xf5\x2a\x1f\x68\x19\x27\x15\x4b\xce\x1d\x95\xf6\x6a\xf0\x22\x15\xdd\xd5\xf7\xe7\xb7\x76\xf8\xac\x21\x58\xf1\xe5\xc4\xb7\xfb\xab\xb1\xa8\x3b\x0f\x36\xfa\x6e\x31\xae\xc8\xcc\x0c\xc6\x8c\xfa\x95\x8a\xcd\xd4\xf3\x4c\xf2\x6b\xcd\xd7\xe9\x2d\x1c\x7d\xfe\x33\x2a\x85\x8c\xf6\xe3\xda\x46\xf6\x78\x4b\xb6\x55\x75\xcf\x94\x8e\xb1\x61\x87\xfd\xea\xf5\xf0\xba\x93\x61\xca\xdc\x7e\xee\x81\x80\x04\xff\xe9\x79\xf9\xeb\xdf\x65\x7d\x07\x09\x75\xff\xd1\xac\xb6\xf7\xbb\x94\xfe\xde\x54\xc8\x28\xfd\x4d\x76\x19\x7a\xbd\xed\x59\x98\x62\xbd\xad\x2a\xf1\xb3\xd1\x2f\xf4\x55\x76\xa5\x19\xf7\x30\x63\x1d\x63\xa9\xbf\x26\x97\xfa\xf3\x9a\x75\x02\x55\x46\x45\xe8\x05\x37\x68\xf0\xbb\x4b\xc1\x72\x04\x9d\x71\xf1\x28\x05\x66\x33\x6d\x54\x24\xd9\xc4\x9f\x51\xbd\x8e\x6f\xf0\x68\xb6\x64\x99\x25\x35\xe1\xb6\x69\x48\x37\x54\x6f\xdf\x73\x57\x8e\xdd\x94\x59\x27\x9b\x50\x79\x0b\x35\x7b\x82\x5e\xcf\xc9\xfe\x91\x87\x94\x49\x7d\x83\x8d\x1b\xff\x1d\x3b\x4b\x5f\x88\x0a\x46\xe1\xe7\xba\x25\x12\x37\x04\xfa\x4f\xef\x1c\x3e\x40\x4c\xa3\x21\x62\x07\x20\x36\x03\x83\x32\x0d\xa3\x12\x4e\x8a\xec\xad\x1e\x46\x5f\xdb\xa6\x88\xc7\x43\x0a\xd1\xdf\x6d\x7b\x18\x7f\xc4\x46\x84\xbf\x27\x0d\x0b\xc7\xe5\x2f\xe1\xcb\xc5\xbf\x02\x00\x00\xff\xff\x2d\x7a\x5b\x8b\xc8\x35\x00\x00")

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
