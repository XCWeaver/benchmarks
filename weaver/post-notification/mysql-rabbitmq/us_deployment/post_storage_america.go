package main

import (
	"context"
	"database/sql"
	"fmt"
	"sync"
	"time"

	"github.com/ServiceWeaver/weaver"
	_ "github.com/go-sql-driver/mysql"
)

// Post_storage_america component.
type Post_storage_america interface {
	GetPost(context.Context, Post_id_obj) (string, error)
	GetConsistencyWindowValues(ctx context.Context) ([]float64, error)
}

// Implementation of the Post_storage_america component.
type post_storage_america struct {
	weaver.Implements[Post_storage_america]
	weaver.WithConfig[post_storage_usOptions]
	db                      *sql.DB
	mu                      sync.Mutex
	consistencyWindowValues []float64
}

type Post_id_obj struct {
	weaver.AutoMarshal
	PostId    string
	WriteTime int64
}

type post_storage_usOptions struct {
	MysqlAddr      string `toml:"mysql_address"`
	MysqlPort      string `toml:"mysql_port"`
	MysqlPassword  string `toml:"mysql_password"`
	MysqlUser      string `toml:"mysql_user"`
	MysqlDatastore string `toml:"mysql_datastore"`
}

func (p *post_storage_america) Init(ctx context.Context) error {
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

func (p *post_storage_america) GetPost(ctx context.Context, postId Post_id_obj) (string, error) {
	logger := p.Logger(ctx)

	logger.Debug("Reading post!", "postId", postId.PostId)
	startTimeMs := time.Now().UnixMilli()

	// Query the database for the value associated with the key
	var post string
	query := "SELECT value FROM posts WHERE k = ?"
	err := p.db.QueryRow(query, postId.PostId).Scan(&post)
	if err == sql.ErrNoRows {
		inconsistencies.Inc()
		logger.Error("post not found")
		return "post not found", err
	} else if err != nil {
		return "", err
	}

	readPostDurationMs.Put(float64(time.Now().UnixMilli() - startTimeMs))
	consistencyWindowMs := float64(time.Now().UnixMilli() - postId.WriteTime)
	consistencyWindow.Put(consistencyWindowMs)
	p.mu.Lock()
	p.consistencyWindowValues = append(p.consistencyWindowValues, consistencyWindowMs)
	p.mu.Unlock()

	return post, nil
}

func (p *post_storage_america) GetConsistencyWindowValues(ctx context.Context) ([]float64, error) {
	logger := p.Logger(ctx)
	logger.Debug("entering GetConsistencyWindowValues")
	p.mu.Lock()
	values := p.consistencyWindowValues
	p.mu.Unlock()
	return values, nil
}
