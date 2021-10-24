# wallertme
wallet alert me


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
BOT_TOKEN=token
DISCORD_WEBHOOK=https://discordapp.com/api/webhooks/idhook

```