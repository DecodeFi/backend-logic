package evm_inspect

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	_ "embed"
)

type EvmInspectClient struct {
	host string
	port int
}

func NewEvmInspectClient() *EvmInspectClient {
	return &EvmInspectClient{
		host: "http://127.0.0.1",
		port: 3000,
	}
}

type Trace struct {
	TraceId     string `json:"trace_id,omitempty"`
	BlockNumber string `json:"block_number,omitempty"`
	TxHash      string `json:"tx_hash"`
	FromAddr    string `json:"from_addr"`
	ToAddr      string `json:"to_addr"`
	StorageAddr string `json:"storage_addr"`
	Calldata    string `json:"calldata"`
	Value       string `json:"value"`
	Action      string `json:"action"`
}

func (e *EvmInspectClient) httpRequest(method string, payload string) ([]byte, error) {
	url := fmt.Sprintf("%s:%d/%s", e.host, e.port, method)

	fmt.Println(url)

	reader := strings.NewReader(payload)

	req, err := http.NewRequest("GET", url, reader)
	if err != nil {
		return nil, err
	}

	req.Header.Add("accept", "application/json")
	req.Header.Add("content-type", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)

	return body, nil
}

/// go:embed mock.json
// var json_resp []byte

func (e *EvmInspectClient) TraceBlock(id string) ([]Trace, error) {

	json_resp, err := e.httpRequest(fmt.Sprintf("trace_block/%s", id), "")

	if err != nil {
		return nil, err
	}

	var res []Trace = nil
	err = json.Unmarshal(json_resp, &res)

	return res, err
}
