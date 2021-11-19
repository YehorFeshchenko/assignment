package models

type Account struct {
	UserId     string
	USDBalance float64
	EURBalance float64
}

type Transaction struct {
	OperationTypeStr string  `json:"type"`
	UserId           string  `json:"user_id"`
	Amount           float64 `json:"amount"`
	TimePlaced       string  `json:"time_placed"`
	CurrencyStr      string  `json:"currency"`
}

type Currency int

const (
	USD Currency = iota
	EUR
)

type OperationType int

const (
	DEPOSIT OperationType = iota
	WITHDRAWAL
)

const FACTOR = 1000
