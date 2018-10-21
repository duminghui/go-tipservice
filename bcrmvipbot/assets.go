// Code generated by go-bindata. DO NOT EDIT.
// sources:
// templates/cmdusages.tmpl (1.996kB)

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

var _templatesCmdusagesTmpl = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xb4\x55\x51\x6b\xdb\x30\x10\x7e\xf7\xaf\x38\xb2\x87\x82\x69\x4d\x9e\x83\x6b\x36\xc2\xa0\x85\x75\x84\xa6\xdd\xd3\x1e\xaa\x44\x67\x47\x9b\x2d\x19\x49\x4e\x1b\x84\xfe\xfb\x90\x2c\x27\x76\x16\x27\x1d\xac\x4f\x6d\x74\xa7\xef\xbe\xef\xbe\xd3\xd9\x18\x8a\x39\xe3\x08\x93\x8a\x30\x3e\xb1\xd6\x18\xe4\xd4\xda\x28\x3a\x44\x1a\x45\x0a\x9c\xc0\x8d\xb5\x91\x31\xc9\xb3\x42\xf9\x80\x5c\x33\xc1\xad\x05\x63\x92\x79\x45\xbf\x93\x0a\xad\x85\xb5\xa8\x2a\xc2\x29\xf8\x0b\x33\x63\x26\x3f\x1d\x62\x64\xcc\x0d\xfc\x05\xba\x65\xf5\x1d\x96\xf5\xf3\x19\xec\x9d\x68\x60\x4d\x38\x34\x0a\x41\x6f\x50\x61\x57\x40\x41\x2e\x64\x64\x8c\xc6\xaa\x2e\x89\x6e\xd1\x02\x52\xe2\x91\x06\x91\x27\x31\x0c\xb2\x1c\x92\x7b\xf5\x40\x38\x29\x50\x86\xd2\x83\x0b\x0b\xc1\xb8\x56\xe3\x80\x8f\xa2\xc4\x0b\xe1\xf7\x40\x2c\xd7\x84\x8f\x67\xcc\x37\x84\x73\x2c\xcf\x60\x84\x8c\x4b\xa5\xbe\x56\xe2\x17\x1b\x84\x3b\x3b\xc6\x8c\xe9\x99\x12\xc7\xb1\x31\xc9\x42\x62\xce\xde\xac\xdd\xb2\x3a\x8e\xe3\x08\x60\xb9\x11\xaf\xce\x17\x79\xa5\xe0\xc5\x98\x64\xb9\xab\x56\xa2\xb4\xf6\x05\x7e\xdc\x2f\xa0\xf6\x8c\x46\xe1\xfb\xed\x3b\x59\xc3\x27\xf4\x0a\x39\x4c\xe9\xce\xce\x42\x0e\xfa\x70\xd3\x73\xda\x61\xf8\xe3\x3b\xc6\xb5\x9b\xf0\x43\x7b\xc2\x68\x27\x10\x06\x7f\x9c\x4f\x8b\x0e\xe9\x67\x47\x24\x83\xb4\xd5\x98\xb5\x2c\xe7\x82\xe7\xac\x00\xbd\x21\x1a\x1c\xe1\x8a\xf0\x5d\xe8\x02\x50\xd1\x9e\x77\x37\xb5\x80\xa6\x2e\x24\xa1\x18\xc1\x1e\x06\xaa\x46\x69\xc8\x6e\x61\x7a\x0d\x2c\x87\x5b\xff\x0f\xc5\x12\xb5\x1b\x7d\xa6\xbc\xfc\x2b\x75\xa9\xb3\x47\x33\x73\x52\x4c\x97\xd3\x32\x0f\x1d\x0e\x87\x57\x0a\x0a\xb6\xc5\x8e\xba\xe7\xfd\xba\x41\xee\x1e\x5e\xed\xc8\x50\xd0\x44\xfd\xbe\x54\xff\xe3\x9c\x18\x14\x80\xf4\xd3\xba\xfd\x3d\xe2\x87\x3a\xf6\x42\x8b\x56\xdf\x91\xa6\xf5\x5e\xbd\x57\xf7\x0f\xb6\x1c\x6e\x5e\x70\xe6\xe3\x5a\xb2\x1f\x4c\xf7\x1c\x8f\x1b\xf1\x85\x52\x10\x12\x54\xb3\xea\x5e\x6b\xef\x7d\xf6\x74\x66\xd3\x19\x10\x4a\x43\xe4\x1a\xd2\xe9\xcc\x5f\xea\x12\xae\xe1\x76\x3a\x63\x05\x17\x12\x47\x25\xf6\xf7\xcc\x7f\x55\xe8\x81\x21\x45\xf7\x27\xeb\xc6\x36\x58\xec\x0f\xdb\x41\x75\xce\xf2\xf7\xec\x9f\xc3\xf7\xe0\x64\xb9\x27\x51\x43\xca\x38\xc5\xb7\x6c\xf0\x46\x1c\xf2\x37\x24\x14\xe5\x4a\x10\x49\x7d\x20\xe4\xb5\x73\xb2\x42\x20\xc0\x9b\x6a\x85\x12\x32\x98\x9e\x5d\x56\xbd\xdd\x3f\xba\x72\x5c\x4e\xc7\xe0\x11\x95\xfb\x14\xd6\x28\xea\x12\x55\xf0\xd1\x25\xf5\xaa\xfc\x09\x00\x00\xff\xff\x7a\x34\xce\xe7\xcc\x07\x00\x00")

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

	info := bindataFileInfo{name: "templates/cmdusages.tmpl", size: 1996, mode: os.FileMode(420), modTime: time.Unix(1540134924, 0)}
	a := &asset{bytes: bytes, info: info, digest: [32]uint8{0x8c, 0xf8, 0x3d, 0x6a, 0xd3, 0xee, 0x97, 0x10, 0x51, 0xa8, 0x46, 0x1c, 0x5c, 0x34, 0x55, 0xcb, 0xf2, 0x7e, 0x96, 0xf5, 0x31, 0xf4, 0x5, 0x47, 0x89, 0x16, 0x63, 0x15, 0x16, 0x56, 0x10, 0xca}}
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
