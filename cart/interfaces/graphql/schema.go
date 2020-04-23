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

var _schemaGraphql = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xec\x5a\xcd\x6e\xdc\x38\x12\xbe\xfb\x29\xd4\xbd\x97\x9e\x20\x3b\x83\xdd\xa3\x6f\x76\x77\x12\x34\x66\xec\x64\x6c\x4f\xf6\x10\x18\x06\x2d\x55\x77\x73\x43\x91\x0a\x49\xd9\x16\x16\xf3\xee\x0b\xfe\x49\x24\x45\x49\x4c\x02\x04\xd8\x9f\x39\x4c\xdc\x64\xb1\x58\x24\x3f\x56\x7d\x55\x94\xec\x1a\x28\xb6\xac\xae\x81\x97\xf0\xb0\x83\x92\x71\x24\xa1\xda\x22\x2e\x8b\x7f\x9d\x15\x45\x51\x94\x88\xcb\xf3\x41\x44\xf5\xac\x74\x47\xe5\x84\x77\x40\xf0\x13\x70\x0c\xe2\xbc\xf8\x14\x08\xee\x22\x91\x6e\x75\xaf\x87\x1e\x61\xdc\x75\xd9\x6d\x59\x05\x9b\xca\xfe\x54\x3f\xce\x8b\x5b\xc9\x31\x3d\xae\x7e\x8a\x0c\x18\x0d\x76\x5a\x2f\x08\xf9\x80\xba\x1a\xa8\xbc\x81\x2f\x2d\xe6\x50\xed\x25\xd4\x22\x1a\xfe\xf0\x81\xe3\xd2\x76\xad\xfa\x45\xde\xb6\x75\x8d\x78\x17\xcb\xda\xe6\xd5\xd9\x9f\x67\x67\xe1\x6e\xf9\xdd\x76\xb3\x2a\x2c\x4a\xd6\x52\x19\xcf\x78\xd1\x34\x04\x43\xb5\x73\xdd\x5a\x58\xb4\x75\xdc\xee\x0d\xd3\x36\x46\x72\xef\xf0\x41\x6e\x11\xaf\x26\xe5\xde\x71\x44\xab\x3b\x26\x11\xf9\x07\x96\xa7\x45\x71\x2d\xe9\x26\x0f\x46\x5c\xd4\xaa\x29\x39\xee\x84\xc4\xd8\xec\x4b\xc6\x08\x20\xba\xea\x35\xa3\x17\x18\x6d\xbb\x6e\x4c\xef\xa3\xdd\x3f\x5c\x9d\x17\xfb\x9d\xd1\x02\x54\x62\xd9\xed\x77\x3d\x0a\x74\xeb\x23\x26\x04\xd3\xe3\x45\x55\x71\x10\xa3\x6d\x36\xad\x5a\xb0\x69\x79\x79\x42\x02\x78\x24\xf3\x01\xb8\x60\xd4\x22\x78\x1a\xb8\x01\x5e\x51\x55\x61\x89\x19\x45\x64\x87\x24\x1a\x4f\xea\x75\x1a\x2b\x1b\x03\xc2\x5b\x20\x50\xaa\xbe\x11\x00\xa3\x7e\xb3\x34\x20\x8c\x1e\xc5\x1d\xbb\x68\xe5\x49\xad\xbe\x54\x10\xff\x43\x2f\x21\xd8\x5f\x14\xf7\xc7\x9b\x84\xcc\xf9\x6c\x59\xdb\x30\xaa\x6e\xd2\x68\x81\x43\x97\x5d\x62\x05\x07\xd4\x12\xb9\x6d\x39\x07\x5a\x76\xa1\x3e\xa9\x70\x82\xcd\x4d\x0a\xf5\xdc\xb9\x1e\xab\x46\xfd\xb9\x35\xd0\xd9\x53\xeb\x28\x1a\xce\xaa\xb6\x94\x71\x33\x16\xc1\x2e\x40\x15\xad\xf2\xd8\x63\x39\x86\xe1\x2a\xc0\xef\x1d\x7a\x49\xa3\xd5\x89\x3d\x6a\xb1\x6b\x98\x10\x40\xa3\xbb\xf5\x29\x75\x77\x5d\xbf\xef\xc2\xbe\xc5\x73\x85\x0e\x6b\xe7\x0d\x52\x33\xdb\x61\xf7\x67\x4e\xe0\x0a\x61\x7a\x7b\xc2\x4d\x83\xe9\xf1\xcd\x15\xc2\x24\x3c\x19\x2c\xde\xd4\x8d\xec\xa2\xad\x3b\x21\xe1\x14\xbf\x65\x7c\xd6\xba\x7e\xdc\x78\x55\xca\x3f\xee\x77\x1b\xac\xff\x59\x5c\xd1\xca\x29\xc8\x1d\xa8\xa4\xfa\x41\xfa\x88\x7e\x97\xdd\xa6\x46\xfc\x33\xc8\x0f\x04\x95\x10\x98\xfa\xba\x78\x42\x1c\x23\x2a\xe3\x05\xec\xa9\x1c\x66\x7e\xf3\x22\x81\x53\x44\x6e\xe0\x00\x0a\xc7\xb0\xe1\x70\x58\xb0\xc0\x8d\xfe\xc8\xda\xf2\x04\xfc\x16\x3d\x61\x7a\x1c\xb9\xcc\xde\x52\x8d\x7a\x48\x38\x96\x07\xd3\x6a\x15\x8a\xb6\x76\xc7\x36\x89\xbc\x50\x46\xf9\xdf\xc9\x40\xd0\x9f\xab\x1b\xb0\x65\x62\xe4\x77\x11\x21\xae\xfb\x0e\x4b\x92\x00\x94\xbb\x0c\xef\x38\x13\x13\x73\x04\x22\x19\x36\x79\xf7\x2b\x4b\x3a\x0c\x3a\xf3\x37\xb7\xbe\x66\x54\x1d\xd2\x0d\x10\x1d\xed\xf3\x06\x7d\xe5\x88\x21\x9e\x0d\x4e\x31\x71\x2f\x7a\x5e\x61\x91\x15\xde\x43\x07\x61\xcd\x29\x2e\xbb\xbb\xae\x81\x8d\x8a\x72\x31\x5a\xe7\xbd\xe7\xe0\xf2\xb6\x27\xc4\x8f\x30\xda\xc4\x07\xdb\x6e\xcd\x1a\x4c\xf7\xbc\x57\xec\x09\x6e\xa0\x46\x98\x62\x7a\x4c\xc9\xa4\x49\x8d\xc7\x8f\x3c\x16\x68\xa9\x54\xb4\x06\x2b\xdc\xdf\x27\xb3\x12\x61\x71\x38\x3b\xe6\xd6\x13\xb2\xe3\x64\xbf\x89\xf1\x5e\xd9\x31\xfd\x2e\xaf\xee\x67\x8d\x77\xf6\x58\xfb\xd1\x0c\x00\x22\x3f\x35\xab\xd6\x37\x39\x43\xb5\xf3\xba\x7b\x7a\x60\x01\x14\x66\x27\xe9\xd7\x98\x31\x43\x99\xa1\xf5\x0e\xbd\x64\x68\x52\x03\x43\x50\x2b\x8a\x7d\x5e\xbc\x25\x0c\xc9\x69\xcd\xe0\x20\x92\xe4\x07\x4a\xe2\xde\x0b\x0d\xe6\x62\xa0\x97\x3b\x6f\xb2\xd8\x2d\xab\x31\x93\x4b\xd1\x3e\xd6\xce\x18\x12\x8b\x3e\x12\xec\x07\x0e\xa2\x7f\xda\xe6\x74\xa8\xd5\x28\xc2\x54\x02\x3f\xa0\x72\x74\x1c\x11\x4d\xb3\xf3\x1e\x91\x84\x67\x94\xe2\x48\x1f\x11\x69\x61\xbc\xbd\xe9\xb5\xec\x0c\xe5\x1a\x4d\x82\xeb\x86\x80\x6a\x12\x3f\xd2\x9c\x51\x4e\xe5\x52\x1a\xfb\x73\x36\xec\xf7\xb9\x60\xf2\xea\xee\xfc\xde\x44\x0a\xe8\xee\xea\x65\xb7\xaf\x36\x36\x05\x98\x4c\xf9\x94\xe0\xd4\x0a\x92\x86\xab\xbb\x37\x61\xbc\xea\xea\xdd\x5b\x12\xbf\x13\x2e\x2d\xd2\xe7\x7b\x85\xe5\x38\xbb\xc0\x6e\xf3\xc8\xed\x12\xb7\xfd\x8a\x68\xfb\x2d\xc1\xf6\xab\x63\xed\x57\x72\x8b\x6f\xa0\x16\x27\x24\x2c\xfa\xe6\xa3\x9b\x7f\xf8\x2e\xba\x05\x4e\x54\xb5\x3c\x33\xfe\xf9\x40\xd8\x73\xd8\x5a\x83\x3c\xb1\x2a\x6c\x2b\x11\xe7\x58\x91\x41\xbf\xd1\x61\xef\x37\x56\xa2\x44\xfe\xb7\x8b\xba\xed\x18\x81\x39\x54\x77\xb8\x86\xf3\x42\xfd\xbf\x2f\x6a\x04\x09\xe6\xe6\x33\x74\x3e\xa3\x08\xf2\xbe\x40\xf2\x57\xe8\x02\x06\xa8\x24\xfe\x12\x89\x79\x7b\x21\xce\x8b\x1a\x35\x9f\x84\x71\x8b\xff\x14\x8c\xfe\x7c\x83\x9e\xaf\x40\x08\x74\x84\x8c\xc1\x57\xa8\x19\xa4\x42\xb3\x3d\xc1\xd8\xfc\x2b\xd4\x8c\x6c\xf7\xc4\xe3\x35\xcc\x9e\xa8\xdb\xce\xc2\x1e\xeb\x38\xa2\xa1\xc5\xb2\x41\x2b\xe0\x32\x2a\x31\x04\x84\x2a\x23\xde\x26\x38\x82\x54\x74\x3c\x34\xa5\x51\xc8\x9d\xba\xb8\x72\xf6\xda\xa3\xe9\xaa\x51\xaa\xd8\xe4\x05\x04\xff\x1a\xed\x69\xa9\xdc\xcb\x04\x19\x08\x3a\x16\xa2\x72\x3c\xe1\x1c\x21\x88\x64\x2d\x2c\x1f\xbb\x2d\xaa\x1b\x84\x8f\x9a\x7d\x6f\x4a\xef\x87\xc7\x12\x72\x96\xf9\x68\x28\xc6\x01\x13\x09\x7c\x8e\x65\x8c\x87\xe7\xac\xad\xa7\xc3\xbe\x81\xa1\x3f\xf0\x92\x88\x22\xec\x22\xe8\x11\x88\x21\x25\x71\x97\x3d\x52\xd7\x39\xcd\xcf\x92\xa3\xb1\xf0\xfc\x70\x5c\x8c\x63\x5c\xbe\xe7\x95\xf2\x50\x96\x0d\x2d\x11\x00\x0f\xb7\x78\x1c\xeb\xfa\x18\x67\xd9\x57\x80\x1f\xdd\x92\x56\xef\x6b\xf5\xab\x7c\x71\xc6\x1e\x79\x5c\x5d\x0e\x68\x46\xe5\x00\xdd\x69\x2b\x02\x57\x13\x25\x03\xdf\xca\x6b\x54\x47\x1d\x82\xb5\xbc\x84\xb8\x72\xf6\x45\x76\x5e\x89\x6a\xd9\x9f\x86\x12\x9a\x6f\x8d\x64\x32\x5d\x78\x9f\xd1\xc5\x93\xc6\xe2\xf6\x78\xcd\x2a\x30\x3d\x12\xd0\x28\x99\xcb\xe9\x07\xa9\xc9\x62\x04\x67\xcf\x4b\x6a\x9c\xc8\x52\x29\x2d\xd7\x31\x0d\xe1\x82\xb3\xe7\xb8\x62\xac\x7f\x4f\x5d\x4a\xe3\x9a\x2d\x9c\x9e\x90\xf4\xee\x45\xfa\x86\x1c\x30\x17\x92\x6a\x10\x4c\xca\x10\x94\x14\x09\xf1\x88\xab\x8a\xc0\xf5\x48\x2a\xa0\xde\xc6\xd9\xcf\xda\x23\x10\x69\xa5\xa5\x06\x93\x32\x92\x03\x24\x96\x36\x96\xb9\xe6\x73\x36\x0f\x18\xb5\xfb\xf6\x1b\xa6\x63\x94\x96\xac\x6e\x10\xed\x46\xd3\x05\xbe\x0d\xcb\xb1\x40\x24\xd3\x30\x21\x7b\xef\x37\x69\xb5\xce\x2c\x67\xf5\x70\x38\x62\xcf\x8f\xa6\xed\x51\x30\xe2\x0b\x36\x1b\x99\x91\xa2\xe0\xc4\x80\x40\x73\x62\x74\x0e\x1d\x50\xeb\xe2\xeb\xa4\xcd\x49\xa0\x9a\xc7\x06\x97\x7c\x2f\xbf\x59\x68\x71\xc5\x80\x24\xc2\x24\x96\xfc\x10\xf6\x3a\xff\x89\x85\xc4\xf4\xb8\x6d\x85\x64\x35\xf0\xc4\x03\xc5\x9b\x84\x48\xda\xdc\x94\x64\xe4\xb3\x67\x96\xd9\x5b\xe6\x12\x30\x24\xe1\xfd\xe1\x12\x73\x79\x8a\x7c\x32\x12\xa2\x61\xdc\x24\xee\xbc\x4b\x77\x5e\xb7\xf5\x63\x4c\xab\x29\x32\x38\xd6\x30\x9c\xdd\xf8\xd0\x89\x5a\x83\xb4\xab\x29\xf5\xda\x2e\xa4\xe4\xf8\xb1\x95\xe0\x11\x57\x0e\x02\xf8\x13\x54\x3a\x5a\x2e\x16\x84\xfa\xda\xdd\x64\x0e\x31\x45\xfa\x72\xca\x2f\xc9\x29\x87\xfa\x64\x72\xce\x39\xfe\xe2\x6a\x7f\x93\xc6\xf6\x04\x24\xe9\xf8\x5d\x09\x71\x32\xf3\xba\x19\x24\x16\x6a\x8b\x1f\x11\xc1\x95\x3e\xc7\x1b\x10\x2d\x71\x8c\xea\x84\x84\x92\x63\xf4\x0d\xe7\x6c\x70\x67\x11\xf7\xee\x05\x6c\x5a\xf2\x2b\x44\xe8\xc1\x9a\x07\x29\xbd\xc2\xbf\xab\x51\x51\x4a\x71\x91\xc1\x0e\xad\x70\xb2\x9c\x98\x90\xf5\xc8\x91\x82\x49\xda\x5d\x4c\x59\x39\x51\xf8\x53\xec\xc5\x22\x6f\xc8\x4b\x99\xfa\xed\x2e\x42\x34\x41\xf8\xa0\xb3\xe0\x87\x1e\xc2\x9c\xe6\x2d\xe3\x0e\xb6\x6b\xdb\xe3\xbc\x53\x71\x50\x7d\x15\x92\x68\x6d\x02\x27\xe3\xb5\xf1\x29\xe6\xbf\x50\xad\xa7\xcf\x68\x1b\x76\xaa\x60\x87\x42\xb4\x06\x55\xee\xd5\xd6\x4d\xf2\xba\x80\xba\x91\x5d\x81\x0f\xfd\xb4\x58\x14\x4f\x6a\xec\xda\x86\x74\xa7\x26\x51\xbd\x79\x50\xd3\x79\x38\xea\xab\x38\xeb\xdb\x13\x7b\x16\x4a\xab\x3c\x41\xc1\xe1\x4b\x0b\x42\x16\xcf\x48\x14\xa2\x2d\x4b\x10\xe2\xd0\x12\xd2\x29\x4a\xa8\x7e\x80\x9d\xab\xff\x39\x30\xab\x89\xa7\x7e\xfb\x4e\xd9\xbf\x04\x78\xf0\xfd\x36\x83\xb3\xa7\x4e\x28\x70\xe7\xf7\x16\x03\xa9\x0a\xd1\x40\x89\x0f\xb8\xf4\x0c\x31\x10\x14\xf6\x18\x95\x94\x06\xef\xb8\x44\xab\x95\xbf\xed\x05\x2c\x1f\x58\xbf\x03\x0a\x1c\x91\x29\x8d\x47\xd3\x3d\xa7\x73\xfe\x62\x0d\x22\x6e\x29\x17\xc5\x67\xe8\x14\x6e\xd4\xf1\xe9\xb9\x8a\xda\xdc\xa0\x9f\x8b\xf7\x07\x09\x54\x65\xe7\x95\x82\x64\x21\x39\xa2\x82\x68\xab\xd6\xb6\x34\x93\x76\x08\xeb\x0b\xb5\x37\xe8\xb3\x42\x9f\x51\xa9\xb3\xb0\x40\xa1\x64\x85\x38\xb1\x67\xf5\x2f\xd0\x4a\xb5\xf1\xe2\xaf\x05\xa6\x45\x89\x04\x14\x94\xf9\xb3\x99\x80\x6b\xf7\xc0\x3e\x9a\xff\x66\xf2\xba\xf9\x1b\x18\xed\xf2\x7f\xd9\x9a\xf5\xb4\xfb\x0a\xa8\xc4\x07\x0c\x5c\xdb\x8b\x8c\x2f\xd1\xd0\xf3\x50\x18\xa6\x62\xe9\xcd\x1a\xfb\xa9\xff\x93\xfc\x45\x92\xef\xa8\xfd\xdf\xce\x97\x65\xfe\x3e\x25\xf3\xbf\x9c\x06\xe8\x14\xc0\x0b\xb7\x29\x99\x8c\x34\x00\xd3\xa6\x95\xd3\x80\xde\xeb\xee\x1c\x54\xff\x40\x50\x67\x60\x3a\x03\xd2\x19\x88\xce\x00\x74\x06\x9e\x33\xe0\x9c\x81\xe6\x0c\x30\x67\x60\x39\x03\xca\x19\x48\xce\x00\x72\x06\x8e\x33\x60\xfc\x3d\x28\x76\x85\x76\x8b\x66\x1f\xc9\xeb\x3f\x28\xfe\xd2\x42\x4f\x4b\x75\x8a\xa1\xa2\x0b\x36\x41\xa1\xd3\x01\xce\xf5\xae\x13\x14\x36\x08\x25\xfd\x5b\xde\x04\x2d\xad\x42\x4b\x62\xc2\x15\x5f\xb7\x9e\x1e\xb6\xa4\x32\x86\xa8\xac\xd0\x46\xdd\x88\x9c\x16\x8f\xe0\xc5\xdc\x13\x16\x91\xd5\x4b\x2f\x04\xeb\xf7\x8d\xc9\x3c\x0b\xf7\x10\x50\x5c\xe9\x77\x23\x17\xb4\xfd\x37\xa4\x9c\x11\xd1\x0b\xd3\xe4\x6b\x72\xb0\x23\x3e\xc5\xff\xb1\x47\xf3\x95\x19\x43\xc0\xd9\xe7\x8e\x45\x98\xd3\xfb\xde\xd3\xc9\x3e\x94\x5e\x70\x6b\xf6\x3f\x7d\x18\xc5\x4c\xd6\x53\x45\x3b\xf5\x1f\x90\xf6\xcc\xdd\x7a\xb7\x21\x06\xad\xd9\xe0\x42\xb4\x78\xf5\xca\x95\xa9\x5e\xbd\xca\x07\x5a\xc6\x49\xc5\x92\x73\x47\xa5\xbd\x1a\xbc\x48\x45\x77\xf5\xfd\xf9\xbd\x1d\x3e\x16\x08\x56\x7c\x3e\xf1\x45\xfc\x6a\x2c\xea\xce\x83\x8d\xbe\x06\x8c\xeb\x1c\xa3\xe9\xaf\x6c\x20\x8d\x2d\xb8\xa8\xaa\x3b\xa6\x54\x6c\x46\x8f\x1f\xfb\xdd\xea\xf5\xf0\x44\xf1\x3a\xbd\x79\x3f\x65\x9a\xbf\x03\x02\x12\xfc\xf7\xd3\xe5\x4f\x58\x97\xf5\xed\x25\xd4\xfd\x97\x9f\xda\xde\xef\x52\xfa\x47\x53\x21\xa3\xf4\x77\xd9\x65\xe8\xf5\xb6\x67\x61\x8a\xf5\x45\x55\x89\x5f\x8c\x7e\xa1\x6f\x8e\xab\x84\xb8\xd7\x05\xeb\x87\x4a\xfd\x49\xb4\xd4\xdf\x88\xac\x13\x08\x30\x2a\x42\xa7\xb3\x41\x83\x9b\x5b\x8a\x4d\xa3\xcf\x9f\xc6\xb5\x9a\x14\xf0\xcc\xb4\x51\x4d\x62\x13\x7f\x0b\xf4\x3a\xbe\x30\xa3\xd9\x92\x55\x8d\xd4\x84\x17\x4d\x43\xba\xa1\x04\xf9\x9e\xbb\x9a\xe2\xa6\xcc\x3a\xd9\x84\xca\x1b\xa8\xd9\x13\xf4\x7a\x8e\xf6\x8f\x3c\xa4\x4c\xea\x1b\x6c\xdc\xf8\x8f\xb1\x59\xfa\x42\x54\x30\x0a\xbf\xd4\x2d\x91\xb8\x21\xd0\x7f\x3f\xe6\xf0\x01\x62\x1a\x0d\x51\x30\x06\xb1\x19\x08\x8b\x69\x18\x55\x4c\x52\xdc\x6a\x75\x3f\xfa\x64\x34\x15\xe7\xef\x53\x88\xfe\x6e\xdb\x43\x77\x2f\x36\x22\xfc\x3d\x69\x58\x38\x2e\x7f\x09\x7f\x9e\xfd\x3b\x00\x00\xff\xff\x76\x09\xbb\x9d\x8d\x34\x00\x00")

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
