package github

var (
	APIRoot = "https://api.github.com"
)

const (
	DeviceCodeURL   = "https://github.com/login/device/code"
	AccessTokenURL  = "https://github.com/login/oauth/access_token"
	DefaultClientID = "BuilderCLI"
	TokenDir        = ".builder"
	TokenFile       = "github.json"
)
