package clients

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type Bank interface {
	Place(context.Context, uuid.UUID, float64, uuid.UUID) error
	Refund(context.Context, uuid.UUID) error
}

type bank struct{}

func NewBank() Bank {
	return &bank{}
}

// This is a mock implementation of the bank service
func (b *bank) Place(ctx context.Context, id uuid.UUID, amount float64, mercahntID uuid.UUID) error {
	time.Sleep(time.Second)
	return nil
}

// This is a mock implementation of the bank service
func (b *bank) Refund(ctx context.Context, id uuid.UUID) error {
	time.Sleep(time.Second)
	return nil
}
