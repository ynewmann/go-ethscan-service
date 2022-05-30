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
		GetBlockByNumber(context.Context, uint64) (*BlockDto, error)
		CallAction(context.Context, string, ...ArgOption) (*BlockPage, error)
	}

	ProxyModule struct {
		api  *api
		Name string
	}
)

func NewProxyModule(api *api) Proxy {
	return &ProxyModule{api, ModuleName}
}

func (p *ProxyModule) GetBlockByNumber(ctx context.Context, number uint64) (*BlockDto, error) {
	page, err := p.CallAction(
		ctx,
		GetBlockByNumber,
		withArg(TagArg, fmt.Sprintf("%x", number)),
		withArg(BooleanArg, "true"),
	)
	if err != nil {
		return nil, err
	}

	return &page.Result, nil
}

func (p *ProxyModule) CallAction(ctx context.Context, action string, args ...ArgOption) (*BlockPage, error) {
	page := &BlockPage{}
	args = append(args, withArg(ModuleArg, p.Name), withArg(ActionArg, action))

	_, err := p.api.doNewRequest(
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
