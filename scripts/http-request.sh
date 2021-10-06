curl --location --request POST 'https://matic-mumbai.chainstacklabs.com/' \
--header 'Content-Type: application/json' \
--data-raw '{
	"jsonrpc":"2.0",
	"method":"net_version",
	"params":[],
	"id":67
}'