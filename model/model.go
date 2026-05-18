package model

import "time"

type Product struct {
	Id           string
	Product_name string
	Quantity     int
	User_Id      string
	Created_at   time.Time
	Updated_at   time.Time
	Price        int
	Description  string
}

type User struct {
	Name       string
	Age        int
	Email      string
	Password   string
	Id         string
	Created_at time.Time
	Admin      bool
}

type Cart struct {
	Quantity   int
	User_id    string
	Product_id string
}
