// Copyright 2018 Axel Springer SE
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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

	if clientStaticPath := getOSPrefixEnv("API_CLIENT_STATIC"); clientStaticPath != nil {
		c.APIClientStaticPath = *clientStaticPath
	}

	if caStagingEnv := getOSPrefixEnv("STAGING"); caStagingEnv != nil {
		c.StagingCA = len(*caStagingEnv) > 0 && *caStagingEnv != "0"
	}
}

// FromParameter read config from application parameter
func (c *Configuration) FromParameter() {
	dbEndpointPtr := flag.String("db-endpoint", "", "DynamoDB endpoint (Required)")
	dbRegionPtr := flag.String("db-region", "", "DynamoDB region (Required)")
	dbKeyPtr := flag.String("db-key", "", "DynamoDB credential key")
	dbSecretPtr := flag.String("db-secret", "", "DynamoDB credential secret")
	dbBootstrapPtr := flag.Bool("bootstrap", false, "Bootstrap the database")

	caStagingEnvPtr := flag.Bool("staging", false, "ca manager will connect the CA staging environment")

	logLevelPtr := flag.String("log-level", "", "Set the log level (info,debug,warning,error,fatal,panic)")
	logFormatterPtr := flag.String("log-formatter", "", "Set the log formatter (text,json)")

	httpListenerPtr := flag.String("http", "", "Set the http listener address")
	httpsListenerPtr := flag.String("https", "", "Set the https listener address")
	apiListenerPtr := flag.String("api", "", "Set the API listener address")

	staticPathPtr := flag.String("client-static", "", "Set the path to api client static files")
	versionPtr := flag.Bool("version", false, "Print the version of the application")
	helpPtr := flag.Bool("help", false, "Print the default usage help dialog")

	flag.Parse()

	if versionPtr != nil && *versionPtr {
		c.Version = true
	}

	if helpPtr != nil && *helpPtr {
		c.Help = true
	}

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

	if caStagingEnvPtr != nil && *caStagingEnvPtr {
		c.StagingCA = *caStagingEnvPtr
	}

	if logLevelPtr != nil && *logLevelPtr != "" {
		c.LogLevel = *logLevelPtr
	}

	if logFormatterPtr != nil && *logFormatterPtr != "" {
		c.LogFormatter = *logFormatterPtr
	}

	if httpListenerPtr != nil && *httpListenerPtr != "" {
		c.HTTPListener = *httpListenerPtr
	}

	if httpsListenerPtr != nil && *httpsListenerPtr != "" {
		c.HTTPSListener = *httpsListenerPtr
	}

	if staticPathPtr != nil && *staticPathPtr != "" {
		c.APIClientStaticPath = *staticPathPtr
	}

	if apiListenerPtr != nil && *apiListenerPtr != "" {
		c.APIListener = *apiListenerPtr
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
