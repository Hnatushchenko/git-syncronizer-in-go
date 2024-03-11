package main

import (
	fileHelper "github.com/hnatushchenko/git-syncronizer/helpers"
	"log"
	"os"
	"path/filepath"
	"slices"
)

func removeFilesRecursively(directoryFullName string) bool {
	fileNamesToIgnore := GetFileNamesToIgnore()
	directoryNamesToIgnore := GetDirectoryNamesToIgnore()
	shouldCurrentDirectoryBeDeleted := true
	entries, err := os.ReadDir(directoryFullName)
	if err != nil {
		log.Printf("Failed to read the directory %s\n", directoryFullName)
		log.Println(err)
		return false
	}

	for _, entry := range entries {
		if entry.IsDir() || slices.Contains(fileNamesToIgnore, entry.Name()) {
			continue
		}

		fileFullName := filepath.Join(directoryFullName, entry.Name())
		err = os.Remove(fileFullName)
		if err != nil {
			log.Printf("Failed to remove the file %s\n", fileFullName)
			log.Println(err)
			continue

		}

		log.Printf("Removed the file %s", fileFullName)
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		if slices.Contains(directoryNamesToIgnore, entry.Name()) {
			shouldCurrentDirectoryBeDeleted = false
			continue
		}

		childDirectoryFullName := filepath.Join(directoryFullName, entry.Name())
		shouldCurrentDirectoryBeDeleted = removeFilesRecursively(childDirectoryFullName)
	}

	if shouldCurrentDirectoryBeDeleted {
		err = os.RemoveAll(directoryFullName)
		if err != nil {
			log.Printf("Failed to remove the current directory %s", directoryFullName)
			log.Println(err)
			return false
		}
	}

	return shouldCurrentDirectoryBeDeleted
}

func copyFilesUsingGoRoutines(sourceDirectoryFullName string,
	destinationDirectoryFullName string) {
	entries, err := os.ReadDir(sourceDirectoryFullName)
	if err != nil {
		log.Printf("Failed to read the directory %s\n", sourceDirectoryFullName)
		log.Println(err)
		return
	}

	for _, entry := range entries {
		sourceFileFullName := filepath.Join(sourceDirectoryFullName, entry.Name())
		newFileFullName := filepath.Join(destinationDirectoryFullName, entry.Name())
		_, err = fileHelper.CopyFile(sourceFileFullName, newFileFullName)
		if err != nil {
			log.Printf("Failed to copy file %s to %s\n", sourceFileFullName, newFileFullName)
			log.Println(err)
		}
		log.Printf("File copied successfully from %s to %s\r\n", sourceFileFullName, newFileFullName)
	}
}
