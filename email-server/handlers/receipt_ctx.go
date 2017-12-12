package handlers

import "github.com/KyleIWS/EmailReceipt/email-server/models"

func NewReceiptCtx(ms *models.MongoStore) *ReceiptCtx {
	return &ReceiptCtx{
		ms: ms,
	}
}

type ReceiptCtx struct {
	ms *models.MongoStore
}
