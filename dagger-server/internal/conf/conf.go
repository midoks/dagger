package conf

import (
	"fmt"
	"log"
	"path/filepath"

	"github.com/pkg/errors"
	"gopkg.in/ini.v1"

	"github.com/midoks/dagger/dagger-server/internal/assets/conf"
	"github.com/midoks/dagger/dagger-server/internal/tools"
)

// Asset is a wrapper for getting conf assets.
func Asset(name string) ([]byte, error) {
	return conf.Asset(name)
}

// AssetDir is a wrapper for getting conf assets.
func AssetDir(name string) ([]string, error) {
	return conf.AssetDir(name)
}

// MustAsset is a wrapper for getting conf assets.
func MustAsset(name string) []byte {
	return conf.MustAsset(name)
}

// File is the configuration object.
var File *ini.File

func Init(customConf string) error {

	File, err := ini.LoadSources(ini.LoadOptions{
		IgnoreInlineComment: true,
	}, conf.MustAsset("conf/app.conf"))

	if err != nil {
		return fmt.Errorf("parse 'conf/app.conf' : %s", err)
	}

	File.NameMapper = ini.TitleUnderscore

	if customConf == "" {
		customConf = filepath.Join(CustomDir(), "conf", "app.conf")
	} else {
		customConf, err = filepath.Abs(customConf)
		if err != nil {
			return errors.Wrap(err, "get absolute path")
		}
	}
	CustomConf = customConf

	if tools.IsFile(customConf) {
		if err = File.Append(customConf); err != nil {
			return errors.Wrapf(err, "append %q", customConf)
		}
	} else {
		log.Println("Custom config ", customConf, " not found. Ignore this warning if you're running for the first time")
	}

	if err = File.Section(ini.DefaultSection).MapTo(&App); err != nil {
		return errors.Wrap(err, "mapping default section")
	}

	// ***************************
	// ----- Log settings -----
	// ***************************
	if err = File.Section("log").MapTo(&Log); err != nil {
		return errors.Wrap(err, "mapping [log] section")
	}

	// ****************************
	// ----- Http settings -----
	// ****************************

	if err = File.Section("http").MapTo(&Http); err != nil {
		return errors.Wrap(err, "mapping [http] section")
	}

	// ****************************
	// ----- Database settings -----
	// ****************************

	if err = File.Section("db").MapTo(&Database); err != nil {
		return errors.Wrap(err, "mapping [database] section")
	}

	// ****************************
	// ----- User settings -----
	// ****************************

	if err = File.Section("user").MapTo(&User); err != nil {
		return errors.Wrap(err, "mapping [user] section")
	}

	// Check run user when the install is locked.
	if Security.InstallLock {
		currentUser, match := CheckRunUser(App.RunUser)
		if !match {
			return fmt.Errorf("user configured to run imail is %q, but the current user is %q", App.RunUser, currentUser)
		}
	}

	return nil
}
