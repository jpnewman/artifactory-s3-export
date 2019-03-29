package models

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

/*
Node - Artifactory node table
*/
type Node struct {
	NodeID           int64          `db:"node_id,key" sqlite:"bigint(20) PRIMARY KEY"`
	NodeType         bool           `db:"node_type" sqlite:"tinyint(4)"`
	Repo             string         `db:"repo" sqlite:"varchar(64)"`
	NodePath         string         `db:"node_path" sqlite:"varchar(1024)"`
	NodeName         string         `db:"node_name" sqlite:"varchar(255)"`
	Depth            int            `db:"depth" sqlite:"tinyint(4)"`
	Created          int64          `db:"created" sqlite:"bigint(20)"`
	CreatedBy        sql.NullString `db:"created_by" sqlite:"varchar(64)"`
	Modified         int64          `db:"modified" sqlite:"bigint(20)"`
	ModifiedBy       sql.NullString `db:"modified_by" sqlite:"varchar(64)"`
	Updated          int64          `db:"updated" sqlite:"bigint(20)"`
	BinLength        uint64         `db:"bin_length" sqlite:"bigint(20)"`
	Sha1Actual       sql.NullString `db:"sha1_actual" sqlite:"bigint(40)"`
	Sha1Original     sql.NullString `db:"sha1_original" sqlite:"bigint(1024)"`
	Md5Actual        sql.NullString `db:"md5_actual" sqlite:"bigint(32)"`
	Md5Original      sql.NullString `db:"md5_original" sqlite:"bigint(1024)"`
	Sha256           sql.NullString `db:"sha256" sqlite:"bigint(64)"`
	RepoPathChecksum sql.NullString `db:"repo_path_checksum" sqlite:"bigint(40)"`
	Uploaded         bool           `db:"uploaded" sqlite:"tinyint(4)"`
	UploadError      string         `db:"upload_error" sqlite:"bigint(1024)"`
	RepoFilePath     string         `db:"repo_file_path" sqlite:"bigint(1024)"`
	RepoFileSize     uint64         `db:"repo_file_size" sqlite:"bigint(20)"`
	RepoFileError    string         `db:"repo_file_error" sqlite:"bigint(1024)"`
}
