package enums

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

type ItemCategory struct {
	value string
}

var (
	Fruits            = ItemCategory{"FRUITS"}
	Vegetables        = ItemCategory{"VEGETABLES"}
	Meat              = ItemCategory{"MEAT"}
	Seafood           = ItemCategory{"SEAFOOD"}
	Dairy             = ItemCategory{"DAIRY"}
	Starch            = ItemCategory{"STARCH"}
	Beverages         = ItemCategory{"BEVERAGES"}
	Snacks            = ItemCategory{"SNACKS"}
	Condiments        = ItemCategory{"CONDIMENTS"}
	Bakery            = ItemCategory{"BAKERY"}
	BakedGoods        = ItemCategory{"BAKED GOODS"}
	CannedGoods       = ItemCategory{"CANNED GOODS"}
	FrozenFoods       = ItemCategory{"FROZEN FOODS"}
	PersonalCare      = ItemCategory{"PERSONAL CARE"}
	HouseholdSupplies = ItemCategory{"HOUSEHOLD SUPPLIES"}
	PetCare           = ItemCategory{"PET CARE"}
	BabyItems         = ItemCategory{"BABY ITEMS"}
	Others            = ItemCategory{"OTHERS"}
)

func (i ItemCategory) String() string {
	return i.value
}

// UnmarshalJSON implements the json.Unmarshaler interface for ItemCategory.
func (i *ItemCategory) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	switch s {
	case Fruits.value:
		*i = Fruits
	case Vegetables.value:
		*i = Vegetables
	case Meat.value:
		*i = Meat
	case Seafood.value:
		*i = Seafood
	case Dairy.value:
		*i = Dairy
	case Starch.value:
		*i = Starch
	case Beverages.value:
		*i = Beverages
	case Snacks.value:
		*i = Snacks
	case Condiments.value:
		*i = Condiments
	case Bakery.value:
		*i = Bakery
	case BakedGoods.value:
		*i = BakedGoods
	case CannedGoods.value:
		*i = CannedGoods
	case FrozenFoods.value:
		*i = FrozenFoods
	case PersonalCare.value:
		*i = PersonalCare
	case HouseholdSupplies.value:
		*i = HouseholdSupplies
	case PetCare.value:
		*i = PetCare
	case BabyItems.value:
		*i = BabyItems
	case Others.value:
		*i = Others
	default:
		return fmt.Errorf("invalid item category value: %s", s)
	}
	return nil
}

// MarshalJSON implements the json.Marshaler interface for ItemCategory.
func (i ItemCategory) MarshalJSON() ([]byte, error) {
	return json.Marshal(i.value)
}

// Scan implements the sql.Scanner interface for ItemCategory.
func (i *ItemCategory) Scan(value interface{}) error {
	if value == nil {
		return fmt.Errorf("item category cannot be null")
	}

	s, ok := value.(string)
	if !ok {
		return fmt.Errorf("cannot scan %T into ItemCategory", value)
	}

	*i = ItemCategory{value: s}
	return nil
}

// Value implements the driver.Valuer interface for ItemCategory.
func (i ItemCategory) Value() (driver.Value, error) {
	return i.value, nil
}
