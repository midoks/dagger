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

var _confAppConf = "\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x5c\x90\x41\x4e\xeb\x30\x10\x86\xf7\x73\x0a\x6b\xba\x7d\x4a\xd3\xa6\xed\x43\x42\x5d\x70\x06\xd8\x55\x51\xe5\x90\x69\x1c\x48\x32\xae\x3d\xa1\x2d\x7b\x96\x5d\x21\xae\x51\xb1\x46\xe2\x38\x14\x71\x0b\xe4\xa4\x15\x88\xa5\x3f\x7b\xfe\x6f\x7e\x6b\x6b\x1b\x5d\x93\x9a\x2b\xcc\x75\x51\x90\x43\x70\x6d\x53\x73\xde\x23\x7a\x40\x00\x80\x45\xc5\x45\x0a\x97\xea\xc6\x90\xaa\xb8\x50\x2b\x76\xb5\x16\x45\xa5\x18\x72\x0a\xef\x3c\x37\xa8\xd8\x29\x14\xda\x0a\xc2\xe9\x7a\x7e\x3e\x3b\x66\x59\x5a\x2d\x26\xa0\x8a\x0b\x8f\x00\x03\x75\x7c\x79\x3d\xee\x0f\x1f\x6f\xcf\x5f\x4f\xfb\xcf\xf7\x03\x2c\xf2\x2c\x05\xd9\xd9\xce\xec\xd7\x55\x29\x94\xa0\x1a\xa8\x7a\xe7\xd7\xd5\xbf\x13\x00\xc3\xbe\x4b\x1e\x8d\xff\x47\x71\x14\x47\x23\x84\xd6\x93\x0b\x28\x68\x10\xac\xf6\x7e\xc3\x2e\xff\x45\xd8\x85\x91\x24\x99\xce\xe0\x6f\xd9\xc1\xad\xd1\xce\x53\x17\xd9\xca\xea\xa2\xce\x26\x08\xd6\xd1\xaa\xdc\x06\x54\xd6\x4b\x04\x29\x6b\x7a\xe4\xa6\x9b\xbb\xf2\xa5\x1e\x5e\x1b\xdd\x14\x46\x97\x41\xd6\x97\xca\xb5\xe8\x61\x9f\x19\xe5\x59\x82\x00\x8b\xb0\x55\x0a\xd4\xe8\xac\xa2\xb9\xc2\xb8\xff\x48\x23\x62\xd3\xf3\x46\x38\x1a\x27\x93\xe9\x4f\x4a\x43\xb2\x61\x77\xdf\x1b\xb9\x0d\x4f\x66\xf1\x77\x00\x00\x00\xff\xff\x2e\xb6\xba\x0d\xa2\x01\x00\x00"

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
