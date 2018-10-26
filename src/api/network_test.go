package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/skycoin/skycoin/src/daemon"
	"github.com/skycoin/skycoin/src/readable"
)

func TestConnection(t *testing.T) {
	tt := []struct {
		name                       string
		method                     string
		status                     int
		err                        string
		addr                       string
		gatewayGetConnectionResult *daemon.Connection
		gatewayGetConnectionError  error
		result                     *readable.Connection
	}{
		{
			name:   "405",
			method: http.MethodPost,
			status: http.StatusMethodNotAllowed,
			err:    "405 Method Not Allowed",
		},
		{
			name:                       "400 - empty addr",
			method:                     http.MethodGet,
			status:                     http.StatusBadRequest,
			err:                        "400 Bad Request - addr is required",
			addr:                       "",
			gatewayGetConnectionResult: nil,
			result:                     nil,
		},
		{
			name:   "200",
			method: http.MethodGet,
			status: http.StatusOK,
			err:    "",
			addr:   "addr",
			gatewayGetConnectionResult: &daemon.Connection{
				ID:           1,
				Addr:         "127.0.0.1",
				LastSent:     99999,
				LastReceived: 1111111,
				Outgoing:     true,
				Introduced:   true,
				Mirror:       9876,
				ListenPort:   9877,
				Height:       1234,
			},
			result: &readable.Connection{
				ID:           1,
				Addr:         "127.0.0.1",
				LastSent:     99999,
				LastReceived: 1111111,
				Outgoing:     true,
				Introduced:   true,
				Mirror:       9876,
				ListenPort:   9877,
				Height:       1234,
			},
		},

		{
			name:                      "500 - GetConnection failed",
			method:                    http.MethodGet,
			status:                    http.StatusInternalServerError,
			err:                       "500 Internal Server Error - GetConnection failed",
			addr:                      "addr",
			gatewayGetConnectionError: errors.New("GetConnection failed"),
		},

		{
			name:   "404",
			method: http.MethodGet,
			status: http.StatusNotFound,
			addr:   "addr",
			err:    "404 Not Found",
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			endpoint := "/api/v1/network/connection"
			gateway := &MockGatewayer{}
			gateway.On("GetConnection", tc.addr).Return(tc.gatewayGetConnectionResult, tc.gatewayGetConnectionError)

			v := url.Values{}
			if tc.addr != "" {
				v.Add("addr", tc.addr)
			}
			if len(v) > 0 {
				endpoint += "?" + v.Encode()
			}
			req, err := http.NewRequest(tc.method, endpoint, nil)
			require.NoError(t, err)

			rr := httptest.NewRecorder()
			handler := newServerMux(defaultMuxConfig(), gateway, &CSRFStore{}, nil)
			handler.ServeHTTP(rr, req)

			status := rr.Code
			require.Equal(t, tc.status, status, "got `%v` want `%v`", status, tc.status)

			if status != http.StatusOK {
				require.Equal(t, tc.err, strings.TrimSpace(rr.Body.String()), "got `%v`| %d, want `%v`",
					strings.TrimSpace(rr.Body.String()), status, tc.err)
			} else {
				var msg *readable.Connection
				err = json.Unmarshal(rr.Body.Bytes(), &msg)
				require.NoError(t, err)
				require.Equal(t, tc.result, msg)
			}
		})
	}
}

func TestConnections(t *testing.T) {
	tt := []struct {
		name                                string
		method                              string
		typ                                 string
		status                              int
		err                                 string
		gatewayGetConnectionsResult         []daemon.Connection
		gatewayGetConnectionsError          error
		gatewayGetOutgoingConnectionsResult []daemon.Connection
		gatewayGetOutgoingConnectionsError  error
		gatewayGetIncomingConnectionsResult []daemon.Connection
		gatewayGetIncomingConnectionsError  error
		result                              Connections
	}{
		{
			name:   "405",
			method: http.MethodPost,
			status: http.StatusMethodNotAllowed,
			err:    "405 Method Not Allowed",
		},

		{
			name:   "400 bad type",
			method: http.MethodGet,
			typ:    "foo",
			status: http.StatusBadRequest,
			err:    "400 Bad Request - invalid type",
		},

		{
			name:   "200 all connections",
			method: http.MethodGet,
			status: http.StatusOK,
			err:    "",
			gatewayGetConnectionsResult: []daemon.Connection{
				{
					ID:           1,
					Addr:         "127.0.0.1",
					LastSent:     99999,
					LastReceived: 1111111,
					Outgoing:     true,
					Introduced:   true,
					Mirror:       9876,
					ListenPort:   9877,
					Height:       1234,
				},
				{
					ID:           3,
					Addr:         "127.1.1.1",
					LastSent:     88888,
					LastReceived: 9999999,
					Outgoing:     false,
					Introduced:   true,
					Mirror:       9877,
					ListenPort:   9877,
					Height:       1235,
				},
			},
			result: Connections{
				Connections: []readable.Connection{
					{
						ID:           1,
						Addr:         "127.0.0.1",
						LastSent:     99999,
						LastReceived: 1111111,
						Outgoing:     true,
						Introduced:   true,
						Mirror:       9876,
						ListenPort:   9877,
						Height:       1234,
					},
					{
						ID:           3,
						Addr:         "127.1.1.1",
						LastSent:     88888,
						LastReceived: 9999999,
						Outgoing:     false,
						Introduced:   true,
						Mirror:       9877,
						ListenPort:   9877,
						Height:       1235,
					},
				},
			},
		},

		{
			name:   "200 outgoing connections",
			method: http.MethodGet,
			typ:    "outgoing",
			status: http.StatusOK,
			err:    "",
			gatewayGetOutgoingConnectionsResult: []daemon.Connection{
				{
					ID:           1,
					Addr:         "127.0.0.1",
					LastSent:     99999,
					LastReceived: 1111111,
					Outgoing:     true,
					Introduced:   true,
					Mirror:       9876,
					ListenPort:   9877,
					Height:       1234,
				},
			},
			result: Connections{
				Connections: []readable.Connection{
					{
						ID:           1,
						Addr:         "127.0.0.1",
						LastSent:     99999,
						LastReceived: 1111111,
						Outgoing:     true,
						Introduced:   true,
						Mirror:       9876,
						ListenPort:   9877,
						Height:       1234,
					},
				},
			},
		},

		{
			name:   "200 incoming connections",
			method: http.MethodGet,
			typ:    "incoming",
			status: http.StatusOK,
			err:    "",
			gatewayGetIncomingConnectionsResult: []daemon.Connection{
				{
					ID:           1,
					Addr:         "127.0.0.1",
					LastSent:     99999,
					LastReceived: 1111111,
					Outgoing:     false,
					Introduced:   true,
					Mirror:       9876,
					ListenPort:   9877,
					Height:       1234,
				},
			},
			result: Connections{
				Connections: []readable.Connection{
					{
						ID:           1,
						Addr:         "127.0.0.1",
						LastSent:     99999,
						LastReceived: 1111111,
						Outgoing:     false,
						Introduced:   true,
						Mirror:       9876,
						ListenPort:   9877,
						Height:       1234,
					},
				},
			},
		},

		{
			name:                       "500 - GetConnections failed",
			method:                     http.MethodGet,
			status:                     http.StatusInternalServerError,
			err:                        "500 Internal Server Error - GetConnections failed",
			gatewayGetConnectionsError: errors.New("GetConnections failed"),
		},

		{
			name:                               "500 - GetOutgoingConnections failed",
			method:                             http.MethodGet,
			typ:                                "outgoing",
			status:                             http.StatusInternalServerError,
			err:                                "500 Internal Server Error - GetOutgoingConnections failed",
			gatewayGetOutgoingConnectionsError: errors.New("GetOutgoingConnections failed"),
		},

		{
			name:                               "500 - GetIncomingConnections failed",
			method:                             http.MethodGet,
			typ:                                "incoming",
			status:                             http.StatusInternalServerError,
			err:                                "500 Internal Server Error - GetIncomingConnections failed",
			gatewayGetIncomingConnectionsError: errors.New("GetIncomingConnections failed"),
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			endpoint := "/api/v1/network/connections"
			gateway := &MockGatewayer{}
			gateway.On("GetConnections").Return(tc.gatewayGetConnectionsResult, tc.gatewayGetConnectionsError)
			gateway.On("GetOutgoingConnections").Return(tc.gatewayGetOutgoingConnectionsResult, tc.gatewayGetOutgoingConnectionsError)
			gateway.On("GetIncomingConnections").Return(tc.gatewayGetIncomingConnectionsResult, tc.gatewayGetIncomingConnectionsError)

			v := url.Values{}
			if tc.typ != "" {
				v.Add("type", tc.typ)
			}

			ve := v.Encode()
			if ve != "" {
				endpoint += "?" + ve
			}

			req, err := http.NewRequest(tc.method, endpoint, nil)
			require.NoError(t, err)

			rr := httptest.NewRecorder()
			handler := newServerMux(defaultMuxConfig(), gateway, &CSRFStore{}, nil)
			handler.ServeHTTP(rr, req)

			status := rr.Code
			require.Equal(t, tc.status, status, "got `%v` want `%v`", status, tc.status)

			if status != http.StatusOK {
				require.Equal(t, tc.err, strings.TrimSpace(rr.Body.String()), "got `%v`| %d, want `%v`",
					strings.TrimSpace(rr.Body.String()), status, tc.err)
			} else {
				var msg Connections
				err = json.Unmarshal(rr.Body.Bytes(), &msg)
				require.NoError(t, err)
				require.Equal(t, tc.result, msg)
			}
		})
	}
}

func TestDefaultConnections(t *testing.T) {
	tt := []struct {
		name                               string
		method                             string
		status                             int
		err                                string
		gatewayGetDefaultConnectionsResult []string
		result                             []string
	}{
		{
			name:   "405",
			method: http.MethodPost,
			status: http.StatusMethodNotAllowed,
			err:    "405 Method Not Allowed",
		},
		{
			name:                               "200",
			method:                             http.MethodGet,
			status:                             http.StatusOK,
			err:                                "",
			gatewayGetDefaultConnectionsResult: []string{"44.33.22.11", "11.44.66.88"},
			result:                             []string{"11.44.66.88", "44.33.22.11"},
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			endpoint := "/api/v1/network/defaultConnections"
			gateway := &MockGatewayer{}
			gateway.On("GetDefaultConnections").Return(tc.gatewayGetDefaultConnectionsResult)

			req, err := http.NewRequest(tc.method, endpoint, nil)
			require.NoError(t, err)

			rr := httptest.NewRecorder()
			handler := newServerMux(defaultMuxConfig(), gateway, &CSRFStore{}, nil)
			handler.ServeHTTP(rr, req)

			status := rr.Code
			require.Equal(t, tc.status, status, "got `%v` want `%v`", status, tc.status)

			if status != http.StatusOK {
				require.Equal(t, tc.err, strings.TrimSpace(rr.Body.String()), "got `%v`| %d, want `%v`",
					strings.TrimSpace(rr.Body.String()), status, tc.err)
			} else {
				var msg []string
				err = json.Unmarshal(rr.Body.Bytes(), &msg)
				require.NoError(t, err)
				require.Equal(t, tc.result, msg)
			}
		})
	}
}

func TestGetTrustConnections(t *testing.T) {
	tt := []struct {
		name                             string
		method                           string
		status                           int
		err                              string
		gatewayGetTrustConnectionsResult []string
		result                           []string
	}{
		{
			name:   "405",
			method: http.MethodPost,
			status: http.StatusMethodNotAllowed,
			err:    "405 Method Not Allowed",
		},
		{
			name:                             "200",
			method:                           http.MethodGet,
			status:                           http.StatusOK,
			err:                              "",
			gatewayGetTrustConnectionsResult: []string{"44.33.22.11", "11.44.66.88"},
			result:                           []string{"11.44.66.88", "44.33.22.11"},
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			endpoint := "/api/v1/network/connections/trust"
			gateway := &MockGatewayer{}
			gateway.On("GetTrustConnections").Return(tc.gatewayGetTrustConnectionsResult)

			req, err := http.NewRequest(tc.method, endpoint, nil)
			require.NoError(t, err)

			rr := httptest.NewRecorder()
			handler := newServerMux(defaultMuxConfig(), gateway, &CSRFStore{}, nil)
			handler.ServeHTTP(rr, req)

			status := rr.Code
			require.Equal(t, tc.status, status, "got `%v` want `%v`", status, tc.status)

			if status != http.StatusOK {
				require.Equal(t, tc.err, strings.TrimSpace(rr.Body.String()), "got `%v`| %d, want `%v`",
					strings.TrimSpace(rr.Body.String()), status, tc.err)
			} else {
				var msg []string
				err = json.Unmarshal(rr.Body.Bytes(), &msg)
				require.NoError(t, err)
				require.Equal(t, tc.result, msg)
			}
		})
	}
}

func TestGetExchgConnection(t *testing.T) {
	tt := []struct {
		name                            string
		method                          string
		status                          int
		err                             string
		gatewayGetExchgConnectionResult []string
		result                          []string
	}{
		{
			name:   "405",
			method: http.MethodPost,
			status: http.StatusMethodNotAllowed,
			err:    "405 Method Not Allowed",
		},
		{
			name:                            "200",
			method:                          http.MethodGet,
			status:                          http.StatusOK,
			err:                             "",
			gatewayGetExchgConnectionResult: []string{"44.33.22.11", "11.44.66.88"},
			result:                          []string{"11.44.66.88", "44.33.22.11"},
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			endpoint := "/api/v1/network/connections/exchange"
			gateway := &MockGatewayer{}
			gateway.On("GetExchgConnection").Return(tc.gatewayGetExchgConnectionResult)

			req, err := http.NewRequest(tc.method, endpoint, nil)
			require.NoError(t, err)

			rr := httptest.NewRecorder()
			handler := newServerMux(defaultMuxConfig(), gateway, &CSRFStore{}, nil)
			handler.ServeHTTP(rr, req)

			status := rr.Code
			require.Equal(t, tc.status, status, "got `%v` want `%v`", status, tc.status)

			if status != http.StatusOK {
				require.Equal(t, tc.err, strings.TrimSpace(rr.Body.String()), "got `%v`| %d, want `%v`",
					strings.TrimSpace(rr.Body.String()), status, tc.err)
			} else {
				var msg []string
				err = json.Unmarshal(rr.Body.Bytes(), &msg)
				require.NoError(t, err)
				require.Equal(t, tc.result, msg)
			}
		})
	}
}
