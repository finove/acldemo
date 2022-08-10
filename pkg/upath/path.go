package upath

import (
	"os"
	"path/filepath"
)

var (
	pwd, appPath string
)

// WorkDirectory return work directory path
func WorkDirectory() string {
	return pwd
}

// AppDirectory return app directory
func AppDirectory() string {
	return filepath.Dir(appPath)
}

func init() {
	pwd, _ = os.Getwd()
	appPath, _ = filepath.Abs(os.Args[0])
}
