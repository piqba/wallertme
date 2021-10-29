# wallertme
wallet alert me is an app off-chain. It's a space where the user can register their wallet and can receive notifications of their transactions with other wallets. Our roadmap you can find here [project board](https://piqba.notion.site/piqba/d0682ea1e7e748969aab454b1a437237?v=e09ed390c7c945bf9288891178f3637a)



## How to interact with app for testing reasson


For the Solana blockchain, we need to have a wallet. How you can get a wallet?

Go to [phatom wallet](https://phantom.app/) install and select Devnet network.

Then if you are new in this ecosystem like me, I recomend you read this posts [solana cli tool](https://docs.solana.com/cli/install-solana-cli-tools) and then follow this steps [Send and Receive Tokens](https://docs.solana.com/cli/transfer-tokens) for testing propourse.

For the Cardano blockchain, we need to have a wallet. How you can get a wallet?

Go to [nami wallet](https://namiwallet.io/) install and select CardanoTestNet network.

Then if you are new in this ecosystem like me, go to [cardano faucets](https://developers.cardano.org/docs/integrate-cardano/testnet-faucet/) and get ADA for testing propourse.



Then when you are completed the previous steps, follow this bash descriptions



### Download walletmectl

Go to last release and download [walletmectl](https://github.com/piqba/wallertme/releases) binary for you OS


### Setup .env file


```bash

touch .env

# put into this file
# Database settings optional at the moment:
DB_SERVER_URL="host=localhost port=5432 user=xxxx password=xxxx dbname=xxxx sslmode=disable"
DB_MAX_CONNECTIONS=100
DB_MAX_IDLE_CONNECTIONS=10
DB_MAX_LIFETIME_CONNECTIONS=2
# REDIS requiered very important
REDIS_URI=localhost:6379
REDIS_PASS=""

# notification providers  very important
SMTP_EMAIL_RECEIVER=test@gmail.com
SMTP_EMAIL_USER=test@gmail.com
SMTP_EMAIL_PASSWORD="password"
BOT_TOKEN=telegram:token


```

### Create wallets.json file

```bash

# its requiered to have the following file wallets.json

touch wallets.json
vim wallets.json
...

# paste this format of json data
# Important this wallets can be found on testnet(cardano) and devnet(solana)

[
    {
        "address": "addr_test1qq6g6s99g9z9w0mlvew28w40lpml9rwfkfgerpkg6g2vpn6dp4cf7k9drrdy0wslarr6hxspcw8ev5ed8lfrmaengneqz34lcx",
        "symbol": "ADA",
        "is_active": true,
        "notifier_service": [
            {
                "name": "telegram",
                "user_id": "xx"
            }
        ],
        "network_type": "CardanoTestNet"
    },
    {
        "address": "addr_test1qq5287luxzj5l4lequrqdp5ln76ver4uls3z0m5ykr5gqsv0vxzrwcq5dmmn9e09rvgttzgrngmpxkguy7220r0u0ljqzuww7g",
        "symbol": "ADA",
        "is_active": true,
        "notifier_service": [
            {
                "name": "telegram",
                "user_id": "xx"
            }
        ],
        "network_type": "CardanoTestNet"
    },
    {
        "address": "9hZaTvCVMcfbheTzebkeGR6Xi2EzMqTtPasbhGoPB94C",
        "symbol": "SOL",
        "is_active": true,
        "notifier_service": [
            {
                "name": "telegram",
                "user_id": "xxx"
            }
        ],
        "network_type": "SolanaDevNet"
    }
]

```
### Install redis for local enviroment

```bash
docker run --name redis -e ALLOW_EMPTY_PASSWORD=yes quay.io/bitnami/redis:latest

```
### Usage CLI walletmectl

```bash
Wallertme ctl is a tool focused on: 
	Send tx data from (SOLANA|CARDANO) blockchain to a queue like (REDIS) streams and then send this information
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
### BB8 module

```bash
Publish Txs data from (SOLANA|CARDANO) blockchains to (REDIS)

Usage:
  walletmectl bb8 [flags]

Flags:
      --source string          select a wallets data source from (json|db) (default "json")
  -h, --help                   help for bb8
      --timer string           select a time duration to watch all txs (default "1s")
      --wallets::name string   select the name of wallet.json file (default "wallets.json")
      --wallets::path string   select the path of wallet.json file (default "/path/<bin file>")
      --watcher                select true|false if you want to run this task periodicaly


```

### R2D2 module
```bash
Subscribe to Txs data topic from (REDIS) and send notifications to this services (telegram|discord|smtp)

Usage:
  walletmectl r2d2 [flags]

Flags:
      --group::name string     select a name for your consumer group
  -h, --help              help for r2d2
      --source string          select a wallets data source from (json|db) (default "json")
      --wallets::name string   select the name of wallet.json file (default "wallets.json")
      --wallets::path string   select the path of wallet.json file (default "/path/<bin file>")

```


### ScreenShots notifications


Telegram notification

![R2D2](/docs/assets/tg.png)

