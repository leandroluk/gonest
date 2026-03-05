package main

import (
	"fmt"

	"github.com/leandroluk/gonest/validator"
	"github.com/leandroluk/gonest/validator/rules"
)

type UpdateUserDto struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Age   int    `json:"age"`
}

var updateUserSchema *validator.Schema[UpdateUserDto]

func init() {
	var dto UpdateUserDto
	builder := validator.NewSchema(&dto)

	// Field names come from json tags automatically!
	builder.Field(&dto.Name, rules.Required[string](), rules.MinLength(2))
	builder.Field(&dto.Email, rules.Required[string](), rules.Email())
	builder.Field(&dto.Age, rules.Min(18), rules.Max(120))

	updateUserSchema = builder.Build()
}

func main() {
	fmt.Println("Schema with Field Pointers & JSON Tags")
	fmt.Println("")

	// Valid
	valid := &UpdateUserDto{
		Name:  "John Doe",
		Email: "john@example.com",
		Age:   30,
	}

	result := updateUserSchema.Validate(valid)
	if result.Valid() {
		fmt.Println("✓ Valid user passed")
	}

	// Invalid
	invalid := &UpdateUserDto{
		Name:  "J",
		Email: "not-email",
		Age:   15,
	}

	result = updateUserSchema.Validate(invalid)
	if result.Invalid() {
		fmt.Printf("\n✗ Invalid user (%d errors):\n", result.Count())
		for _, err := range result.Errors() {
			fmt.Printf("  - %s: %s\n", err.Field(), err.Message())
		}
	}
}
