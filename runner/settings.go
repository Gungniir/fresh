package runner

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/pilu/config"
)

const (
	envSettingsPrefix   = "RUNNER_"
	mainSettingsSection = "Settings"
)

var settings = map[string]string{
	"config_path":       "./runner.conf",
	"root":              ".",
	"workdir":           ".",
	"tmp_path":          "./tmp",
	"build_name":        "runner-build",
	"build_log":         "runner-build-errors.log",
	"valid_ext":         ".go, .tpl, .tmpl, .html",
	"no_rebuild_ext":    ".tpl, .tmpl, .html",
	"ignored":           "assets, tmp",
	"build_delay":       "600",
	"build_root":        ".",
	"use_dlv":           "1",
	"colors":            "1",
	"log_color_main":    "cyan",
	"log_color_build":   "yellow",
	"log_color_runner":  "green",
	"log_color_watcher": "magenta",
	"log_color_app":     "",
}

var colors = map[string]string{
	"reset":          "0",
	"black":          "30",
	"red":            "31",
	"green":          "32",
	"yellow":         "33",
	"blue":           "34",
	"magenta":        "35",
	"cyan":           "36",
	"white":          "37",
	"bold_black":     "30;1",
	"bold_red":       "31;1",
	"bold_green":     "32;1",
	"bold_yellow":    "33;1",
	"bold_blue":      "34;1",
	"bold_magenta":   "35;1",
	"bold_cyan":      "36;1",
	"bold_white":     "37;1",
	"bright_black":   "30;2",
	"bright_red":     "31;2",
	"bright_green":   "32;2",
	"bright_yellow":  "33;2",
	"bright_blue":    "34;2",
	"bright_magenta": "35;2",
	"bright_cyan":    "36;2",
	"bright_white":   "37;2",
}

func logColor(logName string) string {
	settingsKey := fmt.Sprintf("log_color_%s", logName)
	colorName := settings[settingsKey]

	return colors[colorName]
}

func loadEnvSettings() {
	for key, _ := range settings {
		envKey := fmt.Sprintf("%s%s", envSettingsPrefix, strings.ToUpper(key))
		if value := os.Getenv(envKey); value != "" {
			settings[key] = value
		}
	}
}

func loadRunnerConfigSettings() {
	if _, err := os.Stat(configPath()); err != nil {
		return
	}

	mainLog("Loading settings from %s", configPath())
	sections, err := config.ParseFile(configPath(), mainSettingsSection)
	if err != nil {
		mainLog("Failed to load settings %s", err.Error())
		return
	}

	mainLog("Settings loaded %v", sections)

	for key, value := range sections[mainSettingsSection] {
		settings[key] = value
	}
}

func initSettings() {
	loadEnvSettings()
	loadRunnerConfigSettings()
	mainLog("Loaded config: %v", settings)
}

func getenv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}

	return defaultValue
}

func root() string {
	return settings["root"]
}

func buildRoot() string {
	return normalizeDir(settings["build_root"])
}

func tmpPath() string {
	return normalizeDir(settings["tmp_path"])
}

func workDir() string {
	return normalizeDir(settings["workdir"])
}

func buildName() string {
	return settings["build_name"]
}
func buildPath() string {
	p := filepath.Join(tmpPath(), buildName())
	if runtime.GOOS == "windows" && filepath.Ext(p) != ".exe" {
		p += ".exe"
	}
	return p
}

func buildErrorsFileName() string {
	return settings["build_log"]
}

func buildErrorsFilePath() string {
	return filepath.Join(tmpPath(), buildErrorsFileName())
}

func configPath() string {
	return settings["config_path"]
}

func buildDelay() time.Duration {
	value, _ := strconv.Atoi(settings["build_delay"])

	return time.Duration(value)
}

func normalizeDir(dir string) string {
	wd, err := filepath.Abs(dir)
	if err != nil {
		panic(err)
	}

	return wd
}

func isDelve() bool {
	return settings["use_dlv"] == "1"
}
