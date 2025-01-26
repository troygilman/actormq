package client

import (
	"errors"

	"github.com/anthdm/hollywood/actor"
)

func handleResponse[Result any](response *actor.Response) (result Result, err error) {
	resultAny, err := response.Result()
	if err != nil {
		return result, err
	}

	result, ok := resultAny.(Result)
	if !ok {
		return result, errors.New("could not cast result to specified type")
	}

	return result, nil
}
