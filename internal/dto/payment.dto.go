package dto

type CreateTopupRequest struct {
	Amount          int     `json:"amount" binding:"required,gt=0"`
	PaymentMethodID *string `json:"payment_method_id"`
	Note            *string `json:"note"`
}

type CreateTopupResponse struct {
	TrxCode    string `json:"trx_code"`
	Amount     int    `json:"amount"`
	Total      int    `json:"total"`
	Status     string `json:"status"`
	PaymentURL string `json:"payment_url"`
}

// Struct penampung Notification Webhook Callback dari DOKU
type DokuNotificationPayload struct {
	Order struct {
		InvoiceNumber string `json:"invoice_number"`
		Amount        int    `json:"amount"`
	} `json:"order"`
	Transaction struct {
		Status string `json:"status"`
	} `json:"transaction"`
}