package main

import (
	"errors"
	"fmt"

	"github.com/spf13/viper"
)

// Configuration structure
type Configuration struct {
	port                  int
	frontendOrigin        string
	gitlabToken           string
	triggerToken          string
	jwtSecret             []byte
	gitlabURL             string
	gitlabProject         int
	debug                 bool
	prefetchGroups        bool
	groupIds              []string
	commitHistoryDays     int
	updateIntervalSeconds int
}

func readConfig() (*viper.Viper, error) {
	v := viper.New()
	// Use envrironment variables for hosted version with VIPER_ prefix
	v.SetEnvPrefix("viper") // will be uppercased automatically
	// v.BindEnv("gitlabToken")
	// v.BindEnv("triggerToken")
	v.SetDefault("port", 8000)
	v.SetDefault("frontendOrigin", "http://localhost:3000")
	v.SetDefault("gitlabURL", "https://gitlab.cern.ch/api/v4")
	v.SetDefault("gitlabProject", 56283)
	v.SetDefault("debug", true)
	v.SetDefault("prefetchGroups", true)
	v.SetDefault("commitHistoryDays", 90)
	v.SetDefault("updateIntervalSeconds", 300)
	v.SetDefault("groupIds", []string{
		"papers", "notes", "reports",
	})
	v.SetConfigName("config")
	v.AddConfigPath("config")
	v.AutomaticEnv()
	err := v.ReadInConfig()
	if err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found; ignore error if desired
			fmt.Println("Config file not found. Trying to use environment variables.")
			return v, nil
		}
	}
	return v, err
}

func validateAndSetConfig(v1 *viper.Viper) (Configuration, error) {
	var configuration Configuration
	configuration.port = v1.GetInt("port")
	configuration.frontendOrigin = v1.GetString("frontendOrigin")
	configuration.gitlabURL = v1.GetString("gitlabURL")
	configuration.gitlabProject = v1.GetInt("gitlabProject")
	configuration.debug = v1.GetBool("debug")
	configuration.prefetchGroups = v1.GetBool("prefetchGroups")
	configuration.groupIds = v1.GetStringSlice("groupIds")
	configuration.commitHistoryDays = v1.GetInt("commitHistoryDays")
	configuration.updateIntervalSeconds = v1.GetInt("updateIntervalSeconds")
	configuration.gitlabToken = v1.GetString("gitlabToken")

	if configuration.gitlabToken == "" {
		errorMessage := "gitlabToken cannot be empty."
		err := errors.New(errorMessage)
		return configuration, err
	}

	configuration.triggerToken = v1.GetString("triggerToken")
	if configuration.triggerToken == "" {
		errorMessage := "triggerToken cannot be empty."
		err := errors.New(errorMessage)
		return configuration, err
	}

	var jwtSecretString = v1.GetString("jwtSecret")
	if jwtSecretString == "" {
		errorMessage := "jwtSecret cannot be empty."
		err := errors.New(errorMessage)
		return configuration, err
	}
	configuration.jwtSecret = []byte(jwtSecretString)

	fmt.Printf("Reading config for port = %d\n", configuration.port)
	fmt.Printf("Reading config for frontendOrigin = %s\n", configuration.frontendOrigin)
	fmt.Printf("Reading config for gitlabURL = %s\n", configuration.gitlabURL)
	fmt.Printf("Reading config for gitlabProject = %d\n", configuration.gitlabProject)
	fmt.Printf("Reading config for debug = %t\n", configuration.debug)
	fmt.Printf("Reading config for prefetchGroups = %t\n", configuration.prefetchGroups)
	fmt.Printf("Reading config for groupIds = %#v\n", configuration.groupIds)
	fmt.Printf("Reading config for commitHistoryDays = %d\n", configuration.commitHistoryDays)
	fmt.Printf("Reading config for updateIntervalSeconds = %d\n", configuration.updateIntervalSeconds)
	// fmt.Printf("Reading config for gitlabToken = %s\n", configuration.gitlabToken)
	// fmt.Printf("Reading config for triggerToken = %s\n", configuration.triggerToken)
	return configuration, nil
}
