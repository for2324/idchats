#!/bin/bash
abigen --abi bibotreward.abi --pkg contract  --type BBTTradeReward --out bibotreward.go
abigen --abi bbtpledgepool.abi --pkg contract --type BBTPledgePool  --out bbtpledgepool.go
abigen --abi blppledgepool.abi --pkg contract --type BLPPledgePool  --out blppledgepool.go
