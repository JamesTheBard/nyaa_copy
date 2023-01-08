package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"
)

type CacheRecord struct {
	FilterName string    `json:"filter_name"`
	LastUsed   time.Time `json:"last_used"`
	Count      int       `json:"count"`
}

func populateCache() {
	data, err := ioutil.ReadFile(config.CacheLocation)
	if err != nil {
		log.Debug("Cache does not exist, creating new cache.")
		*cache = []CacheRecord{}
	} else {
		log.Debug("Cache exists, loading cache file to memory.")
		err := json.Unmarshal(data, &cache)
		if err != nil {
			log.Fatalf("Could not parse cache file: '%s'.", config.CacheLocation)
		}
	}
	log.Debugf("Cache size: %d entries.", len(*cache))
}

func updateCache(filter_name string) {
	index := -1
	count := 1
	for idx, val := range *cache {
		if val.FilterName == filter_name {
			index = idx
			count = val.Count + 1
			break
		}
	}

	record := CacheRecord{
		FilterName: filter_name,
		LastUsed:   time.Now(),
		Count:      count,
	}

	if index != -1 {
		(*cache)[index] = record
		return
	}
	*cache = append(*cache, record)
}

func deleteCacheEntry(filter_name string) {
	for idx, val := range *cache {
		if val.FilterName == filter_name {
			*cache = append((*cache)[:idx], (*cache)[idx+1:]...)
			return
		}
	}
}

func saveCache() {
	if config.DryRun {
		log.Info("Skipping cache save due to dry-run")
		return
	}
	data, _ := json.MarshalIndent(*cache, "", " ")
	err := os.MkdirAll(filepath.Dir(config.CacheLocation), 0750)
	if err != nil {
		log.Warningf("Could not create directory for cache file: '%s'.", filepath.Dir(config.CacheLocation))
	}
	err = ioutil.WriteFile(config.CacheLocation, data, 0644)
	if err != nil {
		log.Warningf("Could not write cache file to disk: '%s'.", config.CacheLocation)
	}
}

func deleteExpiredFilters() {
	var oldFilters []string
	log.Info("Deleting all cache entries older than 4 months...")
	sixMonthsAgo := time.Now().AddDate(0, -4, 0)
	for _, entry := range *cache {
		if entry.LastUsed.Before(sixMonthsAgo) {
			path, _ := filepath.Abs(entry.FilterName)
			log.Infof("Deleting filter file: '%s'...", path)
			if !config.DryRun {
				err := os.Remove(path)
				if err != nil {
					log.Fatalf("Unable to delete file '%s'!", path)
				}
			} else {
				log.Infof("Not deleting filter '%s' due to dry run!", entry.FilterName)
			}
			oldFilters = append(oldFilters, entry.FilterName)
		}
	}
	for _, filter := range oldFilters {
		log.Infof("Removing cache entry: '%s'...", filter)
		deleteCacheEntry(filter)
	}
	saveCache()
}
