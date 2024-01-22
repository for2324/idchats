package order

func NewOPScanResolver(tag string, chainId int64) ScanResolver {
	return &OPScanResolver{
		ChainId: chainId,
		Coin:    tag,
		EthScanResolver: EthScanResolver{
			ChainId: chainId,
			Coin:    tag,
		},
	}
}

type OPScanResolver struct {
	EthScanResolver
	ChainId int64
	Coin    string
}
