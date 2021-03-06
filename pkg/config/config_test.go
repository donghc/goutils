package config

import (
	_ "embed"
	"fmt"
	"testing"
	"time"
)

type Config struct {
	MySQL struct {
		Read struct {
			Addr string `toml:"addr"`
			User string `toml:"user"`
			Pass string `toml:"pass"`
			Name string `toml:"name"`
		} `toml:"read"`
		Write struct {
			Addr string `toml:"addr"`
			User string `toml:"user"`
			Pass string `toml:"pass"`
			Name string `toml:"name"`
		} `toml:"write"`
		Base struct {
			MaxOpenConn int `toml:"maxOpenConn"`
			MaxIdleConn int `toml:"maxIdleConn"`
		} `toml:"base"`
	} `toml:"mysql"`

	Redis struct {
		Addr         string `toml:"addr"`
		Pass         string `toml:"pass"`
		Db           int    `toml:"db"`
		MaxRetries   int    `toml:"maxRetries"`
		PoolSize     int    `toml:"poolSize"`
		MinIdleConns int    `toml:"minIdleConns"`
	} `toml:"redis"`

	Mail struct {
		Host string `toml:"host"`
		Port int    `toml:"port"`
		User string `toml:"user"`
		Pass string `toml:"pass"`
		To   string `toml:"to"`
	} `toml:"mail"`

	HashIds struct {
		Secret string `toml:"secret"`
		Length int    `toml:"length"`
	} `toml:"hashids"`

	Language struct {
		Local string `toml:"local"`
	} `toml:"language"`
}

var (
	//go embed 对文件夹的支持不是非常友好，当路径中出现特殊字符的时候，会忽略 @ref: fs.ValidPath
	//go:embed dev_configs.toml
	devConfigs []byte
)

func TestLoad(t *testing.T) {
	var (
		p = "D:\\workspace\\github\\goutils\\pkg\\config\\dev_configs.toml"
		c interface{}
	)

	Load(p, &c)
	fmt.Println(c)
	//修改文件，实时监听文件的变化
	time.Sleep(time.Second * 10)
	fmt.Println(c)
}

func TestLoadAndCreate(t *testing.T) {
	var (
		p = "D:\\workspace\\github\\goutils\\pkg\\config\\prod_configs.toml"
		c interface{}
	)

	LoadAndCreate(p, devConfigs, &c)
	fmt.Println(c)
	time.Sleep(time.Second * 20)
	fmt.Println(c)
}
