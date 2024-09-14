package pkg

import "errors"

var (
	ErrTraderAlreadyExist = errors.New("trader already exist")
	ErrCoinAlreadyExist   = errors.New("coin already exist")
	ErrInvalidCoinType    = errors.New("invalid coin type")
	ErrRingAlreadyExist   = errors.New("ring already exist")
)
