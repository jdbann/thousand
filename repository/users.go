package repository

import (
	"context"
	"errors"

	"emailaddress.horse/thousand/form"
	"emailaddress.horse/thousand/models"
	"emailaddress.horse/thousand/repository/queries"
	"github.com/google/uuid"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
)

func (m *Repository) CreateUser(ctx context.Context, form *form.NewUserForm) (models.User, error) {
	dbUser, err := m.queries.CreateUser(ctx, queries.CreateUserParams{
		Email:    form.Email.Value,
		Password: form.Password.Value,
	})
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
			err = models.ErrEmailAlreadyInUse.Cause(err)
		}

		return models.User{}, err
	}

	return models.User{
		ID:    dbUser.ID,
		Email: dbUser.Email,
	}, nil
}

func (m *Repository) GetUser(ctx context.Context, id uuid.UUID) (models.User, error) {
	dbUser, err := m.queries.GetUser(ctx, id)
	if err != nil {
		return models.User{}, err
	}

	return models.User{
		ID:    dbUser.ID,
		Email: dbUser.Email,
	}, nil
}
