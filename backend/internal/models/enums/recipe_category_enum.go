package enums

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

type RecipeCategoryEnum struct {
	value string
}

var (
	Dessert    = RecipeCategoryEnum{"DESSERT"}
	Salad      = RecipeCategoryEnum{"SALAD"}
	MainCourse = RecipeCategoryEnum{"MAIN COURSE"}
	Soup       = RecipeCategoryEnum{"SOUP"}
	Breakfast  = RecipeCategoryEnum{"BREAKFAST"}
	Brunch     = RecipeCategoryEnum{"BRUNCH"}
	Starter    = RecipeCategoryEnum{"STARTER"}
	Sauce      = RecipeCategoryEnum{"SAUCE"}
	Snack      = RecipeCategoryEnum{"SNACK"}
	Beverage   = RecipeCategoryEnum{"BEVERAGE"}
	Vegan      = RecipeCategoryEnum{"VEGAN"}
	Vegetarian = RecipeCategoryEnum{"VEGETARIAN"}
	GlutenFree = RecipeCategoryEnum{"GLUTEN FREE"}
)

func (r RecipeCategoryEnum) String() string {
	return r.value
}

// UnmarshalJSON implements the json.Unmarshaler interface for RecipeCategoryEnum.
func (r *RecipeCategoryEnum) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	switch s {
	case Dessert.value:
		*r = Dessert
	case Salad.value:
		*r = Salad
	case MainCourse.value:
		*r = MainCourse
	case Soup.value:
		*r = Soup
	case Breakfast.value:
		*r = Breakfast
	case Brunch.value:
		*r = Brunch
	case Starter.value:
		*r = Starter
	case Sauce.value:
		*r = Sauce
	case Snack.value:
		*r = Snack
	case Beverage.value:
		*r = Beverage
	case Vegan.value:
		*r = Vegan
	case Vegetarian.value:
		*r = Vegetarian
	case GlutenFree.value:
		*r = GlutenFree
	default:
		return fmt.Errorf("invalid recipe category value: %s", s)
	}
	return nil
}

// MarshalJSON implements the json.Marshaler interface for RecipeCategoryEnum.
func (r RecipeCategoryEnum) MarshalJSON() ([]byte, error) {
	return json.Marshal(r.value)
}

// Scan implements the sql.Scanner interface for RecipeCategoryEnum.
func (r *RecipeCategoryEnum) Scan(value interface{}) error {
	if value == nil {
		return fmt.Errorf("recipe category cannot be null")
	}

	s, ok := value.(string)
	if !ok {
		return fmt.Errorf("cannot scan %T into RecipeCategoryEnum", value)
	}

	*r = RecipeCategoryEnum{value: s}
	return nil
}

// Value implements the driver.Valuer interface for RecipeCategoryEnum.
func (r RecipeCategoryEnum) Value() (driver.Value, error) {
	return r.value, nil
}
