package main

import (
	"fmt"
	fileHelper "github.com/hnatushchenko/git-syncronizer/helpers"
	"log"
	"time"
)

const sourcePath = "C:\\test4\\screens"
const destinationPath = "C:\\test4\\destination"

func main() {
	//validateInputDirectories()
	//log.SetOutput(io.Discard)
	removeFilesRecursively(destinationPath)
	start := time.Now()
	numOfBytes := copyFilesUsingGoRoutines(sourcePath, destinationPath)
	elapsed := time.Since(start)
	fmt.Printf("Coping took %s\r\n", elapsed)
	fmt.Printf("%v bytes copied", numOfBytes)
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
