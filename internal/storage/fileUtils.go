package storage

import "os"

func FileExists(filename string) (bool, error) {
	stats, err := os.Stat(filename)

	if os.IsNotExist(err) {
		return false, nil
	}

	if err != nil {
		return false, err
	}

	return !stats.IsDir(), nil
}

func CreateDirectory(directoryName string) error {
	err := os.Mkdir(directoryName, 0700)

	if err == nil || os.IsExist(err) {
		return nil
	}

	return err
}
