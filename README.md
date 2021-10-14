# Glacier Backup

App to backup your files in [S3 Glacier Deep Archive](https://aws.amazon.com/en/s3/glacier/) storage made using Go.

## Requirements

* Go 1.16 installed in your computer
* AWS account
* Create a S3 bucket
* API key and secret key in the `~/.aws/credentials`

## Setup

1. Execute `make install` to build the application and create the required folders.
2. This command will create:
* An application `glacier-backup` in your `/usr/local/bin` folder
* A folder in `$HOME/.glacier-backup` directory with the `config.yaml` file.
3. Define the configuration params in your `$HOME/.glacier-backup/config.yaml` file
4. Run in your terminal executing:
```sh
glacier-backup
```
5. If you need to stop the process you can use `Ctrl + V` and
the backup process will continue where it stopped.

## Develop

This application has no infrastructure requirements (DB, cache) so to develop you can run
simply the `cmd/main.go` file.

To prevent upload the same files that are already uploaded, the application uses a CSV list
of uploaded files. This list (`uploaded_files.csv`) is stored in the same bucket and it 
contains a list of uploaded files + the uploaded date.

If you executes the command again and some files have the updated_at date after the `uploaded_date` from the CSV file
they will be updated in the remote bucket.

### TODO list

* Add testing
* Add more remote storages (eg. [Cloud Storage](https://cloud.google.com/storage)) implementing the `RemoteFilesRepository interface` and adding
the custom configuration params in the `var/config.yaml` file.