package service

import (
	"context"
	"e-market/auth"
	"e-market/model"
	"errors"
	"log"
	"slices"

	"golang.org/x/crypto/bcrypt"
)

func (s *Service) Signup(ctx context.Context, user model.User) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Println(err)
		return "", err
	}

	user.Password = string(hash)

	err = s.Repo.AddUser(ctx, user)
	if err != nil {
		log.Println(err)
		return "", err
	}

	ss, err := auth.CreateSS(user)
	if err != nil {
		log.Println(err)
		return "", err
	}

	return ss, nil
}

func (s *Service) Login(ctx context.Context, input model.User) (string, error) {
	if input.Email == "" || input.Password == "" {
		log.Println("Invalid Inputs")
		return "", errors.New("Invalid Inputs missing field email or password")
	}

	responses, err := s.Repo.ListUsers(ctx, "email", input.Email)
	if err != nil {
		log.Println(err)
		return "", err
	}
	if len(responses) == 0 {
		log.Println("User not found")
		return "", errors.New("User not found")
	}

	err = bcrypt.CompareHashAndPassword([]byte(responses[0].Password), []byte(input.Password))
	if err != nil {
		log.Println(err)
		return "", err
	}

	ss, err := auth.CreateSS(responses[0])
	if err != nil {
		log.Println(err)
		return "", err
	}

	return ss, nil
}

func (s *Service) UserInfo(ctx context.Context, id string) (model.User, error) {
	responses, err := s.Repo.ListUsers(ctx, "id", id)
	if err != nil {
		log.Print(err)
		return model.User{}, err
	}
	if len(responses) == 0 {
		log.Print(err)
		return model.User{}, errors.New("User not found")
	}

	user := model.User{
		Id:    responses[0].Id,
		Email: responses[0].Email,
		Name:  responses[0].Name,
		Age:   responses[0].Age,
	}
	return user, nil
}

func (s *Service) UpdateUser(ctx context.Context, user model.User, updates map[string]any) error {
	password, exist := func() (string, bool) {
		for key, value := range updates {
			if key == "password" {
				password := value.(string)
				return password, true
			}
		}
		return "", false
	}()

	if exist {
		hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		updates["password"] = string(hash)
	}

	acceptedinputs := []string{
		"name",
		"age",
		"email",
		"password",
	}
	admininputs := []string{
		"id",
		"created_at",
	}

	if user.Admin {
		acceptedinputs = append(acceptedinputs, admininputs...)
	}

	for key := range updates {
		if !slices.Contains(acceptedinputs, key) {
			delete(updates, key)
		}
	}

	err := s.Repo.UpdateUser_V3(ctx, user.Id, updates)
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) Makeadmin(ctx context.Context, user model.User, id string) error {
	instructions := make(map[string]any)
	if !user.Admin {
		return errors.New("Unauthorized")
	}

	if id == "" {
		return errors.New("Invalid url path")
	}

	instructions["admin"] = true

	if err := s.Repo.UpdateUser_V3(ctx, id, instructions); err != nil {
		return err
	}

	return nil
}
