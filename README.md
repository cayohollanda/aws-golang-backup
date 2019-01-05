# AWS Golang Backup Script
A Golang script created for backup of archives to AWS S3.

## Usage
```console
cayohollanda@pc:~$ aws-backup -b bucket-name -p /folder/to/upload -z name-of-the-zip.zip
```
or
```console
cayohollanda@pc:~$ aws-backup --bucket bucket-name --path /folder/to/upload --zip name-of-the-zip.zip
```

## What the script do?
* The script compress the folder which is passed as a parameter, after this, upload the compressed archive to bucket on AWS S3

## How configure the AWS Credentials?
* To configure AWS Credentials, you need to install the AWS CLI package and configure it
```console
cayohollanda@pc:~$ sudo apt-get install awscli
```
* After that, you need to run the command to configure credentials, with this command, you will pass a **Access Key ID**, **Secret Key ID**, **Region name** and **Output Format** (default is **json**)

```console
cayohollanda@pc:~$ aws configure
```

* After configure AWS Credentials, your pc stay prepared to use the script.

## How automate the script to backup all days
One solution to automate the backup of archives is use a [Cron](https://opensource.com/article/17/11/how-use-cron-linux)
* Build a Go script to an executable
```console
cayohollanda@pc:~$ go build main.go
```
* Now, install the cron (if not have installed)
```console
cayohollanda@pc:~$ sudo apt-get install cron
```
* After that, access the file to configure a schedules
```console
cayohollanda@pc:~$ crontab -e
```
* Schedule as you prefer, save and **h4v3 fun**!
