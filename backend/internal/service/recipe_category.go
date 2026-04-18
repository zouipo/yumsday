package service

import "github.com/zouipo/yumsday/backend/internal/model"

type RecipeCategoryService interface {
	validateRecipeCategory(ingredient model.RecipeCategory) error
}
