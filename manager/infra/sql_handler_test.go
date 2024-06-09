package infra

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"touchgift-job-manager/config"
)

// WARNING: ローカルでテストする場合はrdsを起動してから実行してください
func TestNewSQLHandler(t *testing.T) {
	t.Run("正常な接続情報を記述した場合正常にDB接続が確立される", func(t *testing.T) {
		logger := GetLogger()
		handler := NewSQLHandler(logger)
		assert.NotNil(t, handler, "handler should not be nil")
	})
	t.Run("誤った設定を行った際にdbに接続できずpanicエラーを発生させる", func(t *testing.T) {
		logger := GetLogger()
		// 誤ったdb設定を記述
		config.Env.Db.User = "tester"
		config.Env.Db.Password = "password"
		config.Env.Db.Host = "localhost"
		config.Env.Db.Port = 5432
		config.Env.Db.Database = "test"
		config.Env.Db.ConnectTimeout = 30
		config.Env.DriverName = "pgx"

		defer func() {
			// recoverを使用してpanicをキャッチ
			if r := recover(); r == nil {
				t.Errorf("The code did not panic")
			} else {
				t.Logf("DB接続情報が誤っているためDB接続エラー発生: %v", r)
			}
		}()
		NewSQLHandler(logger)
	})
}
