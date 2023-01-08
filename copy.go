package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/dustin/go-humanize"
)

func processFiles() {
	entries, err := os.ReadDir((*config).Directories.Videos)
	if err != nil {
		log.Fatalf("Cannot open download path for processing: '%s'.", (*config).Directories.Videos)
	}
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		wasProcessed := false
		for _, filter := range *filters {
			regex := fmt.Sprintf("%s - (\\d+)", filter.Search)
			regex_compiled := regexp.MustCompile(regex)
			res := regex_compiled.FindStringSubmatch(entry.Name())
			if len(res) == 0 {
				continue
			}
			wasProcessed = true
			newFilename, err := generateFileName(entry.Name(), res[1], &filter)
			if err != nil {
				break
			}
			entryInfo, _ := entry.Info()
			entryPath := filepath.Join(config.Directories.Videos, entry.Name())
			entryDestPartial := filepath.Join(config.Directories.Destination, filter.Destination, fmt.Sprintf("Season %d", filter.Season))
			entryDest := filepath.Join(entryDestPartial, newFilename)
			log.Debugf("Processing info: '%s', '%s', '%s'", entryPath, entryDestPartial, entryDest)
			if !isAllowedExtension(entry.Name()) {
				log.Warningf("File '%s' does not have an allowed file extension, skipping!", entry.Name())
				break
			}
			entryDestInfo, err := os.Stat(entryDest)
			if err == nil {
				if entryDestInfo.Size() == entryInfo.Size() {
					log.Infof("Skipping '%s', already copied.", entry.Name())
					break
				}
			}
			log.Infof("Copying '%s' to '%s' (%s)...", newFilename, entryDestPartial, humanize.Bytes(uint64(entryInfo.Size())))
			err = copyFile(entryPath, entryDest)
			if err != nil {
				log.Fatalf("Could not copy file: '%s' -> '%s'", entryPath, entryDest)
			}
			log.Debugf("Updating cache...")
			updateCache(filter.FilterFile)
		}
		if !wasProcessed {
			log.Warningf("No filter match: '%s'", entry.Name())
		}
	}
}

func isAllowedExtension(sourceFile string) bool {
	currentExt := filepath.Ext(sourceFile)
	for _, ext := range config.AllowedExtensions {
		if ext == currentExt {
			return true
		}
	}
	return false
}

func copyFile(sourceFile string, destFile string) error {
	if config.DryRun {
		log.Infof("Skipping file copy operation due to dry run!")
		return nil
	}
	err := os.MkdirAll(filepath.Dir(destFile), 0755)
	if err != nil {
		return err
	}

	fin, err := os.Open(sourceFile)
	if err != nil {
		return err
	}
	defer fin.Close()

	fout, err := os.Create(destFile)
	if err != nil {
		return err
	}
	defer fout.Close()

	_, err = io.Copy(fout, fin)
	if err != nil {
		return err
	}
	return nil
}

func generateFileName(fileName string, episode string, filter *Filter) (string, error) {
	log.Debugf("Filter name: '%s', '%s'", filter.FilterFile, filter.Search)
	log.Debugf("File name passed: '%s'", fileName)
	log.Debugf("Filter information: '%v'", *filter)
	intEpisode, _ := strconv.Atoi(episode)
	intEpisode += (*filter).Offset
	if intEpisode < 0 {
		log.Warningf("Invalid episode value generated: '%s' [episode %s, offset %d]", fileName, episode, filter.Offset)
		return "", errors.New("invalid episode value generated after offset conversion")
	}
	if filter.Rename != "" {
		// fileName = strings.Replace(fileName, filter.Search, filter.Rename, 1)
		regex := regexp.MustCompile(filter.Search)
		regex.ReplaceAllString(fileName, filter.Rename)
	}

	find := fmt.Sprintf(" - %s", episode)
	replace := fmt.Sprintf(" - %dx%02d", (*filter).Season, intEpisode)
	fileName = strings.Replace(fileName, find, replace, 1)
	log.Debugf("New filename: '%s'.", fileName)
	return fileName, nil
}
