package utils

import (
	"conc/customLog"
	"fmt"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

// GetConfFromEnvFile receives data for the database from the environment file. If successful, returns a non-empty map.
func GetConfFromEnvFile() map[string]string {
	resp := make(map[string]string)
	envFile, err := godotenv.Read(".env")
	if err == nil {
		resp = envFile
	} else {
		customLog.Logging(err)
	}
	return resp
}

// GCRunAndPrintMemory runs a garbage collection and
// if setting the APP_ENV environment variable as "dev" prints currently allocated number of bytes on the heap.
func GCRunAndPrintMemory() {
	debugSet := false
	settings := GetConfFromEnvFile()
	if val, ok := settings["APP_ENV"]; ok && val == "dev" {
		debugSet = true
	}
	if debugSet {
		var stat runtime.MemStats
		runtime.ReadMemStats(&stat)
		fmt.Println(stat.Alloc / 1024)
	}
	if val, ok := settings["GC_MANUAL_RUN"]; ok && val == "true" {
		runtime.GC()
	}
}

// GetEnvByKey returns a string from the env file with the value of the passed string key; if there is no file or such key, then an empty string.
func GetEnvByKey(key string) string {
	mapEnv := GetConfFromEnvFile()
	val, ok := mapEnv[key]
	if ok {
		return val
	} else {
		return ""
	}
}

// CreateDir creates a directory using the passed path string with permissions 0777, returns the directory string and an error.
func CreateDir(dirName string) (string, error) {
	err := os.MkdirAll(dirName, 0777)

	if err != nil {
		customLog.Logging(err)
	}

	return dirName, err
}

// ConcatSlice returns a string from the elements of the passed slice with strings. Separator - space.
func ConcatSlice(strSlice []string) string {
	resp := ""
	if len(strSlice) > 0 {
		var strBuilder strings.Builder
		for _, val := range strSlice {
			strBuilder.WriteString(val)
		}
		resp = strBuilder.String()
		strBuilder.Reset()
	}
	return resp
}

// Duration returns a formatted duration string based on the passed Time.
func Duration(start time.Time) string {
	return fmt.Sprintf("%v\n", time.Since(start))
}
