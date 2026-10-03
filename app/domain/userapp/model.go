package userapp

import (
	"encoding/json"
	"fmt"
	"net/mail"
	"time"

	"github.com/emorydu/service/app/sdk/errs"
	"github.com/emorydu/service/business/domain/userbus"
	"github.com/emorydu/service/business/types/name"
	"github.com/emorydu/service/business/types/password"
	"github.com/emorydu/service/business/types/role"
)

// User represents information about an individual user.
type User struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Email       string   `json:"email"`
	Roles       []string `json:"roles"`
	Department  string   `json:"department"`
	Enabled     bool     `json:"enabled"`
	DateCreated string   `json:"dateCreated"`
	DateUpdated string   `json:"dateUpdated"`
}

// Encode implements the encoder interface.
func (app User) Encode() ([]byte, string, error) {
	data, err := json.Marshal(app)

	return data, "application/json", err
}

func toAppUser(bus userbus.User) User {
	return User{
		ID:          bus.ID.String(),
		Name:        bus.Name.String(),
		Email:       bus.Email.Address,
		Roles:       role.ParseToString(bus.Roles),
		Department:  bus.Department.String(),
		Enabled:     bus.Enabled,
		DateCreated: bus.DateCreated.Format(time.RFC3339),
		DateUpdated: bus.DateUpdated.Format(time.RFC3339),
	}
}

func toAppUsers(users []userbus.User) []User {
	app := make([]User, len(users))
	for i, usr := range users {
		app[i] = toAppUser(usr)
	}

	return app
}

// =============================================================================

// NewUser defines the data needed to add a new user.
type NewUser struct {
	Name            string   `json:"name"`
	Email           string   `json:"email"`
	Roles           []string `json:"roles"`
	Department      string   `json:"department"`
	Password        string   `json:"password"`
	PasswordConfirm string   `json:"passwordConfirm"`
}

// Decode implements the decoder interface.
func (app *NewUser) Decode(data []byte) error {
	return json.Unmarshal(data, app)
}

func toBusNewUser(nu NewUser) (userbus.NewUser, error) {
	var errors errs.FieldErrors

	roles, err := role.ParseMany(nu.Roles)
	if err != nil {
		errors.Add("roles", err)
	}

	addr, err := mail.ParseAddress(nu.Email)
	if err != nil {
		errors.Add("email", err)
	}

	nme, err := name.Parse(nu.Name)
	if err != nil {
		errors.Add("name", err)
	}

	department, err := name.ParseNull(nu.Department)
	if err != nil {
		errors.Add("department", err)
	}

	pass, err := password.ParseConfirm(nu.Password, nu.PasswordConfirm)
	if err != nil {
		errors.Add("password", err)
	}

	if len(errors) > 0 {
		return userbus.NewUser{}, fmt.Errorf("validate: %w", errors.ToError())
	}

	bnu := userbus.NewUser{
		Name:       nme,
		Email:      *addr,
		Roles:      roles,
		Department: department,
		Password:   pass,
	}

	return bnu, nil
}

// =============================================================================

// UpdateUserRole defines the data needed to update a user role.
type UpdateUserRole struct {
	Roles []string `json:"roles"`
}

// Decode implements the decoder interface.
func (app *UpdateUserRole) Decode(data []byte) error {
	return json.Unmarshal(data, app)
}

func toBusUpdateUserRole(uur UpdateUserRole) (userbus.UpdateUser, error) {
	var errors errs.FieldErrors

	var roles []role.Role
	if uur.Roles != nil {
		var err error
		roles, err = role.ParseMany(uur.Roles)
		if err != nil {
			errors.Add("roles", err)
		}
	}

	if len(errors) > 0 {
		return userbus.UpdateUser{}, fmt.Errorf("validate: %w", errors.ToError())
	}

	buu := userbus.UpdateUser{
		Roles: roles,
	}

	return buu, nil
}

// =============================================================================

// UpdateUser defines the data needed to update a user.
type UpdateUser struct {
	Name            *string `json:"name"`
	Email           *string `json:"email"`
	Department      *string `json:"department"`
	Password        *string `json:"password"`
	PasswordConfirm *string `json:"passwordConfirm"`
	Enabled         *bool   `json:"enabled"`
}

// Decode implements the decoder interface.
func (app *UpdateUser) Decode(data []byte) error {
	return json.Unmarshal(data, app)
}

func toBusUpdateUser(uu UpdateUser) (userbus.UpdateUser, error) {
	var errors errs.FieldErrors

	var addr *mail.Address
	if uu.Email != nil {
		var err error
		addr, err = mail.ParseAddress(*uu.Email)
		if err != nil {
			errors.Add("email", err)
		}
	}

	var nme *name.Name
	if uu.Name != nil {
		nm, err := name.Parse(*uu.Name)
		if err != nil {
			errors.Add("name", err)
		}
		nme = &nm
	}

	var department *name.Null
	if uu.Department != nil {
		dep, err := name.ParseNull(*uu.Department)
		if err != nil {
			errors.Add("department", err)
		}
		department = &dep
	}

	var pass *password.Password
	p, err := password.ParseConfirmPointers(uu.Password, uu.PasswordConfirm)
	if err != nil {
		errors.Add("password", err)
	}
	pass = &p

	if len(errors) > 0 {
		return userbus.UpdateUser{}, fmt.Errorf("validate: %w", errors.ToError())
	}

	buu := userbus.UpdateUser{
		Name:       nme,
		Email:      addr,
		Department: department,
		Password:   pass,
		Enabled:    uu.Enabled,
	}

	return buu, nil
}
