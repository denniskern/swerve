package acm

// Config contains the ACM config
type Config struct {
	PebbleCA       string
	UsePebble      bool   `long:"use-pebble" env:"SWERVE_USE_PEBBLE" description:"If set to true, it will use a different http client with insecure settings"`
	UseStage       bool   `long:"use-stage" env:"SWERVE_USE_STAGE" description:"If set to true, it will use a different directory url for letsencrypt"`
	PebbleCAURL    string `long:"pebble-ca-url" env:"SWERVE_PEBBLE_CA_URL" description:"Pebble CA http url, e.g https://localhost:14000/dir"`
	LetsEncryptURL string `long:"letsencrypt-url" env:"SWERVE_LETSENCRYPT_URL" description:"Letsencrypt CA http url, e.g https://localhost:15000"`
}
