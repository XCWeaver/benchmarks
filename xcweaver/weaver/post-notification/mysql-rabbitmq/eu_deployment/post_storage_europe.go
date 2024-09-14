package main

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/ServiceWeaver/weaver"
	_ "github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
)

// Post_storage component.
type Post_storage_europe interface {
	Post(context.Context, string) (Post_id_obj, error)
}

type Post_id_obj struct {
	weaver.AutoMarshal
	PostId    string
	WriteTime int64
}

// Implementation of the Post_storage component.
type post_storage_europe struct {
	weaver.Implements[Post_storage_europe]
	weaver.WithConfig[post_storage_europeOptions]
	db *sql.DB
}

type post_storage_europeOptions struct {
	MysqlAddr      string `toml:"mysql_address"`
	MysqlPort      string `toml:"mysql_port"`
	MysqlPassword  string `toml:"mysql_password"`
	MysqlUser      string `toml:"mysql_user"`
	MysqlDatastore string `toml:"mysql_datastore"`
}

func (p *post_storage_europe) Init(ctx context.Context) error {
	logger := p.Logger(ctx)

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", p.Config().MysqlUser, p.Config().MysqlPassword, p.Config().MysqlAddr, p.Config().MysqlPort, p.Config().MysqlDatastore)
	var err error
	p.db, err = sql.Open("mysql", dsn)
	if err != nil {
		return err
	}
	logger.Info("post storage service at eu running!", "mysql host", p.Config().MysqlAddr, "port", p.Config().MysqlPort, "user", p.Config().MysqlUser, "password", p.Config().MysqlPassword, "datastore", p.Config().MysqlDatastore)

	return nil
}

func (p *post_storage_europe) Post(ctx context.Context, post string) (Post_id_obj, error) {
	logger := p.Logger(ctx)

	id := uuid.New()

	writeStartTimeMs := time.Now().UnixMilli()

	// Prepare the statement and execute the query
	query := "INSERT INTO posts VALUES (?, ?)"
	stmt, err := p.db.Prepare(query)
	if err != nil {
		return Post_id_obj{}, err
	}
	defer stmt.Close()

	_, err = stmt.Exec(id.String(), post)
	writePostDurationMs.Put(float64(time.Now().UnixMilli() - writeStartTimeMs))
	if err != nil {
		logger.Error("Error writing post!", "msg", err.Error())
		return Post_id_obj{}, err
	}

	logger.Debug("Post successfully stored!", "postId", id.String(), "post", post)

	return Post_id_obj{PostId: id.String(), WriteTime: writeStartTimeMs}, nil
}
