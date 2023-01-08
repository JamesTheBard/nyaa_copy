package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/akamensky/argparse"
	"github.com/op/go-logging"
	"gopkg.in/yaml.v3"
)

type ConfigFile struct {
	Directories struct {
		Shows       string `yaml:"shows"`
		Videos      string `yaml:"videos"`
		Destination string `yaml:"destination"`
	}
	AllowedExtensions  []string `yaml:"allowed_extensions"`
	CacheLocation      string
	ShowExpiredFilters bool
	DryRun             bool
}

func getConfig() {
	log.Debugf("Parsing configuration...")
	homeDir, _ := os.UserConfigDir()
	cacheDir, _ := os.UserCacheDir()
	homeDir = filepath.Join(homeDir, "nyaa_copy/config.yaml")
	cacheDir = filepath.Join(cacheDir, "nyaa_copy/cache.json")
	parser := argparse.NewParser("nyaa_copy", "A custom copying program that can relabel videos, add season information, and track the number of times a given filter was used.")
	configFilePath := parser.String("c", "config", &argparse.Options{Required: false, Help: "Use custom configuration file", Default: homeDir})
	cacheFilePath := parser.String("", "cache", &argparse.Options{Required: false, Help: "Use custom cache file", Default: cacheDir})
	showOldFilters := parser.Flag("e", "expired", &argparse.Options{Required: false, Help: "Delete old filter files and exit"})
	dryRun := parser.Flag("d", "dry-run", &argparse.Options{Required: false, Help: "Do not actually copy files, just simulate", Default: false})
	loggingDebug := parser.Flag("", "debug", &argparse.Options{Required: false, Help: "Set logging level to DEBUG"})
	appVersion := parser.Flag("", "version", &argparse.Options{Required: false, Help: "Display version and exit"})
	err := parser.Parse(os.Args)

	if err != nil {
		fmt.Print(parser.Usage(err))
		os.Exit(0)
	}

	if *appVersion {
		fmt.Printf("%s version %s\n", filepath.Base(os.Args[0]), version)
		os.Exit(0)
	}

	if *loggingDebug {
		backend := logging.NewLogBackend(os.Stderr, "", 0)
		debugFormat := logging.MustStringFormatter(
			`%{color}%{time:2006-01-02 15:04:05.000} [%{level:.4s} %{id:03x}] (%{shortfunc})%{color:reset} %{message}`,
		)
		logging.SetLevel(logging.DEBUG, "")
		logging.SetFormatter(debugFormat)
		logging.SetBackend(backend)
	}

	yamlFile, err := ioutil.ReadFile(*configFilePath)
	if err != nil {
		log.Fatalf("Could not read configuration file: '%s'.", *configFilePath)
	}

	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		log.Fatalf("Could not parse configuration file: '%s'.", *configFilePath)
	}

	config.CacheLocation = *cacheFilePath
	config.ShowExpiredFilters = *showOldFilters
	config.DryRun = *dryRun
}
