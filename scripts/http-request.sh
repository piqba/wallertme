curl --location --request POST 'https://matic-mumbai.chainstacklabs.com/' \
--header 'Content-Type: application/json' \
--data-raw '{
	"jsonrpc":"2.0",
	"method":"eth_getTransactionByHash",
	"params":[
		"0x0314f52b94f624695e9035df6f76ba7c0209a57462ec6c9ade577523883fb681"
	],
	"id":1
}'