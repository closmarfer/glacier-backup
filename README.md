# Glacier Backup

App to backup your files in [S3 Glacier Deep Archive](https://aws.amazon.com/en/s3/glacier/) storage or to another folder or Volume in your computer made using Go.

## Requirements

* Go 1.23 installed in your computer
* *Optional (only for S3)*:
  * AWS account
  * A S3 bucket
  * API key and secret key in the `~/.aws/credentials`

## Configuration

The application is configured using environment variables. You must set the following variables before running the application:

### General Configuration

* `GLACIER_BACKUP_PATHS_TO_BACKUP`: A list of absolute paths to the folders you want to backup, separated by `;`.
  * Example: `/Users/me/Documents;/Users/me/Pictures`
* `GLACIER_BACKUP_IGNORED_PATTERNS`: A list of file patterns to ignore, separated by `;`.
  * Example: `*.DS_Store;*.tmp`

### Remote Storage Configuration

Depending on the remote storage you choose (`s3` or `local`), you need to set additional variables.

#### AWS S3 Glacier (`s3`)

* `GLACIER_BACKUP_S3_BUCKETS`: The name of the S3 bucket where files will be stored.
* `GLACIER_BACKUP_S3_REGION`: The AWS region of your bucket (e.g., `us-east-1`, `eu-west-1`).
* `GLACIER_BACKUP_S3_PROFILE`: The AWS profile name from your `~/.aws/credentials` file to use for authentication.

#### Local Storage (`local`)

* `GLACIER_BACKUP_LOCAL_DESTINATION_PATH`: The absolute path where the backup will be stored locally.

## Usage

Run the application using the command line. The general syntax is:

```sh
go run cmd/main.go [remote] [action]
```

### Arguments

1. **Remote**: The storage backend to use. Options are `s3` or `local`.
2. **Action**: The operation to perform.
   * `--backup`: Starts the backup process.
   * `--sizeCount`: Calculates and displays the total size of the files to be backed up.
   * `--cleanRemote`: Cleans up files in the remote storage that are no longer present locally (if applicable/implemented).

### Examples

**Backup to S3:**

```sh
export GLACIER_BACKUP_PATHS_TO_BACKUP="/Users/me/data"
export GLACIER_BACKUP_IGNORED_PATTERNS="*.tmp"
export GLACIER_BACKUP_S3_BUCKETS="my-backup-bucket"
export GLACIER_BACKUP_S3_REGION="us-east-1"
export GLACIER_BACKUP_S3_PROFILE="default"

go run cmd/main.go s3 --backup
```

**Backup to Local Configuration:**

```sh
export GLACIER_BACKUP_PATHS_TO_BACKUP="/Users/me/data"
export GLACIER_BACKUP_IGNORED_PATTERNS="*.tmp"
export GLACIER_BACKUP_LOCAL_DESTINATION_PATH="/Volumes/ExternalDrive/Backup"

go run cmd/main.go local --backup
```

### Stopping and Resuming

If you need to stop the process, you can use `Ctrl + C`. The application will gracefully shut down, ensuring the current state is saved.

Found files are stored in a local SQLite database (`backup.db`). When you run the command again, the backup process will continue from where it left off, skipping files that have already been uploaded (unless they have been modified).

## Develop

This application has no infrastructure requirements (DB, cache) so to develop you can run
simply the `cmd/main.go` file.

To prevent uploading the same files multiple times, the application uses a SQLite database (`backup.db`) that is stored remotely (S3 bucket or local destination, depending on the selected remote) and downloaded at startup.

If you execute the command again and some files have been modified after their last upload time, they will be uploaded again.

### TODO list

* Add more remote storages (eg. [Cloud Storage](https://cloud.google.com/storage)) implementing the `RemoteFilesRepository interface` and documenting the required environment variables.
