package alchemy

import "encoding/json"

type alchemyRequest interface {
	GetPaylod(id int64) []byte
}

type alchemyRequestBase struct {
	Id      int64         `json:"id"`
	Jsonrpc string        `json:"jsonrpc"`
	Method  string        `json:"method"`
	Params  []interface{} `json:"params"`
}

type alchemyResponseBase struct {
	Id      int64  `json:"id"`
	Jsonrpc string `json:"jsonrpc"`
}

var _ alchemyRequest = new(alchemyRequestBase)

func (a *alchemyRequestBase) GetPaylod(id int64) []byte {
	a.Id = id
	res, _ := json.Marshal(a)
	return res
}

func newAlchemyRequestBase(method string) alchemyRequestBase {
	return alchemyRequestBase{Id: -1, Method: method, Jsonrpc: "2.0"}
}
