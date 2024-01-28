package adapter

import (
	"encoding/json"
	"errors"
	"strconv"

	"gopkg.in/resty.v1"

	"it-chain/blockchain"
)

var ErrGetBlockFromPeer = errors.New("Error when getting block from other peer")

type HttpBlockAdapter struct {
}

func NewHttpBlockAdapter() *HttpBlockAdapter {
	return &HttpBlockAdapter{}
}

func (a HttpBlockAdapter) GetLastBlockFromPeer(peer blockchain.Peer) (blockchain.DefaultBlock, error) {

	resp, err := resty.R().
		SetQueryString("height=-1").
		SetHeader("Content-Type", "application/json").
		Get("http://" + peer.ApiGatewayAddress + "/blocks")
	if err != nil {
		return blockchain.DefaultBlock{}, err
	}

	block := blockchain.DefaultBlock{}
	if err := json.Unmarshal(resp.Body(), &block); err != nil {
		return blockchain.DefaultBlock{}, err
	}

	if block.IsEmpty() {
		return blockchain.DefaultBlock{}, ErrGetBlockFromPeer
	}

	return block, nil
}

func (a HttpBlockAdapter) GetBlockByHeightFromPeer(height blockchain.BlockHeight, peer blockchain.Peer) (blockchain.DefaultBlock, error) {

	resp, err := resty.R().
		SetQueryParams(map[string]string{
			"height": strconv.FormatUint(height, 10),
		}).
		SetHeader("Content-Type", "application/json").
		Get("http://" + peer.ApiGatewayAddress + "/blocks")
	if err != nil {
		return blockchain.DefaultBlock{}, err

	}

	block := blockchain.DefaultBlock{}
	if err := json.Unmarshal(resp.Body(), &block); err != nil {
		return blockchain.DefaultBlock{}, err
	}

	if block.IsEmpty() {
		return blockchain.DefaultBlock{}, ErrGetBlockFromPeer
	}

	return block, nil
}
