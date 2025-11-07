package conf

import (
	"fmt"
	"path/filepath"
	
	"github.com/spf13/viper"
)

type CommitPrefix struct {
	Prefix      string `mapstructure:"prefix"`
	Description string `mapstructure:"desc"`
}

const (
	ConfDir  = "./conf"
	ConfFile = "prefix.toml"
)

func LoadPrefixes() ([]CommitPrefix, error) {
	viper.SetConfigFile(filepath.Join(ConfDir, ConfFile))
	
	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("读取配置文件失败: %w", err)
	}
	
	var prefixes []CommitPrefix
	if err := viper.UnmarshalKey("prefixes", &prefixes); err != nil {
		return nil, fmt.Errorf("解析 prefixes 配置失败: %w", err)
	}
	
	if len(prefixes) == 0 {
		return nil, fmt.Errorf("prefix.toml 中未配置任何 prefixes")
	}
	
	return prefixes, nil
}
