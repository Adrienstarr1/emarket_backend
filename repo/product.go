package repo

import (
	"context"
	"e-market/model"
	"fmt"
	"log"
	"strings"
	"time"
)

func (r *Repo) AddProduct_V2(ctx context.Context, instructions map[string]any) error {
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

	sqlString := fmt.Sprintf("INSERT INTO products (%s) VALUES (%s)", cmd, parameter)

	cmdtag, err := r.Pool.Exec(ctx, sqlString, values...)
	if err != nil {
		return err
	}
	if cmdtag.RowsAffected() == 0 {
		return NRA
	}

	return nil
}

func (r *Repo) AddProduct(ctx context.Context, info ...any) error {
	id := info[0]
	name := info[1]
	quantity := info[2]
	user_id := info[3]
	price := info[4]
	cost := info[5]
	description := info[6]
	sql_string := "INSERT INTO products (id, name, quantity, user_id, created_at, price, cost, description) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)"

	_, err := r.Pool.Exec(ctx, sql_string, id, name, quantity, user_id, time.Now(), price, cost, description)
	if err != nil {
		log.Println("Error in storing new product", err)
		return err
	}
	log.Println("Products added successfully")
	return nil
}

func (r *Repo) FindProduct_V2(ctx context.Context, category string, info any) ([]model.Product, error) {
	var response model.Product
	var responselist []model.Product
	sqlString := fmt.Sprintf("SELECT id, name, quantity, user_id, price, description FROM products WHERE %s = $1", category)

	rows, err := r.Pool.Query(ctx, sqlString, info)
	if err != nil {
		return []model.Product{}, err
	}
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&response.Id, &response.Product_name, &response.Quantity, &response.User_Id, &response.Price, &response.Description)
		if err != nil {
			log.Println("problem sending product info \n", err)
			return []model.Product{}, err
		}
		responselist = append(responselist, response)
	}

	return responselist, nil
}

// dynamically update products
// better for handling multiple updates in 1 go ;)
func (r *Repo) UpdateProduct_V2(ctx context.Context, id string, instructions map[string]any) error {
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

	cmdtag, err := r.Pool.Exec(ctx, sqlString, values...)
	if err != nil {
		return err
	}
	if cmdtag.RowsAffected() == 0 {
		return NRA
	}

	return nil
}

func (r *Repo) DeleteProduct(ctx context.Context, id string) error {
	cmdtag, err := r.Pool.Exec(ctx, "DELETE FROM products WHERE id = $1", id)
	if err != nil {
		return err
	}
	if cmdtag.RowsAffected() == 0 {
		return NRA
	}
	return nil
}
