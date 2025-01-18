package config

import (
	"flag"
	"fmt"
	"os"

	"github.com/bsakweson/cqrs-core/pkg/constants"
	"github.com/bsakweson/cqrs-core/pkg/elasticsearch"
	"github.com/bsakweson/cqrs-core/pkg/eventstroredb"
	"github.com/bsakweson/cqrs-core/pkg/logger"
	"github.com/bsakweson/cqrs-core/pkg/mongodb"
	"github.com/bsakweson/cqrs-core/pkg/probes"
	"github.com/bsakweson/cqrs-core/pkg/tracing"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

var configPath string

func init() {
	flag.StringVar(&configPath, "config", "", "ES microservice config path")
}

type Config struct {
	ServiceName      string                         `mapstructure:"serviceName"`
	Logger           *logger.Config                 `mapstructure:"logger"`
	GRPC             GRPC                           `mapstructure:"grpc"`
	Mongo            *mongodb.Config                `mapstructure:"mongo"`
	MongoCollections MongoCollections               `mapstructure:"mongoCollections"`
	Probes           probes.Config                  `mapstructure:"probes"`
	Jaeger           *tracing.Config                `mapstructure:"jaeger"`
	EventStoreConfig eventstroredb.EventStoreConfig `mapstructure:"eventStoreConfig"`
	Subscriptions    Subscriptions                  `mapstructure:"subscriptions"`
	Elastic          elasticsearch.Config           `mapstructure:"elastic"`
	ElasticIndexes   ElasticIndexes                 `mapstructure:"elasticIndexes"`
	Http             Http                           `mapstructure:"http"`
}

type GRPC struct {
	Port        string `mapstructure:"port"`
	Development bool   `mapstructure:"development"`
}

type MongoCollections struct {
	Products string `mapstructure:"products" validate:"required"`
}

type Subscriptions struct {
	PoolSize                   int    `mapstructure:"poolSize" validate:"required,gte=0"`
	ProductPrefix              string `mapstructure:"productPrefix" validate:"required,gte=0"`
	MongoProjectionGroupName   string `mapstructure:"mongoProjectionGroupName" validate:"required,gte=0"`
	ElasticProjectionGroupName string `mapstructure:"elasticProjectionGroupName" validate:"required,gte=0"`
}

type ElasticIndexes struct {
	Products string `mapstructure:"products" validate:"required"`
}

type Http struct {
	Port                string   `mapstructure:"port" validate:"required"`
	Development         bool     `mapstructure:"development"`
	BasePath            string   `mapstructure:"basePath" validate:"required"`
	ProductsPath        string   `mapstructure:"productsPath" validate:"required"`
	DebugErrorsResponse bool     `mapstructure:"debugErrorsResponse"`
	IgnoreLogUrls       []string `mapstructure:"ignoreLogUrls"`
}

func InitConfig() (*Config, error) {
	if configPath == "" {
		configPathFromEnv := os.Getenv(constants.ConfigPath)
		if configPathFromEnv != "" {
			configPath = configPathFromEnv
		} else {
			getwd, err := os.Getwd()
			if err != nil {
				return nil, errors.Wrap(err, "os.Getwd")
			}
			configPath = fmt.Sprintf("%s/config/config.yaml", getwd)
		}
	}

	cfg := &Config{}

	viper.SetConfigType(constants.Yaml)
	viper.SetConfigFile(configPath)

	if err := viper.ReadInConfig(); err != nil {
		return nil, errors.Wrap(err, "viper.ReadInConfig")
	}

	if err := viper.Unmarshal(cfg); err != nil {
		return nil, errors.Wrap(err, "viper.Unmarshal")
	}

	grpcPort := os.Getenv(constants.GrpcPort)
	if grpcPort != "" {
		cfg.GRPC.Port = grpcPort
	}
	mongoURI := os.Getenv(constants.MongoDbURI)
	if mongoURI != "" {
		//cfg.Mongo.URI = "mongodb://host.docker.internal:27017"
		cfg.Mongo.URI = mongoURI
	}
	jaegerAddr := os.Getenv(constants.JaegerHostPort)
	if jaegerAddr != "" {
		cfg.Jaeger.HostPort = jaegerAddr
	}
	eventStoreConnectionString := os.Getenv(constants.EventStoreConnectionString)
	if eventStoreConnectionString != "" {
		cfg.EventStoreConfig.ConnectionString = eventStoreConnectionString
	}
	elasticUrl := os.Getenv(constants.ElasticUrl)
	if elasticUrl != "" {
		cfg.Elastic.URL = elasticUrl
	}

	return cfg, nil
}
