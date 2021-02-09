package connext

// Config struct holds mandatory information how to connect to Connext and configs for operations
type Config struct {
	PublicIdentifier string `json:"publicIdentifier"`
	SignerAddress    string `json:"signerAddress"`
	ChannelAddress   string `json:"channelAddress"`
	Host             string `json:"host"`
	GraceNonceDiff   int    `json:"graceNonceDiff"`
}

const (
	AssetETH  = "0x0000000000000000000000000000000000000000"
	AssetUSDT = "0xdac17f958d2ee523a2206206994597c13d831ec7"

	CurrencyETH  = "ETH"
	CurrencyUSDT = "USDT"
)
