package main

import (
	"bufio"
	"database/sql"
	"flag"
	"fmt"
	"log"
	"os"
	"path"
	"regexp"
	"strings"
	"time"

	awsHelper "github.com/jpnewman/artifactory-s3-export/aws"
	configHelper "github.com/jpnewman/artifactory-s3-export/config"
	dbHelper "github.com/jpnewman/artifactory-s3-export/dbs"
	"github.com/jpnewman/artifactory-s3-export/models"
	"github.com/samonzeweb/godb"

	"github.com/dustin/go-humanize"
	"github.com/golang/glog"

	"github.com/spf13/viper"

	"github.com/aws/aws-sdk-go/aws/session"
)

func usage() {
	usage := `
Usage: -

  artifactory-s3-export
    -stderrthreshold=[INFO|WARN|FATAL]
    -log_dir=[string]
    -updateS3Table
    -dryrun

Arguments: -

`

	fmt.Fprintf(os.Stderr, usage)
	flag.PrintDefaults()
	os.Exit(2)
}

func init() {
}

func queryPackages(args *Args, db *sql.DB, sqliteDb *godb.DB, repo string) {
	glog.Info("Querying Packages...")

	// results, err := db.Query("SELECT node_id, node_type, repo, node_path, node_name, depth, created, created_by, modified, modified_by, updated, bin_length, sha1_actual, sha1_original, md5_actual, md5_original, sha256, repo_path_checksum FROM artdb.nodes WHERE node_name <> '.' AND node_path = '.' AND node_type = 1 AND repo = ? LIMIT 1000;", repo)
	// results, err := db.Query("SELECT node_id, node_type, repo, node_path, node_name, depth, created, created_by, modified, modified_by, updated, bin_length, sha1_actual, sha1_original, md5_actual, md5_original, sha256, repo_path_checksum FROM artdb.nodes WHERE node_name <> '.' AND node_path = '.' AND repo = ? LIMIT 1000;", repo)
	// results, err := db.Query("SELECT node_id, node_type, repo, node_path, node_name, depth, created, created_by, modified, modified_by, updated, bin_length, sha1_actual, sha1_original, md5_actual, md5_original, sha256, repo_path_checksum FROM artdb.nodes WHERE node_name <> '.' AND node_type = 1 AND repo = ? LIMIT 1000;", repo)

	var results *sql.Rows
	var err error
	if viper.IsSet("mysql.select_limit") {
		results, err = db.Query("SELECT node_id, node_type, repo, node_path, node_name, depth, created, created_by, modified, modified_by, updated, bin_length, sha1_actual, sha1_original, md5_actual, md5_original, sha256, repo_path_checksum FROM artdb.nodes WHERE node_name REGEXP '\\.nupkg$' AND node_type = 1 AND repo = ? LIMIT ?;", repo, viper.GetInt64("mysql.select_limit"))
	} else {
		results, err = db.Query("SELECT node_id, node_type, repo, node_path, node_name, depth, created, created_by, modified, modified_by, updated, bin_length, sha1_actual, sha1_original, md5_actual, md5_original, sha256, repo_path_checksum FROM artdb.nodes WHERE node_name REGEXP '\\.nupkg$' AND node_type = 1 AND repo = ?;", repo)
	}
	defer results.Close()

	if err != nil {
		panic(err.Error()) // TODO: Handling error.
	}

	for results.Next() {
		var node models.Node
		err = results.Scan(&node.NodeID,
			&node.NodeType,
			&node.Repo,
			&node.NodePath,
			&node.NodeName,
			&node.Depth,
			&node.Created,
			&node.CreatedBy,
			&node.Modified,
			&node.ModifiedBy,
			&node.Updated,
			&node.BinLength,
			&node.Sha1Actual,
			&node.Sha1Original,
			&node.Md5Actual,
			&node.Md5Original,
			&node.Sha256,
			&node.RepoPathChecksum)

		if err != nil {
			panic(err.Error()) // TODO: Handling error.
		}

		if *args.dryRun == false {
			dbHelper.InsertOrUpdate(sqliteDb, &node)
		}
	}
}

func uploadPackages(args *Args, sqliteDb *godb.DB, awsSession *session.Session, repo string) {
	glog.Info("Uploading Packages...")

	nodes := make([]models.Node, 0, 0)
	err := sqliteDb.SelectFrom("Node").
		Columns("node_id",
			"node_type",
			"repo",
			"node_path",
			"node_name",
			"depth",
			"created",
			"created_by",
			"modified",
			"modified_by",
			"updated",
			"bin_length",
			"sha1_actual",
			"sha1_original",
			"md5_actual",
			"md5_original",
			"sha256",
			"repo_path_checksum",
			"repo_file_path",
			"repo_file_size",
			"repo_file_error").
		Where("repo = ?", repo).
		OrderBy("node_id").
		Do(&nodes)

	if err != nil {
		panic(err.Error()) // TODO: Handling error.
	}

	i := 0
	c := uint64(0)
	for _, node := range nodes {
		i++
		c += node.BinLength
		if node.Sha1Actual.Valid == true {
			node.RepoFilePath = fmt.Sprintf("%s/%s", node.Sha1Actual.String[0:2], node.Sha1Actual.String)

			node.RepoFileSize, err = getFileStats(path.Join(viper.GetString("repo.filestore_path"), node.RepoFilePath))
			if err != nil {
				node.RepoFileError = err.Error()
			} else {
				node.RepoFileError = ""
			}

			if *args.dryRun == false {
				dbHelper.InsertOrUpdate(sqliteDb, &node)
			}

			glog.Infof("  %d: %s - %s/%s (%s)\n", i, node.NodeName, node.Sha1Actual.String[0:2], node.Sha1Actual.String, humanize.Bytes(node.BinLength))
		} else {
			glog.Infof("  %d: %s - NULL/NULL (%s)\n", i, node.NodeName, humanize.Bytes(node.BinLength))
		}

		if node.BinLength == node.RepoFileSize {
			filePath := path.Join(viper.GetString("repo.filestore_path"), node.RepoFilePath)
			s3Key := path.Join(viper.GetString("aws.s3_key"), node.Repo, node.NodeName)

			var s3Obj models.S3Object
			selectErr := sqliteDb.Select(&s3Obj).
				Where("key = ?", s3Key).
				Do()

			if selectErr == sql.ErrNoRows {
				if *args.dryRun == false {
					s3Obj, err = awsHelper.UploadFileToS3(awsSession, filePath, s3Key)
					if err != nil {
						panic(err.Error())
					}
				}
			} else if err != nil {
				panic(err.Error())
			}

			s3Obj.NodeID = node.NodeID
			if *args.dryRun == false {
				dbHelper.InsertOrUpdate(sqliteDb, &s3Obj)
			}
		}
	}

	fmt.Printf("%s (%d) [%s]\n", repo, i, humanize.Bytes(c))
}

func getFileStats(filepath string) (uint64, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		return 0, err
	}

	return uint64(stat.Size()), nil
}

func main() {
	start := time.Now()

	fmt.Println("Export Artifactory to S3...")

	flag.Usage = usage
	args := parseArgs()

	configHelper.LoadConfig("config")

	db := dbHelper.InitMySQLDb(viper.GetString("mysql.connection_string"))
	defer db.Close()

	awsSession := awsHelper.InitAWSSession()
	sqliteDb := dbHelper.InitSqliteDb("./node.db")

	var node models.Node
	dbHelper.CreateTable(sqliteDb, &node)
	var s3Obj models.S3Object
	dbHelper.CreateTable(sqliteDb, &s3Obj)

	defer sqliteDb.Close()

	file, err := os.Open(viper.GetString("repo.list_file"))
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	var lineRegEx = regexp.MustCompile(`(^#.*$|^\s*$)`)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		repo := strings.TrimSpace(scanner.Text())
		if len(repo) == 0 || lineRegEx.MatchString(repo) {
			continue
		}

		glog.Info(repo)

		if *args.updateS3Table && *args.dryRun == false {
			awsHelper.GetS3Objects(awsSession, sqliteDb, repo)
		}

		queryPackages(&args, db, sqliteDb, repo)
		uploadPackages(&args, sqliteDb, awsSession, repo)
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	elapsed := time.Since(start)
	fmt.Printf("Elapse Time %s\n", elapsed)
}
