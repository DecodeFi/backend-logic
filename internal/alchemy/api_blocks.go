package alchemy

import (
	"encoding/json"
)

// API: eth_blockNumber

func NewBlockNumberRequest() alchemyRequest {
	req := newAlchemyRequestBase("eth_blockNumber")
	return &req
}

type BlockNumberResponse struct {
	alchemyResponseBase
	Number string `json:"result"`
}

func (a *AlchemyClient) BlockNumber(req alchemyRequest) (*BlockNumberResponse, error) {
	json_resp, err := a.httpRequest(req)
	if err != nil {
		return nil, err
	}

	var res *BlockNumberResponse = nil
	err = json.Unmarshal(json_resp, &res)
	return res, err
}

// API: eth_getBlockByNumber

func NewBlockRequest(number string) alchemyRequest {
	req := newAlchemyRequestBase("eth_getBlockByNumber")
	req.Params = append(req.Params, number)
	req.Params = append(req.Params, false)
	return &req
}

type BlockResponse struct {
	alchemyResponseBase
	Result struct {
		Hash         string   `json:"hash"`
		Transactions []string `json:"transactions"`
	} `json:"result"`
}

func (a *AlchemyClient) Block(req alchemyRequest) (*BlockResponse, error) {
	json_resp, err := a.httpRequest(req)
	if err != nil {
		return nil, err
	}

	var res *BlockResponse = nil
	err = json.Unmarshal(json_resp, &res)
	return res, err
}
