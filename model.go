package main

type User struct {
	Id          string `json:"id"`
	Email       string `json:"email,omitempty" auth:"username" required:"register, login"`
	Password    string `json:"password,omitempty" auth:"password" required:"register, login"`
	Active      bool   `json:"active" auth:"active"`
	Role        string `json:"role,omitempty"`
	FirstName   string `json:"firstName,omitempty" required:"register"`
	LastName    string `json:"lastName,omitempty" required:"register"`
	Created     int64  `json:"created,omitempty"`
	LastSeen    int64  `json:"lastSeen,omitempty"`
	Phone       string `json:"phone,omitempty"`
	AccountId   string `json:"accountId,omitempty"`
	Primary     bool   `json:"primary"`
	RecoveryPin int    `json:"recoveryPin,omitempty"`
}

type Transaction struct {
	Id          string `json:"id"`
	UserId      string `json:"userId"`
	Title       string `json:"title,omitempty" required:"transcation"`
	Description string `json:"description,omitempty" required:"transcation"`
	BudgetItem  string `json:"budgetItem,omitempty" required:"transcation"`
	Amount      int    `json:"amount,omitempty" required:"transcation"`
	Date        int64  `json:"date,omitempty"`
	AccountId   string `json:"accountId,omitempty"`
}

type BudgetItem struct {
	Id        string `json:"id"`
	AccountId string `json:"accountId"`
	Title     string `json:"title,omitemtpy" required:"item"`
	Planned   string `json:"planned,omitempty" required:"item"`
}

type BudgetGroup struct {
	Id          string       `json:"id"`
	AccountId   string       `json:"accountId"`
	Title       string       `json:"title,omitemtpy" required:"group"`
	BudgetItems []BudgetItem `json:"budgetItems,omitempty"`
}
