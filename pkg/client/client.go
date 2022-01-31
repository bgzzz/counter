package client

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/bgzzz/counter/pkg/model"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

func NewClient(host string) *Client {
	return &Client{
		httpCLient: &http.Client{},
		hostURL:    host,
	}
}

type Client struct {
	hostURL    string
	httpCLient *http.Client
}

func (cl *Client) GetCounterValue() (uint64, error) {
	resp, err := cl.httpCLient.Get(cl.hostURL)
	if err != nil {
		return 0, errors.Wrap(err, "unable to get counter value")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("not expected response code on get: %d",
			resp.StatusCode)
	}

	return getCounterValue(resp.Body)
}

func (cl *Client) IncrementCounterValue() (uint64, error) {
	log.Debug("Incrementing counter value")

	resp, err := cl.httpCLient.Post(cl.hostURL,
		"application/json", nil)
	if err != nil {
		return 0, errors.Wrap(err, "unable to post counter")
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnprocessableEntity {
		b, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return 0,
				errors.Wrap(err,
					"unable to read responce body of maxed increment payload")
		}

		log.Debug(string(b))
		return 0,
			fmt.Errorf("Counter is on maximum and can't be incremented")
	}

	if resp.StatusCode != http.StatusCreated {
		return 0, fmt.Errorf("not expected response code on post: %d",
			resp.StatusCode)
	}

	return getCounterValue(resp.Body)
}

func (cl *Client) DecrementCounterValue() (uint64, error) {
	log.Debug("Decrementing counter value")

	req, err := http.NewRequest("DELETE", cl.hostURL, nil)
	if err != nil {
		return 0, errors.Wrap(err, "unable to create delete request")
	}

	resp, err := cl.httpCLient.Do(req)
	if err != nil {
		return 0, errors.Wrap(err, "unable to make delete request")
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnprocessableEntity {
		log.Info("Counter is on minimun and can't be decremented")
		b, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return 0,
				errors.Wrap(err, "unable to read responce body of minimum decrement payload")
		}

		log.Debug(string(b))
		return 0, fmt.Errorf("Counter is on minimun and can't be decremented")
	}

	if resp.StatusCode != http.StatusOK {
		return 0,
			fmt.Errorf("not expected response code on delete: %d",
				resp.StatusCode)
	}

	return getCounterValue(resp.Body)
}

func getCounterValue(r io.Reader) (uint64, error) {
	log.Debug("Getting counter value")

	b, err := ioutil.ReadAll(r)
	if err != nil {
		return 0, errors.Wrap(err,
			"unable to read responce body of get request")
	}

	var cntr model.CounterRsp
	if err := json.Unmarshal(b, &cntr); err != nil {
		return 0, errors.Wrap(err,
			"unable to marshall counter")
	}

	log.Infof("Counter value is %d", cntr.Counter)
	return cntr.Counter, nil
}
