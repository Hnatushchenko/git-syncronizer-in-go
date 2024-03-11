package main

import (
    fileHelper "github.com/hnatushchenko/git-syncronizer/helpers"
    "log"
    "os"
    "path/filepath"
    "slices"
    "strings"
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

func copyFilesRecursively(currentSourceDirectoryFullName string,
    initialSourceDirectoryFullName string,
    destinationDirectoryFullName string,
    fileNamesToIgnore []string,
    directoryNamesToIgnore []string) {
    entries, err := os.ReadDir(currentSourceDirectoryFullName)
    if err != nil {
        log.Printf("Failed to read the directory %s\n", currentSourceDirectoryFullName)
        log.Println(err)
        return
    }
    
    for _, entry := range entries {
        if entry.IsDir() || slices.Contains(fileNamesToIgnore, entry.Name()) {
            continue
        }
        
        sourceFileFullName := filepath.Join(currentSourceDirectoryFullName, entry.Name())
        newFileFullName := strings.ReplaceAll(sourceFileFullName, initialSourceDirectoryFullName,
            destinationDirectoryFullName)
        newDirectoryFullName := strings.ReplaceAll(currentSourceDirectoryFullName,
            initialSourceDirectoryFullName,
            destinationDirectoryFullName)
        exists, err := fileHelper.FileExists(newDirectoryFullName)
        if err != nil {
            log.Printf("Failed to check the existence of the directory: %s\n", newDirectoryFullName)
            continue
        }
        
        if !exists {
            err := os.MkdirAll(newDirectoryFullName, os.ModePerm)
            log.Printf("Error creating new directory %s\r\n", newDirectoryFullName)
            log.Println(err)
        }
        
        _, err = fileHelper.CopyFile(sourceFileFullName, newFileFullName)
        if err != nil {
            log.Printf("Failed to copy file %s to %s\n", sourceFileFullName, newFileFullName)
            log.Println(err)
        }
        log.Printf("File copied successfully from %s to %s\r\n", sourceFileFullName, newFileFullName)
    }
    
    for _, entry := range entries {
        if !entry.IsDir() || slices.Contains(directoryNamesToIgnore, entry.Name()) {
            continue
        }
        
        childSourceDirectoryFullName := filepath.Join(currentSourceDirectoryFullName, entry.Name())
        copyFilesRecursively(childSourceDirectoryFullName, initialSourceDirectoryFullName,
            destinationDirectoryFullName,
            fileNamesToIgnore,
            directoryNamesToIgnore)
    }
}
