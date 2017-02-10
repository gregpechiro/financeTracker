package main

import "github.com/cagnosolutions/adb"

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
	Id            string `json:"id"`
	UserId        string `json:"userId"`
	Title         string `json:"title,omitempty" required:"transaction"`
	Description   string `json:"description,omitempty" required:"transaction"`
	BudgetItem    string `json:"budgetItem,omitempty" required:"transaction"`
	Amount        int    `json:"amount,omitempty" required:"transaction"`
	Date          int64  `json:"date,omitempty"`
	AccountId     string `json:"accountId,omitempty"`
	SubcategoryId string `json:"subcategoryId"`
	CategoryId    string `json:"categoryId"`
}

type Subcategory struct {
	Id         string `json:"id"`
	CategoryId string `json:"categoryId"`
	AccountId  string `json:"accountId"`
	Title      string `json:"title,omitemtpy" required:"item"`
	Planned    string `json:"planned,omitempty"`
}

type Category struct {
	Id        string `json:"id"`
	AccountId string `json:"accountId"`
	Title     string `json:"title,omitemtpy" required:"group"`
}

type CategoryView struct {
	Category      Category
	Subcategories []Subcategory
}

func getCategoryView(accountId string) []CategoryView {
	var categories []Category
	var categoryViews []CategoryView
	db.TestQuery("category", &categories, adb.Eq("accountId", `"`+accountId+`"`))

	for _, category := range categories {
		var categoryView CategoryView
		var subcategories []Subcategory
		db.TestQuery("subcategory", &subcategories, adb.Eq("categoryId", `"`+category.Id+`"`))
		categoryView.Category = category
		categoryView.Subcategories = subcategories
		categoryViews = append(categoryViews, categoryView)
	}
	return categoryViews

}
