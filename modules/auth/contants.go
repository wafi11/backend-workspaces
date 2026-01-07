package auth

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"unicode"
)

type ErrorMessage string

const (
	// Registration errors
	ErrUsernameExists     ErrorMessage = "username already exists"
	ErrUsernameNotFound   ErrorMessage = "invalid credentials"
	ErrEmailExists        ErrorMessage = "email already exists"
	ErrPhoneNumberExists  ErrorMessage = "phone number already exists"
	ErrWeakPassword       ErrorMessage = "password does not meet security requirements"
	ErrPasswordNotMatch   ErrorMessage = "password not match"
	ErrPasswordHashFailed ErrorMessage = "failed to hash password"
	ErrInvalidUsername    ErrorMessage = "username format is invalid"
	ErrInvalidEmail       ErrorMessage = "email format is invalid"
	ErrInvalidPhoneNumber ErrorMessage = "phone number format is invalid"

	// Login errors
	ErrInvalidCredentials ErrorMessage = "invalid username or password"
	ErrUserNotFound       ErrorMessage = "user not found"
	ErrAccountLocked      ErrorMessage = "account is locked"
	ErrAccountInactive    ErrorMessage = "account is not active"

	// Database errors
	ErrDatabaseConnection ErrorMessage = "database connection failed"
	ErrTransactionFailed  ErrorMessage = "transaction failed"
	ErrQueryFailed        ErrorMessage = "query execution failed"

	// General errors
	ErrInternalServer ErrorMessage = "internal server error"
	ErrUnauthorized   ErrorMessage = "unauthorized access"
	ErrForbidden      ErrorMessage = "forbidden access"
)

func validateUsername(username string) error {
	username = strings.TrimSpace(username)

	// Check if empty
	if username == "" {
		return fmt.Errorf("%s: username cannot be empty", ErrInvalidUsername)
	}

	// Check length (3-30 characters)
	if len(username) < 3 || len(username) > 30 {
		return fmt.Errorf("%s: username must be between 3 and 30 characters", ErrInvalidUsername)
	}

	validUsername := regexp.MustCompile(`^[a-zA-Z0-9_]+$`)
	if !validUsername.MatchString(username) {
		return fmt.Errorf("%s: username can only contain letters, numbers, and underscores", ErrInvalidUsername)
	}

	// Username must start with a letter
	if !unicode.IsLetter(rune(username[0])) {
		return fmt.Errorf("%s: username must start with a letter", ErrInvalidUsername)
	}

	return nil
}

func validateEmail(email string) error {
	email = strings.TrimSpace(email)

	// Check if empty
	if email == "" {
		return fmt.Errorf("%s: email cannot be empty", ErrInvalidEmail)
	}

	// Email regex pattern
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(email) {
		return fmt.Errorf("%s: invalid email format", ErrInvalidEmail)
	}

	// Check length
	if len(email) > 255 {
		return fmt.Errorf("%s: email too long", ErrInvalidEmail)
	}

	return nil
}

func validatePhoneNumber(phoneNumber string) (string, error) {
	phoneNumber = strings.TrimSpace(phoneNumber)

	// Check if empty
	if phoneNumber == "" {
		return "", fmt.Errorf("%s: phone number cannot be empty", ErrInvalidPhoneNumber)
	}

	// Remove common separators for validation
	cleaned := strings.ReplaceAll(phoneNumber, " ", "")
	cleaned = strings.ReplaceAll(cleaned, "-", "")
	cleaned = strings.ReplaceAll(cleaned, "(", "")
	cleaned = strings.ReplaceAll(cleaned, ")", "")

	// Normalize to international format +62
	normalized := normalizePhoneNumber(cleaned)

	// Phone number regex (international format with +62)
	phoneRegex := regexp.MustCompile(`^\+62[0-9]{9,13}$`)
	if !phoneRegex.MatchString(normalized) {
		return "", fmt.Errorf("%s: phone number must be valid Indonesian number", ErrInvalidPhoneNumber)
	}

	return normalized, nil
}

func normalizePhoneNumber(phone string) string {
	phone = strings.TrimPrefix(phone, "+")

	if strings.HasPrefix(phone, "0") {
		phone = "62" + phone[1:]
	}

	if strings.HasPrefix(phone, "62") {
		return "+" + phone
	}

	if strings.HasPrefix(phone, "8") {
		return "+62" + phone
	}

	return "+62" + phone
}

func ValidatePassword(password string) error {
	if len(password) < 8 {
		return fmt.Errorf("%s: password must be at least 8 characters", ErrWeakPassword)
	}

	if len(password) > 72 {
		return fmt.Errorf("%s: password too long (max 72 characters)", ErrWeakPassword)
	}

	var (
		hasUpper   bool
		hasLower   bool
		hasNumber  bool
		hasSpecial bool
	)

	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsNumber(char):
			hasNumber = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}

	if !hasUpper || !hasLower || !hasNumber || !hasSpecial {
		return fmt.Errorf("%s: password must contain uppercase, lowercase, number, and special character", ErrWeakPassword)
	}

	return nil
}

func determineStatusCode(err error) int {
	errMsg := err.Error()

	// Validation errors (400)
	validationErrors := []ErrorMessage{
		ErrInvalidUsername,
		ErrInvalidEmail,
		ErrInvalidPhoneNumber,
		ErrWeakPassword,
	}
	for _, validErr := range validationErrors {
		if strings.Contains(errMsg, string(validErr)) {
			return http.StatusBadRequest
		}
	}

	// Conflict errors (409) - duplicate data
	conflictErrors := []ErrorMessage{
		ErrUsernameExists,
		ErrEmailExists,
		ErrPhoneNumberExists,
	}
	for _, conflictErr := range conflictErrors {
		if strings.Contains(errMsg, string(conflictErr)) {
			return http.StatusConflict
		}
	}

	// Server errors (500)
	serverErrors := []ErrorMessage{
		ErrPasswordHashFailed,
		ErrTransactionFailed,
		ErrQueryFailed,
		ErrDatabaseConnection,
		ErrInternalServer,
	}
	for _, serverErr := range serverErrors {
		if strings.Contains(errMsg, string(serverErr)) {
			return http.StatusInternalServerError
		}
	}

	// Default to 500 for unknown errors
	return http.StatusInternalServerError
}
