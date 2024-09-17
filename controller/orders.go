package controller

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/zkryaev/taskwb-L0/cache"
	"github.com/zkryaev/taskwb-L0/repository/config"
)

type Server struct {
	cfg   config.ConfigApp
	Cache *cache.Cache
}

func New(cfgPath string, cache *cache.Cache) *Server {
	cfg := config.Load(cfgPath)
	return &Server{
		cfg:   cfg.App,
		Cache: cache,
	}
}

func (s *Server) Launch() error {
	http.HandleFunc("/order", s.GetOrderHandler)
	err := http.ListenAndServe(s.cfg.Host+":"+s.cfg.Port, nil)
	if err != nil {
		return fmt.Errorf("failed to launch server: %w", err)
	}
	return nil
}

func (s *Server) GetOrderHandler(w http.ResponseWriter, r *http.Request) {
	orderUID := r.URL.Query().Get("id")
	if order, ok := s.Cache.GetOrder(orderUID); ok {
		orderJSON, _ := json.MarshalIndent(order, "", "    ")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(orderJSON)
	} else {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Order not found"))
	}
}
