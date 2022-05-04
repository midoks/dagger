// Code generated for package conf by go-bindata DO NOT EDIT. (@generated)
// sources:
// ../../../conf/app.conf
package conf

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

func bindataRead(data, name string) ([]byte, error) {
	gz, err := gzip.NewReader(strings.NewReader(data))
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

// Mode return file modify time
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

var _confAppConf = "\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x5c\x90\x4d\x4e\xc3\x30\x10\x85\xf7\x73\x0a\x6b\xba\x45\xe9\x4f\x5a\x40\x42\x5d\x70\x06\xd8\x55\x51\xe5\x90\x69\x1c\x48\x3c\xae\x3d\xa5\x2d\x7b\x96\x5d\x21\xae\x51\xb1\x46\xe2\x38\x14\x71\x0b\xe4\xa4\x15\x88\xa5\x3f\x6b\xbe\x37\x6f\xb4\x73\x56\x37\xa4\xa6\x0a\x0b\x5d\x96\xe4\x11\xfc\xca\x36\x5c\x74\x88\x1e\x11\x00\x60\x56\x73\x99\xc1\x95\xba\x35\xa4\x6a\x2e\xd5\x82\x7d\xa3\x45\x51\x25\x86\xbc\xc2\xfb\xc0\x16\x15\x7b\x85\x42\x1b\x41\x38\x7e\x4f\x4f\x6f\xcf\x2c\x73\xa7\xc5\x44\x54\x73\x19\x10\xa0\xa7\x0e\xaf\x6f\x87\xdd\xfe\xf3\xfd\xe5\xfb\x79\xf7\xf5\xb1\x87\x59\x91\x67\x20\x5b\xd7\x26\x87\x65\x5d\x09\xa5\xa8\x7a\xaa\xd9\x86\x65\x7d\x76\x04\x60\x38\xb4\xe6\xe1\xe8\x22\x19\x24\x83\x64\x88\xb0\x0a\xe4\x23\x8a\x31\x08\x4e\x87\xb0\x66\x5f\xfc\x21\xec\xe3\x48\x9a\x4e\xce\xe1\x7f\xd9\xde\x9d\xd1\x3e\x50\xab\x5c\xc9\xe2\xb2\xc9\xc7\x08\xce\xd3\xa2\xda\x44\x54\x35\x73\x04\xa9\x1a\x7a\x62\xdb\xce\x5d\x87\x4a\xf7\x6f\x8c\xb6\xa5\xd1\x55\x0c\xeb\x4a\x15\x5a\x74\xbf\x73\x26\x45\x9e\x22\xc0\x2c\x6e\x95\x01\x59\x9d\xd7\x34\x55\x38\xe8\x0e\x69\x44\x5c\x76\xda\x08\x87\xa3\x74\x3c\xf9\xb5\x58\x92\x35\xfb\x07\xfc\x09\x00\x00\xff\xff\x07\x79\xd7\xc1\x95\x01\x00\x00"

func confAppConfBytes() ([]byte, error) {
	return bindataRead(
		_confAppConf,
		"conf/app.conf",
	)
}

func confAppConf() (*asset, error) {
	bytes, err := confAppConfBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "conf/app.conf", size: 0, mode: os.FileMode(0), modTime: time.Unix(0, 0)}
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
	"conf/app.conf": confAppConf,
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
	"conf": &bintree{nil, map[string]*bintree{
		"app.conf": &bintree{confAppConf, map[string]*bintree{}},
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
