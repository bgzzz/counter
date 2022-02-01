package server

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/bgzzz/counter/pkg/model"
	"github.com/stretchr/testify/assert"
)

func TestCounterServer(t *testing.T) {

	tests := []struct {
		Srv                 *Server
		SrvState            uint64
		ReqMethod           string
		ExpectedStatusCode  int
		ExpectedResult      uint64
		ExpectedErrorString string
	}{
		{
			Srv:                 NewServer(8080),
			SrvState:            0,
			ReqMethod:           http.MethodGet,
			ExpectedStatusCode:  http.StatusOK,
			ExpectedResult:      0,
			ExpectedErrorString: "",
		},
		{
			Srv:                 NewServer(8080),
			SrvState:            0,
			ReqMethod:           http.MethodPost,
			ExpectedStatusCode:  http.StatusCreated,
			ExpectedResult:      1,
			ExpectedErrorString: "",
		},
		{
			Srv:                 NewServer(8080),
			SrvState:            MaxUint64,
			ReqMethod:           http.MethodPost,
			ExpectedStatusCode:  http.StatusUnprocessableEntity,
			ExpectedResult:      0,
			ExpectedErrorString: "unable to increment, counter has reached its maximum value",
		},
		{
			Srv:                 NewServer(8080),
			SrvState:            1,
			ReqMethod:           http.MethodDelete,
			ExpectedStatusCode:  http.StatusOK,
			ExpectedResult:      0,
			ExpectedErrorString: "",
		},
		{
			Srv:                 NewServer(8080),
			SrvState:            0,
			ReqMethod:           http.MethodDelete,
			ExpectedStatusCode:  http.StatusUnprocessableEntity,
			ExpectedResult:      0,
			ExpectedErrorString: "unable to decrement, counter has reached its minimum value",
		},
		{
			Srv:                 NewServer(8080),
			SrvState:            0,
			ReqMethod:           http.MethodPut,
			ExpectedStatusCode:  http.StatusMethodNotAllowed,
			ExpectedResult:      0,
			ExpectedErrorString: fmt.Sprintf("method %s is not supported", http.MethodPut),
		},
	}

	for i, test := range tests {
		tt := test
		t.Run(fmt.Sprintf("server table test %d", i),
			func(t *testing.T) {
				t.Parallel()
				req := httptest.NewRequest(tt.ReqMethod,
					"/api/v1/counter", nil)

				w := httptest.NewRecorder()

				tt.Srv.counter = tt.SrvState

				tt.Srv.handleCounter(w, req)
				res := w.Result()
				defer res.Body.Close()

				assert.Equal(t, tt.ExpectedStatusCode, res.StatusCode)

				data, err := ioutil.ReadAll(res.Body)
				if err != nil {
					t.Errorf("expected error to be nil got %v", err)
				}

				if tt.ExpectedErrorString != "" {
					return
				}

				var cntrRsp model.CounterRsp
				if err := json.Unmarshal(data, &cntrRsp); err != nil {
					t.Errorf("expected error(unmarshal) to be nil got %v", err)
				}

				assert.Equal(t, tt.ExpectedResult, cntrRsp.Counter)
			})
	}

}
