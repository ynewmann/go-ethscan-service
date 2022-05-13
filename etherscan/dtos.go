package etherscan

type ErrorPage struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Result  string `json:"result"`
}

type BlockPage struct {
	JsonRpc string   `json:"jsonrpc"`
	ID      uint64   `json:"id"`
	Result  BlockDto `json:"result"`
}

type TransactionDto struct {
	BlockHash        string `json:"blockHash"`
	BlockNumber      string `json:"blockNumber"`
	From             string `json:"from"`
	Gas              string `json:"gas"`
	GasPrice         string `json:"gasPrice"`
	Hash             string `json:"hash"`
	Input            string `json:"input"`
	Nonce            string `json:"nonce"`
	To               string `json:"to"`
	TransactionIndex string `json:"transactionIndex"`
	Value            string `json:"value"`
	Type             string `json:"type"`
	V                string `json:"v"`
	R                string `json:"r"`
	S                string `json:"s"`
}

type BlockDto struct {
	Difficulty       string           `json:"difficulty"`
	ExtraData        string           `json:"extraData"`
	GasLimit         string           `json:"gasLimit"`
	GasUsed          string           `json:"gasUsed"`
	Hash             string           `json:"hash"`
	LogsBloom        string           `json:"logsBloom"`
	Miner            string           `json:"miner"`
	MixHash          string           `json:"mixHash"`
	Nonce            string           `json:"nonce"`
	Number           string           `json:"number"`
	ParentHash       string           `json:"parentHash"`
	ReceiptsRoot     string           `json:"receiptsRoot"`
	Sha3Uncles       string           `json:"sha3Uncles"`
	Size             string           `json:"size"`
	StateRoot        string           `json:"stateRoot"`
	Timestamp        string           `json:"timestamp"`
	TotalDifficulty  string           `json:"totalDifficulty"`
	Transactions     []TransactionDto `json:"transactions"`
	TransactionsRoot string           `json:"transactionsRoot"`
	Uncles           []string         `json:"uncles"`
}
