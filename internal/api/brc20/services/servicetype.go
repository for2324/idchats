package services

type BalanceListItem struct {
	Token            string `json:"token"`
	TokenType        string `json:"tokenType"`
	Balance          string `json:"balance"`
	AvailableBalance string `json:"availableBalance"`
	TransferBalance  string `json:"transferBalance"`
}
type TransferAbleInscript struct {
	InscriptionId     string `json:"inscriptionId"`
	Ticker            string `json:"ticker"`
	Amount            string `json:"amount"`
	InscriptionNumber int    `json:"inscriptionNumber"`
	UtxoHash          string `json:"utxoHash"`
	Vout              int    `json:"vout"`
}
type UnspendUtxo struct {
	TxId         string        `json:"txId"`
	OutputIndex  int           `json:"outputIndex"`
	Satoshis     int64         `json:"satoshis"`
	ScriptPk     string        `json:"scriptPk"`
	AddressType  int           `json:"addressType"`
	Inscriptions []interface{} `json:"inscriptions"`
}
