package gui

// Network-related information for the GUI
import (
	"fmt"
	"net/http"
	"sort"

	wh "github.com/skycoin/skycoin/src/util/http" //http,json helpers
)

// InfoResponse encapsulates useful information from the ...
type InfoResponse struct {
	StatusListEnable []string `json:"list_status_enable"`
	// StatusListDisable      []string `json:"list_status_disable"`
	DefaultConnectionCount int `json:"default_connection_count"`
	OpenConnectionCount    int `json:"open_connection_count"`
}

func connectionHandler(gateway Gatewayer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			wh.Error405(w)
			return
		}

		addr := r.FormValue("addr")
		if addr == "" {
			wh.Error400(w, "addr is required")
			return
		}

		c := gateway.GetConnection(addr)
		if c == nil {
			wh.Error404(w)
			return
		}

		wh.SendJSONOr500(logger, w, c)
	}
}

func connectionsHandler(gateway Gatewayer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			wh.Error405(w)
			return
		}

		wh.SendJSONOr500(logger, w, gateway.GetConnections())
	}
}

func defaultConnectionsHandler(gateway Gatewayer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			wh.Error405(w)
			return
		}

		conns := gateway.GetDefaultConnections()
		sort.Strings(conns)

		wh.SendJSONOr500(logger, w, conns)
	}
}

func trustConnectionsHandler(gateway Gatewayer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			wh.Error405(w)
			return
		}

		conns := gateway.GetTrustConnections()
		sort.Strings(conns)

		wh.SendJSONOr500(logger, w, conns)
	}
}

func exchgConnectionsHandler(gateway Gatewayer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			wh.Error405(w)
			return
		}

		conns := gateway.GetExchgConnection()
		sort.Strings(conns)

		wh.SendJSONOr500(logger, w, conns)
	}
}

// TODO Function to obtain the issue info # 1049
func infoHandler(gateway Gatewayer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			wh.Error405(w)
			return
		}

		listdefault := gateway.GetDefaultConnections()
		listconnections := gateway.GetConnections().Connections

		connectionCount := len(listconnections)
		defaultconnectionCount := len(listdefault)
		var enable []string
		countenable := 0

		for _, tmpdefault := range listdefault {

			for _, tmpconnection := range listconnections {

				if tmpdefault == string(tmpconnection.Addr) {
					enable = append(enable, tmpdefault)
					countenable++
				}
			}
		}

		resp := &InfoResponse{
			StatusListEnable:       enable,
			DefaultConnectionCount: defaultconnectionCount,
			OpenConnectionCount:    connectionCount,
		}

		fmt.Println(resp)

		if resp == nil {
			wh.Error404(w)
			return
		}

		wh.SendJSONOr500(logger, w, &resp)

	}
}
