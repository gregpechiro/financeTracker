package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"
)

func genId() string {
	return strconv.Itoa(int(time.Now().UnixNano()))
}

func OrderCategories(m map[string]struct{}) []string {
	var ss []string
	for k, _ := range m {
		ss = append(ss, k)
	}
	sort.Strings(ss)
	return ss
}

func NormalIzeString(s string) string {
	return strings.Join(strings.Fields(strings.ToLower(s)), " ")
}

func PrettyDate(ts int64) string {
	if ts == 0 {
		return ""
	}
	t := time.Unix(ts, 0)
	return t.Format("1/2/2006")
}

func PrettyDateTime(ts int64) string {
	if ts == 0 {
		return ""
	}
	t := time.Unix(ts, 0)
	return t.Format("1/2/2006 3:04 PM")
}

func isIncome(ammount float32) bool {
	return ammount > float32(0)
}

func ajaxResponse(w http.ResponseWriter, v interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	b, err := json.Marshal(v)
	if err != nil {
		b = []byte(`{"error":true,output:"Error"}`)
	}
	fmt.Fprint(w, string(b))
}

func toJson(v interface{}) string {
	b, err := json.Marshal(v)
	if err != nil {
		return ""
	}
	return string(b)
}

func toBase64Json(v interface{}) string {
	b, err := json.Marshal(v)
	if err != nil {
		return ""
	}
	return base64.StdEncoding.EncodeToString(b)
}

func defaultUsers() {
	admin := User{
		Id:        "0",
		Role:      "ADMIN",
		FirstName: "Admin",
		LastName:  "Temporary",
		Email:     "admin@temp.com",
		Password:  "admin",
		Active:    true,
		Created:   time.Now().Unix(),
		LastSeen:  time.Now().Unix(),
	}

	db.Set("user", "0", admin)

	fmt.Printf("\nTemporary admin credentials:\n\n\tEmail:\t\t%s\n\tPassword:\t%s\n\n", admin.Email, admin.Password)
}
