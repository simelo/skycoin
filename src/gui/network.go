package gui

// Network-related information for the GUI
import (
	"net/http"
	"sort"

	wh "github.com/skycoin/skycoin/src/util/http" //http,json helpers
)

// Status encapsulates useful information from the ...
type Status struct {
	address string `json:"ip:port"`
	status  bool   `json:"is_conections"`
}

// InfoResponse encapsulates useful information from the ...
type InfoResponse struct {
	StatusList             []Status `json:"list_status"`
	DefaultConnectionCount int      `json:"default_connection_count"`
	OpenConnectionCount    int      `json:"open_connection_count"`
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
		list := make([]Status, 0)
		var tmpstruct Status

		// var tmpstatus bool

		for _, tmpdefault := range listdefault {

			for _, tmpconnection := range listconnections {

				if tmpdefault == string(tmpconnection.Addr) {
					tmpstatus = true
				}

			}

			tmpstruct = Status{address: tmpdefault, status: tmpstatus}

			list = append(list, tmpstruct)

		}

		resp := &InfoResponse{
			StatusList:             list,
			DefaultConnectionCount: defaultconnectionCount,
			OpenConnectionCount:    connectionCount,
		}

		if resp == nil {
			wh.Error404(w)
			return
		}

		wh.SendJSONOr500(logger, w, &resp)
	}
}
