package main

import (
    fileHelper "github.com/hnatushchenko/git-syncronizer/helpers"
    "log"
)

const sourcePath = "C:\\Users\\gnatu\\mydisk\\L35\\repos\\test2\\luminaryone"
const destinationPath = "C:\\Users\\gnatu\\mydisk\\L35\\src2\\test"

func main() {
    validateInputDirectories()
    removeFilesRecursively(destinationPath)
    copyFilesRecursively(sourcePath, sourcePath, destinationPath, GetFileNamesToIgnore(), GetDirectoryNamesToIgnore())
}

func validateInputDirectories() {
    exists, err := fileHelper.FileExists(sourcePath)
    if err != nil {
        log.Println("Failed to check the existence of the source directory.")
        log.Fatalln(err)
    }
    if !exists {
        log.Fatalln("The source directory does not exist!")
    }
    
    exists, err = fileHelper.FileExists(destinationPath)
    if err != nil {
        log.Println("Failed to check the existence of the destination directory.")
        log.Fatalln(err)
    }
    if !exists {
        log.Fatalln("The destination directory does not exist!")
    }
}
