package infra

import (
	"github.com/aws/aws-sdk-go/aws/ec2metadata"
	"github.com/aws/aws-sdk-go/aws/session"
)

type region struct {
	logger *Logger
	region *string
}

type Region interface {
	Get() *string
}

func NewRegion(logger *Logger) Region {
	return &region{
		logger: logger,
	}
}

func (r region) Get() *string {
	if r.region == nil {
		metaSession, err := session.NewSession()
		if err != nil {
			r.logger.Error().Err(err).Msg("Failed to new metaSession")
			return nil
		}
		metaClient := ec2metadata.New(metaSession)
		if !metaClient.Available() {
			r.logger.Error().Err(err).Msg("Unavailable ec2 metadata")
			return nil
		}
		region, err := metaClient.Region()
		if err != nil {
			r.logger.Error().Err(err).Msg("Failed to get region")
			return nil
		}
		r.region = &region
	}
	return r.region
}
