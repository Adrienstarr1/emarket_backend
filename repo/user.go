package repo

import (
	"context"
	"e-market/model"
	"fmt"
	"log"
	"strings"
)

func (r *Repo) AddUser(ctx context.Context, user model.User) error {
	_, err := r.Pool.Exec(ctx, "INSERT INTO users (id, email, age, name, password, created_at, admin) VALUES ($1, $2, $3, $4, $5, $6, $7)",
		user.Id, user.Email, user.Age, user.Name, user.Password, user.Created_at, user.Admin)
	if err != nil {
		log.Println(err)
		return err
	}
	log.Println("User added")
	return nil
}

func (r *Repo) ListUsers(ctx context.Context, category string, info any) ([]model.User, error) {
	var response model.User
	responselist := make([]model.User, 0)
	query := fmt.Sprintf("SELECT name, age, email, password, id, created_at, admin FROM users WHERE %s = $1", category)

	rows, err := r.Pool.Query(ctx, query, info)
	if err != nil {
		return []model.User{}, err
	}
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&response.Name, &response.Age, &response.Email, &response.Password, &response.Id, &response.Created_at, &response.Admin)
		if err != nil {
			return []model.User{}, err
		}
		responselist = append(responselist, response)
	}

	return responselist, nil
}

// dynamically update users
// better for handling multiple updates in 1 go ;)
func (r *Repo) UpdateUser_V3(ctx context.Context, id string, instructions map[string]any) error {
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

	cmdtag, err := r.Pool.Exec(ctx, sqlString, values...)
	if err != nil {
		return err
	}
	if cmdtag.RowsAffected() == 0 {
		return NRA
	}

	return nil
}

func (r *Repo) DeleteUser(ctx context.Context, id string) error {
	cmdtag, err := r.Pool.Exec(ctx, "DELETE FROM users WHERE id = $1", id)
	if err != nil {
		return err
	}
	if cmdtag.RowsAffected() == 0 {
		return NRA
	}
	return nil
}
