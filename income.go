package lknpd

import "time"

const (
	defaultPaymentType = "CASH"
)

type Income struct {
	Name     string  `json:"name"`
	Amount   float64 `json:"amount"`
	Quantity uint    `json:"quantity"`
}

type IncomeClient struct {
	ContactPhone *string          `json:"contactPhone"`
	DisplayName  *string          `json:"displayName"`
	IncomeType   IncomeClientType `json:"incomeType"`
	INN          *string          `json:"inn"`
}

type IncomeClientType string

var (
	// Физическому лицу
	IncomeClientTypeFromIndividual IncomeClientType = "FROM_INDIVIDUAL"
	// Юридическому лицу. Необходимо указать INN и DisplayName
	IncomeClientTypeFromLegalEntity IncomeClientType = "FROM_LEGAL_ENTITY"
)

type CreateIncomeRequest struct {
	Client                           IncomeClient `json:"client"`
	IgnoreMaxTotalIncomeRestrictions bool         `json:"ignoreMaxTotalIncomeRestrictions"`
	OperationTime                    time.Time    `json:"operationTime"`
	PaymentType                      string       `json:"paymentType"`
	RequestTime                      time.Time    `json:"requestTime"`
	Services                         []Income     `json:"services"`
	TotalAmount                      string       `json:"totalAmount"`
}

type CreateIncomeResponse struct {
	ApprovedReceiptUUID string `json:"approvedReceiptUuid"`
}
