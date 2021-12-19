package main

import (
	"crypto/sha512"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"time"
)

// Eq return true if the two file path are equals, false otherwise and return an error if it can't read the file
func Eq(f0Path, f1Path string) (bool, error) {
	hash0, err := getHash(f0Path)
	if err != nil {
		return false, err
	}
	hash1, err := getHash(f1Path)
	if err != nil {
		return false, err
	}

	return hash0 == hash1, nil
}

// Duplicates return a slice of slice of equal files in a directory
func Duplicates(dirPath string) ([][]string, error) {
	files, err := ioutil.ReadDir(dirPath)
	if err != nil {
		return [][]string{}, err
	}

	// List of equals files we'll return
	var equalsFiles [][]string
	// Checked files to don't re-check duplicated files
	var checkedFiles []string

	// Check all files in the directory
	for _, file := range files {
		currentFilePath := filepath.Join(dirPath, file.Name())

		// Don't compare an already checked file or directory
		if contain(checkedFiles, currentFilePath) || file.IsDir() {
			continue
		}

		duplicates, _ := FindSame(dirPath, currentFilePath)
		checkedFiles = append(checkedFiles, currentFilePath)

		if len(duplicates) > 0 {
			// Add duplicates files to the list of duplication
			equalsFiles = append(
				equalsFiles,
				append([]string{currentFilePath}, duplicates...),
			)
			// Add duplicated files to checked files
			checkedFiles = append(checkedFiles, duplicates...)
		}
	}

	return equalsFiles, nil
}

// FindSame return same files as `filePath` in the `dirPath` directory
func FindSame(dirPath string, filePath string) ([]string, error) {
	files, err := ioutil.ReadDir(dirPath)
	if err != nil {
		return []string{}, err
	}

	var sameFiles []string
	for _, f := range files {
		fPath := filepath.Join(dirPath, f.Name())

		// Don't compare file to directory or himself
		if f.IsDir() || fPath == filePath {
			continue
		}

		// Compare files
		eq, err := Eq(filePath, fPath)
		if err != nil {
			// Ignore possible errors and don't add the file
			continue
		}

		if eq {
			sameFiles = append(sameFiles, fPath)
		}
	}
	return sameFiles, nil
}

func contain(s []string, f string) bool {
	for _, file := range s {
		if file == f {
			return true
		}
	}
	return false
}

func getHash(filePath string) (string, error) {
	f, err := os.Open(filePath)
	defer f.Close()
	if err != nil {
		return "", err
	}

	h := sha512.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", h.Sum(nil)), nil
}

// WatchFiles in the `srcdir` then move them in the `dstdir` if it doesn't exist
// Wait `duration` between each scan
func WatchFiles(srcdir, dstdir string, duration time.Duration) {
	for {
		// List files in the `srcdir`
		sourceFiles, err := ioutil.ReadDir(srcdir)
		if err != nil {
			log.Println(err)
		}

		// Check all files
		for _, f := range sourceFiles {
			fPath := filepath.Join(srcdir, f.Name())
			sameFiles, err := FindSame(dstdir, fPath)
			if err != nil {
				log.Println(err)
			}
			// The file already exist
			if len(sameFiles) != 0 {
				err := os.Remove(fPath)
				if err != nil {
					log.Println(err)
				}
				continue
			}

			// Move the file in the `dstdir`
			// Don't use os.Rename becaust it cause an `invalid cross-device link` when renaming on differents partitions
			file, err := os.ReadFile(fPath)
			if err != nil {
				log.Println(err)
				continue
			}
			err = os.WriteFile(filepath.Join(dstdir, f.Name()), file, 0644)
			if err != nil {
				log.Println(err)
				continue
			}
		}

		time.Sleep(duration)
	}
}
