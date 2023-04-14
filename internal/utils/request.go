package utils

import (
	"errors"
	"io"
	"net/http"
	"strconv"
)

func SendGetRequest(requestURL string) (resBody []byte, err error) {
	req, err := http.NewRequest(http.MethodGet, requestURL, nil)
	if err != nil {
		return nil, err
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	if res.StatusCode != 200 {
		return nil, errors.New("status code is " + strconv.Itoa(res.StatusCode))
	}

	resBody, err = io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	return resBody, nil

}
