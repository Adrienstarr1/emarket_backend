package service

import (
	"context"
	"errors"
)

func (s *Service) Add2cart(ctx context.Context, userId, productId string, instructions map[string]any) error {
	instructions["user_id"] = userId
	instructions["product_id"] = productId

	res, err := s.Repo.FindCart(ctx, userId)
	if err != nil {
		return err
	}

	exist := func() bool {
		for _, r := range res {
			if productId == r.Product_id {
				return true
			}
		}
		return false
	}()

	switch exist {
	case true:
		err = s.Repo.UpdateCart(ctx, userId, productId, 1)
		if err != nil {
			return err
		}
	case false:
		err = s.Repo.Add2cart(ctx, instructions)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *Service) DeleteCart(ctx context.Context, userId, productId string, instructions map[string]any) error {
	res, err := s.Repo.FindCart(ctx, userId)
	if err != nil {
		return err
	}

	exist, available := func() (bool, bool) {
		for _, r := range res {
			if productId == r.Product_id && r.Quantity > 0 {
				return true, true
			} else if productId == r.Product_id && r.Quantity <= 0 {
				return true, false
			}
		}
		return false, false
	}()

	switch exist {
	case true:
		switch available {
		case true:
			err = s.Repo.UpdateCart(ctx, userId, productId, 0)
			if err != nil {
				return err
			}
		case false:
			err = s.Repo.DeleteCart(ctx, userId, productId)
			if err != nil {
				return err
			}
		}
	case false:
		return errors.New("Cart doesn't exist")
	}
	return nil
}
