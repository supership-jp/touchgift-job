package config

import (
	"github.com/kelseyhightower/envconfig"
	"log"
	"time"
)

type App struct {
	LogLevel string `envconfig:"LOG_LEVEL" default:"info"`
}

type Server struct {
	Port            string        `envconfig:"PORT" default:"8080"`
	AdminPort       string        `envconfig:"ADMIN_PORT" default:"8081"`
	ShutdownTimeout time.Duration `envconfig:"SHUTDOWN_TIMEOUT" default:"1m"`
	MetricsPath     string        `envconfig:"METRICS_PATH" default:"/metrics"`
}

type Db struct {
	DriverName      string        `envconfig:"DB_DRIVER_NAME" default:"mysql"`
	User            string        `envconfig:"DB_USER" default:"touchgift"`
	Password        string        `envconfig:"DB_PASSWORD" default:"test"`
	Database        string        `envconfig:"DB_DATABASE" default:"touchgift"`
	Host            string        `envconfig:"DB_HOST" default:"localhost"`
	Port            int           `envconfig:"DB_PORT" default:"3306"`
	ConnectTimeout  int           `envconfig:"DB_CONNECT_TIMEOUT_SEC" default:"60"`
	MaxOpenConns    int           `envconfig:"DB_MAX_OPEN_CONNS" default:"5"`
	MaxIdleConns    int           `envconfig:"DB_MAX_IDLE_CONNS" default:"5"`
	ConnMaxLifetime time.Duration `envconfig:"DB_CONN_MAX_LIFETIME" default:"1h"`
}

type DeliveryStart struct {
	TaskInterval       time.Duration `envconfig:"DELIVERY_START_TASK_INTERVAL" default:"1m"`
	TaskLimit          int           `envconfig:"DELIVERY_START_WORKER_TASK_LIMIT" default:"10"` // 1回のSQLで取得する数
	NumberOfConcurrent int           `envconfig:"DELIVERY_START_WORKER_NUMBER_OF_CONCURRENT" default:"5"`
	NumberOfQueue      int           `envconfig:"DELIVERY_START_WORKER_NUMBER_OF_QUEUE" default:"5"`
}

type DynamoDB struct {
	EndPoint                string `envconfig:"DYNAMODB_ENDPOINT" default:"http://localhost:4566"` // デフォルトはローカル用
	TableNamePrefix         string `envconfig:"TABLE_NAME_PREFIX" default:""`                      // デフォルトはローカル/CI用
	DeliveryTableName       string `envconfig:"CAMPAIGN_TABLE_NAME" default:"campaign_data"`
	CreativeTableName       string `envconfig:"TOUCH_POINT_TABLE_NAME" default:"touch_point_data"`
	DeliveryBudgetTableName string `envconfig:"CONTENTS_TABLE_NAME" default:"contents_data"`
}

type SQS struct {
	Region                    string `envconfig:"AWS_REGION" default:"us-east-1"`
	EndPoint                  string `envconfig:"SQS_ENDPOINT" default:"http://localhost:4566"` // デフォルトはローカル用
	DeliveryOperationQueueURL string `envconfig:"SQS_DELIVERY_OPERATION_QUEUE_URL" default:"http://localhost:4566/000000000000/delivery-touchgift-operation"`
	DeliveryControlQueueURL   string `envconfig:"SQS_DELIVERY_CONTROL_QUEUE_URL" default:"http://localhost:4566/000000000000/delivery-touchgift-control"`
	VisibilityTimeoutSeconds  int64  `envconfig:"SQS_VISIBILITY_TIMEOUT_SECONDS" default:"60"` // 取得したメッセージを処理する時間(これを過ぎると別のアプリがメッセージを取得してしまう)
	WaitTimeSeconds           int64  `envconfig:"SQS_WAIT_TIME_SECONDS" default:"20"`          // SQSからメッセージを取得する待ち時間
	MaxMessages               int64  `envconfig:"SQS_MAX_MESSAGES" default:"10"`               // 一度に取得するメッセージ数
}

var Env = EnvConfig{}

type EnvConfig struct {
	RegionFromEC2Metadata bool `envconfig:"REGION_FROM_EC2METADATA" default:"false"`
	App
	DeliveryStart
	Server
	SQS
	Db
	DynamoDB
}

func init() {
	err := envconfig.Process("", &Env)
	if err != nil {
		log.Fatalf("Fail to load env config : %v", err)
	}
}
