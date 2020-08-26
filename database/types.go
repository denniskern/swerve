package database

// Database is the API to the database
type Database struct {
	Service DynamoDBAdapter
	Config  Config
}

// Config contains the databases config
type Config struct {
	Key            string `long:"dyno-aws-key" env:"SWERVE_DYNO_AWS_KEY" default:"0" required:"false" description:"AWS access key for dynamodb"`
	Secret         string `long:"dyno-aws-sec" env:"SWERVE_DYNO_AWS_SECRET" default:"0" required:"false" description:"AWS secret key for dynamodb"`
	Region         string `long:"dyno-aws-region" env:"SWERVE_DYNO_AWS_REGION" required:"false" description:"AWS region for dynamodb" default:"eu-central-1"`
	Endpoint       string `long:"dyno-endpoint" env:"SWERVE_DYNO_ENDPOINT" required:"false" description:"Endpoint of dynamodb" default:"http://localhost:8000"`
	DefaultUserPW  string `long:"dyno-default-user-pw" env:"SWERVE_DYNO_DEFAULT_PW" required:"false" description:"Default PW for the admin user"`
	DefaultUser    string `long:"dyno-default-user" env:"SWERVE_DYNO_DEFAULT_USER" required:"false" description:"Default PW for the admin user" default:"admin"`
	TableRedirects string `long:"dyno-tbl-redirects" env:"SWERVE_DYNO_TABLE_REDIRECTS" required:"false" description:"Table name for redirects" default:"Swerve_Redirects"`
	TableCertCache string `long:"dyno-tbl-certcache" env:"SWERVE_DYNO_TABLE_CERTCACHE" required:"false" description:"Table name for cert cache" default:"Swerve_CertCache"`
	TableUsers     string `long:"dyno-tbl-users" env:"SWERVE_DYNO_TABLE_USERS" required:"false" description:"Table name for users" default:"Swerve_Users"`
	Bootstrap      bool   `long:"dyno-bootstrap" env:"SWERVE_DYNO_BOOTSTRAP" required:"false" description:"Create tables and default user on startup"`
}

// Redirect is the redirect entry model
type Redirect struct {
	RedirectFrom string  `json:"redirect_from"`
	Description  string  `json:"description"`
	RedirectTo   string  `json:"redirect_to"`
	Promotable   bool    `json:"promotable"`
	Code         int     `json:"code"`
	Created      int     `json:"created"`
	Modified     int     `json:"modified"`
	CPathMaps    *[]byte `json:"cpath_map,omitempty"`
}

// CertCacheEntry contains a certificate for the domain
type CertCacheEntry struct {
	Key   string `json:"domain"`
	Value string `json:"cert"`
}

// User contains a users cerdentials
type User struct {
	Name string `json:"username"`
	Pwd  string `json:"pwd"`
}
