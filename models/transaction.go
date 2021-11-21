package models

type Account struct {
	User_id    string
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
