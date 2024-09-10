package config

import (
	"path/filepath"
	"strings"
)

type LogPaths struct {
	BrowserLog string
	ClientLog  string
	RuntimeLog string
  TypeIntervalLog string
}

// JSONPaths holds the paths for the JSON output of each type of log.
type JSONPaths struct {
	BrowserLogJSON string
	ClientLogJSON  string
	RuntimeLogJSON string
  TypeIntervalJSON string
}

func DeriveJSONPaths(baseDir string, logPaths LogPaths) JSONPaths {
	return JSONPaths{
		BrowserLogJSON: replaceExtension(logPaths.BrowserLog, ".json"),
		ClientLogJSON:  replaceExtension(logPaths.ClientLog, ".json"),
		RuntimeLogJSON: replaceExtension(logPaths.RuntimeLog, ".json"),
    TypeIntervalJSON: replaceExtension(logPaths.TypeIntervalLog, ".json"),
	}
}

func replaceExtension(filePath, newExt string) string {
	return strings.TrimSuffix(filePath, filepath.Ext(filePath)) + newExt
}
