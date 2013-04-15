package msg

type Deposit struct {
	AccountNumber string
	Amount        int
}

type Withdraw struct {
	AccountNumber string
	Amount        int
}

type Transfer struct {
	FromAccount string
	ToAccount   string
	Amount      int
}

type Balance struct {
	AccountNumber string
}
