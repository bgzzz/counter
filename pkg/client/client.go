package Client

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"path"

	"github.com/bgzzz/counter/pkg/model"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

const protocolSchemaTemplate = "http://%s"

func NewClient(host string) *Client {
	return &Client{
		hostURL: fmt.Sprintf(protocolSchemaTemplate,
			path.Join(host,
				"api",
				model.APIVersion,
				"counter")),
	}
}

type Client struct {
	hostURL string
	timeout string
}

func (cl *Client) GetCounterValue() error {
	resp, err := http.Get(cl.hostURL)
	if err != nil {
		return errors.Wrap(err, "unable to get counter value")
	}

	if resp.StatusCode != http.StatusAccepted {
		return fmt.Errorf("not expected response code on get: %d",
			resp.StatusCode)
	}

	return getCounterValue(resp.Body)
}

func (cl *Client) IncrementCounterValue() error {
	log.Debug("Incrementing counter value")

	resp, err := http.Post(cl.hostURL,
		"application/json", nil)
	if err != nil {
		return errors.Wrap(err, "unable to post counter")
	}

	if resp.StatusCode == http.StatusUnprocessableEntity {
		log.Info("Counter is on maximum and can't be incremented")
		b, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return errors.Wrap(err, "unable to read responce body of maxed increment payload")
		}

		log.Debug(string(b))
		return nil
	}

	if resp.StatusCode != http.StatusCreated {
		return errors.Wrapf(err,
			"not expected response code on post: %d", resp.StatusCode)
	}

	return getCounterValue(resp.Body)
}

func (cl *Client) DecrementCounterValue() error {
	log.Debug("Decrementing counter value")

	client := &http.Client{}
	req, err := http.NewRequest("DELETE", cl.hostURL, nil)
	if err != nil {
		return errors.Wrap(err, "unable to create delete request")
	}

	resp, err := client.Do(req)
	if err != nil {
		return errors.Wrap(err, "unable to make delete request")
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnprocessableEntity {
		log.Info("Counter is on minimun and can't be decremented")
		b, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return errors.Wrap(err, "unable to read responce body of minimum decrement payload")
		}

		log.Debug(string(b))
		return nil
	}

	return getCounterValue(resp.Body)
}

func getCounterValue(r io.Reader) error {
	log.Debug("Getting counter value")

	b, err := ioutil.ReadAll(r)
	if err != nil {
		return errors.Wrap(err,
			"unable to read responce body of get request")
	}

	var cntr model.CounterRsp
	if err := json.Unmarshal(b, &cntr); err != nil {
		return errors.Wrap(err,
			"unable to marshall counter")
	}

	log.Infof("Counter value is %d", cntr.Counter)
	return nil
}
