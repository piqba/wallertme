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
# Kafka
```