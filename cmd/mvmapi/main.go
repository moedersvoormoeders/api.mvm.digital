package main

import (
	"flag"
	"fmt"
	"strings"

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

func initConfig() {
	viper.AutomaticEnv()
}

func main() {
	flag.Parse()
	cobra.OnInitialize(initConfig)
	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)

	err := rootCmd.Execute()
	if err != nil {
		glog.Error(err)
	}
}

func printLogo() {
	fmt.Println(strings.Join([]string{
		"                                _ _       _ _        _ ",
		"                               | (_)     (_) |      | |",
		"  _ __ _____   ___ __ ___    __| |_  __ _ _| |_ __ _| |",
		" | '_ ` _ \\ \\ / / '_ ` _ \\  / _` | |/ _` | | __/ _` | |",
		" | | | | | \\ V /| | | | | || (_| | | (_| | | || (_| | |",
		" |_| |_| |_|\\_/ |_| |_| |_(_)__,_|_|\\__, |_|\\__\\__,_|_|",
		"                                     __/ |             ",
		"                                    |___/              ",
	}, "\n"))
}
