package useful_servise

import (
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"math"
	"testing"
	"time"

	"github.com/valyala/fasthttp"
)

func TestHandler(t *testing.T) {
	us := NewUsefulService()
	handler := newHandler(us)

	t.Run("BadRoute", func(t *testing.T) {
		req := fasthttp.AcquireRequest()
		defer fasthttp.ReleaseRequest(req)
		req.URI().SetPath("/bad")

		var ctx fasthttp.RequestCtx
		ctx.Init(req, nil, nil)
		handler.Handle(&ctx)
		assert.EqualValues(t, fasthttp.StatusNotFound, ctx.Response.StatusCode())
		assert.EqualValues(t, ErrBadRoute.Error(), string(ctx.Response.Body()))
	})

	t.Run("ValidRoute", func(t *testing.T) {
		req := fasthttp.AcquireRequest()
		defer fasthttp.ReleaseRequest(req)
		req.URI().SetPath("/api/blocks/111/total")

		var ctx fasthttp.RequestCtx
		ctx.Init(req, nil, nil)
		handler.Handle(&ctx)
		require.EqualValues(t, fasthttp.StatusOK, ctx.Response.StatusCode())

		res := &TotalResponse{}
		err := json.Unmarshal(ctx.Response.Body(), res)
		assert.NoError(t, err)
		assert.NotNil(t, res)
	})
}

func TestIntegration(t *testing.T) {
	expectedResutls := []struct {
		Block uint32
		Resp  *TotalResponse
	}{
		{11509797, &TotalResponse{155, 2.285405}},
		{11508993, &TotalResponse{241, 1130.987085}},
		{109789, &TotalResponse{1, 4.99877}},
	}

	us := NewUsefulService(WithConfig(&Config{Address: ":4444"}))
	go func() {
		err := us.Start()
		assert.NoError(t, err, err)
	}()

	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)

	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(resp)

	for _, exp := range expectedResutls {
		if us.cfg.ApiKey == "" {
			// to avoid rate error because ApiKey is not set
			time.Sleep(time.Second * 5)
		}

		req.SetHost(fmt.Sprintf("localhost%s", us.cfg.Address))
		req.URI().SetPath(fmt.Sprintf("/api/blocks/%d/total", exp.Block))

		err := fasthttp.Do(req, resp)
		require.NoError(t, err, err)

		result := &TotalResponse{}
		err = json.Unmarshal(resp.Body(), result)
		assert.NoError(t, err, err)
		require.NotNil(t, result)
		assert.Equal(t, fasthttp.StatusOK, resp.StatusCode())

		if result.Amount != 0 {
			// round to provided value
			result.Amount = math.Round(result.Amount*1000000) / 1000000
		}
		assert.EqualValues(t, exp.Resp, result)
	}
}
