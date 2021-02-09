# connext-cli
[![PRs Welcome](https://img.shields.io/badge/PRs-welcome-brightgreen.svg)](#contributing)

A helper tool for easy operations with your active Connext server node, which uses some of the endpoints from the [Connext Node API](https://docs.connext.network/reference/nodeAPI/).
Download the binary from release, adapt the config and get started.

### Get started
1) Log on your VM where the active Connext server node is running
2) New folder ``mkdir connext`` and go in new directory``cd connext``
3) Download binary from release ``wget https://github.com/cwntr/connext-cli/releases/download/v1.0.0/cncli``
4) Make executable ``chmod +x cncli``
5) Adapt config: currently only supports an absolute fixed file path at `"cfg.json"` on same directory as the `cncli`

**example config**
 ```
 {
  "publicIdentifier":"vector8PtDJ7CYRDLG85Bdez9wbGyer8xdREVbEcvgHX4zqHugcFzEfY",
  "signerAddress":"0x2Ee926f87FBbeE0Ab083F2194E524fFC0fe4aaC0",
  "channelAddress":"0x4f50E45fAF15AA7eeF139f129ce1b8B0f7D4A998",
  "host":"http://localhost:8001",
  "graceNonceDiff": 5
 }
 ```
_graceNonceDiff_ : only used for deleting active transfers (pending / stuck), the channel nonce difference what you consider "old" transfers are safe to be deleted.

### Commands
```
COMMANDS:
   delete-active-transfers, dt  delete all active transfers for all assets
   get-active-transfers, t      gets all active transfers for current channel
   get-channel-state, s         gets all channels for public identifier and channel-id
   get-channels, c              gets all channels for public identifier
   help, h                      Shows a list of commands or help for one command
```

**example usage**

Running:
`` ./cncli get-channel-state``

Will output you the channel state with grouped balances
```
-- INFO (GET-CHANNEL-STATE)
 pub-id: vector8PtDJ7CYRDLG85Bdez9wbGyer8xdREVbEcvgHX3zqHugcFzEfY
 channel-address: 0x4f60E45fAF15AA8eeF139f129ce1b8B0f7D4A998
-- INFO

-- RAW
 {AssetIds:[0x0000000000000000000000000000000000000000 0xdac17f958d2ee523a2206206994597c13d831ec7] Balances:[{Amount:[14999680000000000 15000320000000000] To:[0xeDb1EBFaf2413b8C250bFDD0f58B0bcfeF54F980 0x2Ee926f87FBbeE0Ab083F4194E524fFC0fe4aaC0]} {Amount:[101955070 14522547] To:[0xeDb1EBFaf2413b8C250bFDD0f58B0bcfeF54F980 0x2Ee926f87FBbeE0Ab083F4194E524fFC0fe4aaC0]}] ChannelAddress:0x4f60E45fAF15AA8eeF139f129ce1b8B0f7D4A998}
-- RAW

-- BALANCES
 balances found: 2
 localUSDT: 14.522547 USDT
 remoteUSDT: 101.955070 USDT
 localETH: 0.015000319999999999 ETH
 remoteETH: 0.014999680000000000 ETH
-- BALANCES
```

Commands will be improved / new commands be added over time. Just a starting point for basic operations.