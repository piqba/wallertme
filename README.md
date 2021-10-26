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



# its requiered to have the following file wallets.json

touch wallets.json
vim wallets.json
...

# paste this format of json data 
# Important this wallets can be found on testnet(cardano) and devnet(solana)

[
    {
        "address": "addr_test1qq6g6s99g9z9w0mlvew28w40lpml9rwfkfgerpkg6g2vpn6dp4cf7k9drrdy0wslarr6hxspcw8ev5ed8lfrmaengneqz34lcx",
        "lastTx": "5ef1187f5e125090675a3c2d2d2cee359aaf6941df625db598ec996ab1011f55",
        "symbol": "ADA"
    },
    {
        "address": "addr_test1qq5287luxzj5l4lequrqdp5ln76ver4uls3z0m5ykr5gqsv0vxzrwcq5dmmn9e09rvgttzgrngmpxkguy7220r0u0ljqzuww7g",
        "lastTx": "5ef1187f5e125090675a3c2d2d2cee359aaf6941df625db598ec996ab1011f55",
        "symbol": "ADA"
    },
    {
        "address": "9hZaTvCVMcfbheTzebkeGR6Xi2EzMqTtPasbhGoPB94C",
        "lastTx": "3EDaSfApwCzkHcZdLBnMdDAyo9aVV9KaxCxSdmcMuJoq4sAoedb7ziHwBwBDe2jNxjnzZC5oAb9YFfGiHSs6taGu",
        "symbol": "SOL"
    }
]

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
TELEGRAM_USER_ID=telegramID(only a number)

```