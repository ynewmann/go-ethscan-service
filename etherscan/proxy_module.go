package etherscan

import (
	"context"
	"fmt"
	"net/http"
)

const (
	ModuleName = "proxy"

	// common proxy args
	ActionArg  = "action"
	TagArg     = "tag"
	BooleanArg = "boolean"
	IndexArg   = "index"
	TxhashAtg  = "txhash"
	AddressAtg = "Address"

	GetBlockByNumber = "eth_getBlockByNumber"
)

type (
	Proxy interface {
		Moduler

		GetBlockByNumber(context.Context, uint64) (*BlockDto, error)
		CallAction(context.Context, string, ...ArgOption) (*BlockPage, error)
	}

	ProxyModule struct {
		*Module
	}
)

func NewProxyModule(client *apiClient) Proxy {
	return &ProxyModule{NewModule(client, ModuleName)}
}

func (p *ProxyModule) GetBlockByNumber(ctx context.Context, number uint64) (*BlockDto, error) {
	page, err := p.CallAction(
		ctx,
		GetBlockByNumber,
		WithArg(TagArg, fmt.Sprintf("%x", number)),
		WithArg(BooleanArg, "true"),
	)
	if err != nil {
		return nil, err
	}

	return &page.Result, nil
}

func (p *ProxyModule) CallAction(ctx context.Context, action string, args ...ArgOption) (*BlockPage, error) {
	page := &BlockPage{}
	args = append(args, WithArg(ModuleArg, p.Name), WithArg(ActionArg, action))

	_, err := p.client.doNewRequest(
		ctx,
		http.MethodGet,
		nil,
		page,
		args...,
	)
	if err != nil {
		return nil, err
	}

	return page, nil
}
