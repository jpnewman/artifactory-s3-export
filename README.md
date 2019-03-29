
# Artifactory S3 Export

A program to export NUGET packages from Artifactory repositories to an AWS S3 Bucket.

## Configuration

Create ```config.yaml``` file with content: -

~~~
mysql:
  connection_string: root:<MYSQL_PASSWORD>@tcp(localhost:3306)/artdb
  select_limit: 10
repo:
  list_file: repos.txt
  filestore_path: /var/opt/jfrog/artifactory/data/filestore
aws:
  access_key: <AWS_ACCCESS_KEY>
  secret_key: <AWS_SECRET_KEY>
  aws_region: eu-west-1
  s3_bucket: <AWS_S3_BUCKET>
  s3_key: Artifactory-backups

~~~

Where: -

- ```<MYSQL_PASSWORD``` is the MySQL password.
- ```<AWS_ACCCESS_KEY>``` is the AWS Access Key.
- ```<AWS_SECRET_KEY>``` is the AWS Secret Key.
- ```<AWS_S3_BUCKET>``` is the AWS S3 Bucket Name.

Notes: -

- If ```select_limit:``` is defined it'll limit the number of records return from the MySQL query.

## Repository list

File ```repos.txt``` should contain a list of Artifactory repositories, one per line.

## Go Modules

~~~
export GO111MODULE=on
~~~

## Go Build

~~~
go build
~~~

### For Linux, on Mac OS X

Cross Compile sqlite3 <https://github.com/mattn/go-sqlite3/issues/384>

~~~
brew install FiloSottile/musl-cross/musl-cross
~~~

~~~
CC=x86_64-linux-musl-gcc CXX=x86_64-linux-musl-g++ GOARCH=amd64 GOOS=linux CGO_ENABLED=1 go build -ldflags "-linkmode external -extldflags -static" -o artifactory-s3-export-linux
~~~

## File Info

### Mac OS X

~~~
file artifactory-s3-export
~~~

> Results

~~~
artifactory-s3-export: Mach-O 64-bit executable x86_64
~~~

### For Linux, on Mac OS X

~~~
file artifactory-s3-export-linux
~~~

> Results

~~~
artifactory-s3-export-linux: ELF 64-bit LSB executable, x86-64, version 1 (SYSV), statically linked, not stripped
~~~

## Run

~~~
./artifactory-s3-export | tee artifactory-s3-export.log
~~~

> Info log level

~~~
./artifactory-s3-export -stderrthreshold=INFO | tee artifactory-s3-export.log
~~~
