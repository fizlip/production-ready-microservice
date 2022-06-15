package main

import (
  "crypto/sha256"
  "context"
  "errors"
)

type Service interface {
  CreateReaction(ctx context.Context, r Reaction) ([]byte,error)
}

type Reaction struct {
  Address string    `json:id`
  PostID  string    `json:postId`
}

var (
  ErrorAlreadyExists = errors.New("Reaction already exists")
)

type reactionService struct{}

func (reactionService) CreateReaction(ctx context.Context, r Reaction) ([]byte, error){
  return encryptSHA256(r.Address),nil
}

func encryptSHA256(data string) []byte{
  hash := sha256.Sum256([]byte(data))
  return hash[:]
}

