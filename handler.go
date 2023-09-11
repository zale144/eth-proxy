package main

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"

	"github.com/ethereum/go-ethereum/common"
	"github.com/julienschmidt/httprouter"
	"go.uber.org/zap"
)

// Handler is the HTTP handler for the service
type Handler struct {
	client ethClient
	log    *zap.Logger
}

type ethClient interface {
	GetBalance(ctx context.Context, account common.Address) (*big.Int, error)
	BlockNumber(ctx context.Context) (uint64, error)
}

// NewHandler creates a new HTTP handler
func NewHandler(client ethClient, log *zap.Logger) (Handler, error) {
	return Handler{
		client: client,
		log:    log,
	}, nil
}

// GetBalance @Summary Show the Ethereum balance
// @Description Get the Ethereum balance for the given address
// @Accept  json
// @Produce  json
// @Param   address     path    string     true        "Ethereum Address"
// @Success 200  {object}  main.GetBalance.getBalanceResponse
// @Failure 400 {string} string "Invalid address"
// @Failure 500 {string} string "Failed to get balance"
// @Router  /eth/balance/{address} [get]
func (h Handler) GetBalance(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	address := params.ByName("address")

	if !common.IsHexAddress(address) {
		http.Error(w, "Invalid address", http.StatusBadRequest)
		return
	}

	account := common.HexToAddress(address)

	balance, err := h.client.GetBalance(r.Context(), account)
	if err != nil {
		ethBalanceRequests.WithLabelValues("failed").Inc()
		h.log.Error("Failed to get balance", zap.Error(err))
		http.Error(w, "Failed to get balance", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	enc := json.NewEncoder(w)
	enc.SetIndent("", "\t")

	type getBalanceResponse struct {
		Balance string `json:"balance" example:"1000000000"`
	}

	response := getBalanceResponse{
		Balance: balance.String(),
	}

	if err = enc.Encode(response); err != nil {
		ethBalanceRequests.WithLabelValues("failed").Inc()
		h.log.Error("Failed to encode response", zap.Error(err))
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}

	ethBalanceRequests.WithLabelValues("success").Inc()
}

// Healthy @Summary Show the health status
// @Description Show the health status
// @Accept  json
// @Produce  json
// @Success 200  {string} string "OK"
// @Router  /healthy [get]
func (h Handler) Healthy(w http.ResponseWriter, _ *http.Request, _ httprouter.Params) {
	w.WriteHeader(http.StatusOK)
	_, _ = fmt.Fprintf(w, "OK")
}

// Ready @Summary Show the readiness status
// @Description Show the readiness status
// @Accept  json
// @Produce  json
// @Success 200  {string} string "OK"
// @Failure 503 {string} string "Service not ready"
// @Router  /ready [get]
func (h Handler) Ready(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	_, err := h.client.BlockNumber(r.Context())
	if err != nil {
		h.log.Error("Failed to get block number", zap.Error(err))
		http.Error(w, "Service not ready", http.StatusServiceUnavailable)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, _ = fmt.Fprintf(w, "OK")
}
