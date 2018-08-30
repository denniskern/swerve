package configuration

import (
	"flag"
	"os"
	"strings"
)

const (
	envPrefix = "SWERVE_"
)

// getOSPrefixEnv get os env
func getOSPrefixEnv(s string) *string {
	if e := strings.TrimSpace(os.Getenv(envPrefix + s)); len(e) > 0 {
		return &e
	}

	return nil
}

// FromEnv read the config from envs
func (c *Configuration) FromEnv() {
	if api := getOSPrefixEnv("API"); api != nil {
		c.APIListener = *api
	}

	if httpListener := getOSPrefixEnv("HTTP"); httpListener != nil {
		c.HTTPListener = *httpListener
	}

	if httpsListener := getOSPrefixEnv("HTTPS"); httpsListener != nil {
		c.HTTPSListener = *httpsListener
	}

	if dbEndpoint := getOSPrefixEnv("DB_ENDPOINT"); dbEndpoint != nil {
		c.DynamoDB.Endpoint = *dbEndpoint
	}

	if dbRegion := getOSPrefixEnv("DB_REGION"); dbRegion != nil {
		c.DynamoDB.Region = *dbRegion
	}

	if dbKey := getOSPrefixEnv("DB_KEY"); dbKey != nil {
		if dbSecret := getOSPrefixEnv("DB_SECRET"); dbSecret != nil {
			c.DynamoDB.Key = *dbKey
			c.DynamoDB.Secret = *dbSecret
		}
	}

	if dbBootstrap := getOSPrefixEnv("BOOTSTRAP"); dbBootstrap != nil {
		c.Bootstrap = len(*dbBootstrap) > 0 && *dbBootstrap != "0"
	}

	if logLevel := getOSPrefixEnv("LOG_LEVEL"); logLevel != nil {
		c.LogLevel = *logLevel
	}

	if logFormatter := getOSPrefixEnv("LOG_FORMATTER"); logFormatter != nil {
		c.LogFormatter = *logFormatter
	}
}

// FromParameter read config from application parameter
func (c *Configuration) FromParameter() {
	dbEndpointPtr := flag.String("db-endpoint", "", "DynamoDB endpoint (Required)")
	dbRegionPtr := flag.String("db-region", "", "DynamoDB region (Required)")
	dbKeyPtr := flag.String("db-key", "", "DynamoDB credential key")
	dbSecretPtr := flag.String("db-secret", "", "DynamoDB credential secret")
	dbBootstrapPtr := flag.Bool("bootstrap", false, "Bootstrap the database")

	logLevelPtr := flag.String("log-level", "", "Set the log level (info,debug,warning,error,fatal,panic)")
	logFormatterPtr := flag.String("log-formatter", "", "Set the log formatter (text,json)")

	flag.Parse()

	if dbEndpointPtr != nil && *dbEndpointPtr != "" {
		c.DynamoDB.Endpoint = *dbEndpointPtr
	}

	if dbRegionPtr != nil && *dbRegionPtr != "" {
		c.DynamoDB.Region = *dbRegionPtr
	}

	if dbKeyPtr != nil && dbSecretPtr != nil && *dbKeyPtr != "" && *dbSecretPtr != "" {
		c.DynamoDB.Key = *dbKeyPtr
		c.DynamoDB.Secret = *dbSecretPtr
	}

	if dbBootstrapPtr != nil && *dbBootstrapPtr {
		c.Bootstrap = *dbBootstrapPtr
	}

	if logLevelPtr != nil && *logLevelPtr != "" {
		c.LogLevel = *logLevelPtr
	}

	if logFormatterPtr != nil && *logFormatterPtr != "" {
		c.LogFormatter = *logFormatterPtr
	}
}

// NewConfiguration creates a new instance
func NewConfiguration() *Configuration {
	return &Configuration{
		HTTPListener:  ":8080",
		HTTPSListener: ":8081",
		APIListener:   ":8082",
		LogFormatter:  "text",
		LogLevel:      "debug",
		Bootstrap:     false,
	}
}
