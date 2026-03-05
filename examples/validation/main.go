package main

import (
	"context"
	"fmt"

	"github.com/leandroluk/gonest/validator"
	"github.com/leandroluk/gonest/validator/rules"
)

// ========================================
// DTOs with Type-Safe Validation
// ========================================

type CreateUserDto struct {
	Email           string
	Password        string
	PasswordConfirm string
	Age             int
	Username        string
	Website         string
}

// Define validators as package-level variables (compiled once)
var (
	emailValidator = validator.Field[string]("email").
			Is(rules.Required[string]()).
			Is(rules.Email()).
			Is(rules.MaxLength(255))

	passwordValidator = validator.Field[string]("password").
				Is(rules.Required[string]()).
				Is(rules.MinLength(8)).
				Is(rules.MaxLength(100)).
				Is(rules.HasUpperCase()).
				Is(rules.HasLowerCase()).
				Is(rules.HasDigit())

	ageValidator = validator.Field[int]("age").
			Is(rules.Required[int]()).
			Is(rules.Min(18)).
			Is(rules.Max(120))

	usernameValidator = validator.Field[string]("username").
				Is(rules.Required[string]()).
				Is(rules.MinLength(3)).
				Is(rules.MaxLength(20)).
				Is(rules.AlphaNumeric())

	websiteValidator = validator.Field[string]("website").
				Is(rules.URL())
)

// Validate method using validators
func (dto *CreateUserDto) Validate() *validator.ValidationResult {
	result := validator.NewValidationResult()

	// Validate each field
	if err := emailValidator.Check(dto.Email); err != nil {
		result.AddError(err)
	}

	if err := passwordValidator.Check(dto.Password); err != nil {
		result.AddError(err)
	}

	if err := ageValidator.Check(dto.Age); err != nil {
		result.AddError(err)
	}

	if err := usernameValidator.Check(dto.Username); err != nil {
		result.AddError(err)
	}

	if dto.Website != "" { // Optional field
		if err := websiteValidator.Check(dto.Website); err != nil {
			result.AddError(err)
		}
	}

	// Cross-field validation
	if dto.Password != dto.PasswordConfirm {
		result.AddError(validator.NewFieldError(
			"passwordConfirm",
			"match",
			"Passwords must match",
		))
	}

	return result
}

// ========================================
// Using Schema Builder
// ========================================

type UpdateUserDto struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Age   int    `json:"age"`
}

var updateUserSchema *validator.Schema[UpdateUserDto]

func init() {
	var dto UpdateUserDto
	builder := validator.NewSchema(&dto)

	builder.Field(&dto.Name,
		rules.MinLength(2),
	)

	builder.Field(&dto.Email,
		rules.Required[string](),
		rules.Email(),
	)

	builder.Field(&dto.Age,
		rules.Range(18, 120),
	)

	updateUserSchema = builder.Build()
}

// ========================================
// Custom Validators
// ========================================

// Custom validator: check if email domain is allowed
func allowedEmailDomain(allowedDomains []string) validator.Validator[string] {
	return func(value string) *validator.FieldError {
		for _, domain := range allowedDomains {
			if len(value) > len(domain) && value[len(value)-len(domain):] == domain {
				return nil
			}
		}

		err := validator.NewFieldError(
			"",
			"email_domain",
			"Email domain not allowed",
		)
		err.WithParam("allowed_domains", allowedDomains)
		return err
	}
}

// Async validator: check if email already exists (simulated)
func uniqueEmail() validator.ContextValidator[string] {
	return func(ctx context.Context, value string) *validator.FieldError {
		// Simulate database check
		existingEmails := []string{"taken@example.com", "used@example.com"}

		for _, email := range existingEmails {
			if value == email {
				return validator.NewFieldError(
					"",
					"unique",
					"Email is already taken",
				)
			}
		}

		return nil
	}
}

// ========================================
// Main Demo
// ========================================

func main() {
	ctx := context.Background()

	fmt.Println("========================================")
	fmt.Println("GoNest Type-Safe Validation Examples")
	fmt.Println("10 comprehensive examples")
	fmt.Println("========================================")

	// ========================================
	// Example 1: Valid DTO
	// ========================================
	fmt.Println("")
	fmt.Println("1. Testing VALID DTO...")
	validDto := &CreateUserDto{
		Email:           "john@example.com",
		Password:        "SecurePass123",
		PasswordConfirm: "SecurePass123",
		Age:             25,
		Username:        "johndoe",
		Website:         "https://example.com",
	}

	result := validDto.Validate()
	if result.Valid() {
		fmt.Println("   ✓ Validation PASSED")
	} else {
		fmt.Println("   ✗ Validation FAILED:")
		for _, err := range result.Errors() {
			fmt.Printf("     - %s\n", err.Error())
		}
	}
	fmt.Println()

	// ========================================
	// Example 2: Invalid DTO
	// ========================================
	fmt.Println("2. Testing INVALID DTO...")
	invalidDto := &CreateUserDto{
		Email:           "invalid-email",
		Password:        "weak",
		PasswordConfirm: "different",
		Age:             15,
		Username:        "ab",
		Website:         "not-a-url",
	}

	result = invalidDto.Validate()
	if result.Valid() {
		fmt.Println("   ✓ Validation PASSED")
	} else {
		fmt.Printf("   ✗ Validation FAILED (%d errors):\n", result.Count())
		for _, err := range result.Errors() {
			fmt.Printf("     - %s: %s (code: %s)\n", err.Field(), err.Message(), err.Code())
		}
	}
	fmt.Println()

	// ========================================
	// Example 3: Field-by-Field Validation
	// ========================================
	fmt.Println("3. Testing Individual Fields...")

	testCases := []struct {
		name     string
		value    string
		expected bool
	}{
		{"Valid email", "test@example.com", true},
		{"Invalid email", "not-an-email", false},
		{"Empty email", "", false},
	}

	for _, tc := range testCases {
		err := emailValidator.Check(tc.value)
		passed := err == nil
		status := "✓"
		if !passed {
			status = "✗"
		}
		fmt.Printf("   %s %s: %v", status, tc.name, passed)
		if err != nil {
			fmt.Printf(" (%s)", err.Message())
		}
		fmt.Println()
	}
	fmt.Println()

	// ========================================
	// Example 4: Schema Validation
	// ========================================
	fmt.Println("4. Testing Schema Validation...")
	updateDto := &UpdateUserDto{
		Name:  "John",
		Email: "john@example.com",
		Age:   30,
	}

	schemaResult := updateUserSchema.Validate(updateDto)
	if schemaResult.Valid() {
		fmt.Println("   ✓ Schema validation PASSED")
	} else {
		fmt.Println("   ✗ Schema validation FAILED:")
		for _, err := range schemaResult.Errors() {
			fmt.Printf("     - %s\n", err.Error())
		}
	}
	fmt.Println()

	// ========================================
	// Example 5: Custom Validators
	// ========================================
	fmt.Println("5. Testing Custom Validators...")

	domainValidator := validator.Field[string]("email").
		Is(rules.Required[string]()).
		Is(allowedEmailDomain([]string{"@company.com", "@partner.com"}))

	emails := []string{
		"user@company.com",
		"admin@partner.com",
		"external@gmail.com",
	}

	for _, email := range emails {
		err := domainValidator.Check(email)
		if err == nil {
			fmt.Printf("   ✓ %s: Allowed\n", email)
		} else {
			fmt.Printf("   ✗ %s: %s\n", email, err.Message())
		}
	}
	fmt.Println()

	// ========================================
	// Example 6: Async Validation
	// ========================================
	fmt.Println("6. Testing Async Validation...")

	asyncValidator := validator.Field[string]("email").
		Is(rules.Required[string]()).
		Is(rules.Email()).
		IsAsync(uniqueEmail())

	testEmails := []string{
		"new@example.com",
		"taken@example.com",
	}

	for _, email := range testEmails {
		err := asyncValidator.CheckAsync(ctx, email)
		if err == nil {
			fmt.Printf("   ✓ %s: Available\n", email)
		} else {
			fmt.Printf("   ✗ %s: %s\n", email, err.Message())
		}
	}
	fmt.Println()

	// ========================================
	// Example 7: JSON Error Output
	// ========================================
	fmt.Println("7. JSON Error Format...")
	jsonErrors := result.ToJSON()
	fmt.Printf("   %+v\n", jsonErrors)
	fmt.Println()

	// ========================================
	// Example 8: Numeric Validators
	// ========================================
	fmt.Println("8. Testing Numeric Validators...")

	priceValidator := validator.Field[float64]("price").
		Is(rules.Required[float64]()).
		Is(rules.Positive[float64]()).
		Is(rules.Max(9999.99))

	prices := []float64{10.50, -5.00, 10000.00}
	for _, price := range prices {
		err := priceValidator.Check(price)
		if err == nil {
			fmt.Printf("   ✓ $%.2f: Valid\n", price)
		} else {
			fmt.Printf("   ✗ $%.2f: %s\n", price, err.Message())
		}
	}
	fmt.Println()

	// ========================================
	// Example 9: Strong Password Validator
	// ========================================
	fmt.Println("9. Testing Strong Password Validator...")

	strongPasswordValidator := validator.Field[string]("password").
		Is(rules.StrongPassword())

	passwords := []string{
		"Weak123",    // No special char
		"weak123!",   // No uppercase
		"WEAK123!",   // No lowercase
		"WeakPass!",  // No digit
		"Strong123!", // Valid
	}

	for _, pwd := range passwords {
		err := strongPasswordValidator.Check(pwd)
		if err == nil {
			fmt.Printf("   ✓ '%s': Strong password\n", pwd)
		} else {
			fmt.Printf("   ✗ '%s': %s\n", pwd, err.Message())
		}
	}
	fmt.Println()

	// ========================================
	// Example 10: Boolean Validators
	// ========================================
	fmt.Println("10. Testing Boolean Validators...")

	acceptTermsValidator := validator.Field[bool]("acceptTerms").
		Is(rules.MustAccept())

	isAdultValidator := validator.Field[bool]("isAdult").
		Is(rules.IsTrue())

	// Test accept terms
	if err := acceptTermsValidator.Check(true); err == nil {
		fmt.Println("   ✓ Terms accepted: valid")
	}

	if err := acceptTermsValidator.Check(false); err != nil {
		fmt.Printf("   ✗ Terms not accepted: %s\n", err.Message())
	}

	// Test is adult
	if err := isAdultValidator.Check(true); err == nil {
		fmt.Println("   ✓ Is adult: valid")
	}

	if err := isAdultValidator.Check(false); err != nil {
		fmt.Printf("   ✗ Not adult: %s\n", err.Message())
	}
	fmt.Println()

	// ========================================
	// Summary
	// ========================================
	fmt.Println("========================================")
	fmt.Println("Summary:")
	fmt.Println("✓ Type-safe validation")
	fmt.Println("✓ Composable validators")
	fmt.Println("✓ Custom validators")
	fmt.Println("✓ Async validation")
	fmt.Println("✓ Schema-based validation")
	fmt.Println("✓ Boolean validation")
	fmt.Println("✓ Detailed error messages")
	fmt.Println("========================================")
}
