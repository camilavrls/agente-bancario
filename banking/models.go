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
