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
	Server
	SQS
}

func init() {
	err := envconfig.Process("", &Env)
	if err != nil {
		log.Fatalf("Fail to load env config : %v", err)
	}
}
