package utils

import (
	"os"
	"path/filepath"
	"runtime"
)

func CreatePathIfNotExist(filePath string) error {
	_, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		err = os.MkdirAll(filePath, 0666)
		if err != nil {
			return err
		}
	}
	if err != nil {
		return err
	}
	return nil
}

func PathExist(filePath string) (bool, error) {
	_, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}

func ReadFile(filePath string) ([]byte, error) {
	return os.ReadFile(filePath)
}

func WriteFile(filePath string, data []byte, perm os.FileMode) error {
	err := CreatePathIfNotExist(filepath.Dir(filePath))
	if err != nil {
		return err
	}
	err = os.WriteFile(filePath, data, perm)
	if err != nil {
		return err
	}
	return nil
}

func OpenFile(filePath string, flag int, perm os.FileMode) (*os.File, error) {
	err := CreatePathIfNotExist(filepath.Dir(filePath))
	if err != nil {
		return nil, err
	}
	file, err := os.OpenFile(filePath, flag, perm)
	if err != nil {
		return nil, err
	}
	return file, nil
}

func GetRootPath() string {
	rootPath, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	switch runtime.GOOS {
	case "linux":
		rootPath = filepath.Join(rootPath, ".ouroboros")
	case "windows":
		rootPath = filepath.Join(rootPath, "AppData", "Local", "Ouroboros")
	}
	return rootPath
}
