package transaction

import "github.com/deuna/payment/infraestructure/db/repos"

type PlaceRequest struct {
	Amount float64 `json:"amount"`
}

type RefundRequest struct {
	ID string `json:"id"`
}

type ListTxResponse struct {
	Txs []repos.Tx `json:"transactions"`
}
