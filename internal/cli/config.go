package cli

import "github.com/spf13/viper"

type App struct {
	GithubUser      string
	GithubToken     string
	GithubRepo      string
	GithubConfigDir string
	GithubValueDir  string
}

func ParseConfig() (*App, error) {
	v := viper.New()
	v.AddConfigPath("$HOME/.config/quimby")
	v.SetConfigName("config")
	v.SetConfigType("yaml")

	v.SetEnvPrefix("quimby")
	v.AutomaticEnv()

	cfg := &App{}
	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, err
		}
	}

	if err := v.Unmarshal(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}
