package main

import (
	"flag"

	"github.com/golang/glog"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

var (
	// Used for flags.
	cfgFile     string
	userLicense string

	rootCmd = &cobra.Command{
		Use:   "mvmapi",
		Short: "mvmapi is the API server for mvm.digital",
		Long:  "mvmapi is the API server for mvm.digital",
	}
)

func init() {
	cobra.OnInitialize(initConfig)
	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)
}

func initConfig() {
	viper.AutomaticEnv()
}

func main() {
	flag.Parse()
	err := rootCmd.Execute()
	if err != nil {
		glog.Error(err)
	}
}
