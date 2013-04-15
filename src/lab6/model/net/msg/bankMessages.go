package msg

import (
	"encoding/gob"
)

func init() {
	gob.Register(Deposit{})
	gob.Register(Withdraw{})
	gob.Register(Transfer{})
	gob.Register(Balance{})
}

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
