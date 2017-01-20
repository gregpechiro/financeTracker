package main

type User struct {
	Id          string `json:"id"`
	Email       string `json:"email,omitempty" auth:"username"`
	Password    string `json:"password,omitempty" auth:"password"`
	Active      bool   `json:"active" auth:"active"`
	Role        string `json:"role,omitempty"`
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
	Title       string `json:"title,omitempty"`
	Description string `json:"description,omitempty"`
	Category    string `json:"category,omitempty"`
	Amount      int    `json:"amoun,omitemptyt"`
	Date        int64  `json:"date,omitempty"`
	AccountId   string `json:"accountId,omitempty"`
}

type Category struct {
	Id        string `json:"id"`
	AccountId string `json:"accountId"`
	Title     string `json:"title,omitemtpy"`
}
