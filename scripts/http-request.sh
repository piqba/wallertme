curl --location --request POST 'localhost:8545/' \
--header 'Content-Type: application/json' \
--data-raw '{
	"jsonrpc":"2.0",
	"method":"eth_getBlockByNumber",
	"params":[
		"latest",
		true
	],
	"id":1
}'