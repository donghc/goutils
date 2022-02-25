package config

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/donghc/goutils/pkg/file"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

func Load(path string, c interface{}) {
	dir := filepath.Dir(path)
	filename := filepath.Base(path)
	ext := filepath.Ext(path)

	viper.SetConfigType(strings.ReplaceAll(ext, ".", ""))
	//不包含扩展名
	viper.SetConfigName(filename)
	//路径
	viper.AddConfigPath(dir)
	_, ok := file.IsExists(path)
	if !ok {
		panic(errors.New(fmt.Sprintf("file not fount : %v", path)))
	}

	bys, err := ioutil.ReadFile(path)
	if os.IsNotExist(err) {
		panic(err)
	}

	var r io.Reader = bytes.NewReader(bys)
	if err := viper.ReadConfig(r); err != nil {
		panic(err)
	}
	if err := viper.Unmarshal(&c); err != nil {
		panic(err)
	}

	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		if err := viper.Unmarshal(&c); err != nil {
			panic(err)
		}
	})
}

func LoadAndCreate(path string, bys []byte, c interface{}) {
	dir := filepath.Dir(path)
	filename := filepath.Base(path)
	ext := filepath.Ext(path)

	viper.SetConfigType(strings.ReplaceAll(ext, ".", ""))
	//不包含扩展名
	viper.SetConfigName(filename)
	//路径
	viper.AddConfigPath(dir)
	_, ok := file.IsExists(path)
	if !ok {
		if err := os.MkdirAll(filepath.Dir(path), 0766); err != nil {
			panic(err)
		}
		f, err := os.Create(path)
		if err != nil {
			panic(err)
		}
		defer f.Close()

		var r io.Reader = bytes.NewReader(bys)
		if err := viper.ReadConfig(r); err != nil {
			panic(err)
		}

		if err := viper.WriteConfig(); err != nil {
			panic(err)
		}
	}

	if err := viper.Unmarshal(&c); err != nil {
		panic(err)
	}

	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		if err := viper.Unmarshal(&c); err != nil {
			panic(err)
		}
	})
}
