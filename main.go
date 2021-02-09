package main

import (
	"encoding/json"
	"fmt"
	"github.com/cwntr/connext-cli/connext"
	"github.com/urfave/cli/v2"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
)

const (
	CmdGetChannels                 = "get-channels"
	CmdGetChannelState             = "get-channel-state"
	CmdGetActiveTransfers          = "get-active-transfers"
	CmdGetDeleteAllActiveTransfers = "delete-active-transfers"

	ConfigFile = "cfg.json"
)

var (
	app *cli.App
	cfg connext.Config
)

func info() {
	app.Name = "Connext CLI"
	app.Usage = "Connext HTTP API wrapper"
	app.Version = "1.0.0"
}

func main() {
	err := readConfig()
	if err != nil {
		fmt.Printf("unable to read cfg.json file, err:%v\n", err)
		panic("cannot parse config")
	}
	connext.SetHost(cfg.Host)

	app = &cli.App{}
	info()
	commands()

	sort.Sort(cli.FlagsByName(app.Flags))
	sort.Sort(cli.CommandsByName(app.Commands))

	err = app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func commands() {
	app.Commands = []*cli.Command{
		{
			Name:    CmdGetChannels,
			Aliases: []string{"c"},
			Usage:   "gets all channels for public identifier",
			Action: func(c *cli.Context) error {
				fmt.Printf("-- INFO (%s) \n", strings.ToUpper(CmdGetChannels))
				fmt.Printf(" pub-id: %s  \n", cfg.PublicIdentifier)
				fmt.Printf("-- INFO \n\n")

				ch, err := connext.GetChannels(cfg.PublicIdentifier)
				if err != nil {
					fmt.Printf("[%s] err: %v", CmdGetChannels, err)
					return nil
				}

				fmt.Printf("-- CHANNELS \n")
				for i, c := range ch {
					fmt.Printf(" #%d [pub-id: %s] [channel: %s] \n", i+1, cfg.PublicIdentifier, c)
				}
				fmt.Printf("-- CHANNELS \n")

				return nil
			},
		},
		{
			Name:    CmdGetChannelState,
			Aliases: []string{"s"},
			Usage:   "gets all channels for public identifier and channel-id",
			Action: func(c *cli.Context) error {
				fmt.Printf("-- INFO (%s) \n", strings.ToUpper(CmdGetChannelState))
				fmt.Printf(" pub-id: %s \n", cfg.PublicIdentifier)
				fmt.Printf(" channel-address: %s \n", cfg.ChannelAddress)
				fmt.Printf("-- INFO \n\n")
				ch, err := connext.GetVectorChannel(cfg.PublicIdentifier, cfg.ChannelAddress)
				if err != nil {
					fmt.Printf("[%s] err: %v", CmdGetChannelState, err)
					return nil
				}
				fmt.Printf("-- RAW \n")
				fmt.Printf(" %+v\n", ch)
				fmt.Printf("-- RAW \n\n")

				var localUSDT int64 = 0
				var remoteUSDT int64 = 0
				var localETH int64 = 0
				var remoteETH int64 = 0
				for i, r := range ch.Balances {
					if ch.AssetIds[i] == connext.AssetUSDT {
						if r.To[0] == cfg.SignerAddress {
							amt, _ := strconv.ParseInt(r.Amount[0], 10, 64)
							localUSDT += amt
							amt, _ = strconv.ParseInt(r.Amount[1], 10, 64)
							remoteUSDT += amt
						} else {
							amt, _ := strconv.ParseInt(r.Amount[1], 10, 64)
							localUSDT += amt
							amt, _ = strconv.ParseInt(r.Amount[0], 10, 64)
							remoteUSDT += amt
						}
					} else if ch.AssetIds[i] == connext.AssetETH {
						if r.To[0] == cfg.SignerAddress {
							amt, _ := strconv.ParseInt(r.Amount[0], 10, 64)
							localETH += amt
							amt, _ = strconv.ParseInt(r.Amount[1], 10, 64)
							remoteETH += amt
						} else {
							amt, _ := strconv.ParseInt(r.Amount[1], 10, 64)
							localETH += amt
							amt, _ = strconv.ParseInt(r.Amount[0], 10, 64)
							remoteETH += amt
						}
					}
				}
				fmt.Printf("-- BALANCES \n")
				fmt.Printf(" balances found: %d \n", len(ch.Balances))
				fmt.Printf(" localUSDT: %.6f USDT \n", float64(localUSDT)/1e6)
				fmt.Printf(" remoteUSDT: %.6f USDT \n", float64(remoteUSDT)/1e6)
				fmt.Printf(" localETH: %.18f ETH \n", float64(localETH)/1e18)
				fmt.Printf(" remoteETH: %.18f ETH \n", float64(remoteETH)/1e18)
				fmt.Printf("-- BALANCES \n")
				return nil
			},
		},
		{
			Name:    CmdGetActiveTransfers,
			Aliases: []string{"t"},
			Usage:   "gets all active transfers for current channel",
			Action: func(c *cli.Context) error {
				fmt.Printf("-- INFO (%s) \n", strings.ToUpper(CmdGetActiveTransfers))
				fmt.Printf(" pub-id: %s \n", cfg.PublicIdentifier)
				fmt.Printf(" channel-address: %s \n", cfg.ChannelAddress)
				fmt.Printf("-- INFO \n\n")
				transfers, err := connext.GetActiveTransfers(cfg.PublicIdentifier, cfg.ChannelAddress)
				if err != nil {
					fmt.Printf("[%s] err: %v", CmdGetActiveTransfers, err)
					return nil
				}

				fmt.Printf("-- RAW \n")
				fmt.Printf(" %+v\n", transfers)
				fmt.Printf("-- RAW \n\n")

				var incomingUSDT int64 = 0
				var outgoingUSDT int64 = 0
				var incomingETH int64 = 0
				var outgoingETH int64 = 0
				for _, r := range transfers {
					if r.AssetID == connext.AssetUSDT {
						if r.Balance.To[0] == cfg.PublicIdentifier {
							amt, _ := strconv.ParseInt(r.Balance.Amount[0], 10, 64)
							incomingUSDT += amt
						} else {
							amt, _ := strconv.ParseInt(r.Balance.Amount[0], 10, 64)
							outgoingUSDT += amt
						}
					} else if r.AssetID == connext.AssetETH {
						if r.Balance.To[0] == cfg.PublicIdentifier {
							amt, _ := strconv.ParseInt(r.Balance.Amount[0], 10, 64)
							incomingETH += amt
						} else {
							amt, _ := strconv.ParseInt(r.Balance.Amount[0], 10, 64)
							outgoingETH += amt
						}
					}
				}
				fmt.Printf("-- TRANSFERS \n")
				fmt.Printf(" active transfers: %d \n", len(transfers))
				fmt.Printf(" incomingUSDT: %.6f USDT \n", float64(incomingUSDT)/1e6)
				fmt.Printf(" outgoingUSDT: %.6f USDT \n", float64(outgoingUSDT)/1e6)
				fmt.Printf(" incomingETH: %.18f ETH \n", float64(incomingETH)/1e18)
				fmt.Printf(" outgoingETH: %.18f ETH \n", float64(outgoingETH)/1e18)
				fmt.Printf("-- TRANSFERS \n")
				return nil
			},
		},
		{
			Name:    CmdGetDeleteAllActiveTransfers,
			Aliases: []string{"dt"},
			Usage:   "delete all active transfers for all assets",
			Action: func(c *cli.Context) error {
				myPubId := cfg.PublicIdentifier
				myChanAddr := cfg.ChannelAddress

				fmt.Printf("-- INFO (%s) \n", strings.ToUpper(CmdGetDeleteAllActiveTransfers))
				fmt.Printf(" pub-id: %s \n", myPubId)
				fmt.Printf(" channel-address: %s \n", myChanAddr)

				transfers, err := connext.GetActiveTransfers(myPubId, myChanAddr)
				if err != nil {
					fmt.Printf("[%s] err: %v", CmdGetDeleteAllActiveTransfers, err)
					return nil
				}
				fmt.Printf(" grace-nonce-difference: %d (if current channel nonce - transfer nonce > this value, will cancel it)", cfg.GraceNonceDiff)
				fmt.Printf(" active transfers: %d \n", len(transfers))
				fmt.Printf("-- INFO \n\n")

				channelInfo, err := connext.GetVectorChannel(myPubId, myChanAddr)
				if err != nil {
					fmt.Printf("[%s] error while connext.GetVectorChannel, err: %v", CmdGetDeleteAllActiveTransfers, err)
					return nil
				}

				fmt.Printf("-- CANCELLING TRANFERS \n")
				for i, c := range transfers {
					asset := getAssetByAssetId(c.AssetID)
					if channelInfo.Nonce-c.ChannelNonce < cfg.GraceNonceDiff {
						fmt.Printf(" #%d [%s] skipping transfer due to nonce check [transferId: %s] [transfer-nonce: %d] [chan-nonce:%d] \n", i+1, asset, c.TransferID, c.ChannelNonce, channelInfo.Nonce)
						continue
					}

					if c.Meta.SenderIdentifier == "" {
						fmt.Printf(" #%d [%s] delete start NO META DATA [transferId: %s] [transfer-nonce: %d] [chan-nonce:%d] \n", i+1, asset, c.TransferID, c.ChannelNonce, channelInfo.Nonce)
						err := connext.CancelTransfer(myPubId, myChanAddr, c.TransferID)
						if err != nil {
							fmt.Printf(" #%d [%s] delete error [transferId: %s] err: %+v \n", i+1, asset, c.TransferID, err)
							continue
						}
					}
					if len(c.Meta.Path) > 0 && c.Meta.Path[0].Recipient == myPubId {
						fmt.Printf(" #%d [%s] delete start [transferId: %s] [transfer-nonce: %d] [chan-nonce:%d] \n", i+1, asset, c.TransferID, c.ChannelNonce, channelInfo.Nonce)
						err := connext.CancelTransfer(myPubId, myChanAddr, c.TransferID)
						if err != nil {
							fmt.Printf(" #%d [%s] delete error [transferId: %s] err: %+v \n", i+1, asset, c.TransferID, err)
						} else {
							fmt.Printf(" #%d [%s] delete end [transferId: %s] \n", i+1, asset, c.TransferID)
						}
					}
				}
				fmt.Printf("-- CANCELLING TRANFERS \n")
				return nil
			},
		},
	}
}

func getAssetByAssetId(assetId string) string {
	switch assetId {
	case connext.AssetUSDT:
		return connext.CurrencyUSDT
	case connext.AssetETH:
		return connext.CurrencyETH
	default:
		return "unknown-asset"
	}
}

// Reads the entire config, "cfg.json" is hardcoded and must be placed on same level as the application binary
func readConfig() error {
	file, err := os.Open(ConfigFile)
	if err != nil {
		fmt.Printf("can't open config file: %v", err)
		return err
	}
	defer file.Close()
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&cfg)
	if err != nil {
		fmt.Printf("can't decode config JSON: %v", err)
		return err
	}
	return nil
}
