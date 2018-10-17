// Code generated by go-bindata. DO NOT EDIT.
// sources:
// templates/cmdmsgs.tmpl (5.322kB)
// templates/cmdusages.tmpl (4.371kB)

package main

import (
	"bytes"
	"compress/gzip"
	"crypto/sha256"
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
	bytes  []byte
	info   os.FileInfo
	digest [sha256.Size]byte
}

type bindataFileInfo struct {
	name    string
	size    int64
	mode    os.FileMode
	modTime time.Time
}

func (fi bindataFileInfo) Name() string {
	return fi.name
}
func (fi bindataFileInfo) Size() int64 {
	return fi.size
}
func (fi bindataFileInfo) Mode() os.FileMode {
	return fi.mode
}
func (fi bindataFileInfo) ModTime() time.Time {
	return fi.modTime
}
func (fi bindataFileInfo) IsDir() bool {
	return false
}
func (fi bindataFileInfo) Sys() interface{} {
	return nil
}

var _templatesCmdmsgsTmpl = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xac\x58\xdb\x6e\xe3\xbc\x11\xbe\xd7\x53\x4c\xb3\x17\x6e\x8d\xc4\x0f\x60\x2c\x16\x70\x0e\x45\x53\xac\xd3\x60\x93\xb4\xd8\xab\x35\x2d\x8d\x6d\x76\x29\x52\x25\xa9\x38\x86\xc0\x3e\xfb\x8f\x21\x29\x59\x8e\x2d\x59\xf9\x77\xb1\x37\x6b\x6a\x66\xbe\x99\x6f\x0e\x1c\xa6\xaa\x32\x5c\x71\x89\x70\xb1\x65\x42\xa0\x9d\x33\x2e\x2d\x4a\x26\x53\xbc\x80\x2b\xe7\x92\xaa\x9a\x38\x07\xff\xf1\x1f\x21\xdf\x7f\x4d\xaa\xea\x0a\x50\x66\xce\x25\x49\xcb\x08\xb7\x9b\x4c\xb3\xed\x9c\xcb\x59\xae\x4a\x69\xef\xb4\x6e\xec\xbc\x18\xd4\x73\x94\x96\x2b\xe9\x1c\xd4\xa2\x90\x73\xc9\xf3\x32\x07\xe6\x15\x80\x1b\x58\x54\xd5\x64\xce\x49\xa8\xaa\x26\x4f\xbb\x7c\xa9\x84\x73\x8b\x7e\xc4\x7f\x33\xc1\x33\x66\x71\x96\x65\xba\x1b\x94\x2c\x93\x84\x73\x0b\x02\x92\xca\xb6\x21\x46\x06\x58\x96\x69\x34\xa6\x1f\xeb\x5a\xd9\x0f\xc1\x70\x09\x4b\x65\x47\x06\x02\xc7\x97\xb0\x53\x25\xa4\x4c\x42\x69\xd0\xcb\x3e\x6a\x5c\xf1\x37\xe7\x0a\x8e\x0b\x48\x55\x9e\x33\x99\x81\x55\xb0\xe6\xaf\x08\x46\xe5\xa8\x24\xb6\x1d\xed\x77\x2f\x30\xff\xa0\xec\x9d\x54\xe5\x7a\xd3\xed\x26\xb9\x91\x29\x39\xb2\xb0\x61\xaf\x08\x28\x55\x69\x37\x75\x1e\xac\x6a\x52\xd4\x0f\xf7\x54\xa6\x29\x1a\xd3\x83\xd1\xa4\x9a\x68\xf1\xd6\x0f\x53\x4b\x58\x2d\xca\x12\xfb\xb6\x42\x9c\x92\xc4\xf3\xdb\xdf\x11\x0f\x85\x13\x7f\x7c\xf7\x56\xbc\x68\xe1\x9c\xff\x71\x7f\xeb\x5c\xf2\xa0\xb6\x04\xa6\x61\xc9\x04\x15\xe8\xf4\xb4\xd7\x4b\x26\x82\x0b\xed\xf2\xfe\x7e\x5e\x2f\xc3\x42\x19\x6e\xef\xe5\x4a\x75\x44\xea\x8d\x1c\x56\x53\x54\xaa\xab\x0a\xb8\x99\x26\xed\x40\x4f\x22\x15\x1c\x83\x87\x73\x2e\xbb\x73\x37\x3f\xec\x9a\xd3\x2d\x03\x4c\x08\xb5\x45\x5f\x4b\x4b\x84\x8c\x1b\xab\xf9\xb2\xb4\x1d\xdd\x5b\x70\x6c\xaa\xe6\x5c\xfb\x9e\xa8\x9d\x75\xbb\x76\x5a\x58\x85\x40\x66\xd0\x17\x7b\x5d\xdb\xed\xa2\x8f\x24\x2d\x7c\xc1\xa3\x7d\x4f\x5a\x8f\xab\x8f\xa8\x0a\x81\x6d\x1f\x9d\x83\x07\x05\x85\x3f\x3f\x0a\x1a\x0a\x8e\x97\xf0\xac\x77\xc0\xd6\x8c\x4b\xd8\x6e\x50\xd6\xb2\x4c\x23\x28\x29\xb8\x1c\x40\xcd\x1d\x4b\x7b\x9a\x8a\xf8\x47\x99\x9d\xac\xf3\x38\x72\x22\x59\x07\x2c\x91\xdc\x37\x4c\x91\xbf\xa2\xbe\x89\xaa\xc1\xb7\x6e\x06\x08\xc7\xcb\xfe\x83\xb7\xea\xf9\xa8\x2a\x29\xee\xbe\xd6\xfb\x33\xc8\xed\x9e\x9f\x0a\x25\x04\x2f\x54\x31\x85\xff\xbf\xfb\x57\x55\x93\x1b\xc5\xe5\x03\xcb\xa9\x8b\x1f\x39\xbe\x97\x98\xa6\x4c\x66\x3b\xea\xb9\xe8\x1e\x91\x7b\xe4\x62\x55\x69\x26\xd7\x08\x8d\x9f\x26\x88\x50\xff\x7b\xf7\xaa\x8a\xaf\x60\xf2\xb4\x51\xdb\x99\x10\xa1\x2c\x9c\x83\xbf\x32\x21\x7a\xa2\xfb\x5b\x54\x3e\x1d\x65\xba\x61\x52\xa2\xf8\x57\x81\x9a\x59\xa5\x7b\xa6\x5c\x3c\xe2\x2b\xc0\xff\xc1\xa4\x56\x80\x0b\x96\x65\xb5\xf8\x05\x5c\x38\x37\xcb\x32\xcc\xc0\x04\x43\xab\x52\x88\x5d\x40\x16\x06\xdb\x62\xdf\x30\x57\x34\xf6\x8f\xe4\xbc\x87\x8b\x76\x7f\x8f\x4c\xdd\x51\x54\x58\x5b\x60\xa9\xa5\x1b\x83\x4b\xb0\x1b\xa4\x76\x0b\x31\x98\x30\xd3\x22\x87\x37\xf1\x10\x9c\x83\xcf\x9f\x3c\x8b\x5f\xaa\x8a\xbc\x70\x8e\x06\x46\xa3\xd5\x4b\x4f\xc1\x91\x76\x85\x39\x93\x6c\x8d\xfa\x2b\x37\xb1\x00\xc7\xe3\x71\x3c\x9b\x8e\xc7\xe3\x04\x60\x3c\x26\xaa\xa6\xfe\xff\x00\x57\x50\x55\xff\x55\x5c\xc2\x24\x4a\x19\xb8\xb8\xa4\xa8\xbd\xe4\x37\x25\xb0\x53\x92\x3e\x36\xd2\x03\x7c\xea\x19\xd5\xc4\x95\x41\xfd\x8a\x1a\xf2\xe8\x6c\x52\x55\x16\xf3\x42\x30\xdb\x11\xdc\xe4\x0c\x6a\x0f\xdc\x53\x80\xba\x51\x72\xc5\xd7\x87\xb9\x50\x5c\x86\x63\x13\xb8\x6b\x65\xb7\xa6\x2f\xcc\xc9\xe9\x78\x0c\xad\xa1\xe9\xbf\xcc\x42\xba\x9b\x2c\x7b\x11\x6a\x84\x3a\xc5\xd4\x18\x81\xc3\x26\xe9\xbe\xc4\xea\x74\xcf\xf6\xd9\x3e\x4e\xf6\x30\x42\x8e\x64\xee\xde\x52\x51\x66\x38\x8c\xb4\x10\x4e\xf7\x30\x2d\xfc\x77\x7f\x5d\x5c\x2b\x5b\x07\xef\x97\x2a\x8d\x3e\x83\x19\xac\x94\xa6\x72\xa7\x15\x6b\x00\x56\x60\x97\x66\xf9\x1b\x37\x3d\xf7\x5b\xeb\x6e\xb3\x1b\x4e\x7d\xc6\xe5\xc8\x80\xf1\xea\xb0\x38\xbf\x98\x1e\x46\xd8\x0f\x56\x5f\x8b\xad\x70\xf7\xb1\xa6\xbe\x3e\x4a\xcd\xfc\xac\xb1\xca\x07\x7b\xe0\x00\xac\x18\x17\x98\xfd\x25\x39\x54\xdc\x30\x03\x4b\x44\xd9\x58\x88\x64\x7d\xd4\xf9\xfe\x05\xef\x57\x7d\xef\x98\x73\xc7\xce\x50\x45\xf5\x34\x59\x13\x24\x5b\x0a\x3c\x4c\xd7\x34\xb9\xda\x4f\x93\x80\x3c\x68\x90\x1c\x94\x72\x1c\x6e\xf1\x0c\xfc\x34\x8a\x23\x6e\x6f\x3c\x7e\x1d\x3c\xaa\xa2\xfc\xb0\x51\x85\x11\x5a\x7b\xe8\x5f\xea\xbd\x59\x69\xd5\x23\xd3\x2c\x37\x5f\x71\xbf\x64\xde\x1e\xed\x73\x4d\x72\xbd\xec\x7b\x48\xb2\xf2\x62\xd8\x1a\xcf\x83\xc5\xf9\xd3\x40\xfd\xb3\x34\x16\x24\x62\x46\x5b\xd7\x0e\xe8\x7d\xf3\xf9\x53\x1c\x45\x5f\x4e\xe1\x78\x6f\xa3\x95\x81\x98\x94\x85\x06\x70\x4e\x80\x3e\x34\x8f\x45\x1c\x76\xe3\x90\xe6\x40\x90\x27\xcb\x6c\x69\x1a\x98\xf0\x13\x72\x42\x5b\x22\x8c\xc2\x4e\x39\x02\xa5\x61\xc4\x84\x18\x75\x42\x06\xc5\x81\xa0\xf7\xd2\xa2\x7e\x65\xa2\x35\x51\x4e\x5b\xad\x05\x07\xda\x3d\xb7\xf4\x7f\x0e\xfb\xfd\x97\x26\x3e\x06\xb2\xcc\x97\xa8\x3b\x1d\x08\x16\x3f\x04\xdf\xf7\xec\xf9\x9d\x38\xbf\xe1\x81\xcc\xb2\x8c\x96\x5a\xb2\x09\x96\x99\x9f\xc9\x03\x55\xb4\x9f\x6f\x65\x7e\x6a\xdd\x5e\x5c\xfa\x85\x9c\xac\x75\xbc\x33\xeb\x7e\xd9\xa5\x02\x9f\x79\x8e\x67\x53\xdc\x48\x0e\x2d\x9e\xbb\x37\x8b\x32\x1b\x64\x7b\xb6\xb2\xa8\x49\x72\xa8\xe9\xda\x60\x4d\x0a\x11\x14\x2e\xa7\xb6\xde\x89\x2e\x7a\x77\xc5\xc4\x87\xcb\x7b\x7e\x0f\x2e\x8b\x7e\x02\xa3\xc5\xdb\xf9\x87\x6c\x26\xa9\xe0\xe9\x4f\xd0\x48\x9b\xb4\x92\x30\x65\x5a\xab\xad\xf9\x91\x52\x2a\x51\xa7\x42\xa5\x3f\xb7\xdc\xe0\x94\xb2\xaf\x71\xa5\xd1\x6c\xc2\x7e\xe0\x4d\x71\xb9\x52\x47\x26\x8c\x55\xc5\x8f\x65\x69\xad\x92\x5e\x8d\x7e\x03\x0d\x55\x1d\x16\x7d\xaf\xfe\x4c\xd5\xd3\x17\x0f\x09\x3c\xc4\xca\x0d\x2f\x84\xfd\xa5\xf1\xfd\xcc\x5b\x5c\x97\x12\x7c\xbc\xf4\x00\x23\x47\x69\xeb\xf7\x49\xe6\x78\x7f\xeb\xdc\x65\x2b\x84\xe0\x54\x06\xcb\x1d\x5c\x2b\x9b\x5c\xd7\x7f\x14\x09\xc8\xd3\xae\x77\x64\xf2\x22\xfd\x05\xac\x73\xcc\xa0\x25\xdb\x3a\x3e\xa9\xd6\x3b\xe5\x4c\x08\x34\xbe\xa2\x9e\xf7\x34\x9b\xda\xcd\x96\xfe\x1f\x01\x00\x00\xff\xff\x9c\x2e\x3e\xe5\xca\x14\x00\x00")

func templatesCmdmsgsTmplBytes() ([]byte, error) {
	return bindataRead(
		_templatesCmdmsgsTmpl,
		"templates/cmdmsgs.tmpl",
	)
}

func templatesCmdmsgsTmpl() (*asset, error) {
	bytes, err := templatesCmdmsgsTmplBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "templates/cmdmsgs.tmpl", size: 5322, mode: os.FileMode(420), modTime: time.Unix(1539663312, 0)}
	a := &asset{bytes: bytes, info: info, digest: [32]uint8{0x4e, 0x65, 0xaf, 0x65, 0x2e, 0xfb, 0xe7, 0x21, 0x9c, 0xdf, 0x75, 0xbd, 0x55, 0x36, 0x44, 0x57, 0xf7, 0x4e, 0x80, 0xa, 0x77, 0x11, 0x9, 0xd, 0xdb, 0xd9, 0xbe, 0x8, 0x18, 0x52, 0x72, 0x54}}
	return a, nil
}

var _templatesCmdusagesTmpl = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xbc\x57\x5b\x6f\xdb\x46\x13\x7d\xd7\xaf\x18\xf8\x43\xe0\xaf\x84\x45\x38\xc8\x4b\x21\x28\x42\x9c\x20\x45\x0c\x54\x8d\x91\xfb\xa5\x05\xb8\x22\x87\xd2\xa2\x7b\x61\x77\x97\xb2\x0d\x46\xff\xbd\x18\x72\x29\x71\x49\xd3\xf2\x83\x53\x3f\xd1\x7b\x99\x39\x33\x73\xe6\xcc\xaa\xaa\x32\xcc\xb9\x42\x38\x91\x8c\xab\x93\xdd\xae\xaa\x50\x65\xbb\xdd\x64\x72\xd8\x29\x2d\x5b\xe3\x09\x4c\x77\xbb\x49\x55\xc5\x1f\x2d\x9a\x25\x2a\xc7\xb5\xda\xed\xa0\xaa\xe2\x57\x32\xfb\x83\x49\xdc\xed\x20\xd5\x52\x32\x95\x41\x7d\x61\x56\x55\x27\x7f\x92\xc5\x49\x55\x4d\x61\x60\xb4\xe0\xf8\xb1\x6b\x97\xe7\x10\x5f\xda\xf7\x1b\x7d\x5d\x2f\xbf\xe1\xca\x11\x18\x87\xb2\x10\xcc\x1d\x50\xc4\xe0\x31\xd6\xf7\xa2\x28\xaa\xaa\xf8\xca\x60\xce\x6f\x76\xbb\x82\x23\x7c\x7f\x61\x30\x45\xbe\x45\x13\xc7\xf1\x5f\x30\x67\x52\x97\xca\x2d\xa2\x28\x9a\x00\x48\xae\xb8\x2c\x25\x34\x8b\x33\xba\xca\x71\xc9\x7d\x20\xef\x6f\xe5\x4a\x8b\x31\xc0\x2b\x26\x3a\x80\x43\xc7\x2b\x26\x1a\x07\x6b\x74\xb0\x62\x82\xa9\x14\xbd\x93\xbb\x8d\x65\x58\x68\xcb\xdd\xa8\x41\xbf\x7f\x30\xea\x17\x80\x65\x99\x41\x6b\xef\xb6\x7a\xcd\xdd\x26\x33\xec\xfa\xd1\x13\xdb\x1a\x86\xb9\xf7\xbf\x38\x9a\xd8\xcf\xfe\xca\x30\xbb\x00\xee\x26\x47\xe2\x47\x61\xb8\x72\x39\x9c\x3c\x59\x3f\x79\x72\x02\xf1\x87\x9b\xdf\x10\xaf\xd0\xa4\x48\x10\x41\xe7\x7b\x27\xa0\xcd\xde\x45\x55\x35\x07\x1f\x58\xb6\x74\xc3\x94\x42\xf1\xe8\x29\xf1\x76\xeb\x8c\xfc\x30\x28\xf5\x16\x17\x30\xff\x9f\x5f\x8e\xe3\xd8\x67\x86\x65\x19\xa1\x6f\x4e\x00\x4b\x1d\xdf\x22\xb4\x97\x73\x6d\x20\xe9\x84\x90\x9c\xda\xb6\x89\x46\x4a\xbc\x41\x51\x7c\xbc\xa7\x1f\x6f\x75\x09\x29\x53\x50\x5a\x04\xb7\x41\x8b\x7b\x7b\x43\x67\x93\x6e\xc8\x87\x7e\x8c\x6b\xbb\x87\x9d\x03\xf1\xfb\x3b\x21\x8b\xfb\xbb\x3d\x36\x52\x32\xf7\xb9\x5f\x32\xc5\xd6\x68\x7c\x0c\x87\x3b\x61\xb5\xe2\xbb\x3a\xfc\xa2\x74\x1a\xa6\x16\xe6\x5a\x09\xae\xf0\x07\x13\x62\x01\x53\x0e\x73\xae\x1c\x9a\x2d\xa3\xff\xd8\x81\x39\xd3\x14\xe6\xe9\x6d\x2a\xf0\x03\x97\x68\x17\x30\x45\x98\xb3\xdc\xa1\x71\x5c\x76\x2b\xb6\x80\xf9\x0b\xa3\x05\xfa\xb2\x51\xfe\x92\xa1\xe3\x04\xd6\xba\xee\x47\xa9\x0d\x02\x57\xb9\x36\x92\x51\xea\xbb\xe5\x6a\x83\xfc\xc4\x8b\xb7\xaa\x0e\x71\xda\xf9\x0b\xe2\xdd\xf2\x22\xc8\xde\xc1\xc4\xdd\x92\xb9\x64\x5c\xf9\xdc\x3d\x0a\xa3\xab\x2a\x7e\xa9\x5d\x57\x3f\xa5\xaf\x4c\xc8\xeb\x17\xa5\x45\xf3\xa3\x4d\xd0\x80\xd6\xed\x25\xa2\xd8\x15\xc7\x97\x7a\x44\xf6\x7c\x04\x8d\xbf\x9f\x14\x40\x51\x7f\xc3\xdc\xd6\x34\x5f\xc0\xbc\x59\x68\x60\xa7\x5a\xe5\x7c\xbd\xdf\x3c\xdd\x8f\xac\xe6\xd0\x04\x82\x06\x4a\x86\xd6\x05\xb7\x2e\x01\xd7\x90\xa0\xb1\x56\x1a\xb6\x12\x08\x8d\xc9\x91\xb6\xf5\x81\x5f\xaa\x5c\x07\x9a\x3f\x74\x40\x9c\xaa\xa1\xda\x8d\xbe\x06\x8b\x66\x8b\xa6\x85\x4d\x7b\xf7\xda\xff\x9d\x5b\x77\xc4\x3e\x05\x50\xdb\xa7\x0f\x60\x42\x84\x51\xa4\x9a\xab\x53\x4b\xc1\x1c\x8b\xe5\xf5\x4d\x2a\xca\xec\x71\x86\xf8\x10\x26\x36\xd6\x43\x1a\x7e\xaf\x09\x48\x73\xfd\x0e\x0e\x76\x6f\xd1\xb1\xfb\xd1\xbf\x39\xa2\xa3\x5f\x47\x75\x74\xd6\x97\xcd\x5e\x61\xfb\x32\x38\x2c\xcd\xc8\x89\xa0\x2b\x46\xce\x84\xbd\x3f\x72\x28\xac\x4c\x7c\x8f\x9c\x90\xa4\x5d\x31\xc3\xe4\x7b\xc7\x5c\x69\x03\xea\xf4\x24\x76\x16\x45\xf0\x56\x89\x5b\xea\x6f\xe2\x7f\xb3\x05\x05\xea\x42\xa0\xa5\x3a\x10\x97\xda\x7f\x13\x5d\x50\x22\x67\x19\xe6\xac\x14\xce\x9f\x4e\x8e\x03\xb9\xf4\x22\x1e\x42\xe9\xaa\x3b\x01\x99\x7f\xd9\x7c\x95\xdf\xec\xe2\xff\x5f\x60\xa3\x4b\x03\x5f\xe9\x75\x50\x3a\x84\x6f\x60\x31\xd5\x2a\x83\x5f\x80\xda\x47\x68\xb5\x06\x12\x7a\x42\x4c\xc0\xb5\x4a\xf1\x0c\x12\x59\x5a\x07\x0b\x78\xfa\xeb\xb9\x3d\xeb\x21\xa5\xb5\x64\x02\x00\xf0\xfa\x86\xc9\x42\xe0\x0c\x92\xa7\x9b\x04\xb8\x05\xad\xb0\x76\x77\x46\x2b\xcf\xce\xa5\x5f\x6c\x20\x90\x88\x3c\x3b\xf7\x38\xc6\xb9\xb7\x0f\xf4\xa2\x9e\x4f\x61\x98\x7e\x6c\xd1\xa3\x6a\x06\xc9\x92\x40\xae\x10\x18\xa8\x52\xae\xd0\xc0\xe2\x39\x8c\xbc\x58\x13\x80\x28\x8a\x12\x83\xff\x94\x86\x63\x96\xd0\x0c\x3b\x8a\xe0\xd5\x7e\x2a\x86\x28\xc2\x81\x49\xe9\x7e\xa3\xaf\x49\xe2\x6f\xeb\x54\x5a\xca\x65\x3d\x84\xaf\xf8\x3e\x97\x01\x4c\x38\x3f\x83\x7e\x56\xcf\x1f\x50\xfb\x0b\x1a\xcb\xe4\x35\xc4\x13\x0c\xec\x16\x4e\x50\xd9\x16\x0d\xd4\xe7\x20\x35\x48\xed\xe0\x98\xfd\xfb\x0c\x2c\x93\x08\xcc\x1e\x08\xd4\xcc\xfd\x55\xe9\x5a\x1a\x3c\xa7\xba\x35\xb4\xb1\x03\xe0\xcf\x88\x0d\xc7\x53\x39\x78\x64\x46\xd1\xe1\x65\x41\x98\x3f\x6f\x78\xba\x69\xdf\x7d\xa7\xd6\x37\x4a\x2d\x31\x34\x49\x8a\x3a\x95\x3d\xdf\x69\x69\x0c\x2a\xd7\xde\x4a\xea\x1a\x5f\xe6\x90\x73\x21\x80\x3b\xd8\x67\xde\x91\x44\x11\x42\x7a\x33\xd3\x67\x3b\xd5\xea\xa6\x7b\x10\x19\xde\x69\x11\xa6\xdd\xbf\x86\x1a\xec\x9a\x98\x4e\x86\xfd\x6a\x17\x77\x8d\xaa\x87\x9c\xa4\x80\x0e\x9e\xc1\xf2\x71\x30\xfe\x8c\x9f\x8a\xff\xcd\x43\xf2\x82\x86\x94\x42\x60\xe4\x8d\xd2\x45\xb4\x9c\x00\xf4\x74\xfb\x6e\x25\x8e\xeb\x9f\x4b\x63\x47\x7b\x5a\x79\xff\xe1\x40\x6f\xee\x3f\x3a\x10\x86\x23\x96\x7b\x6d\x7b\xc4\xf8\xe0\x85\x3f\x7e\xb6\x43\xca\xd1\x29\x76\x78\x41\x0f\x8b\xbc\xe5\x45\x53\x04\xe2\x0a\x8d\x72\x73\x6a\x9b\x47\xdd\xab\x77\x4b\x52\xcd\x4f\x97\x57\x50\x68\xae\x5c\x57\xb0\xff\x0d\x00\x00\xff\xff\xb1\xcc\x6c\xee\x13\x11\x00\x00")

func templatesCmdusagesTmplBytes() ([]byte, error) {
	return bindataRead(
		_templatesCmdusagesTmpl,
		"templates/cmdusages.tmpl",
	)
}

func templatesCmdusagesTmpl() (*asset, error) {
	bytes, err := templatesCmdusagesTmplBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "templates/cmdusages.tmpl", size: 4371, mode: os.FileMode(420), modTime: time.Unix(1539800149, 0)}
	a := &asset{bytes: bytes, info: info, digest: [32]uint8{0x9b, 0x4e, 0x9a, 0xff, 0x50, 0x15, 0xa5, 0x19, 0x7d, 0xc9, 0x62, 0x54, 0xe2, 0x5e, 0x56, 0xad, 0x87, 0x50, 0x3e, 0x45, 0xd8, 0x47, 0x44, 0x8f, 0xd, 0x1f, 0xe4, 0x7a, 0x1f, 0xbe, 0x32, 0x7d}}
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

// AssetString returns the asset contents as a string (instead of a []byte).
func AssetString(name string) (string, error) {
	data, err := Asset(name)
	return string(data), err
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

// MustAssetString is like AssetString but panics when Asset would return an
// error. It simplifies safe initialization of global variables.
func MustAssetString(name string) string {
	return string(MustAsset(name))
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

// AssetDigest returns the digest of the file with the given name. It returns an
// error if the asset could not be found or the digest could not be loaded.
func AssetDigest(name string) ([sha256.Size]byte, error) {
	canonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[canonicalName]; ok {
		a, err := f()
		if err != nil {
			return [sha256.Size]byte{}, fmt.Errorf("AssetDigest %s can't read by error: %v", name, err)
		}
		return a.digest, nil
	}
	return [sha256.Size]byte{}, fmt.Errorf("AssetDigest %s not found", name)
}

// Digests returns a map of all known files and their checksums.
func Digests() (map[string][sha256.Size]byte, error) {
	mp := make(map[string][sha256.Size]byte, len(_bindata))
	for name := range _bindata {
		a, err := _bindata[name]()
		if err != nil {
			return nil, err
		}
		mp[name] = a.digest
	}
	return mp, nil
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
	"templates/cmdmsgs.tmpl": templatesCmdmsgsTmpl,

	"templates/cmdusages.tmpl": templatesCmdusagesTmpl,
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
// then AssetDir("data") would return []string{"foo.txt", "img"},
// AssetDir("data/img") would return []string{"a.png", "b.png"},
// AssetDir("foo.txt") and AssetDir("notexist") would return an error, and
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
	"templates": &bintree{nil, map[string]*bintree{
		"cmdmsgs.tmpl":   &bintree{templatesCmdmsgsTmpl, map[string]*bintree{}},
		"cmdusages.tmpl": &bintree{templatesCmdusagesTmpl, map[string]*bintree{}},
	}},
}}

// RestoreAsset restores an asset under the given directory.
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
	return os.Chtimes(_filePath(dir, name), info.ModTime(), info.ModTime())
}

// RestoreAssets restores an asset under the given directory recursively.
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
