package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestRouter(t *testing.T) {
	log := zap.NewNop()
	type args struct {
		h        Handler
		endpoint string
	}
	tests := []struct {
		name          string
		args          args
		wantCode      int
		wantBody      string
		skipBodyCheck bool
	}{
		{
			name: "/healthy OK",
			args: args{
				h:        Handler{},
				endpoint: "/healthy",
			},
			wantCode: http.StatusOK,
			wantBody: "OK",
		}, {
			name: "/ready OK",
			args: args{
				h: Handler{
					client: mockClient{},
					log:    log,
				},
				endpoint: "/ready",
			},
			wantCode: http.StatusOK,
			wantBody: "OK",
		}, {
			name: "/ready not ready",
			args: args{
				h: Handler{
					client: mockClient{
						err: errors.New("client error"),
					},
					log: log,
				},
				endpoint: "/ready",
			},
			wantCode: http.StatusServiceUnavailable,
			wantBody: "Service not ready\n",
		}, {
			name: "/metrics",
			args: args{
				h: Handler{
					log: log,
				},
				endpoint: "/metrics",
			},
			wantCode:      http.StatusOK,
			skipBodyCheck: true,
		}, {
			name: "/eth/balance OK",
			args: args{
				h: Handler{
					client: mockClient{
						balance: big.NewInt(100),
					},
					log: log,
				},
				endpoint: "/eth/balance/0xfe3b557e8fb62b89f4916b721be55ceb828dbd73",
			},
			wantCode: http.StatusOK,
			wantBody: fmt.Sprintf("{\n\t\"balance\": \"%d\"\n}\n", 100),
		}, {
			name: "/eth/balance bad address",
			args: args{
				h: Handler{
					client: mockClient{},
					log:    log,
				},
				endpoint: "/eth/balance/bad-address",
			},
			wantCode: http.StatusBadRequest,
			wantBody: "Invalid address\n",
		}, {
			name: "/eth/balance client error",
			args: args{
				h: Handler{
					client: mockClient{
						err: errors.New("client error"),
					},
					log: log,
				},
				endpoint: "/eth/balance/0xfe3b557e8fb62b89f4916b721be55ceb828dbd73",
			},
			wantCode: http.StatusInternalServerError,
			wantBody: "Failed to get balance\n",
		}, {
			name: "not found",
			args: args{
				h: Handler{
					client: mockClient{},
					log:    log,
				},
				endpoint: "/not-exists",
			},
			wantCode: http.StatusNotFound,
			wantBody: "404 page not found\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rls := NewRateLimiterStore(1, 1)
			r := Router(rls, tt.args.h)
			ts := httptest.NewServer(r)
			defer ts.Close()

			resp, err := http.Get(ts.URL + tt.args.endpoint)
			if err != nil {
				t.Fatal(err)
			}

			body, err := io.ReadAll(resp.Body)
			_ = resp.Body.Close()
			if err != nil {
				t.Fatal(err)
			}

			assert.Equalf(t, tt.wantCode, resp.StatusCode, "Response status code mismatch. Have: %d, want: %d.", resp.StatusCode, tt.wantCode)
			if !tt.skipBodyCheck {
				assert.Equalf(t, tt.wantBody, string(body), "Response body mismatch. Have: %s, want: %s.", string(body), tt.wantBody)
			}
		})
	}
}

type mockClient struct {
	balance *big.Int
	err     error
}

func (m mockClient) GetBalance(_ context.Context, _ common.Address) (*big.Int, error) {
	return big.NewInt(100), m.err
}

func (m mockClient) BlockNumber(_ context.Context) (uint64, error) {
	return 0, m.err
}
