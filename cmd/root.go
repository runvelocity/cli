/*
Copyright © 2023 Utibeabasi Umanah utibeabasiumanah6@gmail.com

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"log"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile    string
	managerUrl string
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "velocity",
	Short: "Create and execute self-hosted serverless functions with firecracker",
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.velocity.yaml)")
	rootCmd.PersistentFlags().StringVar(&managerUrl, "manager-url", "https://velocity-manager.fly.dev", "default is localhost:8000")

}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		configPath := home + "/.velocity"

		if err := os.MkdirAll(configPath, os.ModePerm); err != nil {
			log.Fatalln("Error creating directories:", err)
		}

		_, err = os.Stat(configPath + "/config.json")
		if err != nil {
			_, err = os.Create(configPath + "/config.json")
			if err != nil {
				log.Fatalln("A fatal error occured: " + err.Error())
			}
		}

		viper.AddConfigPath(configPath)
		viper.SetConfigType("json")
		viper.SetConfigName("config")
		viper.SetDefault("ManagerUrl", managerUrl)
		_ = viper.SafeWriteConfig()
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found; ignore error if desired
		} else {
			// Config file was found but another error was produced
			log.Fatalln("An error occured while reading config at file", viper.ConfigFileUsed())
		}
	}
}
