package banking

type CustomerProfile struct {
	ID          string
	Name        string
	Segment     string
	CreditScore int
}

type CardLimit struct {
	CustomerID      string
	CurrentLimit    int
	AvailableLimit  int
	MaxAllowedLimit int
}

type AccountBalance struct {
	CustomerID   string
	BalanceCents int
}

type PixTransaction struct {
	ID             string
	FromCustomerID string
	ToPixKey       string
	AmountCents    int
	Status         string
}
