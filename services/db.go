package services

import (
	"context"
	"north-api/ent"
	"north-api/ent/migrate"

	"entgo.io/ent/dialect/sql"
	_ "github.com/go-sql-driver/mysql"
)

var Db *ent.Client

// OpenDb 打开数据库链接
func OpenDb(backend, dsn string) error {
	drv, err := sql.Open(backend, dsn)
	if err != nil {
		return err
	}
	// db := drv.DB()
	// db.SetMaxIdleConns(10)
	// db.SetMaxOpenConns(100)
	// db.SetConnMaxLifetime(time.Hour)
	Db = ent.NewClient(ent.Driver(drv))
	return nil
}

// CloseDb 关闭数据库链接
func CloseDb() error {
	return Db.Close()
}

// MigrateDb 数据库迁移
func MigrateDb() error {
	ctx := context.Background()
	err := Db.Schema.Create(
		ctx,
		migrate.WithDropIndex(true),
		migrate.WithDropColumn(true),
	)
	return err
}
