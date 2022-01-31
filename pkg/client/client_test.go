package client

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/bgzzz/counter/pkg/model"
	"github.com/stretchr/testify/assert"
)

func TestClientGet(t *testing.T) {

	tests := []struct {
		Srv            *httptest.Server
		ExpectedResult uint64
		ExpectedErr    error
		ClientFunction func(client *Client) (uint64, error)
	}{
		{

			Srv: httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter,
				req *http.Request) {
				assert.Equal(t, http.MethodGet, req.Method)
				getHandler(t, rw, 0, http.StatusOK)
			})),
			ExpectedResult: 0,
			ExpectedErr:    nil,
			ClientFunction: func(client *Client) (uint64, error) {
				return client.GetCounterValue()
			},
		},
		{
			Srv: httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter,
				req *http.Request) {
				assert.Equal(t, http.MethodGet, req.Method)
				http.Error(rw,
					"something went wrong",
					http.StatusInternalServerError)
			})),
			ExpectedResult: 0,
			ExpectedErr: fmt.Errorf("not expected response code on get: %d",
				http.StatusInternalServerError),
			ClientFunction: func(client *Client) (uint64, error) {
				return client.GetCounterValue()
			},
		},
		{
			Srv: httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter,
				req *http.Request) {
				assert.Equal(t, http.MethodGet, req.Method)
				rw.WriteHeader(http.StatusOK)
				fmt.Fprintf(rw, "wrong content")
			})),
			ExpectedResult: 0,
			ExpectedErr:    fmt.Errorf("unable to marshall counter: invalid character 'w' looking for beginning of value"),
			ClientFunction: func(client *Client) (uint64, error) {
				return client.GetCounterValue()
			},
		},
		{

			Srv: httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter,
				req *http.Request) {
				assert.Equal(t, http.MethodPost, req.Method)
				getHandler(t, rw, 0, http.StatusCreated)
			})),
			ExpectedResult: 0,
			ExpectedErr:    nil,
			ClientFunction: func(client *Client) (uint64, error) {
				return client.IncrementCounterValue()
			},
		},
		{

			Srv: httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter,
				req *http.Request) {
				assert.Equal(t, http.MethodPost, req.Method)
				http.Error(rw,
					"unable to increment, counter has reached its maximum value",
					http.StatusUnprocessableEntity)
			})),
			ExpectedResult: 0,
			ExpectedErr:    fmt.Errorf("Counter is on maximum and can't be incremented"),
			ClientFunction: func(client *Client) (uint64, error) {
				return client.IncrementCounterValue()
			},
		},
		{

			Srv: httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter,
				req *http.Request) {
				assert.Equal(t, http.MethodPost, req.Method)
				http.Error(rw,
					"something",
					http.StatusBadGateway)
			})),
			ExpectedResult: 0,
			ExpectedErr: fmt.Errorf("not expected response code on post: %d",
				http.StatusBadGateway),
			ClientFunction: func(client *Client) (uint64, error) {
				return client.IncrementCounterValue()
			},
		},
		{

			Srv: httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter,
				req *http.Request) {
				assert.Equal(t, http.MethodDelete, req.Method)
				getHandler(t, rw, 0, http.StatusOK)
			})),
			ExpectedResult: 0,
			ExpectedErr:    nil,
			ClientFunction: func(client *Client) (uint64, error) {
				return client.DecrementCounterValue()
			},
		},
		{

			Srv: httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter,
				req *http.Request) {
				assert.Equal(t, http.MethodDelete, req.Method)
				http.Error(rw,
					"something",
					http.StatusUnprocessableEntity)
			})),
			ExpectedResult: 0,
			ExpectedErr:    fmt.Errorf("Counter is on minimun and can't be decremented"),
			ClientFunction: func(client *Client) (uint64, error) {
				return client.DecrementCounterValue()
			},
		},
	}

	for i, test := range tests {
		tt := test
		t.Run(fmt.Sprintf("client get table test %d", i),
			func(t *testing.T) {
				t.Parallel()
				defer tt.Srv.Close()

				client := NewClient(tt.Srv.URL)
				client.httpCLient = tt.Srv.Client()

				val, err := tt.ClientFunction(client)

				assert.Equal(t, tt.ExpectedResult, val)

				// this is done to align with errors containing stack
				var evaluatedErr error
				if err != nil {
					evaluatedErr = fmt.Errorf(err.Error())
				}
				assert.Equal(t, tt.ExpectedErr, evaluatedErr)
			})
	}

}

func getHandler(t *testing.T, rw http.ResponseWriter, state uint64,
	status int) {
	b, err := json.Marshal(model.CounterRsp{
		Counter: state,
	})
	if err != nil {
		t.Fatalf("unable to marshal: %v", err)
	}

	rw.WriteHeader(status)
	if _, err := rw.Write(b); err != nil {
		t.Fatalf("unable to write rsp: %v", err)
	}
}
