package repository

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

type UserResponse struct {
	Name       string
	Age        int
	Email      string
	Password   string
	Id         string
	Created_at time.Time
}

type ProductResponse struct {
	Id           string
	Product_name string
	Quantity     int
	User_Id      string
	Created_at   time.Time
	Updated_at   time.Time
	Price        int
	Cost         int
	Description  string
}

type CartResponse struct {
	Quantity   int
	User_id    string
	Product_id string
}

var NRA = errors.New("No rows affected prolem with entry")

func Connect(ctx context.Context) (*pgxpool.Pool, error) {
	if err := godotenv.Load(); err != nil {
		log.Panicf("Problem loading from .env\n%v", err)
	}
	connString := os.Getenv("DB_URL")
	pool, err := pgxpool.New(ctx, connString)
	if err != nil {
		log.Fatal("Unable to connect to database:", err)
		return nil, err
	}

	if err = pool.Ping(ctx); err != nil {
		log.Println("DB is not connecting")
		return nil, err
	}
	return pool, nil
}

func AddUser(ctx context.Context, pool *pgxpool.Pool, info ...any) error {
	id := info[4]
	email := info[2]
	age := info[1]
	name := info[0]
	password := info[3]

	_, err := pool.Exec(ctx, "INSERT INTO users (id, email, age, name, password, created_at) VALUES ($1, $2, $3, $4, $5, $6)", id, email, age, name, password, time.Now())
	if err != nil {
		log.Println(err)
		return err
	}
	log.Println("User added")
	return nil
}

func AddProduct(ctx context.Context, pool *pgxpool.Pool, info ...any) error {
	id := info[0]
	name := info[1]
	quantity := info[2]
	user_id := info[3]
	price := info[4]
	cost := info[5]
	description := info[6]
	sql_string := "INSERT INTO products (id, name, quantity, user_id, created_at, price, cost, description) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)"

	_, err := pool.Exec(ctx, sql_string, id, name, quantity, user_id, time.Now(), price, cost, description)
	if err != nil {
		log.Println("Error in storing new product", err)
		return err
	}
	log.Println("Products added successfully")
	return nil
}

func FindProduct_V2(ctx context.Context, pool *pgxpool.Pool, category string, info any) ([]ProductResponse, error) {
	var response ProductResponse
	var responselist []ProductResponse
	sqlString := fmt.Sprintf("SELECT id, name, quantity, user_id, price, description FROM products WHERE %s = $1", category)

	rows, err := pool.Query(ctx, sqlString, info)
	if err != nil {
		return []ProductResponse{}, err
	}

	for rows.Next() {
		err = rows.Scan(&response.Id, &response.Product_name, &response.Quantity, &response.User_Id, &response.Price, &response.Description)
		if err != nil {
			log.Println("problem sending product info \n", err)
			return []ProductResponse{}, err
		}
		responselist = append(responselist, response)
	}

	return responselist, nil
}

func ListUsers(ctx context.Context, pool *pgxpool.Pool, category string, info any) ([]UserResponse, error) {
	var response UserResponse
	responselist := make([]UserResponse, 0)
	query := fmt.Sprintf("SELECT name, age, email, password, id, created_at FROM users WHERE %s = $1", category)

	rows, err := pool.Query(ctx, query, info)
	if err != nil {
		return []UserResponse{}, err
	}

	for rows.Next() {
		err = rows.Scan(&response.Name, &response.Age, &response.Email, &response.Password, &response.Id, &response.Created_at)
		if err != nil {
			return []UserResponse{}, err
		}
		responselist = append(responselist, response)
	}

	return responselist, nil
}

// dynamically update users
// better for handling multiple updates in 1 go ;)
func UpdateUser_V3(ctx context.Context, pool *pgxpool.Pool, id string, instructions map[string]any) error {
	cmds := make([]string, 0) // store commands i.e the column and the paramatized change
	values := make([]any, 0)  // store users paramatized channges-

	// counter is used to keep track of the changes we make
	// it is also used to keep track of the ${number} when make sql commands
	counter := 1
	for key, value := range instructions {
		str := fmt.Sprintf("%s = $%d", key, counter) // format string to make sql commands e.g 'age' = 45
		cmds = append(cmds, str)                     // store the command in slice
		values = append(values, value)               //store its corresponding value in slice
		counter++                                    // increment counter for next input
	}

	// join commands to together while segregating them with a comma
	cmd := strings.Join(cmds, ",")

	// id is the final parameter
	// it doesnt need to increment counter since it is always going to be the last input
	sqlString := fmt.Sprintf("UPDATE users SET %s WHERE id = $%d", cmd, counter)

	// id is always the last input
	values = append(values, id)

	cmdtag, err := pool.Exec(ctx, sqlString, values...)
	if err != nil {
		return err
	}
	if cmdtag.RowsAffected() == 0 {
		return NRA
	}

	return nil
}

// dynamically update products
// better for handling multiple updates in 1 go ;)
func UpdateProduct_V2(ctx context.Context, pool *pgxpool.Pool, id string, instructions map[string]any) error {
	cmds := make([]string, 0) // store commands i.e the column and the paramatized change
	values := make([]any, 0)  // store users paramatized channges-

	// counter is used to keep track of the changes we make
	// it is also used to keep track of the ${number} when make sql commands
	counter := 1
	for key, value := range instructions {
		str := fmt.Sprintf("%s = $%d", key, counter) // format string to make sql commands e.g 'age' = 45
		cmds = append(cmds, str)                     // store the command in slice
		values = append(values, value)               // store its corresponding value in slice
		counter++                                    // increment counter for next input
	}

	// join commands to together while segregating them with a comma
	cmd := strings.Join(cmds, ",")

	// id is the final parameter
	// it doesnt need to increment counter since it is always going to be the last input
	sqlString := fmt.Sprintf("UPDATE products SET %s WHERE id = $%d", cmd, counter)

	// id is always the last input
	values = append(values, id)

	cmdtag, err := pool.Exec(ctx, sqlString, values...)
	if err != nil {
		return err
	}
	if cmdtag.RowsAffected() == 0 {
		return NRA
	}

	return nil
}

func DeleteUser(ctx context.Context, pool *pgxpool.Pool, id string) error {
	if cmdtag, err := pool.Exec(ctx, "DELETE FROM users WHERE id = $1", id); err != nil ||
		cmdtag.RowsAffected() == 0 {
		return err
	}
	return nil
}

func DeleteProduct(ctx context.Context, pool *pgxpool.Pool, id string) error {
	if cmdtag, err := pool.Exec(ctx, "DELETE FROM users WHERE id = $1", id); err != nil ||
		cmdtag.RowsAffected() == 0 {
		return err
	}
	return nil
}

func Add2cart(ctx context.Context, pool *pgxpool.Pool, instructions map[string]any) error {
	cmds := make([]string, 0)
	values := make([]any, 0)
	parameters := make([]string, 0)
	counter := 1

	for key, value := range instructions {
		cmds = append(cmds, key)
		values = append(values, value)
		parameters = append(parameters, fmt.Sprintf("$%d", counter))
		counter++
	}

	cmd := strings.Join(cmds, ",")
	parameter := strings.Join(parameters, ",")

	sqlString := fmt.Sprintf("INSERT INTO cart (%s) VALUES (%s)", cmd, parameter)

	cmdtag, err := pool.Exec(ctx, sqlString, values...)
	if err != nil {
		return err
	}
	if cmdtag.RowsAffected() == 0 {
		return NRA
	}

	return nil

}

func FindCart(ctx context.Context, pool *pgxpool.Pool, user_id string) ([]CartResponse, error) {
	var response CartResponse
	responselist := make([]CartResponse, 0)
	query := "SELECT user_id, product_id, quantity FROM cart WHERE user_id = $1"

	rows, err := pool.Query(ctx, query, user_id)
	if err != nil {
		return []CartResponse{}, err
	}

	for rows.Next() {
		err = rows.Scan(&response.User_id, &response.Product_id, &response.Quantity)
		if err != nil {
			return []CartResponse{}, err
		}
		responselist = append(responselist, response)
	}

	return responselist, nil
}

func UpdateCart(ctx context.Context, pool *pgxpool.Pool, user_id, product_id string) error {
	cmdtag, err := pool.Exec(ctx, "UPDATE cart SET quantity = quantity + 1 WHERE user_id = $1 AND product_id = $2", user_id, product_id)
	if err != nil {
		return err
	}
	if cmdtag.RowsAffected() == 0 {
		return NRA
	}
	return nil
}
