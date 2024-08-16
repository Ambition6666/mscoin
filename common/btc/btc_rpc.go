package btc

import (
	"common/tools"
	"encoding/json"
	"errors"
	"log"
)

var apiUrl = "http://127.0.0.1:18332"
var auth = "Basic Yml0Y29pbjoxMjM0NTY="

type ListUnspentInfo struct {
	Id     string              `json:"id"`
	Error  string              `json:"error"`
	Result []ListUnspentResult `json:"result"`
}

type ListUnspentResult struct {
	Txid          string  `json:"txid"`
	Vout          int     `json:"vout"`
	Address       string  `json:"address"`
	Amount        float64 `json:"amount"`
	Confirmations int     `json:"confirmations"`
}

func ListUnspent(min, max int, addresses []string) ([]ListUnspentResult, error) {
	params := make(map[string]any)
	params["id"] = "mscoin"
	params["method"] = "listunspent"
	params["jsonrpc"] = "1.0"
	params["params"] = []any{min, max, addresses}
	headers := make(map[string]string)
	headers["Authorization"] = auth
	bytes, err := tools.PostWithHeader(apiUrl, params, headers, "")
	if err != nil {
		log.Println(err)
		return nil, err
	}
	var result ListUnspentInfo
	json.Unmarshal(bytes, &result)
	if result.Error == "" {
		return result.Result, nil
	}
	return nil, errors.New(result.Error)
}

type Input struct {
	Txid string `json:"txid"`
	Vout int    `json:"vout"`
}
type CreateRawTransactionInfo struct {
	Id     string `json:"id"`
	Error  string `json:"error"`
	Result string `json:"result"`
}

// CreateRawTransaction values=[{"address":amount},{"address":amount}]
func CreateRawTransaction(inputs []Input, values []map[string]any) (string, error) {
	params := make(map[string]any)
	params["id"] = "mscoin"
	params["method"] = "createrawtransaction"
	params["jsonrpc"] = "1.0"
	params["params"] = []any{inputs, values}
	headers := make(map[string]string)
	headers["Authorization"] = auth
	bytes, err := tools.PostWithHeader(apiUrl, params, headers, "")
	if err != nil {
		log.Println(err)
		return "", err
	}
	var result CreateRawTransactionInfo
	json.Unmarshal(bytes, &result)
	if result.Error == "" {
		return result.Result, nil
	}
	return "", errors.New(result.Error)
}

type SignRawTransactionWithWalletInfo struct {
	Id     string                             `json:"id"`
	Error  string                             `json:"error"`
	Result SignRawTransactionWithWalletResult `json:"result"`
}
type SignRawTransactionWithWalletResult struct {
	Hex      string `json:"hex"`
	Complete bool   `json:"complete"`
}

// SignRawTransactionWithWallet values=[{"address":amount}]
func SignRawTransactionWithWallet(hexTxid string) (*SignRawTransactionWithWalletResult, error) {
	params := make(map[string]any)
	params["id"] = "mscoin"
	params["method"] = "signrawtransactionwithwallet"
	params["jsonrpc"] = "1.0"
	params["params"] = []any{hexTxid}
	headers := make(map[string]string)
	headers["Authorization"] = auth
	bytes, err := tools.PostWithHeader(apiUrl, params, headers, "")
	if err != nil {
		log.Println(err)
		return nil, err
	}
	var result SignRawTransactionWithWalletInfo
	json.Unmarshal(bytes, &result)
	if result.Error == "" {
		return &result.Result, nil
	}
	return nil, errors.New(result.Error)
}

type SendRawTransactionInfo struct {
	Id     string `json:"id"`
	Error  string `json:"error"`
	Result string `json:"result"`
}

// SendRawTransaction 将交易广播到网络中 返回txid
func SendRawTransaction(signHex string) (string, error) {
	params := make(map[string]any)
	params["id"] = "mscoin"
	params["method"] = "sendrawtransaction"
	params["jsonrpc"] = "1.0"
	params["params"] = []any{signHex, 0} //0代表任意手续费 前面创建交易的时候 一定要算好手续费
	headers := make(map[string]string)
	headers["Authorization"] = auth
	bytes, err := tools.PostWithHeader(apiUrl, params, headers, "")
	if err != nil {
		log.Println(err)
		return "", err
	}
	var result SendRawTransactionInfo
	json.Unmarshal(bytes, &result)
	if result.Error == "" {
		return result.Result, nil
	}
	return "", errors.New(result.Error)
}
