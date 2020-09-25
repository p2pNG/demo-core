package utils

import (
	"go.uber.org/zap"
	"os"
	"path/filepath"
	"runtime"
)

// HomeDir returns the best guess of the current user's home
// directory from environment variables. If unknown, "." (the
// current directory) is returned instead, except GOOS=android,
// which returns "/sdcard".
func HomeDir() string {
	home := homeDirUnsafe()
	if home == "" && runtime.GOOS == "android" {
		home = "/sdcard"
	}
	if home == "" {
		home = "."
	}
	return home
}

// homeDirUnsafe is a low-level function that returns
// the user's home directory from environment
// variables. Careful: if it cannot be determined, an
// empty string is returned. If not accounting for
// that case, use HomeDir() instead; otherwise you
// may end up using the root of the file system.
func homeDirUnsafe() string {
	home := os.Getenv("HOME")
	if home == "" && runtime.GOOS == "windows" {
		drive := os.Getenv("HOMEDRIVE")
		path := os.Getenv("HOMEPATH")
		home = drive + path
		if drive == "" || path == "" {
			home = os.Getenv("USERPROFILE")
		}
	}
	if home == "" && runtime.GOOS == "plan9" {
		home = os.Getenv("home")
	}
	return home
}

// AppConfigDir returns the directory where to store user's config.
//
// If XDG_CONFIG_HOME is set, it returns: $XDG_CONFIG_HOME/p2pNG.
// Otherwise, os.UserConfigDir() is used; if successful, it appends
// "p2pNG" to the path.
// If it returns an error, the fallback path "./p2pNG" is returned.
//
// The config directory is not guaranteed to be different from
// AppDataDir().
//
// Unlike os.UserConfigDir(), this function prefers the
// XDG_CONFIG_HOME env var on all platforms, not just Unix.
//
// Ref: https://specifications.freedesktop.org/basedir-spec/basedir-spec-latest.html
func AppConfigDir() string {
	if basedir := os.Getenv("XDG_CONFIG_HOME"); basedir != "" {
		return filepath.Join(basedir, "p2pNG")
	}
	basedir, err := os.UserConfigDir()
	if err != nil {
		Log().Warn("unable to determine directory for user configuration; falling back to current directory", zap.Error(err))
		return "./p2pNG"
	}

	return filepath.Join(basedir, "p2pNG")
}

// AppDataDir returns a directory path that is suitable for storing
// application data on disk. It uses the environment for finding the
// best place to store data, and appends a "p2pNG" subdirectory.
//
// For a base directory path:
// If XDG_DATA_HOME is set, it returns: $XDG_DATA_HOME/p2pNG; otherwise,
// on Windows it returns: %AppData%/p2pNG,
// on Mac: $HOME/Library/Application Support/p2pNG,
// on Plan9: $home/lib/p2pNG,
// on Android: $HOME/p2pNG,
// and on everything else: $HOME/.local/share/p2pNG.
//
// If a data directory cannot be determined, it returns "./p2pNG"
// (this is not ideal, and the environment should be fixed).
//
// The data directory is not guaranteed to be different from AppConfigDir().
//
// Ref: https://specifications.freedesktop.org/basedir-spec/basedir-spec-latest.html
func AppDataDir() string {
	if basedir := os.Getenv("XDG_DATA_HOME"); basedir != "" {
		return filepath.Join(basedir, "p2pNG")
	}
	switch runtime.GOOS {
	case "windows":
		appData := os.Getenv("AppData")
		if appData != "" {
			return filepath.Join(appData, "p2pNG")
		}
	case "darwin":
		home := homeDirUnsafe()
		if home != "" {
			return filepath.Join(home, "Library", "Application Support", "p2pNG")
		}
	case "plan9":
		home := homeDirUnsafe()
		if home != "" {
			return filepath.Join(home, "lib", "p2pNG")
		}
	case "android":
		home := homeDirUnsafe()
		if home != "" {
			return filepath.Join(home, "p2pNG")
		}
	default:
		home := homeDirUnsafe()
		if home != "" {
			return filepath.Join(home, ".local", "share", "p2pNG")
		}
	}
	return "./p2pNG"
}
