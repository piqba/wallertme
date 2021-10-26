# wallertme
wallet alert me

## Usage 

```bash
Wallertme ctl is a tool focused on: 
	Send tx data from (SOLANA|CARDANO) blockchain to a queue like (REDIS|KAFKA) and then send this information
	to DISCORD|TELEGRAM|SMTP

Usage:
  walletmectl [command]

Available Commands:
  bb8         Publish Txs data from (SOLANA|CARDANO) blockchains to (REDIS)
  completion  generate the autocompletion script for the specified shell
  help        Help about any command
  r2d2        Subscribe to Txs data topic from (REDIS) and send notifications to this services (telegram|discord|smtp)
  version     Print the version number of wallertmectl

Flags:
  -h, --help   help for walletmectl

Use "walletmectl [command] --help" for more information about a command.

```
## BB8 module

```bash
Publish Txs data from (SOLANA|CARDANO) blockchains to (REDIS)

Usage:
  walletmectl bb8 [flags]

Flags:
      --exporter string        select a exporter to send data (default "redis")
  -h, --help                   help for bb8
      --timer string           select a time duration to watch all txs (default "1s")
      --wallets::name string   select the name of wallet.json file (default "wallets.json")
      --wallets::path string   select the path of wallet.json file (default "/path/<bin file>")
      --watcher                select true|false if you want to run this task periodicaly

```

## R2D2 module
```bash
Subscribe to Txs data topic from (REDIS) and send notifications to this services (telegram|discord|smtp)

Usage:
  walletmectl r2d2 [flags]

Flags:
  -h, --help              help for r2d2
      --notifier string   select a provider to send notifications (default "telegram")

```
## setup env file

```bash

touch .env

# put into this file
# Database settings:
DB_SERVER_URL="host=localhost port=5432 user=xxxx password=xxxx dbname=xxxx sslmode=disable"
DB_MAX_CONNECTIONS=100
DB_MAX_IDLE_CONNECTIONS=10
DB_MAX_LIFETIME_CONNECTIONS=2
# REDIS
REDIS_URI=localhost:6379
REDIS_PASS=""
# KAFKA_HOST="localhost:9092"
# notification providers
SMTP_EMAIL_RECEIVER=test@gmail.com
SMTP_EMAIL_USER=test@gmail.com
SMTP_EMAIL_PASSWORD="password"
BOT_TOKEN=telegram:token
DISCORD_WEBHOOK=https://discordapp.com/api/webhooks/idhook

```