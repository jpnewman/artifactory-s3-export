package models

/*
S3Object - AWS S3 object table
*/
type S3Object struct {
	ID          int64  `db:"id,key,auto" sqlite:"INTEGER PRIMARY KEY AUTOINCREMENT"`
	NodeID      int64  `db:"node_id" sqlite:"bigint(20)"`
	Key         string `db:"key" sqlite:"VARCHAR(1024) NOT NULL UNIQUE"`
	Size        uint64 `db:"size" sqlite:"bigint(20)"`
	UploadError string `db:"upload_file_error" sqlite:"VARCHAR(1024)"`
}
