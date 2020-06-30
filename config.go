package main

import (
	"errors"
	"fmt"
	"log"

	"github.com/manifoldco/promptui"
	"github.com/spf13/viper"
)

var configKeys = []configItem{
	{
		Name:    "path",
		Prompt:  "Local path to store repos",
		Default: "repos",
	},
	{
		Name:   "ghUser",
		Prompt: "GitHub Username",
	},
	{
		Name:   "ghOrg",
		Prompt: "GitHub Organization (optional)",
	},
	{
		Name:     "ghKey",
		Prompt:   "GitHub API Token",
		Secret:   true,
		Validate: ghKeyValidate,
	},
	{
		Name:   "glUser",
		Prompt: "GitLab Username",
	},
	{
		Name:     "glKey",
		Prompt:   "GitLab API Token",
		Secret:   true,
		Validate: glKeyValidate,
	},
	{
		Name:   "bbUser",
		Prompt: "Bitbucket Username",
	},
	{
		Name:     "bbKey",
		Prompt:   "Bitbucket API Token",
		Secret:   true,
		Validate: bbKeyValidate,
	},
}

type configItem struct {
	Name     string
	Prompt   string
	Default  string
	Secret   bool
	Validate func(string) error
}

func init() {
	fmt.Println("Config stuff...")
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			err = viper.WriteConfigAs("config.yaml")
			if err != nil {
				log.Panic(err)
			}
		} else {
			log.Panic(err)
		}
	}
	err := checkForConfigValues()
	if err != nil {
		log.Panic(err)
	}
}

func checkForConfigValues() error {
	for _, k := range configKeys {
		if v := viper.Get(k.Name); v != nil {
			fmt.Printf("Key \"%v\" already configured, skipping\n", k.Name)
		} else {
			prompt := promptui.Prompt{
				Label: k.Prompt,
			}
			if k.Default != "" {
				prompt.Default = k.Default
			}
			if k.Secret == true {
				prompt.Mask = '*'
			}
			if k.Validate != nil {
				prompt.Validate = k.Validate
			}
			result, err := prompt.Run()
			if err != nil {
				return err
			}

			viper.Set(k.Name, result)
			err = viper.WriteConfig()
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func ghKeyValidate(key string) error {
	if len(key) != 40 {
		return errors.New("GitHub Tokens must be 40 characters")
	}
	return nil
}

func glKeyValidate(key string) error {
	if len(key) != 20 {
		return errors.New("GitLab Tokens must be 20 characters")
	}
	return nil
}

func bbKeyValidate(key string) error {
	if len(key) != 20 {
		return errors.New("Bitbucket Tokens must be 20 characters")
	}
	return nil
}
