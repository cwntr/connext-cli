package connext

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httputil"
	"time"
)

const (
	ReqHeaderContentType = "Content-Type"
	ReqHeaderJson        = "application/json"
)

var nodeURL string

func SetHost(h string) {
	nodeURL = h
}

func GetChannels(publicIdentifier string) ([]string, error) {
	resp, err := http.Get(fmt.Sprintf("%s/%s/channels", nodeURL, publicIdentifier))
	var resData []string
	defer resp.Body.Close()
	err = getJson(resp, &resData)
	return resData, err
}

func GetVectorChannel(publicIdentifier string, channelId string) (VectorChannel, error) {
	resp, err := http.Get(fmt.Sprintf("%s/%s/channels/%s", nodeURL, publicIdentifier, channelId))
	var resData VectorChannel
	defer resp.Body.Close()
	err = getJson(resp, &resData)
	return resData, err
}

func GetActiveTransfers(publicIdentifier string, channelId string) ([]ActiveTransfer, error) {
	resp, err := http.Get(fmt.Sprintf("%s/%s/channels/%s/active-transfers", nodeURL, publicIdentifier, channelId))
	var resData []ActiveTransfer
	defer resp.Body.Close()
	err = getJson(resp, &resData)
	return resData, err
}

const PreImageDelete = "0x0000000000000000000000000000000000000000000000000000000000000000"

func getHttpClient() *http.Client {
	tr := &http.Transport{
		MaxIdleConns:       10,
		IdleConnTimeout:    30 * time.Second,
		DisableCompression: true,
	}
	return &http.Client{Transport: tr}
}

var isDebug = false

type CancelTransferResponse struct {
	ChannelAddress string `json:"channelAddress"`
	TransferID     string `json:"transferId"`
}

func CancelTransfer(publicIdentifier string, channelId string, transferId string) error {
	c := getHttpClient()
	o := TransferResolve{}
	o.TransferID = transferId
	o.PublicIdentifier = publicIdentifier
	o.ChannelAddress = channelId
	o.TransferResolver.PreImage = PreImageDelete

	request, err := json.Marshal(o)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/transfers/resolve", nodeURL), bytes.NewBuffer(request))
	if err != nil {
		return err
	}

	req.Header.Add(ReqHeaderContentType, ReqHeaderJson)
	resp, err := c.Do(req)
	if err != nil {
		return err
	}
	if isDebug {
		fmt.Printf("[CancelTransfer] response status: %d \n", resp.StatusCode)
		dump, err := httputil.DumpResponse(resp, true)
		if err != nil {
			fmt.Printf("[CancelTransfer] err while dumping response, err: %v \n", err)
		}
		fmt.Printf("[CancelTransfer] response body: %s \n", string(dump))
	}
	var resData CancelTransferResponse
	err = getJson(resp, &resData)
	return err
}

type TransferResolve struct {
	PublicIdentifier string `json:"publicIdentifier"`
	ChannelAddress   string `json:"channelAddress"`
	TransferID       string `json:"transferId"`
	TransferResolver struct {
		PreImage string `json:"preImage"`
	} `json:"transferResolver"`
}

type TransferBalance struct {
	Amount []string `json:"amount"`
	To     []string `json:"to"`
}

type TransferState struct {
	Expiry   interface{} `json:"expiry,omitempty"`
	LockHash string      `json:"lockHash,omitempty"`
}

type ActiveTransfer struct {
	InDispute             bool            `json:"inDispute"`
	ChannelFactoryAddress string          `json:"channelFactoryAddress"`
	AssetID               string          `json:"assetId"`
	ChainID               int             `json:"chainId"`
	ChannelAddress        string          `json:"channelAddress"`
	Balance               TransferBalance `json:"balance"`
	Initiator             string          `json:"initiator"`
	Responder             string          `json:"responder"`
	InitialStateHash      string          `json:"initialStateHash"`
	TransferDefinition    string          `json:"transferDefinition"`
	InitiatorIdentifier   string          `json:"initiatorIdentifier"`
	ResponderIdentifier   string          `json:"responderIdentifier"`
	ChannelNonce          int             `json:"channelNonce"`
	TransferEncodings     []string        `json:"transferEncodings"`
	TransferID            string          `json:"transferId"`
	TransferState         TransferState   `json:"transferState,omitempty"`
	TransferTimeout       string          `json:"transferTimeout,omitempty"`
	Meta                  struct {
		SenderIdentifier string `json:"senderIdentifier,omitempty"`
		RequireOnline    bool   `json:"requireOnline,omitempty"`
		RoutingID        string `json:"routingId,omitempty"`
		Path             []struct {
			Recipient        string `json:"recipient,omitempty"`
			RecipientChainID int    `json:"recipientChainId,omitempty"`
			RecipientAssetID string `json:"recipientAssetId,omitempty"`
		} `json:"path"`
	} `json:"meta,omitempty"`
}

type VectorChannel struct {
	AssetIds []string `json:"assetIds"`
	Balances []struct {
		Amount []string `json:"amount"`
		To     []string `json:"to"`
	} `json:"balances"`
	ChannelAddress     string   `json:"channelAddress"`
	MerkleRoot         string   `json:"merkleRoot"`
	ProcessedDepositsA []string `json:"processedDepositsA"`
	ProcessedDepositsB []string `json:"processedDepositsB"`
	DefundNonces       []string `json:"defundNonces"`
	NetworkContext     struct {
		ChainID                 int    `json:"chainId"`
		ChannelFactoryAddress   string `json:"channelFactoryAddress"`
		TransferRegistryAddress string `json:"transferRegistryAddress"`
		ProviderURL             string `json:"providerUrl"`
	} `json:"networkContext"`
	Nonce           int    `json:"nonce"`
	Alice           string `json:"alice"`
	AliceIdentifier string `json:"aliceIdentifier"`
	Bob             string `json:"bob"`
	BobIdentifier   string `json:"bobIdentifier"`
	Timeout         string `json:"timeout"`
	LatestUpdate    struct {
		AssetID string `json:"assetId"`
		Balance struct {
			Amount []string `json:"amount"`
			To     []string `json:"to"`
		} `json:"balance"`
		ChannelAddress string `json:"channelAddress"`
		Details        struct {
			Balance struct {
				To     []string `json:"to"`
				Amount []string `json:"amount"`
			} `json:"balance"`
			MerkleProofData      []string `json:"merkleProofData"`
			MerkleRoot           string   `json:"merkleRoot"`
			TransferDefinition   string   `json:"transferDefinition"`
			TransferTimeout      string   `json:"transferTimeout"`
			TransferID           string   `json:"transferId"`
			TransferEncodings    []string `json:"transferEncodings"`
			TransferInitialState struct {
				LockHash string `json:"lockHash"`
				Expiry   string `json:"expiry"`
			} `json:"transferInitialState"`
			Meta struct {
				RoutingID     string `json:"routingId"`
				Hello         string `json:"hello"`
				RequireOnline bool   `json:"requireOnline"`
			} `json:"meta"`
		} `json:"details"`
		FromIdentifier string `json:"fromIdentifier"`
		Nonce          int    `json:"nonce"`
		AliceSignature string `json:"aliceSignature"`
		BobSignature   string `json:"bobSignature"`
		ToIdentifier   string `json:"toIdentifier"`
		Type           string `json:"type"`
	} `json:"latestUpdate"`
	InDispute bool `json:"inDispute"`
}

// getJson will parse the http response to the target struct
func getJson(response *http.Response, target interface{}) error {
	defer response.Body.Close()
	return json.NewDecoder(response.Body).Decode(target)
}
