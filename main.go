package main

import (
	"os"

	"github.com/op/go-logging"
)

var cache *[]CacheRecord
var config *ConfigFile
var filters *[]Filter

var log = logging.MustGetLogger("main")
var format = logging.MustStringFormatter(
	`%{color}%{time:2006-01-02 15:04:05.000} [%{level:.4s} %{id:03x}]%{color:reset} %{message}`,
)

var version = "1.1.4"

func main() {
	backend := logging.NewLogBackend(os.Stderr, "", 0)
	backendFormatter := logging.NewBackendFormatter(backend, format)
	logging.SetBackend(backendFormatter)
	logging.SetFormatter(format)
	logging.SetLevel(logging.INFO, "")

	getConfig()
	log.Infof("Starting run...")

	log.Infof("Loading cache from '%s'.", config.CacheLocation)
	populateCache()

	if config.ShowExpiredFilters {
		deleteExpiredFilters()
		return
	}

	log.Infof("Loading filters from '%s'.", config.Directories.Shows)
	loadFilters()

	log.Infof("Processing downloaded files...")
	processFiles()

	log.Infof("Writing cache to file...")
	saveCache()
}
