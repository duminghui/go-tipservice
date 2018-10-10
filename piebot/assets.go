// Code generated by go-bindata. DO NOT EDIT.
// sources:
// templates/cmdmsgs.tmpl (3.559kB)
// templates/cmdusages.tmpl (2.032kB)

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

var _templatesCmdmsgsTmpl = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xa4\x56\xc1\x6e\xe3\x36\x10\xbd\xeb\x2b\xa6\xee\x21\xad\x91\xf8\x03\x84\xa2\x85\x37\x4d\xd1\x05\xea\x74\x91\x6c\x5a\xf4\x26\xda\x1c\xdb\x53\x50\xa4\x4a\x52\x71\x02\x82\xfd\xf6\x05\x49\x49\xa1\x13\x4b\xc9\x66\x91\x8b\x43\x3d\x72\xde\xbc\x37\x33\xa4\x73\x1c\xb7\x24\x11\x66\x07\x26\x04\xda\x15\x23\x69\x51\x32\xb9\xc1\x19\x5c\x78\x5f\x38\xb7\xf0\x1e\xfe\x8e\x1f\xa1\x7e\xfa\x5a\x38\x77\x01\x28\xb9\xf7\x45\x91\x1d\x42\x76\xcf\x35\x3b\xac\x48\x2e\x6b\xd5\x4a\x7b\xa5\xf5\x70\xce\x9d\x41\xbd\x42\x69\x49\x49\xef\xa1\x87\x42\x4d\x92\xea\xb6\x06\x16\x37\x00\x19\xa8\x9c\x5b\xac\x28\x80\x9c\x5b\xdc\x3e\xd6\x6b\x25\xbc\xaf\xa6\x23\xfe\xc5\x04\x71\x66\x71\xc9\xb9\x1e\x0f\x1a\x4e\x0e\x08\xef\xab\x10\x48\x2a\x9b\x87\x38\x33\xc0\x38\xd7\x68\xcc\x74\xac\x0f\xca\x7e\x55\x18\x92\xb0\x56\xf6\xcc\x40\xd2\xf8\x1c\x1e\x55\x0b\x1b\x26\xa1\x35\x18\xb1\x9f\x34\x6e\xe9\xc1\xfb\x86\xb0\x82\x8d\xaa\x6b\x26\x39\x58\x05\x3b\xba\x47\x30\xaa\x46\x25\x31\x27\x3a\x4d\x2f\x29\x7f\xad\xec\x95\x54\xed\x6e\x3f\x4e\x33\xd0\xe0\x4a\x9e\x59\xd8\xb3\x7b\x04\x94\xaa\xb5\xfb\xde\x07\xab\x06\x8b\x8a\xaa\xaa\x3e\x30\x11\x6c\xef\xbe\x96\x21\xc1\xf8\xeb\xd8\xa4\xe2\x4e\x6e\x94\xdc\x92\xae\x91\x67\xd0\x6c\xf5\xd4\xae\xaa\x7a\xc5\xdc\xdb\x76\xb3\x41\x63\x26\xd2\x18\xaa\x69\x84\x58\x48\x27\x73\xa5\xb0\x0f\x5b\xc4\x32\x20\x3e\x3f\xfc\x86\xf8\x2c\x8b\xb8\x7c\xf5\xd0\xdc\x69\xe1\x7d\xfc\xe7\xe3\xaf\x63\xaa\xaf\x99\x48\x01\x47\xd8\xfd\xa3\x5a\x0d\xeb\xa4\x5e\x19\x94\xbc\xec\xa5\x28\xc7\xc8\xe6\x2a\x46\xd0\x7b\xf5\xe3\xd8\x28\x43\xf6\xa3\xdc\xaa\x29\x76\xc7\x2d\xd0\x6d\xea\x5b\x01\xc8\x94\x45\x2e\xdd\xc9\x48\x0d\x61\x62\xb6\x22\x39\x5e\x70\xab\xe3\x56\x3f\xdd\xe7\xc0\x84\x50\x07\x8c\x0d\xb0\x46\xe0\x64\xac\xa6\x75\x6b\x47\x46\x4e\x43\x38\x94\xfa\x6b\x33\xe7\x44\xc1\xef\xf2\x82\xcf\x62\x35\x02\x99\xc1\xd8\xa1\x7d\x43\xe6\x9d\xda\x89\x54\xc5\x2e\x45\xfb\x5c\xb4\x13\x1d\xf3\x16\xb3\x73\xec\x7b\x3d\x8f\x7a\x7c\x42\xd5\x08\xcc\x85\xf0\x1e\xae\x15\x34\x71\xfd\x85\xb2\xd0\x10\x9e\xc3\x67\xfd\x08\x6c\xc7\x48\xc2\x61\x8f\xb2\xc7\x32\x8d\xa0\xa4\x20\xf9\x06\xfd\xaf\xd8\x66\x62\xdc\x04\xf6\x28\x4f\x66\xd3\x0f\xe3\xce\x91\x23\x2b\x02\xee\x06\x37\x48\xf7\xa8\x2f\xbb\xad\x89\xdb\xc8\x98\x6e\x08\x43\x9c\x88\xfd\x9d\xa6\x1a\x33\xe4\x3d\x35\x31\xde\x13\x39\x1f\x55\xa5\x50\x42\x50\xa3\x9a\x12\xfe\x7f\xf6\xe7\xdc\xe2\x52\x91\xbc\x66\x75\x18\x3e\x81\xc7\x33\x44\xb9\x61\x92\x3f\x96\xc5\x40\x2f\x88\xfb\x82\xa2\x73\x9a\xc9\x1d\xc2\xc0\xd3\x24\x48\x18\x5b\x91\x9e\x73\xb4\x85\xc5\xed\x5e\x1d\x96\x42\xa4\xb2\xf0\x1e\x7e\x60\x42\x4c\x64\xf7\x63\xb7\xf9\x74\x96\x9b\x3d\x93\x12\xc5\x9f\x0d\x6a\x66\x95\x9e\x18\xce\xdd\x12\x6d\x01\xff\x83\x45\xbf\x01\x66\x8c\xf3\x1e\x3e\x83\x99\xf7\x4b\xce\x91\x83\x49\x07\x6d\x5b\x21\x1e\x53\x64\x61\x30\x87\xdd\x60\xad\xc2\x85\xf8\x02\x17\x19\x56\x79\x77\x9c\x99\xbe\x6d\x43\x61\x1d\x80\x6d\x6c\xb8\x4b\x49\x82\xdd\x63\xe8\xe9\x94\x83\x29\xe3\x01\x9d\x86\x97\xdd\x22\x78\x0f\x3f\x7d\x1f\x55\xfc\xd9\xb9\xc0\xc2\xfb\x30\x95\x86\x5d\x93\xf2\x34\x84\xe1\x15\xb5\x62\x92\xed\x50\xff\x41\xa6\x2b\xc0\xf9\x7c\xde\xad\x95\xf3\xf9\xbc\x00\x98\xcf\x83\x54\x65\xfc\x0d\x70\x01\xce\xfd\xab\x48\xc2\xa2\x43\x19\x98\x9d\x87\xac\x23\xf2\x46\x09\x1c\x45\x86\x8f\x03\xfa\x0d\x9c\x26\xee\x83\xa0\x95\x41\x7d\x8f\x1a\xea\x8e\x6c\xe1\x9c\xc5\xba\x11\xcc\x8e\x24\xb7\x78\x25\xea\x44\xb8\xdb\x14\x2a\x5e\x89\xbb\x63\x2f\x14\xc9\xb4\x6c\x92\x76\x99\xbb\xbd\x7c\x69\x18\x97\xf3\x39\x64\x93\x39\x7e\x59\x26\xbb\x07\x97\x23\x24\x34\x42\x6f\x71\x68\x8c\xa4\xe1\x60\x7a\x2c\xb1\xde\xee\xe5\x93\xdb\x2f\xcd\xfe\x56\x41\x12\xd5\x67\xd3\xb9\x89\x8b\x50\xfd\x12\x9f\x8c\x1a\xa3\x0b\x1c\xb6\x4a\x87\x92\x0d\x0f\xc8\x37\x9c\x99\x14\x0a\xf3\xf8\x81\xcc\xc4\x45\x98\x5d\x82\x76\x4f\xa1\x57\x48\x9e\x19\x30\x71\x3b\x54\xaf\x3f\xbb\x8f\x33\x99\x0e\xd6\xdf\x9f\x7d\x86\x99\x59\xe1\xb5\x1b\x3c\x6e\x35\x8b\xf3\xc2\xaa\x98\xec\x11\x01\xd8\x32\x12\xc8\xbf\x2b\x8e\x37\xee\x99\x81\x35\xa2\x1c\x4e\xe8\xc4\xfa\x5a\xf2\xd3\x6f\xcb\x6f\xe5\x3e\x32\xab\x5e\x92\x09\x95\x33\xd1\x28\x43\x92\x6c\x2d\xf0\xd8\xae\xb2\xb8\x78\x9a\x08\x29\x72\x36\x0c\x52\xbc\x2f\x01\x00\x00\xff\xff\xed\xbd\x10\xbb\xe7\x0d\x00\x00")

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

	info := bindataFileInfo{name: "templates/cmdmsgs.tmpl", size: 3559, mode: os.FileMode(420), modTime: time.Unix(1539146459, 0)}
	a := &asset{bytes: bytes, info: info, digest: [32]uint8{0xd1, 0x38, 0xcb, 0x9e, 0x5, 0x5f, 0xc3, 0xbd, 0x6c, 0x3d, 0x19, 0x8d, 0xa5, 0x69, 0xfb, 0xbb, 0x95, 0x31, 0xb2, 0xd8, 0x66, 0x49, 0xf5, 0x5f, 0x86, 0xeb, 0x91, 0x94, 0x51, 0x62, 0xbb, 0x95}}
	return a, nil
}

var _templatesCmdusagesTmpl = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xac\x94\x6d\x6b\xdb\x30\x10\xc7\xdf\xe7\x53\x1c\x19\xa5\x10\xa8\x3e\x40\x09\x59\xd9\x60\xb4\xb0\x8c\x42\x57\xc6\xd8\x06\x51\xec\x73\x22\xd0\x83\x91\xe4\xa4\x45\xd5\x77\x1f\x7a\x70\x62\x3b\x4e\xb7\x17\x79\x67\x7c\xba\xdf\xfd\xef\xd1\xb9\x12\x2b\x26\x11\xa6\x82\x32\x39\xf5\xde\x39\x94\xa5\xf7\x93\xa3\xa1\x31\x74\x83\x53\xb8\x89\x3f\xc9\xb3\x41\xbd\x44\x69\x99\x92\xde\x83\x73\xe4\xb3\x28\xbf\x51\x81\xde\x43\xa1\x84\xa0\xb2\x84\xe8\x70\xeb\xdc\xf4\x77\x00\x4e\x9c\xbb\x81\xc4\xec\x40\x6b\x86\xcf\x5d\x2e\xab\x80\x3c\x98\xa7\xad\xda\xc7\xdf\xf7\x4c\xda\xa0\xc5\xa2\xa8\x39\xb5\x47\x15\x04\xb2\xc4\xe8\x37\x9b\xcd\x9c\x23\x8f\x1a\x2b\xf6\xe2\x7d\xcd\x10\x7e\xdd\x69\x2c\x90\xed\x50\x13\x42\xfe\xc0\x9c\x0a\xd5\x48\xbb\x98\xcd\x66\x13\x80\xf0\xc0\x39\xf2\xf4\x2a\xd6\x8a\x7b\x0f\x56\x41\xfb\xda\x4c\x00\x04\x93\x4c\x34\x02\x92\xcf\x6d\x20\x33\x5c\xb2\x9c\x67\xeb\x35\x9e\xcf\x9a\xf2\x4e\x3e\x7d\x5d\x6b\xca\x53\xfc\x0d\x5a\x58\x53\x4e\x65\x81\x39\xc8\x38\xac\xc4\x5a\x19\x66\xcf\x02\xb3\xfd\x08\xcd\x3f\x80\x96\xa5\x46\x63\xc6\xa9\x7b\x66\xb7\xa5\xa6\xfb\x8b\xd7\xbd\x05\xc3\x3c\xc7\x5f\x0c\xea\x7e\x5a\xd8\x1f\xd9\xe5\xb4\xba\x00\xf6\xa5\xc2\x30\x3e\xb5\x66\xd2\x56\x30\xbd\xda\x5c\x5d\x4d\x81\x7c\x7f\xf9\x82\xf8\x88\xba\xc0\x20\x11\x54\x75\x08\x02\x4a\x1f\x42\x38\x97\x1e\xfe\x67\xdb\x8a\x2d\x95\x12\xf9\xc5\x4b\x92\xb9\xb1\x22\x6f\x1a\x85\xda\xe1\x02\xe6\x1f\xf2\x6f\x42\x48\xae\x0c\x2d\xcb\xa0\x3e\xbd\x00\x5a\x58\xb6\x43\x68\x9d\x2b\xa5\x61\xd5\x49\x61\x05\x2b\xd3\x2e\xd9\x99\x1e\x6f\x91\xd7\xcf\xef\xec\xeb\xab\x6a\xa0\xa0\x12\x1a\x83\x60\xb7\x68\xf0\xc0\x3b\x8d\x36\xe9\xe6\x7c\xdc\x57\x12\xb9\x47\xcb\x71\xf2\x87\x96\xfe\x18\x0f\xad\x83\x71\x0c\xd5\x8c\x39\xa5\xf2\x2f\xa9\xa4\x1b\xd4\xa1\xc2\xe9\x8c\x74\x5d\xfb\x5d\x23\xdd\xa3\x35\x7e\x6a\x96\x94\xc9\x4c\xbc\x48\xab\x3f\x86\x4b\x22\x12\x70\xd0\xe3\xbb\xc6\xa0\x7e\xbb\xd3\x8a\xe3\x62\xa4\xc5\xad\x53\xa8\xf6\x23\xc3\x4f\xea\xcc\x09\xc8\xa2\xd3\x3c\x5d\x4e\x73\x1d\x79\x30\x37\xb1\xc9\x0b\x98\xa7\x1f\x49\x69\xa1\x64\xc5\x36\x07\xe3\xf5\xe1\xa0\xa7\x47\x13\xe8\x8d\xcf\x2a\x02\x39\x33\x76\x15\x0e\x69\x38\x43\x09\xd0\x68\xba\xe6\x08\x89\x72\x66\x4e\x73\x7a\x0f\xb2\x52\xbd\x2b\x17\x99\x4c\x56\x2a\x0a\x32\x5b\xb5\x07\x83\x7a\x87\xba\x15\x17\x6c\xef\x22\xbf\x32\x63\x4f\x91\x41\x66\x44\x86\x0f\xa0\x9c\xf7\xb5\x16\x8a\xc9\x6b\x13\x24\xff\x4b\xf1\xfd\xc9\x82\x79\x0f\x3f\xcf\x6e\xd5\xed\x70\x89\xfa\x59\x0f\x56\xe2\x24\x85\x71\x7b\x77\x28\xc6\x5f\xf4\x66\xbd\xb7\x17\x7f\x03\x00\x00\xff\xff\x51\xc0\x53\xed\xf0\x07\x00\x00")

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

	info := bindataFileInfo{name: "templates/cmdusages.tmpl", size: 2032, mode: os.FileMode(420), modTime: time.Unix(1539059061, 0)}
	a := &asset{bytes: bytes, info: info, digest: [32]uint8{0x95, 0xfc, 0x82, 0x1c, 0x68, 0x4d, 0xeb, 0x2d, 0x55, 0x8e, 0xd2, 0xab, 0xf7, 0x92, 0x32, 0x48, 0xff, 0x13, 0xab, 0xdd, 0x64, 0xe5, 0x66, 0x80, 0x84, 0x58, 0x45, 0x28, 0x5f, 0x98, 0xd6, 0x9d}}
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
