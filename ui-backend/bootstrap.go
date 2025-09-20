package main

import (
	"context"
	"os"

	"github.com/tommzn/go-config"
	"github.com/tommzn/go-log"
	"github.com/tommzn/go-secrets"
)

func bootstrap() log.Logger {

	secretsManager := newSecretsManager()
	conf := loadConfig()
	ctx := context.Background()
	logger := newLogger(conf, secretsManager, ctx)
	return logger
}

func loadConfig() config.Config {

	configSource, err := config.NewS3ConfigSourceFromEnv()
	if err != nil {
		panic(err)
	}

	conf, err := configSource.Load()
	if err != nil {
		panic(err)
	}
	return conf
}

func newSecretsManager() secrets.SecretsManager {
	secretsManager := secrets.NewDockerecretsManager("/run/secrets/token")
	return secretsManager
}

func newLogger(conf config.Config, secretsMenager secrets.SecretsManager, ctx context.Context) log.Logger {
	logger := log.NewLoggerFromConfig(conf, secretsMenager)
	logContextValues := make(map[string]string)
	logContextValues[log.LogCtxNamespace] = "utte-universe-ui"
	if node, ok := os.LookupEnv("K8S_NODE_NAME"); ok {
		logContextValues[log.LogCtxK8sNode] = node
	}
	if pod, ok := os.LookupEnv("K8S_POD_NAME"); ok {
		logContextValues[log.LogCtxK8sPod] = pod
	}
	logger.WithContext(log.LogContextWithValues(ctx, logContextValues))
	return logger
}
