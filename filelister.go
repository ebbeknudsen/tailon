package main

import (
	"os"
	"path"
	"path/filepath"
	"time"
)

type ListEntry struct {
	Path    string    `json:"path"`
	Alias   string    `json:"alias"`
	Size    int64     `json:"size"`
	ModTime time.Time `json:"mtime"`
	Exists  bool      `json:"exists"`
}

func fileInfo(path string) *ListEntry {
	entry := ListEntry{}
	entry.Path = path

	info, err := os.Stat(path)
	if !os.IsNotExist(err) {
		entry.Exists = true
		entry.Size = info.Size()
		entry.ModTime = info.ModTime()
	}

	return &entry
}

var AllFiles map[string]bool

func createListing(filespecs []FileSpec) map[string][]*ListEntry {
	AllFiles = make(map[string]bool)
	res := make(map[string][]*ListEntry)

	for _, spec := range filespecs {
		group := "__default__"
		if spec.Group != "" {
			group = spec.Group
		}

		switch spec.Type {
		case "file":
			entry := fileInfo(spec.Path)
			if spec.Alias != "" {
				entry.Alias = spec.Alias
			} else {
				entry.Alias = entry.Path
			}
			res[group] = append(res[group], entry)
			AllFiles[entry.Path] = true
		case "glob":
			matches, _ := filepath.Glob(spec.Path)
			for _, match := range matches {
				entry := fileInfo(match)
				if spec.Alias != "" {
					entry.Alias = path.Join(spec.Alias, path.Base(entry.Path))
				} else {
					cwd, _ := os.Getwd()
					rel, _ := filepath.Rel(cwd, entry.Path)
					entry.Alias = rel
				}
				res[group] = append(res[group], entry)
				AllFiles[entry.Path] = true
			}
		}
	}

	return res
}

func fileAllowed(path string) bool {
	_, ok := AllFiles[path]
	return ok
}
