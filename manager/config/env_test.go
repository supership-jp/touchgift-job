package config

import (
	"github.com/kelseyhightower/envconfig"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
	"time"
)

func TestEnvConfig(t *testing.T) {
	t.Run("環境変数をセットし、EnvConfigが期待通りに設定を読み込むかをテストする", func(t *testing.T) {
		os.Setenv("DB_DRIVER_NAME", "postgres")
		os.Setenv("DB_USER", "testuser")
		os.Setenv("DB_PASSWORD", "testpassword")
		os.Setenv("DB_DATABASE", "testdb")
		os.Setenv("DB_HOST", "localhost")
		os.Setenv("DB_PORT", "5432")
		os.Setenv("DB_CONNECT_TIMEOUT_SEC", "30")
		os.Setenv("DB_MAX_OPEN_CONNS", "10")
		os.Setenv("DB_MAX_IDLE_CONNS", "10")
		os.Setenv("DB_CONN_MAX_LIFETIME", "30m")

		// 環境変数の読み込み
		var env EnvConfig
		err := envconfig.Process("", &env)
		assert.Nil(t, err)

		// 結果の検証
		assert.Equal(t, "postgres", env.Db.DriverName)
		assert.Equal(t, "testuser", env.Db.User)
		assert.Equal(t, "testpassword", env.Db.Password)
		assert.Equal(t, "testdb", env.Db.Database)
		assert.Equal(t, "localhost", env.Db.Host)
		assert.Equal(t, 5432, env.Db.Port)
		assert.Equal(t, 30, env.Db.ConnectTimeout)
		assert.Equal(t, 10, env.Db.MaxOpenConns)
		assert.Equal(t, 10, env.Db.MaxIdleConns)
		assert.Equal(t, time.Duration(30)*time.Minute, env.Db.ConnMaxLifetime)

		// 環境変数のクリーンアップ
		os.Unsetenv("DB_DRIVER_NAME")
		os.Unsetenv("DB_USER")
		os.Unsetenv("DB_PASSWORD")
		os.Unsetenv("DB_DATABASE")
		os.Unsetenv("DB_HOST")
		os.Unsetenv("DB_PORT")
		os.Unsetenv("DB_CONNECT_TIMEOUT_SEC")
		os.Unsetenv("DB_MAX_OPEN_CONNS")
		os.Unsetenv("DB_MAX_IDLE_CONNS")
		os.Unsetenv("DB_CONN_MAX_LIFETIME")
	})
	t.Run("環境変数を設定しなければdefaultの値を返すことを確認する", func(t *testing.T) {
		var env EnvConfig
		err := envconfig.Process("", &env)
		assert.Nil(t, err)

		assert.Equal(t, "mysql", env.Db.DriverName)
		assert.Equal(t, "touchgift", env.Db.User)
		assert.Equal(t, "test", env.Db.Password)
		assert.Equal(t, "touchgift-db", env.Db.Database)
		assert.Equal(t, "localhost", env.Db.Host)
		assert.Equal(t, 3306, env.Db.Port)
		assert.Equal(t, 60, env.Db.ConnectTimeout)
		assert.Equal(t, 5, env.Db.MaxOpenConns)
		assert.Equal(t, 5, env.Db.MaxIdleConns)
		assert.Equal(t, time.Duration(60)*time.Minute, env.Db.ConnMaxLifetime)

	})

}
