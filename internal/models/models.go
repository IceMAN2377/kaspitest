package models

type User struct {
	Name  string `json:"name" db:"name"`
	IIN   string `json:"iin" db:"iin"`
	Phone string `json:"phone" db:"phone"`
}

type IINInfo struct {
	Correct   bool   `json:"correct"`
	Gender    string `json:"gender"`
	Birthdate string `json:"birthdate"`
}
