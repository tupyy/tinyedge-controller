package configuration

import (
	"os"
	"time"

	"github.com/spf13/viper"
)

const (
	prefix = "TINYEDGE"

	defaultCertificateTTL = 3600 * 24 * 365 * time.Second
)

type Configuration struct {
	BaseDomain            string
	DefaultCertificateTTL time.Duration
	VaultAddress          string
	VaultApproleRoleID    string
	VaultAppRoleSecretID  string
	VaultSecretMountPath  string
	PostgresAddress       string
	PostgresUser          string
	PostgresPassword      string
	PostgresDB            string
}

func GetConfiguration() Configuration {
	return Configuration{
		BaseDomain:            "home.net",
		DefaultCertificateTTL: defaultCertificateTTL,
		VaultAddress:          "http://localhost:8200",
		VaultApproleRoleID:    "app-role-id",
		VaultSecretMountPath:  "tinyedge",
		VaultAppRoleSecretID:  os.Getenv("VAULT_SECRET_ID"),
		PostgresAddress:       "localhost:5432",
		PostgresUser:          "postgres",
		PostgresPassword:      "postgres",
		PostgresDB:            "tinyedge",
	}
}

// func ParseConfiguration(confFile string) error {
// 	// setDefaults()

// 	// viper.SetEnvPrefix(prefix)
// 	// viper.AutomaticEnv() // read in environment variables that match

// 	// if len(confFile) == 0 {
// 	// 	return errors.New("no config file specified")
// 	// }

// 	// content, err := os.ReadFile(confFile)
// 	// if err != nil {
// 	// 	return err
// 	// }
// 	// if err := json.Unmarshal(content, &configuration); err != nil {
// 	// 	return err
// 	// }

// 	// zap.S().Infow("configuration read", "file", confFile)
// 	// return nil
// }

func setDefaults() {
	viper.SetDefault("defaultCertificateTTL", defaultCertificateTTL)
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
