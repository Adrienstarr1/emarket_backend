package repo

import (
	"context"
	"e-market/model"
	"errors"
	"fmt"
	"strings"
)

func (r *Repo) Add2cart(ctx context.Context, instructions map[string]any) error {
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

	cmdtag, err := r.Pool.Exec(ctx, sqlString, values...)
	if err != nil {
		return err
	}
	if cmdtag.RowsAffected() == 0 {
		return NRA
	}

	return nil

}

func (r *Repo) FindCart(ctx context.Context, user_id string) ([]model.Cart, error) {
	var response model.Cart
	responselist := make([]model.Cart, 0)
	query := "SELECT user_id, product_id, quantity FROM cart WHERE user_id = $1"

	rows, err := r.Pool.Query(ctx, query, user_id)
	if err != nil {
		return []model.Cart{}, err
	}
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&response.User_id, &response.Product_id, &response.Quantity)
		if err != nil {
			return []model.Cart{}, err
		}
		responselist = append(responselist, response)
	}

	return responselist, nil
}

func (r *Repo) UpdateCart(ctx context.Context, user_id, product_id string, operation int) error {
	var sqlstring string

	switch operation {
	case 0:
		sqlstring = "UPDATE cart SET quantity = CASE WHEN quantity > 0 THEN quantity - 1 ELSE 0 END WHERE user_id = $1 AND product_id = $2 "
	case 1:
		sqlstring = "UPDATE cart SET quantity = quantity + 1 WHERE user_id = $1 AND product_id = $2"
	default:
		return errors.New("Invalid operation: Operations must be 1 or 0")
	}

	cmdtag, err := r.Pool.Exec(ctx, sqlstring, user_id, product_id)
	if err != nil {
		return err
	}
	if cmdtag.RowsAffected() == 0 {
		return NRA
	}
	return nil
}

func (r *Repo) DeleteCart(ctx context.Context, user_id, product_id string) error {
	if cmdtag, err := r.Pool.Exec(ctx, "DELETE FROM cart WHERE user_id = $1 AND product_id = $2", user_id, product_id); err != nil ||
		cmdtag.RowsAffected() == 0 {
		return err
	}
	return nil
}
