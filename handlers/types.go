package handlers

type SigninRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type CreateAccountRequest struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

// type TransferRequest struct {
// 	ToAccount int `json:"toAccount"`
// 	Amount    int `json:"amount"`
// }
