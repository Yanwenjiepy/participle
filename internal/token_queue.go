package internal

import (
	"errors"
	"github.com/enriquebris/goconcurrentqueue"
)

var tokenQueue *goconcurrentqueue.FixedFIFO

func New(tokenNum int) {
	tokenQueue = goconcurrentqueue.NewFixedFIFO(tokenNum)
}

func GetQueueLen() int {
	qLen := tokenQueue.GetLen()
	return qLen
}

func GetQueueCap() int {
	qCap := tokenQueue.GetCap()
	return qCap
}

func GetToken() (string, error) {
	tokenUnknown, err := tokenQueue.Dequeue()
	if err != nil {
		return "", err
	}

	token, ok := tokenUnknown.(string)
	if !ok {
		ErrType := errors.New("token isn't string")
		return "", ErrType
	}

	return token, nil
}

func AddToken(token string) error {
	err := tokenQueue.Enqueue(token)
	if err != nil {
		return err
	}

	return nil
}