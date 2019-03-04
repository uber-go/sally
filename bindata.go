// Code generated by go-bindata. DO NOT EDIT.
// sources:
// templates/index.html
// templates/package.html
package main

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

var _templatesIndexHtml = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xb4\x54\xc1\x6e\xdb\x30\x0c\xbd\xef\x2b\x38\x21\xbd\xad\xe2\x30\xec\x54\x28\xd9\x61\x19\x7a\x19\xb0\xa0\xdd\x0e\x3b\x2a\x12\x13\xab\x91\xa5\x40\xa2\xdd\x05\x85\xff\x7d\x88\x9d\x36\xad\x33\xbb\x01\x86\xfa\x62\xd3\x7c\xef\x91\x22\x29\xaa\xf7\xf3\x1f\x5f\x7f\xfe\x5e\x7c\x83\x82\x4b\x3f\x7b\xa7\xba\x17\x00\x80\x2a\x48\xdb\xee\xb3\x35\xbd\x0b\x1b\x48\xe4\xa7\x22\xf3\xce\x53\x2e\x88\x58\x40\x91\x68\x35\x15\x05\xf3\x36\x5f\x21\x1a\x1b\xee\xb2\x34\x3e\x56\x76\xe5\x75\x22\x69\x62\x89\xfa\x4e\xff\x41\xef\x96\x19\xf3\x86\x3c\x71\x0c\xf8\x49\x7e\x94\x9f\x9f\x4c\x59\xba\x20\x4d\xce\x02\xf0\x10\x1a\x8f\xb1\xd5\x32\xda\xdd\xb3\x34\xac\xab\xc1\x78\x9d\xf3\x54\x98\x18\x58\xbb\x40\x49\x1c\xfd\x7d\x4c\x8a\xf7\x3d\x6f\x8b\x60\xbd\xf4\xf4\x88\xa9\x2e\x57\x95\xf7\x97\xf7\xce\x72\xf1\x0f\x70\x47\x78\x59\x8d\x53\x7f\x1a\x76\x1e\x04\x66\x0b\x6d\x36\x7a\x4d\x0a\xb9\x78\x1d\x7c\x1b\xab\x64\xce\xc4\xce\xa3\xa9\x4a\x0a\xac\xd9\xc5\x30\x4e\x51\x38\x94\xe9\x9e\x37\x78\x46\xc5\x2f\xdb\xd0\x7f\x1e\x1e\x20\xe9\xb0\x26\x98\x6c\x68\xf7\x01\x26\xb5\xf6\x15\xc1\xd5\x14\xe4\xe1\xd0\x19\x9a\x66\x8c\x3d\x71\xe5\x36\x26\x5e\x68\x2e\xf6\xb4\x6d\x72\x81\x57\x20\x2e\x6a\xbc\xa8\x05\x4c\xe4\xaf\x9b\xef\xad\xf6\x98\xcc\x19\x4d\xb0\xb3\x5e\xac\xa6\x51\xc8\x23\x8d\x7d\xa4\x8d\x02\x5a\x90\x3e\xdc\x05\xc4\x7d\x88\xb6\x00\xf2\x86\xb6\x11\x9a\x46\xcc\x4e\x7e\x29\xd4\xaf\x04\x7d\x93\xb4\xe4\x75\xb4\xd1\xdc\x52\xaa\x29\x41\xd3\xe0\x49\x35\x06\xe6\xff\x44\xd6\x95\x6b\xc8\xc9\x9c\xad\xfb\x25\xb3\xe6\x2a\xcb\x5c\xaf\x05\x68\xcf\x53\x71\x1d\xe7\xd1\x3c\x5d\xf9\xf1\x52\xfc\x57\xad\x86\x47\x1e\xba\xd9\xa3\x60\x87\xc6\x4a\xe1\xc0\xe0\x2b\x6c\x57\x48\x6f\xf1\xa0\x75\xf5\xb3\x5d\x75\x34\x15\x76\x32\x0a\xbb\x05\xfb\x37\x00\x00\xff\xff\x64\x04\xa6\x05\x78\x05\x00\x00")

func templatesIndexHtmlBytes() ([]byte, error) {
	return bindataRead(
		_templatesIndexHtml,
		"templates/index.html",
	)
}

func templatesIndexHtml() (*asset, error) {
	bytes, err := templatesIndexHtmlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "templates/index.html", size: 1400, mode: os.FileMode(420), modTime: time.Unix(1551736378, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _templatesPackageHtml = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x8c\xd1\x3f\x4f\xf3\x30\x10\x06\xf0\xfd\xfd\x14\xf7\x86\xb9\x35\x33\x38\x5d\x0a\x62\xa9\xa0\xaa\x60\x60\x34\xc9\xd3\xd8\x92\xed\x0b\xf6\xa5\x12\x8a\xf2\xdd\x51\x70\xa5\x82\x28\x82\xc9\x7f\xf5\xbb\x47\x77\xfa\xff\xcd\xc3\xfa\xf1\x79\x7b\x4b\x56\x82\x5f\xfd\xd3\x65\x21\x22\xd2\x16\xa6\x2d\xdb\x8f\x63\x80\x18\x8a\x26\xa0\xae\x3a\x5e\xb8\xd0\x73\x92\x8a\x1a\x8e\x82\x28\x75\x35\x8e\xb4\x5c\x9b\xc8\xd1\x35\xc6\x3f\xed\x36\x34\x4d\xd4\x39\x21\x2b\xd2\xe7\x2b\xa5\xe6\xf7\x1d\x7a\xa6\x69\xaa\x7e\x54\x33\x0f\xa9\xc1\x2f\xea\x19\xf1\xdc\x9d\x92\x04\xa8\x60\xb2\x20\x8d\xaa\x75\xe9\x8f\xbf\xd4\xb8\x77\x1e\xd3\xc5\x66\xf4\x2e\xe2\x7b\xd8\x19\x59\xe0\x75\x70\x87\xba\x4a\xd8\x27\x64\xfb\x29\xf0\xe5\x35\x0d\xc9\xd7\x73\x85\x3b\x6e\xb9\x29\x99\x8f\x88\x56\xa7\x9e\xea\x17\x6e\xdf\x4e\xf6\x3d\x8b\x75\xb1\x23\x61\xca\x00\x59\x24\x2c\x69\xeb\x61\x32\x48\x1b\xb2\x09\xfb\xd2\x8d\x2f\x6a\xe0\x03\xc8\x78\x8e\x9d\x56\x66\xb5\x3c\x16\x29\xb2\x56\x65\x96\xef\x01\x00\x00\xff\xff\x69\xc7\x6b\xa3\xe3\x01\x00\x00")

func templatesPackageHtmlBytes() ([]byte, error) {
	return bindataRead(
		_templatesPackageHtml,
		"templates/package.html",
	)
}

func templatesPackageHtml() (*asset, error) {
	bytes, err := templatesPackageHtmlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "templates/package.html", size: 483, mode: os.FileMode(420), modTime: time.Unix(1546544917, 0)}
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
	"templates/index.html": templatesIndexHtml,
	"templates/package.html": templatesPackageHtml,
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
	"templates": &bintree{nil, map[string]*bintree{
		"index.html": &bintree{templatesIndexHtml, map[string]*bintree{}},
		"package.html": &bintree{templatesPackageHtml, map[string]*bintree{}},
	}},
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

