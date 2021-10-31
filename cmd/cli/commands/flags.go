package commands

const (
	flagConsumerGroup = "wallets::groupname"
	flagSource        = "wallets::source"
	flagTimer         = "wallets::timer"
	flagWalletsName   = "wallets::name"
	flagWalletsPath   = "wallets::path"
	flagWatcher       = "wallets::watcher"
)

var (
	BuildTime    string
	Version      string
	VersionHash  string
	OtelNameBb8  = "bb8"
	OtelNameR2D2 = "r2d2"
	OtelVersion  = "v0.4.0"
	OtelNameEnv  = "dev"
)
