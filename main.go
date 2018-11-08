/*
    - Created by: http://github.com/cayohollanda

	- funtions used to compress the folder for zip: zipWriter() and addFiles()
	  found on this link -> https://stackoverflow.com/questions/37869793/how-do-i-zip-a-directory-containing-sub-directories-or-files-in-golang

	- configure your aws credentials
	- if not know, install the AWS CLI
	- after, run the command: aws configure
	- after configure, you can run the script

	- if necessary, change the 'Region' constant with the
	  region of your bucket

	- to run: go run main.go name-of-bucket /absolute/path/to/folder/ filename.zip

	- you can use cron to automate your backups
	  with cron, you can personalize the exactly time
	  or time interval that you want between one backup
	  and another
*/

package main

import (
	"archive/zip"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

// Region need to be changed if necessary
const Region = "us-east-1"

func zipWriter(path, filename string) {
	baseFolder := path

	// Get a Buffer to Write To
	outFile, err := os.Create(filename)
	if err != nil {
		fmt.Println(err)
	}
	defer outFile.Close()

	// Create a new zip archive.
	w := zip.NewWriter(outFile)

	// Add some files to the archive.
	addFiles(w, baseFolder, "")

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Make sure to check the error on Close.
	err = w.Close()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func addFiles(w *zip.Writer, basePath, baseInZip string) {
	// Open the Directory
	files, err := ioutil.ReadDir(basePath)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	for _, file := range files {
		fmt.Println(basePath + file.Name())
		if !file.IsDir() {
			dat, err := ioutil.ReadFile(basePath + file.Name())
			if err != nil {
				fmt.Println(err)
			}

			// Add some files to the archive.
			f, err := w.Create(baseInZip + file.Name())
			if err != nil {
				fmt.Println(err)
			}
			_, err = f.Write(dat)
			if err != nil {
				fmt.Println(err)
			}
		} else if file.IsDir() {

			// Recurse
			newBase := basePath + file.Name() + "/"
			fmt.Println("Recursing and Adding SubDir: " + file.Name())
			fmt.Println("Recursing and Adding SubDir: " + newBase)

			addFiles(w, newBase, file.Name()+"/")
		}
	}
}

/*
  Function used to upload the zip in bucket

  @param filename Receive a name of file to be uploaded
  @param bucket   Receive a name of bucket where the file will go
*/
func uploadArchive(filename, bucket string) {
	// maybe need to change the region
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(Region),
	})
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	newUpload := s3manager.NewUploader(sess)

	file, err := os.Open(filename)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	result, err := newUpload.Upload(&s3manager.UploadInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(filepath.Base(filename)),
		Body:   file,
	})

	if err != nil {
		fmt.Println("An error appeared: ", err)
		os.Exit(1)
	}

	fmt.Println("Upload successfully! Path of archive:", result.Location)
}

func main() {

	if len(os.Args) != 4 {
		fmt.Println("Correct syntax: go run main.go name-of-bucket absolute/volumes/path/ filename.zip")
		os.Exit(1)
	}

	var (
		bucket   = os.Args[1]
		path     = os.Args[2]
		filename = os.Args[3]
	)

	t := time.Now()

	fmt.Println("=== STARTING NEW BACKUP ====")
	fmt.Println("Time: " + t.String())

	fmt.Println("Compressing files...")
	zipWriter(path, filename)
	fmt.Println("Compressed files!")

	fmt.Printf("Uploading %s...\n", filename)
	uploadArchive(filename, bucket)

}
