package service

import (
	"context"
	"e-market/model"
	"slices"
	"time"

	"github.com/google/uuid"
)

func (s *Service) CreateProduct(ctx context.Context, product model.Product, user model.User) error {
	product.Id = uuid.New().String()
	product.User_Id = user.Id
	err := s.Repo.AddProduct(ctx, product.Id, product.Product_name, product.Quantity, product.User_Id, product.Price, product.Description)
	if err != nil {
		return err
	}

	return nil
}
func (s *Service) UpdateProduct(ctx context.Context, id string, updates map[string]any) error {
	acceptedinputs := []string{
		"name",
		"quantity",
		"price",
		"description",
	}

	for key := range updates {
		if !slices.Contains(acceptedinputs, key) {
			delete(updates, key)
		}
	}

	updates["updated_at"] = time.Now()

	err := s.Repo.UpdateUser_V3(ctx, id, updates)
	if err != nil {
		return err
	}

	return nil
}
func (s *Service) GetProduct(ctx context.Context, id string) (model.Product, error) {
	response, err := s.Repo.FindProduct_V2(ctx, "id", id)
	if err != nil || len(response) == 0 {
		return model.Product{}, err
	}

	product := model.Product{
		Id:           response[0].Id,
		User_Id:      response[0].User_Id,
		Product_name: response[0].Product_name,
		Quantity:     response[0].Quantity,
		Price:        response[0].Price,
		Description:  response[0].Description,
	}

	return product, nil
}
