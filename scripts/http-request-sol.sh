
# curl https://api.devnet.solana.com -X POST -H "Content-Type: application/json" -d '
#   {
#     "jsonrpc": "2.0",
#     "id": 1,
#     "method": "getSignaturesForAddress",
#     "params": [
#       "CsZcvDrJZGpMu3Y36qK5wYkVEFUcMKEdF3BXqgDWGkSv",
#       {
#         "limit": 1
#       }
#     ]
#   }
# '
# tx sender account
# curl https://api.devnet.solana.com -X POST -H "Content-Type: application/json" -d '
#   {
#     "jsonrpc": "2.0",
#     "id": 1,
#     "method": "getTransaction",
#     "params": [
#       "46Q15x1R2FtTMU4nENCYA4zwVoVPiwsRWEimYwebVfSYDAirLaYkz28SHnshvY1uPDeYsmvSc6HSWofx6xmjhWV4",
#       "json"
#     ]
#   }
# '


# 9hZaTvCVMcfbheTzebkeGR6Xi2EzMqTtPasbhGoPB94C
# curl https://api.devnet.solana.com -X POST -H "Content-Type: application/json" -d '
#   {
#     "jsonrpc": "2.0",
#     "id": 1,
#     "method": "getSignaturesForAddress",
#     "params": [
#       "9hZaTvCVMcfbheTzebkeGR6Xi2EzMqTtPasbhGoPB94C",
#       {
#         "limit": 1
#       }
#     ]
#   }
# '

curl https://api.devnet.solana.com -X POST -H "Content-Type: application/json" -d '
  {
    "jsonrpc": "2.0",
    "id": 1,
    "method": "getTransaction",
    "params": [
      "3EDaSfApwCzkHcZdLBnMdDAyo9aVV9KaxCxSdmcMuJoq4sAoedb7ziHwBwBDe2jNxjnzZC5oAb9YFfGiHSs6taGu",
      "json"
    ]
  }
'