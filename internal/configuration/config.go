package configuration

import (
	"time"

	"github.com/cristalhq/aconfig"
)

const (
	prefix = "TINYEDGE"
)

type Configuration struct {
	BaseDomain            string `default:"home.net" usage:"base domain"`
	DefaultCertificateTTL int64  `default:"31536000"`
	VaultAddress          string `usage:"vault address"`
	VaultApproleRoleID    string `default:"app-role-id"`
	VaultAppRoleSecretID  string
	VaultSecretMountPath  string `default:"tinyedge"`
	PostgresAddress       string `default:"localhost:5432"`
	PostgresUser          string `default:"postgres"`
	PostgresPassword      string `default:"postgres"`
	PostgresDB            string `default:"tinyedge"`
}

func (c Configuration) GetCertificateTTL() time.Duration {
	return time.Duration(c.DefaultCertificateTTL) * time.Second
}

func GetConfiguration() Configuration {
	var cfg Configuration
	loader := aconfig.LoaderFor(&cfg, aconfig.Config{
		SkipFiles: true,
		EnvPrefix: prefix,
	})
	if err := loader.Load(); err != nil {
		panic(err)
	}
	return cfg
}

// func GetLogLevel() zapcore.Level {
// 	switch configuration.LogLevel {
// 	case "TRACE":
// 		return zap.DebugLevel
// 	case "DEBUG":
// 		return zap.DebugLevel
// 	case "INFO":
// 		return zap.InfoLevel
// 	case "WARN":
// 		return zap.WarnLevel
// 	case "ERROR":
// 		return zap.ErrorLevel
// 	}

// 	return zap.DebugLevel
// }
