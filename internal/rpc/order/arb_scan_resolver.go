package order

func NewARBScanResolver(tag string, chainId int64) ScanResolver {
	return &ARBScanResolver{
		ChainId: chainId,
		Coin:    tag,
		EthScanResolver: EthScanResolver{
			ChainId: chainId,
			Coin:    tag,
		},
	}
}

type ARBScanResolver struct {
	EthScanResolver
	ChainId int64
	Coin    string
}

type BlockResultJson struct {
	Jsonrpc string `json:"jsonrpc"`
	ID      int    `json:"id"`
	Result  struct {
		Timestamp string `json:"timestamp"`
	}
}
