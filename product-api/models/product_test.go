package models

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

// Test Product struct creation
func TestProductCreation(t *testing.T) {
	product := Product{
		ID:    1,
		Name:  "Test Product",
		Price: 100,
		Stock: 10,
	}

	assert.Equal(t, uint(1), product.ID)
	assert.Equal(t, "Test Product", product.Name)
	assert.Equal(t, 100, product.Price)
	assert.Equal(t, 10, product.Stock)
}

// Test Product with zero values
func TestProductZeroValues(t *testing.T) {
	product := Product{}

	assert.Equal(t, uint(0), product.ID)
	assert.Equal(t, "", product.Name)
	assert.Equal(t, 0, product.Price)
	assert.Equal(t, 0, product.Stock)
}

// Test Product price validation
func TestProductPriceValidation(t *testing.T) {
	testCases := []struct {
		name  string
		price int
		valid bool
	}{
		{"Positive price", 100, true},
		{"Zero price", 0, true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			p := Product{Price: tc.price}
			// Validate price is >= 0
			assert.GreaterOrEqual(t, p.Price, 0, tc.name)
		})
	}
}

// Test Product stock validation
func TestProductStockValidation(t *testing.T) {
	testCases := []struct {
		name  string
		stock int
	}{
		{"In stock", 50},
		{"Low stock", 1},
		{"Out of stock", 0},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			p := Product{Stock: tc.stock}
			assert.GreaterOrEqual(t, p.Stock, 0, tc.name)
		})
	}
}
