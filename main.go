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
	"flag"
	"fmt"
	"io/ioutil"
	"log"
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
	checkErr(err)
	defer outFile.Close()

	// Create a new zip archive.
	w := zip.NewWriter(outFile)

	// Add some files to the archive.
	addFiles(w, baseFolder, "")

	// Make sure to check the error on Close.
	err = w.Close()
	checkErr(err)
}

func addFiles(w *zip.Writer, basePath, baseInZip string) {
	// Open the Directory
	files, err := ioutil.ReadDir(basePath)
	checkErr(err)

	for _, file := range files {
		log.Println("[INFO] " + basePath + file.Name())
		if !file.IsDir() {
			dat, err := ioutil.ReadFile(basePath + file.Name())
			checkErr(err)

			// Add some files to the archive.
			f, err := w.Create(baseInZip + file.Name())
			checkErr(err)

			_, err = f.Write(dat)
			checkErr(err)
		} else if file.IsDir() {

			// Recurse
			newBase := basePath + file.Name() + "/"
			log.Println("[INFO] Recursing and Adding SubDir: " + file.Name())
			log.Println("[INFO] Recursing and Adding SubDir: " + newBase)

			addFiles(w, newBase, file.Name()+"/")
		}
	}
}

/*
  Function used to upload the zip in bucket
*/
func uploadArchive(filename, bucket string) {
	// maybe need to change the region
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(Region),
	})
	checkErr(err)

	newUpload := s3manager.NewUploader(sess)

	file, err := os.Open(filename)
	checkErr(err)

	result, err := newUpload.Upload(&s3manager.UploadInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(filepath.Base(filename)),
		Body:   file,
	})
	checkErr(err)

	log.Println("[INFO] Upload successfully! Path of archive:", result.Location)
}

func checkErr(err error) {
	if err != nil {
		log.Printf("[ERROR] %s", err)
		os.Exit(1)
	}
}

func main() {
	var (
		bucketNamePtr     = flag.String("bucket", "", "Bucket name")
		folderToUploadPtr = flag.String("path", "", "Folder to upload")
		zipNamePtr        = flag.String("zip", "", "Name of the compressed file")
		bucket            string
		path              string
		filename          string
	)

	flag.StringVar(bucketNamePtr, "b", "", "Bucket name")
	flag.StringVar(folderToUploadPtr, "p", "", "Folder to upload")
	flag.StringVar(zipNamePtr, "z", "", "Name of the compressed file")
	flag.Parse()

	bucket = *bucketNamePtr
	path = *folderToUploadPtr
	filename = *zipNamePtr

	if bucket == "" || path == "" || filename == "" {
		log.Println("[WARNING] Correct syntax: aws-backup -b name-of-bucket -p absolute/volumes/path/ -f filename.zip")
		os.Exit(1)
	}

	t := time.Now()

	fmt.Println("=== STARTING NEW BACKUP ====")
	log.Println("[INFO] Time: " + t.String())

	log.Println("[INFO] Compressing files...")
	zipWriter(path, filename)
	log.Println("[INFO] Compressed files!")

	log.Printf("[INFO] Uploading %s...\n", filename)
	uploadArchive(filename, bucket)

}
