package alchemy

import (
	"io"
	"net/http"
	"strings"
	"sync/atomic"
)

type AlchemyClient struct {
	gateway string
	apiKey  string

	id atomic.Int64
}

func NewAlchemyClient() *AlchemyClient {
	return &AlchemyClient{
		gateway: "https://eth-mainnet.g.alchemy.com/v2/",
		apiKey:  "_mSXeLBO_vTQ9KMdPQ0LkRryUY0sdvtM",
	}
}

func (a *AlchemyClient) httpRequest(rq alchemyRequest) ([]byte, error) {
	url := a.gateway + a.apiKey

	id := a.id.Add(1)
	payload := string(rq.GetPaylod(id))

	reader := strings.NewReader(payload)

	req, err := http.NewRequest("POST", url, reader)
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
