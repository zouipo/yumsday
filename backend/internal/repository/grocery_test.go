package repository

import (
	"testing"

	"github.com/zouipo/yumsday/backend/internal/pkg/utils"
)

func TestNewGroceryRepository(t *testing.T) {
	db := utils.SetUpTestDB(t)
	defer db.Close()

	repo := NewGroceryRepository(db)
	if repo == nil {
		t.Fatal("expected non-nil repository, got nil")
	}
	if repo.db == nil {
		t.Fatal("expected non-nil database connection, got nil")
	}
}

func TestGroceryRepositoryHasItem(t *testing.T) {
	db := utils.SetUpTestDB(t)
	defer db.Close()

	repo := NewGroceryRepository(db)

	tests := []struct {
		name     string
		itemID   int64
		expected bool
	}{
		{
			name:     "item present in groceries",
			itemID:   1,
			expected: true,
		},
		{
			name:     "existing item not present in groceries",
			itemID:   3,
			expected: false,
		},
		{
			name:     "unknown item",
			itemID:   -1,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual, err := repo.HasItem(tt.itemID)
			if err != nil {
				t.Fatalf("didn't expected error, got %v", err)
			}

			if actual != tt.expected {
				t.Fatalf("expected %v, got %v", tt.expected, actual)
			}
		})
	}
}
