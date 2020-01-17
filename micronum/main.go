package micronum

import (
  "sort"
	"context"
  "strconv"
	"encoding/json"
	"errors"
  "log"
	"net/http"

  "gomicro/worker"
  "gomicro/workload/fib"

	"github.com/go-kit/kit/endpoint"
	httptransport "github.com/go-kit/kit/transport/http"
)

type numberService struct{}
type NumberService interface {
	FibonacciCon(string) ([]int, error)
}
func (numberService) FibonacciCon(s string) ([]int, error) {
  layers, err := strconv.Atoi(s)
  response    := make([]int, layers)
  log.Print("|-: jobs layers: ", layers)
  if err == nil {
    jobs    := make(chan int, layers)
    results := make(chan int, layers)
    log.Print("|-> setup done.")
    for w := 0; w < layers; w++ { go worker.Spawn(jobs, results, fib.Run) }
    for i := 0; i < layers; i++ { jobs <- i }
    close(jobs)
    log.Print("|-> jobs enqueue done.")
    for j := 0; j < layers; j++ { response[j] = <- results }
    log.Print("|-> jobs results done.")
    sort.Ints(response)
    return response, nil
  } else {
    log.Print("layer cast error.")
    return nil, errors.New("empty string")
  }
}

type fibonacciConRequest struct {
	F string `json:"f"`
}
type fibonacciConResponse struct {
	V   []int `json:"v"`
	Err string `json:"err,omitempty"` // errors don't define JSON marshaling
}

func decodefibonacciConRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request fibonacciConRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil { return nil, err }
	return request, nil
}
func encodeResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	return json.NewEncoder(w).Encode(response)
}
func makefibonacciConEndpoint(svc NumberService) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(fibonacciConRequest)
		v, err := svc.FibonacciCon(req.F)
		if err != nil { return fibonacciConResponse{v, err.Error()}, nil }
		return fibonacciConResponse{v, ""}, nil
	}
}

func InitializeFibc() *httptransport.Server {
	svc := numberService{}
	fibonacciConHandler := httptransport.NewServer(
		makefibonacciConEndpoint(svc),
		decodefibonacciConRequest,
		encodeResponse,
	)
  return fibonacciConHandler
}
