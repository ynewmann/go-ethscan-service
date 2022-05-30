package useful_servise

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math/big"
	"net/http"

	"github.com/valyala/fasthttp"
)

const (
	blockRoutePattern      = "\\/api\\/blocks\\/[0-9]*"
	blockTotalRoutePattern = blockRoutePattern + "\\/total"

	ZeroHex = "0x0"
)

var ErrBadRoute = errors.New("bad route")

type (
	Handler interface {
		Handle(ctx *fasthttp.RequestCtx)
	}

	TotalResponse struct {
		Transactions uint64  `json:"transactions"`
		Amount       float64 `json:"amount"`
	}
)

func newMiddleware(handler Handler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		log.Printf("Client: %s, %s Request %s", ctx.RemoteAddr(), ctx.Method(), ctx.URI())
		handler.Handle(ctx)
	}
}

type handler struct {
	us *UsefulService
}

func newHandler(us *UsefulService) *handler {
	return &handler{
		us: us,
	}
}

func (h *handler) Handle(ctx *fasthttp.RequestCtx) {
	path := ctx.Path()

	switch {
	case strictMatch(blockRoutePattern, path):
		id, err := getBlockNumber(path)
		if err != nil {
			ctx.Error("Cannot extract block number: "+err.Error(), http.StatusBadRequest)
			return
		}

		h.getBlock(ctx, id)
	case strictMatch(blockTotalRoutePattern, path):
		id, err := getBlockNumber(path)
		if err != nil {
			ctx.Error("Cannot extract block number: "+err.Error(), http.StatusBadRequest)
			return
		}

		h.getTotalInBlock(ctx, id)
	default:
		// path is not supported
		ctx.Error(ErrBadRoute.Error(), fasthttp.StatusNotFound)
		return
	}
}

func (h *handler) getBlock(ctx *fasthttp.RequestCtx, num uint64) {
	block, err := h.us.api.Proxy().GetBlockByNumber(ctx, num)
	if err != nil {
		ctx.Error("Cannot get block: "+err.Error(), http.StatusBadGateway)
		return
	}

	body, err := json.Marshal(block)
	if err != nil {
		ctx.Error("Cannot create JSON: "+err.Error(), http.StatusInternalServerError)
		return
	}

	ctx.Response.Header.SetContentType("application/json")
	ctx.Response.SetBody(body)
}

func (h *handler) getTotalInBlock(ctx *fasthttp.RequestCtx, num uint64) {
	hexNum := fmt.Sprintf("%x", num)
	if h.us.cfg.UseCache {
		entry, err := h.us.totalAmountCache.Get(hexNum)
		if err != nil {
			log.Printf("Error while getting %d from cache: %s\n", num, err)
		}

		if entry != nil {
			ctx.Response.Header.SetContentType("application/json")
			ctx.Response.SetBody(entry.Data())
			return
		}
	}

	block, err := h.us.api.Proxy().GetBlockByNumber(ctx, num)
	if err != nil {
		ctx.Error("Cannot get block: "+err.Error(), http.StatusBadGateway)
		return
	}

	totalWei := new(big.Float).SetUint64(0)
	for _, trx := range block.Transactions {
		if trx.Value == ZeroHex {
			continue
		}

		val, _ := new(big.Float).SetString(trx.Value)
		totalWei.Add(totalWei, val)
	}

	tr := &TotalResponse{
		Transactions: uint64(len(block.Transactions)),
		Amount:       weiToEth(totalWei),
	}

	body, err := json.Marshal(tr)
	if err != nil {
		ctx.Error("Cannot create JSON: "+err.Error(), http.StatusInternalServerError)
		return
	}

	ctx.Response.Header.SetContentType("application/json")
	ctx.Response.SetBody(body)

	if h.us.cfg.UseCache {
		err = h.us.totalAmountCache.Put(hexNum, body)
		if err != nil {
			log.Printf("Error while saving new entry to totalAmountCache: %s\n", err)
		}
	}
}
