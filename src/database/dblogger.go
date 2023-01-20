package database

import (
	"context"
	"fmt"
	"log"
	"time"

	"gorm.io/gorm/logger"
)

type GormLogger struct{}

func (L *GormLogger) LogMode(l logger.LogLevel) logger.Interface {
	newLogger := *L

	return &newLogger
}
func (L *GormLogger) Info(ctx context.Context, format string, obj ...any) {
	log.Printf(fmt.Sprintf("[gorm info] %s\n", format), obj...)
}
func (L *GormLogger) Warn(ctx context.Context, format string, obj ...any) {
	log.Printf(fmt.Sprintf("[gorm warn] %s\n", format), obj...)
}
func (L *GormLogger) Error(ctx context.Context, format string, obj ...any) {
	log.Printf(fmt.Sprintf("[gorm error] %s\n", format), obj...)
}
func (L *GormLogger) Trace(
	ctx context.Context,
	begin time.Time,
	fc func() (sql string, rowsAffected int64),
	err error,
) {

	sql, rows := fc()

	trimmedSql := sql
	if len(trimmedSql) > 1000 {
		trimmedSql = trimmedSql[:1000] + "..."
	}
	log.Printf("[gorm trace] Query: %s", trimmedSql)
	log.Printf("[gorm trace] Affected rows: %d\n", rows)

	if err != nil {
		log.Printf("[gorm trace] Error: %q\n", err)
	}
}
