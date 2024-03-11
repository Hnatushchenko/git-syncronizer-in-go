package main

import (
	fileHelper "github.com/hnatushchenko/git-syncronizer/helpers"
	"log"
	"os"
	"path/filepath"
	"slices"
	"sync"
)

var resultChan = make(chan int64, 1500)
var waitGroup sync.WaitGroup

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

	return shouldCurrentDirectoryBeDeleted
}

func copyFilesUsingGoRoutines(sourceDirectoryFullName string,
	destinationDirectoryFullName string) int64 {
	entries, err := os.ReadDir(sourceDirectoryFullName)
	if err != nil {
		log.Fatal(err)
	}

	var totalBytes int64 = 0
	for _, entry := range entries {
		sourceFileFullName := filepath.Join(sourceDirectoryFullName, entry.Name())
		newFileFullName := filepath.Join(destinationDirectoryFullName, entry.Name())
		waitGroup.Add(1)
		go copyFile(sourceFileFullName, newFileFullName)
	}

	waitGroup.Wait()
	close(resultChan)
	for byteResult := range resultChan {
		log.Printf("Receiving %d bytes", byteResult)
		totalBytes += byteResult
	}

	return totalBytes
}

func copyFile(sourceFileFullName string,
	newFileFullName string,
) {
	numberOfBytes, _ := fileHelper.CopyFile(sourceFileFullName, newFileFullName)
	log.Printf("Sending bytes: %v", numberOfBytes)
	resultChan <- numberOfBytes
	waitGroup.Done()
}
