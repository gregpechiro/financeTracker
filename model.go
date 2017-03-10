package main

type User struct {
	Id         string              `json:"id"`
	Email      string              `json:"email,omitempty" auth:"username" required:"register, login, account"`
	Password   string              `json:"password,omitempty" auth:"password" required:"register, login"`
	Active     bool                `json:"active" auth:"active"`
	Role       string              `json:"role,omitempty"`
	FirstName  string              `json:"firstName,omitempty" required:"register, account"`
	LastName   string              `json:"lastName,omitempty" required:"register, account"`
	Created    int64               `json:"created,omitempty"`
	LastSeen   int64               `json:"lastSeen,omitempty"`
	AccountId  string              `json:"accountId,omitempty"`
	Primary    bool                `json:"primary"`
	Categories map[string]struct{} `json:"categories,omitempty"`
	People     map[string]struct{} `json:"people,omitempty"`
	Phone      string              `json:"phone,omitempty"`
}

type Transaction struct {
	Id                 string  `json:"id"`
	UserId             string  `json:"userId"`
	AccountId          string  `json:"accountId,omitempty"`
	DateAdded          int64   `json:"dateAdded,omitempty"`
	Date               int64   `json:"date,omitempty"`
	Title              string  `json:"title,omitempty" required:"transaction"`
	Description        string  `json:"description,omitempty"`
	Ammount            float32 `json:"ammount,omitempty" required:"transaction"`
	Who                string  `json:"who,omitempty"`
	Category           string  `json:"category,omitempty"`
	SecondaryCategory  string  `json:"secondaryCategory,omitempty"`
	QuickTransactionId string  `json:"quickTransactionId,omitempty"`
}

/*type Category struct {
	Id            string `json:"id"`
	AccountId     string `json:"accountId"`
	Title         string `json:"title,omitemtpy"`
	TransactionId string `json:"transactionId, omitempty"`
}*/
