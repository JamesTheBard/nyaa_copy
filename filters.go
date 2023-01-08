package main

import (
	"io/ioutil"
	"path/filepath"

	"github.com/creasty/defaults"
	"gopkg.in/yaml.v3"
)

type Filter struct {
	Search      string `yaml:"search"`
	Destination string `yaml:"destination"`
	Season      int    `default:"1" yaml:"season"`
	Offset      int    `default:"0" yaml:"offset"`
	Rename      string `yaml:"rename"`
	FilterFile  string
}

func (s *Filter) UnmarshalYAML(unmarshal func(interface{}) error) error {
	defaults.Set(s)
	type plain Filter
	if err := unmarshal((*plain)(s)); err != nil {
		return err
	}
	return nil
}

func loadFilters() {
	files, err := ioutil.ReadDir(config.Directories.Shows)
	if err != nil {
		log.Fatalf("Could not read filter directory: '%s'.", config.Directories.Shows)
	}

	var newFilters []Filter
	for _, file := range files {
		fullPath := filepath.Join(config.Directories.Shows, file.Name())
		data, err := ioutil.ReadFile(fullPath)
		if err != nil {
			log.Warningf("Could not parse filter file: '%s'.", fullPath)
			continue
		}
		var filter Filter
		yaml.Unmarshal(data, &filter)
		filter.FilterFile = filepath.Join(config.Directories.Shows, file.Name())
		newFilters = append(newFilters, filter)
	}
	log.Debugf("Current filters: '%v'", filters)
	filters = &newFilters
	log.Debugf("Filter set size: %d", len(newFilters))
}
