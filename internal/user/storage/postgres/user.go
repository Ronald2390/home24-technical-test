package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"home24-technical-test/internal/user/model"
	"home24-technical-test/internal/user/public"

	"home24-technical-test/pkg/appcontext"
	"home24-technical-test/pkg/data"

	"github.com/jmoiron/sqlx"
)

// PostgresStorage implements the user repository service interface
type PostgresStorage struct {
	db *sqlx.DB
}

// FindByID get user by userId
func (s *PostgresStorage) FindByID(ctx context.Context, userID int) (*model.User, error) {
	user := &model.User{}

	rows, err := s.db.NamedQuery(`
	SELECT 
		"id", "name", "email", "address", "password", "createdBy" 
	FROM
		"user"
	WHERE
		"id" = :id`, map[string]interface{}{
		"id": userID,
	})
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, data.ErrNotFound
		}
		return nil, err
	}

	for rows.Next() {
		err = rows.StructScan(user)
		if err != nil {
			return nil, err
		}
	}

	return user, nil
}

// FindByEmail get user by email
func (s *PostgresStorage) FindByEmail(ctx context.Context, email string) (*model.User, error) {
	user := &model.User{}

	rows, err := s.db.NamedQuery(`
	SELECT 
		"id", "name", "email", "address", "password", "createdBy" 
	FROM
		"user"
	WHERE
		"email" = :email`, map[string]interface{}{
		"email": email,
	})
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, data.ErrNotFound
		}
		return nil, err
	}

	for rows.Next() {
		err = rows.StructScan(user)
		if err != nil {
			return nil, err
		}
	}

	return user, nil
}

// Insert inserts an user
func (s *PostgresStorage) Insert(ctx context.Context, singleUser *model.User) error {
	rows, err := s.db.NamedQuery(`
	INSERT INTO 
		"user" ("name","email","address","password","createdBy", "createdAt", "updatedAt", "updatedBy")
	VALUES
		(:name, :email, :address, :password, :createdBy, now(), :createdBy, now())
	RETURNING
		"id", "name","email","address","password","createdBy", "createdAt", "updatedAt", "updatedBy"`, singleUser)
	if err != nil {
		return err
	}

	for rows.Next() {
		err = rows.StructScan(singleUser)
		if err != nil {
			return err
		}
	}

	return nil
}

// Delete user data
func (s *PostgresStorage) Delete(ctx context.Context, userID int) error {
	_, err := s.db.Queryx(`
	UPDATE "user" 
	SET
		"deletedAt" = NOW(),
		"deletedBy" = :deletedBy
	WHERE
		"id" = ?`, userID)
	if err != nil {
		return err
	}

	return nil
}

// Update user data
func (s *PostgresStorage) Update(ctx context.Context, updatedUser *model.User) error {
	updatedUser.UpdatedAt = time.Now()
	updatedUser.UpdatedBy = appcontext.UserID(ctx)

	fmt.Printf("Query: %s", `
	UPDATE "user" 
	SET
		"name" = :name,
		"email" = :email,
		"address" = :address,
		"password" = :password,
		"updatedAt" = :updatedAt,
		"updatedBy" = :updatedBy
	WHERE
		"id" = :id
	RETURNING
		"id", "name","email","address","password","createdBy", "createdAt", "updatedAt", "updatedBy"`)
	fmt.Printf("arg: %v", map[string]interface{}{
		"id":        updatedUser.ID,
		"name":      updatedUser.Name,
		"email":     updatedUser.Email,
		"address":   updatedUser.Address,
		"password":  updatedUser.Password,
		"updatedAt": updatedUser.UpdatedAt,
		"updatedBy": updatedUser.UpdatedBy,
	})

	rows, err := s.db.NamedQuery(`
	UPDATE "user" 
	SET
		"name" = :name,
		"email" = :email,
		"address" = "address,
		"password" = :password,
		"updatedAt" = :updatedAt,
		"updatedBy" = :updatedBy
	WHERE
		"id" = :id
	RETURNING
		"id", "name","email","address","password","createdBy", "createdAt", "updatedAt", "updatedBy"`,
		map[string]interface{}{
			"id":        updatedUser.ID,
			"name":      updatedUser.Name,
			"email":     updatedUser.Email,
			"address":   updatedUser.Address,
			"password":  updatedUser.Password,
			"updatedAt": updatedUser.UpdatedAt,
			"updatedBy": updatedUser.UpdatedBy,
		})
	if err != nil {
		return err
	}

	for rows.Next() {
		err = rows.StructScan(updatedUser)
		if err != nil {
			return err
		}
	}

	return nil
}

//FindAll is stand for
func (s *PostgresStorage) FindAll(ctx context.Context, params *public.FindAllUsersParams) ([]*model.User, error) {
	users := []*model.User{}

	where := `"deletedAt" IS NULL`
	if params.Email != "" {
		where += ` AND "email" ILIKE :email`
	}
	if params.Name != "" {
		where += ` AND "name" ILIKE :name`
	}
	if params.Search != "" {
		where += ` AND "name" ILIKE :search`
	}
	if params.Page != 0 && params.Limit != 0 {
		where = fmt.Sprintf(`%s ORDER BY "id" DESC LIMIT :limit OFFSET :offset`, where)
	} else {
		where = fmt.Sprintf(`%s ORDER BY "id" DESC`, where)
	}

	rows, err := s.db.NamedQuery(fmt.Sprintf(`
	SELECT 
		"id", "name", "email", "address", "password", "createdBy" 
	FROM
		"user"
	WHERE
		%s`, where), map[string]interface{}{
		"limit":  params.Limit,
		"email":  params.Email,
		"name":   params.Name,
		"search": "%" + params.Search + "%",
		"offset": ((params.Page - 1) * params.Limit),
	})
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, data.ErrNotFound
		}
		return nil, err
	}

	for rows.Next() {
		user := &model.User{}
		err = rows.StructScan(user)
		if err != nil {
			return nil, err
		}

		users = append(users, user)
	}

	return users, nil
}

// NewPostgresStorage creates new user repository service
func NewPostgresStorage(
	db *sqlx.DB,
) *PostgresStorage {
	return &PostgresStorage{
		db: db,
	}
}
