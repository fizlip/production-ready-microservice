package main

import (
  "context"
  "net/url"
  "strings"

  "github.com/go-kit/kit/endpoint"
  httptransport "github.com/go-kit/kit/transport/http"
)

type Endpoints struct {
  CreateReactionEndpoint endpoint.Endpoint
}

func MakeServerEndpoints(s Service) Endpoints {
  return Endpoints{
    CreateReactionEndpoint: MakeCreateReactionEndpoint(s),
  }
}

func MakeClientEndpoints(instance string) (Endpoints, error){
  if !strings.HasPrefix(instance, "http"){
    instance = "http://" + instance
  }

  tgt, err := url.Parse(instance)
  if err != nil {
    return Endpoints{}, err
  }

  tgt.Path = ""

  options := []httptransport.ClientOption{}

  return Endpoints{
    CreateReactionEndpoint: httptransport.NewClient("POST", tgt, encodeCreateReactionRequest, decodeCreateReactionResponse, options...).Endpoint(),
  }, nil
}

func (e Endpoints) CreateReaction(ctx context.Context, r Reaction) ([]byte,error) {
  request := createReactionRequest{Reaction: r}
  response, err := e.CreateReactionEndpoint(ctx, request)
  if err != nil{
    return []byte(""), err
  }

  resp := response.(createReactionResponse)
  return resp.Hash, nil
}

func MakeCreateReactionEndpoint(s Service) endpoint.Endpoint {
  return func(ctx context.Context, request interface{}) (response interface{}, err error){
    req := request.(createReactionRequest)
    h, nil := s.CreateReaction(ctx, req.Reaction.Address)
    return createReactionResponse{Hash: h}, nil
  }
}

type createReactionRequest struct {
  Reaction Reaction 
}

type createReactionResponse struct {
  Hash []byte `json:"hash"`
}

func (r createReactionResponse) error() []byte{return r.Hash}
