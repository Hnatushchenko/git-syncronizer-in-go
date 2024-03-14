package main

import (
	fileHelper "github.com/hnatushchenko/git-syncronizer/helpers"
	"log"
	"os"
	"path/filepath"
	"slices"
	"sync"
)

var resultChan = make(chan int64)
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

	page := 0
	pageSize := 50
	recordsProcessed := 0
	accumulationChannel := make(chan int64, 1)
	go accumulateResults(accumulationChannel)
	for ok := true; ok; ok = recordsProcessed > 0 {
		recordsProcessed = 0
		for i := page * pageSize; i < page*pageSize+pageSize; i++ {

			if i >= len(entries) {
				break
			}

			name := entries[i].Name()
			sourceFileFullName := filepath.Join(sourceDirectoryFullName, name)
			newFileFullName := filepath.Join(destinationDirectoryFullName, name)
			waitGroup.Add(1)
			go copyFile(sourceFileFullName, newFileFullName)
			recordsProcessed++
		}
		waitGroup.Wait()
		page++
	}

	waitGroup.Wait()
	close(resultChan)

	result := <-accumulationChannel
	return result
}

func copyFile(sourceFileFullName string,
	newFileFullName string,
) {
	numberOfBytes, _ := fileHelper.CopyFile(sourceFileFullName, newFileFullName)
	log.Printf("Sending bytes: %v", numberOfBytes)
	resultChan <- numberOfBytes
	waitGroup.Done()
}

func accumulateResults(accumulationChannel chan<- int64) {
	var totalBytes int64 = 0
	for byteResult := range resultChan {
		log.Printf("Receiving %d bytes", byteResult)
		totalBytes += byteResult
	}
	accumulationChannel <- totalBytes
}
