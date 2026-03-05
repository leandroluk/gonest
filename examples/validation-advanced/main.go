package main

import (
	"context"
	"fmt"
	"time"

	"github.com/leandroluk/gonest/validator"
	"github.com/leandroluk/gonest/validator/rules"
)

// ========================================
// Example 1: Array Validation
// ========================================

type TagsDto struct {
	Tags []string
}

var tagsValidator = validator.Field[[]string]("tags").
	Is(rules.ArrayNotEmpty[string]()).
	Is(rules.ArrayMinSize[string](1)).
	Is(rules.ArrayMaxSize[string](10)).
	Is(rules.ArrayUnique[string]()).
	Is(rules.ArrayEvery(func(tag string) bool {
		return len(tag) >= 2 && len(tag) <= 20
	}, "Each tag must be between 2-20 characters"))

// ========================================
// Example 2: Date Validation
// ========================================

type EventDto struct {
	StartDate time.Time
	EndDate   time.Time
	Birthdate time.Time
}

var (
	startDateValidator = validator.Field[time.Time]("startDate").
				Is(rules.DateFuture())

	birthdateValidator = validator.Field[time.Time]("birthdate").
				Is(rules.DatePast()).
				Is(rules.DateMinAge(18)).
				Is(rules.DateMaxAge(120))
)

// ========================================
// Example 3: Nested Struct Validation
// ========================================

type Address struct {
	Street  string
	City    string
	ZipCode string
}

// Compile-time interface check
var _ rules.Validatable = (*Address)(nil)

func (a *Address) Validate() *validator.ValidationResult {
	result := validator.NewValidationResult()

	streetVal := validator.Field[string]("street").
		Is(rules.Required[string]()).
		Is(rules.MinLength(5))
	if err := streetVal.Check(a.Street); err != nil {
		result.AddError(err)
	}

	cityVal := validator.Field[string]("city").
		Is(rules.Required[string]()).
		Is(rules.MinLength(2))
	if err := cityVal.Check(a.City); err != nil {
		result.AddError(err)
	}

	zipCodeVal := validator.Field[string]("zipCode").
		Is(rules.Required[string]()).
		Is(rules.Pattern(`^\d{5}(-\d{4})?$`))
	if err := zipCodeVal.Check(a.ZipCode); err != nil {
		result.AddError(err)
	}

	return result
}

type UserWithAddressDto struct {
	Name    string
	Address *Address
}

var addressValidator = validator.Field[*Address]("address").
	Is(rules.ValidStructPtr[Address]())

// ========================================
// Example 4: Async Validation
// ========================================

// Simulated database
var existingEmails = map[string]bool{
	"taken@example.com": true,
	"used@example.com":  true,
}

func checkEmailUnique(ctx context.Context, email string) (bool, error) {
	// Simulate database check
	return !existingEmails[email], nil
}

var asyncEmailValidator = validator.Field[string]("email").
	Is(rules.Required[string]()).
	Is(rules.Email()).
	IsAsync(rules.AsyncUnique(checkEmailUnique, "Email"))

// ========================================
// Example 5: Comparison Validators
// ========================================

type ChangePasswordDto struct {
	OldPassword        string
	NewPassword        string
	ConfirmNewPassword string
}

func (dto *ChangePasswordDto) Validate() *validator.ValidationResult {
	result := validator.NewValidationResult()

	// Old password required
	oldPwdVal := validator.Field[string]("oldPassword").
		Is(rules.Required[string]())
	if err := oldPwdVal.Check(dto.OldPassword); err != nil {
		result.AddError(err)
	}

	// New password must be strong and different from old
	newPwdVal := validator.Field[string]("newPassword").
		Is(rules.Required[string]()).
		Is(rules.StrongPassword()).
		Is(rules.DifferentFrom(dto.OldPassword)).
		WithMessage("New password must be different from old password")
	if err := newPwdVal.Check(dto.NewPassword); err != nil {
		result.AddError(err)
	}

	// Confirm must match new password
	confirmVal := validator.Field[string]("confirmNewPassword").
		Is(rules.Required[string]()).
		Is(rules.SameAs(func() string { return dto.NewPassword }, "newPassword"))
	if err := confirmVal.Check(dto.ConfirmNewPassword); err != nil {
		result.AddError(err)
	}

	return result
}

// ========================================
// Example 6: Conditional Validation
// ========================================

type ProductDto struct {
	Type     string  // "physical" or "digital"
	Weight   float64 // Required for physical
	FileSize int64   // Required for digital
	Price    float64
}

func (dto *ProductDto) Validate() *validator.ValidationResult {
	result := validator.NewValidationResult()

	// Type validation
	typeVal := validator.Field[string]("type").
		Is(rules.Required[string]()).
		Is(rules.OneOf([]string{"physical", "digital"}))
	if err := typeVal.Check(dto.Type); err != nil {
		result.AddError(err)
	}

	// Conditional: Weight required only for physical products
	weightVal := validator.Field[float64]("weight").
		Is(rules.When(
			func(w float64) bool { return dto.Type == "physical" },
			rules.Positive[float64](),
		))
	if dto.Type == "physical" {
		if err := weightVal.Check(dto.Weight); err != nil {
			result.AddError(err)
		}
	}

	// Conditional: FileSize required only for digital products
	fileSizeVal := validator.Field[int64]("fileSize").
		Is(rules.When(
			func(fs int64) bool { return dto.Type == "digital" },
			rules.Positive[int64](),
		))
	if dto.Type == "digital" {
		if err := fileSizeVal.Check(dto.FileSize); err != nil {
			result.AddError(err)
		}
	}

	// Price always required
	priceVal := validator.Field[float64]("price").
		Is(rules.Required[float64]()).
		Is(rules.Positive[float64]())
	if err := priceVal.Check(dto.Price); err != nil {
		result.AddError(err)
	}

	return result
}

// ========================================
// Example 7: Array with Element Validation
// ========================================

type BulkCreateDto struct {
	Items []int
}

var bulkItemsValidator = validator.Field[[]int]("items").
	Is(rules.ArrayNotEmpty[int]()).
	Is(rules.ArrayMinSize[int](1)).
	Is(rules.ArrayMaxSize[int](100)).
	Is(rules.ArrayEach(rules.Range(1, 1000)))

// ========================================
// Main Demo
// ========================================

func main() {
	ctx := context.Background()

	fmt.Println("========================================")
	fmt.Println("GoNest Advanced Validation Examples")
	fmt.Println("========================================")

	// ========================================
	// Test 1: Array Validation
	// ========================================
	fmt.Println("")
	fmt.Println("1. Array Validation")
	fmt.Println("-------------------")

	validTags := TagsDto{Tags: []string{"golang", "backend", "api"}}
	if err := tagsValidator.Check(validTags.Tags); err == nil {
		fmt.Println("✓ Valid tags:", validTags.Tags)
	}

	duplicateTags := TagsDto{Tags: []string{"golang", "golang", "api"}}
	if err := tagsValidator.Check(duplicateTags.Tags); err != nil {
		fmt.Printf("✗ Duplicate tags error: %s\n", err.Message())
	}

	longTags := TagsDto{Tags: []string{"a"}}
	if err := tagsValidator.Check(longTags.Tags); err != nil {
		fmt.Printf("✗ Tag too short: %s\n", err.Message())
	}
	fmt.Println()

	// ========================================
	// Test 2: Date Validation
	// ========================================
	fmt.Println("2. Date Validation")
	fmt.Println("------------------")

	futureDate := time.Now().AddDate(0, 1, 0)
	if err := startDateValidator.Check(futureDate); err == nil {
		fmt.Println("✓ Future date valid")
	}

	pastDate := time.Now().AddDate(-1, 0, 0)
	if err := startDateValidator.Check(pastDate); err != nil {
		fmt.Printf("✗ Past date for future field: %s\n", err.Message())
	}

	birthdate := time.Now().AddDate(-25, 0, 0)
	if err := birthdateValidator.Check(birthdate); err == nil {
		fmt.Println("✓ Birthdate valid (25 years old)")
	}

	tooYoung := time.Now().AddDate(-15, 0, 0)
	if err := birthdateValidator.Check(tooYoung); err != nil {
		fmt.Printf("✗ Too young: %s\n", err.Message())
	}
	fmt.Println()

	// ========================================
	// Test 3: Nested Struct Validation
	// ========================================
	fmt.Println("3. Nested Struct Validation")
	fmt.Println("---------------------------")

	validUser := UserWithAddressDto{
		Address: &Address{
			Street:  "123 Main St",
			City:    "Springfield",
			ZipCode: "12345",
		},
	}

	if err := addressValidator.Check(validUser.Address); err == nil {
		fmt.Println("✓ Valid nested address")
	}

	invalidUser := UserWithAddressDto{
		Address: &Address{
			Street:  "St",
			City:    "C",
			ZipCode: "123",
		},
	}

	if err := addressValidator.Check(invalidUser.Address); err != nil {
		fmt.Printf("✗ Invalid address: %s\n", err.Message())
	}
	fmt.Println()

	// ========================================
	// Test 4: Async Validation
	// ========================================
	fmt.Println("4. Async Validation")
	fmt.Println("-------------------")

	newEmail := "new@example.com"
	if err := asyncEmailValidator.CheckAsync(ctx, newEmail); err == nil {
		fmt.Printf("✓ Email '%s' is available\n", newEmail)
	}

	takenEmail := "taken@example.com"
	if err := asyncEmailValidator.CheckAsync(ctx, takenEmail); err != nil {
		fmt.Printf("✗ Email '%s': %s\n", takenEmail, err.Message())
	}
	fmt.Println()

	// ========================================
	// Test 5: Comparison & Cross-Field
	// ========================================
	fmt.Println("5. Comparison & Cross-Field")
	fmt.Println("---------------------------")

	validChange := &ChangePasswordDto{
		OldPassword:        "OldPass123!",
		NewPassword:        "NewPass456!",
		ConfirmNewPassword: "NewPass456!",
	}

	result := validChange.Validate()
	if result.Valid() {
		fmt.Println("✓ Password change valid")
	}

	invalidChange := &ChangePasswordDto{
		OldPassword:        "OldPass123!",
		NewPassword:        "OldPass123!", // Same as old
		ConfirmNewPassword: "Different!",
	}

	result = invalidChange.Validate()
	if result.Invalid() {
		fmt.Printf("✗ Password change errors (%d):\n", result.Count())
		for _, err := range result.Errors() {
			fmt.Printf("  - %s: %s\n", err.Field(), err.Message())
		}
	}
	fmt.Println()

	// ========================================
	// Test 6: Conditional Validation
	// ========================================
	fmt.Println("6. Conditional Validation")
	fmt.Println("-------------------------")

	physical := &ProductDto{
		Type:   "physical",
		Weight: 1.5,
		Price:  29.99,
	}

	result = physical.Validate()
	if result.Valid() {
		fmt.Println("✓ Physical product valid")
	}

	digital := &ProductDto{
		Type:     "digital",
		FileSize: 1024000,
		Price:    9.99,
	}

	result = digital.Validate()
	if result.Valid() {
		fmt.Println("✓ Digital product valid")
	}

	invalidPhysical := &ProductDto{
		Type:  "physical",
		Price: 19.99,
		// Missing weight!
	}

	result = invalidPhysical.Validate()
	if result.Invalid() {
		fmt.Printf("✗ Invalid physical product: %d errors\n", result.Count())
	}
	fmt.Println()

	// ========================================
	// Test 7: Array Element Validation
	// ========================================
	fmt.Println("7. Array Element Validation")
	fmt.Println("---------------------------")

	validBulk := &BulkCreateDto{
		Items: []int{10, 50, 100, 500},
	}

	if err := bulkItemsValidator.Check(validBulk.Items); err == nil {
		fmt.Println("✓ Bulk items valid")
	}

	invalidBulk := &BulkCreateDto{
		Items: []int{10, 5000, 100}, // 5000 > max(1000)
	}

	if err := bulkItemsValidator.Check(invalidBulk.Items); err != nil {
		fmt.Printf("✗ Bulk items invalid: %s\n", err.Message())
	}
	fmt.Println()

	// ========================================
	// Summary
	// ========================================
	fmt.Println("========================================")
	fmt.Println("Summary:")
	fmt.Println("✓ Array validation (size, unique, each)")
	fmt.Println("✓ Date validation (past, future, age)")
	fmt.Println("✓ Nested struct validation")
	fmt.Println("✓ Async validation (DB checks)")
	fmt.Println("✓ Comparison (same as, different from)")
	fmt.Println("✓ Conditional validation (when/unless)")
	fmt.Println("✓ Array element validation")
	fmt.Println("========================================")
}
