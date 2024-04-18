package models

import "fmt"

// TransactionType is a custom type for transaction types
type TransactionType int64

const (
	Initiate TransactionType = iota + 1
	Authorize
	InProgress
	Complete
	Rejected
	Cancelled
	Refunded
)

var transactionTypes = map[TransactionType]string{
	Initiate:   "Initiate",
	Authorize:  "Authorize",
	InProgress: "In Progress",
	Complete:   "Complete",
	Rejected:   "Rejected",
	Cancelled:  "Cancelled",
	Refunded:   "Refunded",
}

// String returns the string representation of the transaction type
func (t TransactionType) String() string {
	return transactionTypes[t]
}

// Value returns the value of the transaction type
func (t TransactionType) Value() int64 {
	return int64(t)
}

// ParseTransactionType returns the transaction type based on the value
func ParseTransactionType(value int64) (*TransactionType, error) {
	if value > int64(len(transactionTypes)) {
		return nil, fmt.Errorf("there is not any transaction type with that ID: %v", value)
	}

	tranType := TransactionType(value)

	_, ok := transactionTypes[tranType]
	if !ok {
		return nil, fmt.Errorf("there is not any transaction type with that ID: %v", value)
	}

	return &tranType, nil
}
