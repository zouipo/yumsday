package service

import "github.com/zouipo/yumsday/backend/internal/model"

type IngredientService interface {
	validateIngredient(ingredient model.Ingredient) error
}
