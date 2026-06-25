# Testing Guide for Product API

## Overview
This project includes unit tests for controllers and models. Tests are organized in the same directory as the code being tested, with `_test.go` suffix.

## Running Tests

### Run all tests in the project
```bash
go test ./...
```

### Run tests in a specific package
```bash
go test ./controllers
go test ./models
```

### Run tests with verbose output
```bash
go test -v ./...
```

### Run tests with coverage
```bash
go test -cover ./...
```

### Generate coverage report
```bash
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

### Run a specific test
```bash
go test -run TestCreateProductSingle ./controllers
```

## Test Structure

### Controller Tests (`controllers/product_controller_test.go`)
Tests HTTP endpoints and request/response handling:
- `TestCreateProductSingle` - Create single product
- `TestCreateProductBulk` - Bulk create multiple products
- `TestCreateProductInvalidJSON` - Handle invalid input
- `TestGetProducts` - Retrieve all products
- `TestGetProductByIDNotFound` - Handle not found scenario
- `TestUpdateProduct` - Update existing product
- `TestDeleteProduct` - Delete product

### Model Tests (`models/product_test.go`)
Tests data structures and validation:
- `TestProductCreation` - Create product struct
- `TestProductZeroValues` - Handle empty values
- `TestProductPriceValidation` - Validate price
- `TestProductStockValidation` - Validate stock

## Testing Best Practices

### 1. Use Table-Driven Tests
```go
testCases := []struct {
    name  string
    input string
    want  int
}{
    {"case1", "value1", 10},
    {"case2", "value2", 20},
}

for _, tc := range testCases {
    t.Run(tc.name, func(t *testing.T) {
        // test code
    })
}
```

### 2. Use Assertions from testify
```go
import "github.com/stretchr/testify/assert"

assert.Equal(t, expected, actual)
assert.NotEqual(t, expected, actual)
assert.True(t, condition)
assert.Nil(t, err)
assert.Error(t, err)
```

### 3. Test Error Cases
Always test both success and failure scenarios:
- Valid input
- Invalid input
- Edge cases
- Not found scenarios
- Server errors

### 4. Use httptest for HTTP Testing
```go
router := gin.New()
router.POST("/api/products", CreateProduct)

req, _ := http.NewRequest("POST", "/api/products", body)
w := httptest.NewRecorder()
router.ServeHTTP(w, req)

assert.Equal(t, http.StatusCreated, w.Code)
```

## File Naming Convention

- Test files: `*_test.go`
- Benchmark files: `*_bench_test.go`
- Example files: `example_test.go`

## Example: Writing a New Test

```go
func TestMyFunction(t *testing.T) {
    // Arrange
    input := "test data"
    expected := "expected result"

    // Act
    result := MyFunction(input)

    // Assert
    assert.Equal(t, expected, result)
}
```

## Integration Tests

For testing with actual database, create integration tests:

```go
// +build integration

func TestCreateProductWithDB(t *testing.T) {
    // Connect to test database
    db := setupTestDB()
    defer db.Close()
    
    // Test code with real database
}
```

Run integration tests separately:
```bash
go test -tags=integration ./...
```

## CI/CD Integration

Add to your pipeline:
```bash
go test -v -coverprofile=coverage.out ./...
go tool cover -func=coverage.out
```

## Tips

1. **Keep tests focused** - Each test should verify one thing
2. **Use meaningful names** - Test names should describe what they test
3. **Mock external dependencies** - Don't rely on real database in unit tests
4. **Clean up after tests** - Use `defer` to cleanup resources
5. **Run tests frequently** - Run before every commit

## Troubleshooting

### Tests fail with import errors
```bash
go mod tidy
go mod download
```

### Race conditions in tests
```bash
go test -race ./...
```

### Tests hang or timeout
```bash
go test -timeout 30s ./...
```
