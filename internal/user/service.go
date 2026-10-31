package user

import (
	"context"
	"errors"
	"net/mail"
	"regexp"
	"strings"
	"unicode"

	"github.com/google/uuid"
)

var (
	ErrInvalidUsername = errors.New("invalid username")
	ErrInvalidEmail    = errors.New("invalid email")
)

type Service interface {
	Create(ctx context.Context, input CreateUserInput) (*User, error)
	Read(ctx context.Context, id uuid.UUID) (*User, error)
	Update(ctx context.Context, id uuid.UUID, input UpdateUserInput) (*User, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type UserService struct {
	repo Repository
}

func NewService(repo Repository) UserService {
	return UserService{repo: repo}
}

func (s *UserService) Create(ctx context.Context, input CreateUserInput) (*User, error) {
	if !IsValidUsername(input.Username) {
		return nil, ErrInvalidUsername
	}
	if !IsValidEmail(input.Email) {
		return nil, ErrInvalidEmail
	}
	err := ValidatePassword(input.Password)
	if err != nil {
		return nil, err
	}
	if input.Preferences != nil {
		prefsInput, err := ValidatePreferences(*input.Preferences)
		if err != nil {
			return nil, err
		}
		*input.Preferences = prefsInput
	}
	return s.repo.Create(ctx, input)
}

func (s *UserService) Read(ctx context.Context, id uuid.UUID) (*User, error) {
	userFromCtx, ok := FromContext(ctx)
	if !ok {
		return nil, ErrUserNotFoundInCtx
	}
	userRead, err := s.repo.Read(ctx, id)
	if err != nil {
		return nil, err
	}
	if !userFromCtx.IsAdmin && userRead.ID != userFromCtx.ID {
		return nil, ErrForbidden
	}
	return userRead, nil
}

func (s *UserService) Update(
	ctx context.Context,
	id uuid.UUID,
	input UpdateUserInput,
) (*User, error) {
	// implicitly check if user has sufficient permissions to update
	_, err := s.Read(ctx, id)
	if err != nil {
		return nil, err
	}

	if input.Email != nil && !IsValidEmail(*input.Email) {
		return nil, ErrInvalidEmail
	}
	if input.Password != nil {
		err := ValidatePassword(*input.Password)
		if err != nil {
			return nil, err
		}
	}
	if input.Preferences != nil {
		prefsInput, err := ValidatePreferences(*input.Preferences)
		if err != nil {
			return nil, err
		}
		*input.Preferences = prefsInput
	}

	a := ValidateAdmin(ctx)
	if a == nil && input.IsAdmin != nil {
		return nil, ErrForbidden
	}

	return s.repo.Update(ctx, id, input)
}

func (s *UserService) Delete(ctx context.Context, id uuid.UUID) error {
	// implicitly check if user has sufficient permissions to update
	_, err := s.Read(ctx, id)
	if err != nil {
		return err
	}
	return s.repo.Delete(ctx, id)
}

func IsValidUsername(username string) bool {
	username = strings.TrimSpace(username)
	usernameRegex := regexp.MustCompile(`^[a-z][a-z0-9]{4,14}$`)
	return usernameRegex.MatchString(username)
}

func IsValidEmail(email string) bool {
	email = strings.TrimSpace(email)
	if email == "" {
		return false
	}

	addr, err := mail.ParseAddress(email)
	if err != nil {
		return false
	}

	// Ensure the parsed address matches exactly what was input
	// (mail.ParseAddress can successfully extract "user@domain.com" from "Name <user@domain.com")
	if addr.Address != email {
		return false
	}

	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return false
	}
	domain := parts[1]

	// A minimal regex to ensure the domain looks sane (has a dot and a TLD)
	// This catches typos like user@gmail or user@com
	domainRegex := regexp.MustCompile(`^[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return domainRegex.MatchString(domain)
}

func ValidatePassword(password string) error {
	if len(password) < 8 || len(password) > 64 {
		return errors.New("password must be between 8 and 64 chars long")
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
		case unicode.IsDigit(char):
			hasNumber = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}

	if !hasUpper {
		return errors.New("password must contain at least one uppercase letter")
	}
	if !hasLower {
		return errors.New("password must have at least one lowercase letter")
	}
	if !hasNumber {
		return errors.New("password must contain at least one number")
	}
	if !hasSpecial {
		return errors.New("password must contain at least one special character")
	}

	return nil
}

func ValidatePreferences(prefs PreferencesInput) (PreferencesInput, error) {
	if prefs.Mode != nil {
		switch *prefs.Mode {
		case "light", "dark", "system":
		case "":
			*prefs.Mode = "system"
		default:
			return prefs, errors.New(
				"invalid theme selection; must be 'light', 'dark', or 'system'",
			)
		}
	} else {
		mode := "system"
		prefs.Mode = &mode
	}

	if prefs.ItemsPerPage != nil {
		if *prefs.ItemsPerPage == 0 {
			*prefs.ItemsPerPage = 20 // Default fallback if not provided
		} else if *prefs.ItemsPerPage < 10 || *prefs.ItemsPerPage > 100 {
			return prefs, errors.New("items per page must be between 10 and 100")
		}
	} else {
		itemsPerPage := 20
		prefs.ItemsPerPage = &itemsPerPage
	}

	if prefs.Notifications == nil {
		notifs := false
		prefs.Notifications = &notifs
	}

	return prefs, nil
}
