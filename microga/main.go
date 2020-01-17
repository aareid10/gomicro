package microga

import (
	"context"
	"encoding/json"
	"errors"
  "log"
	"net/http"

  // "gomicro/worker"
  "gomicro/workload/ga"

	"github.com/go-kit/kit/endpoint"
	httptransport "github.com/go-kit/kit/transport/http"
)

type geneticService struct{}
type GeneticService interface {
	GeneticAlgo(string) (int, error)
}
func (geneticService) GeneticAlgo(s string) (string, error) {
  origin := []byte(s)
  if s == "" {
		return "", errors.New("empty string")
	} else {
    log.Print("|-: Ga orign: ", s, origin)
    return ga.Run(origin), nil
  }
}

type GeneticAlgoRequest struct {
	P string `json:"p"`
}
type GeneticAlgoResponse struct {
	V   string `json:"v"`
	Err string `json:"err,omitempty"` // errors don't define JSON marshaling
}

func decodeGeneticAlgoRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request GeneticAlgoRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil { return nil, err }
	return request, nil
}
func encodeResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	return json.NewEncoder(w).Encode(response)
}
func makeGeneticAlgoEndpoint(svc geneticService) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(GeneticAlgoRequest)
		v, err := svc.GeneticAlgo(req.P)
		if err != nil { return GeneticAlgoResponse{v, err.Error()}, nil }
		return GeneticAlgoResponse{v, ""}, nil
	}
}

func InitializeGAc() *httptransport.Server {
	svc := geneticService{}
	GeneticAlgoHandler := httptransport.NewServer(
		makeGeneticAlgoEndpoint(svc),
		decodeGeneticAlgoRequest,
		encodeResponse,
	)
  return GeneticAlgoHandler
}
