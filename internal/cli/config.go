package cli

import "github.com/spf13/viper"

type Config struct {
	GithubUser      string
	GithubToken     string
	GithubRepo      string
	GithubConfigDir string
}

func ParseConfig() (*Config, error) {
	v := viper.New()
	v.AddConfigPath("$HOME/.config/quimby")
	v.SetConfigName("config")
	v.SetConfigType("yaml")

	v.SetEnvPrefix("quimby")
	v.AutomaticEnv()

	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, err
		}
	}

	cfg := &Config{}
	if err := v.Unmarshal(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}
