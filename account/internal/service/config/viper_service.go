package config

import (
	"fmt"

	"github.com/spf13/viper"

	"git.paygear.ir/giftino/account/internal/account/domain"
)

type config struct {
	path string
	ext  string
	name string
}

func New(path, ext, name string) domain.ConfigurationService {
	return &config{
		path: path,
		ext:  ext,
		name: name,
	}
}

func (c config) GetConfig() (*domain.Config, error) {
	viper.AddConfigPath(c.path)
	viper.SetConfigName(c.name)
	err := viper.ReadInConfig()

	if err != nil {
		fmt.Printf("%v", err)
	}

	conf := &domain.Config{}
	err = viper.Unmarshal(conf)
	if err != nil {
		return conf, err
	}

	return conf, err
}
