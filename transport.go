package main

import (
  "bytes"
  "context"
  "errors"
  "io/ioutil"
  "net/http"
  "encoding/json"

  "github.com/gorilla/mux"

  "github.com/go-kit/kit/log"
  "github.com/go-kit/kit/transport"
  httptransport "github.com/go-kit/kit/transport/http"

)

var (
  ErrBadRouting = errors.New("inconsistent mapping between route and handler (programmer error)")
)

func MakeHTTPHandler(s Service, logger log.Logger) http.Handler{
  r := mux.NewRouter()
  e := MakeServerEndpoints(s)
  options := []httptransport.ServerOption{
    httptransport.ServerErrorHandler(transport.NewLogErrorHandler(logger)),
    httptransport.ServerErrorEncoder(encodeError),
  }

  r.Methods("POST").Path("/profiles/").Handler(httptransport.NewServer(
    e.CreateReactionEndpoint,
    decodeCreateReactionRequest,
    encodeResponse,
    options...,
  ))
  return r
}

func decodeCreateReactionRequest(_ context.Context, r *http.Request) (request interface{}, err error){
  var req createReactionRequest
  if e := json.NewDecoder(r.Body).Decode(&req.Reaction); e != nil{
    return nil, e
  }
  return req, nil
}

func encodeCreateReactionRequest(ctx context.Context, req *http.Request, request interface{}) error {
  req.URL.Path = "/api/create"
  return encodeRequest(ctx, req, request)
}

func decodeCreateReactionResponse(_ context.Context, resp *http.Response) (interface{}, error){
  var response createReactionResponse
  err := json.NewDecoder(resp.Body).Decode(&response)
  return response, err
}

type errorer interface {
  error() error
}

func encodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error{
  if e, ok := response.(errorer); ok && e.error() != nil {
    encodeError(ctx, e.error(), w)
    return nil
  }
  w.Header().Set("Content-Type", "application/json; charset=utf-8")
  return json.NewEncoder(w).Encode(response)

}

func encodeRequest(ctx context.Context, req *http.Request, request interface{}) error{
  var buf bytes.Buffer
  err := json.NewEncoder(&buf).Encode(request)
  if err != nil{
    return err
  }
  req.Body = ioutil.NopCloser(&buf)
  return nil
}

func encodeError(_ context.Context, err error, w http.ResponseWriter){
  if err == nil {
    panic("encodeError with nil error")
  }

  w.Header().Set("Content-Type", "application/json; charset=utf-8")
  w.WriteHeader(codeFrom(err))
  json.NewEncoder(w).Encode(map[string]interface{}{
    "error": err.Error(),
  })
}

func codeFrom(err error) int {
  switch err {
    default:
      return http.StatusInternalServerError
  }
}
