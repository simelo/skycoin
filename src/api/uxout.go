package api

import (
	"net/http"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/readable"
	wh "github.com/skycoin/skycoin/src/util/http"
)

// URI: /api/v1/uxout
// Method: GET
// Args:
//	uxid: unspent output ID hash
// Returns an unspent output by ID
func uxOutHandler(gateway Gatewayer) http.HandlerFunc {

	// swagger:operation GET /api/v1/uxout uxout
	//
	// Returns an unspent output by ID.
	//
	// ---
	//
	// produces:
	// - application/json
	// parameters:
	// - name: uxid
	//   in: query
	//   description: uxid to filter by
	//   type: string
	// responses:
	//   200:
	//     description: Response for endpoint /api/v1/uxout
	//     schema:
	//       properties:
	//         uxid:
	//           type: string
	//         time:
	//           type: integer
	//           format: int64
	//         src_block_seq:
	//           type: integer
	//           format: int64
	//         src_tc:
	//           type: string
	//         owner_address:
	//           type: string
	//         coins:
	//           type: integer
	//           format: integer
	//         hours:
	//           type: integer
	//           format: int64
	//         spent_block_seq:
	//           type: integer
	//           format: in64
	//         spent_tx:
	//           type: string
	//   default:
	//     $ref: '#/responses/genericError'


	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			wh.Error405(w)
			return
		}

		uxid := r.FormValue("uxid")
		if uxid == "" {
			wh.Error400(w, "uxid is empty")
			return
		}

		id, err := cipher.SHA256FromHex(uxid)
		if err != nil {
			wh.Error400(w, err.Error())
			return
		}

		uxout, err := gateway.GetUxOutByID(id)
		if err != nil {
			wh.Error400(w, err.Error())
			return
		}

		if uxout == nil {
			wh.Error404(w, "")
			return
		}

		wh.SendJSONOr500(logger, w, readable.NewSpentOutput(uxout))
	}
}

// URI: /api/v1/address_uxouts
// Method: GET
// Args:
//	address
// Returns the historical, spent outputs associated with an address
func addrUxOutsHandler(gateway Gatewayer) http.HandlerFunc {

	// swagger:operation GET /api/v1/address_uxouts addressUxouts
	//
	// Returns the historical, spent outputs associated with an address
	//
	// ---
	//
	// produces:
	// - application/json
	// parameters:
	// - name: address
	//   in: query
	//   description: address to filter by
	//   type: string
	// responses:
	//   200:
	//     description: Response for endpoint /api/v1/address_uxouts
	//     schema:
	//       type: array
	//       items:
	//         properties:
	//           uxid:
	//             type: string
	//           time:
	//             type: integer
	//             format: int64
	//           src_block_seq:
	//             type: integer
	//             format: int64
	//           src_tc:
	//             type: string
	//           owner_address:
	//             type: string
	//           coins:
	//             type: integer
	//             format: integer
	//           hours:
	//             type: integer
	//             format: int64
	//           spent_block_seq:
	//             type: integer
	//             format: in64
	//           spent_tx:
	//             type: string
	//   default:
	//     $ref: '#/responses/genericError'

	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			wh.Error405(w)
			return
		}

		addr := r.FormValue("address")
		if addr == "" {
			wh.Error400(w, "address is empty")
			return
		}

		cipherAddr, err := cipher.DecodeBase58Address(addr)
		if err != nil {
			wh.Error400(w, err.Error())
			return
		}

		uxs, err := gateway.GetSpentOutputsForAddresses([]cipher.Address{cipherAddr})
		if err != nil {
			wh.Error400(w, err.Error())
			return
		}

		ret := make([]readable.SpentOutput, 0)
		for _, u := range uxs {
			ret = append(ret, readable.NewSpentOutputs(u)...)
		}

		wh.SendJSONOr500(logger, w, ret)
	}
}
