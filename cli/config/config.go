package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/spf13/cobra"
)

// Args expose available global args for cli command
type Args struct {
	ConfigPath      string
	ShareConfigPath string
}

// GlobalConfig expose available global config for cli command
type GlobalConfig struct {
	LogLevel string `yaml:"LogLevel" env:"LOG_LEVEL" env-default:"info" env-description:"Defines logger's log level'"`
}

// ProcessArgs processes and handles CLI arguments
func ProcessArgs(cfg interface{}, a *Args, cmd *cobra.Command) {
	configFlag := "config"
	cmd.PersistentFlags().StringVarP(&a.ConfigPath, configFlag, "c", "./config/config.yaml", "Path to configuration file")
	_ = cmd.MarkFlagRequired(configFlag)

	shareConfigFlag := "share-config"
	cmd.PersistentFlags().StringVarP(&a.ShareConfigPath, shareConfigFlag, "s", "", "Path to local share configuration file")
	_ = cmd.MarkFlagRequired(shareConfigFlag)

	envHelp, _ := cleanenv.GetDescription(cfg, nil)
	cmd.SetUsageTemplate(envHelp + "\n" + cmd.UsageTemplate())

}
