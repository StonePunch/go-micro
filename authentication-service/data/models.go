package data

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"time"

	"golang.org/x/crypto/bcrypt"
)

const dbTimeout = time.Second * 3

var client *sql.DB

// Models is the type for this package. Any model that is included as a member
// in this type is available throughout the application, anywhere that the app
// variable is used, provided that the model is also added in the New function
type Models struct {
	User User
}

// User is the structure which represents one user from the database
type User struct {
	ID        int       `json:"id"`
	Email     string    `json:"email"`
	FirstName string    `json:"first_name,omitempty"`
	LastName  string    `json:"last_name,omitempty"`
	Password  string    `json:"-"`
	Active    int       `json:"active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// New creates an instance of the data package. It returns the type
// Models, which embeds all the types needed for the application
func New(dbPool *sql.DB) Models {
	client = dbPool

	return Models{
		User: User{},
	}
}

// GetAll returns a slice of all users, sorted by last name
func (u *User) GetAll() ([]*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	sql := `
	SELECT
		id,
		email,
		first_name,
		last_name,
		password,
		user_active,
		created_at,
		updated_at
	FROM
		users
	ORDER BY
		last_name
	`

	rows, err := client.QueryContext(ctx, sql)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*User

	// TODO: Check possible errors from rows.Next()
	for rows.Next() {
		var user User
		err := rows.Scan(
			&user.ID,
			&user.Email,
			&user.FirstName,
			&user.LastName,
			&user.Password,
			&user.Active,
			&user.CreatedAt,
			&user.UpdatedAt,
		)
		if err != nil {
			log.Println("Error scanning", err)
			return nil, err
		}

		users = append(users, &user)
	}

	return users, nil
}

// GetByEmail returns one user by email
func (u *User) GetByEmail(email string) (*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	sql := `
	SELECT
		id,
		email,
		first_name,
		last_name,
		password,
		user_active,
		created_at,
		updated_at
	FROM
		users
	WHERE
		email = $1
	`

	var user User
	row := client.QueryRowContext(ctx, sql, email)

	err := row.Scan(
		&user.ID,
		&user.Email,
		&user.FirstName,
		&user.LastName,
		&user.Password,
		&user.Active,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

// GetOne returns one user by id
func (u *User) GetOne(id int) (*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	sql := `
	SELECT
		id,
		email,
		first_name,
		last_name,
		password,
		user_active,
		created_at,
		updated_at
	FROM
		users
	WHERE
		id = $1
	`

	var user User
	row := client.QueryRowContext(ctx, sql, id)

	err := row.Scan(
		&user.ID,
		&user.Email,
		&user.FirstName,
		&user.LastName,
		&user.Password,
		&user.Active,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

// Update updates one user using the information stored in the receiver
func (u *User) Update() error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	sql := `
	UPDATE users
	SET
		email = $1,
		first_name = $2,
		last_name = $3,
		user_active = $4,
		updated_at = $5
	WHERE
		id = $6
	`

	_, err := client.ExecContext(ctx, sql,
		u.Email,
		u.FirstName,
		u.LastName,
		u.Active,
		time.Now(),
		u.ID,
	)

	if err != nil {
		return err
	}

	return nil
}

// Delete deletes one user using the ID stored in the receiver
func (u *User) Delete() error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	sql := `
	DELETE FROM users
	WHERE
		id = $1
	`

	_, err := client.ExecContext(ctx, sql, u.ID)
	if err != nil {
		return err
	}

	return nil
}

// DeleteByID deletes one user by id
func (u *User) DeleteByID(id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	sql := `
	DELETE FROM users
	WHERE
		id = $1
	`

	_, err := client.ExecContext(ctx, sql, id)
	if err != nil {
		return err
	}

	return nil
}

// Insert adds a new user into the database and returns the id of the newly
// inserted row
func (u *User) Insert(user User) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), 12)
	if err != nil {
		return 0, err
	}

	var newID int
	sql := `
	INSERT INTO users
		(email, first_name, last_name, password, user_active, created_at, updated_at)
	VALUES
		($1, $2, $3, $4, $5, $6, $7)
	RETURNING
		id
	`

	err = client.QueryRowContext(ctx, sql,
		user.Email,
		user.FirstName,
		user.LastName,
		hashedPassword,
		user.Active,
		time.Now(),
		time.Now(),
	).Scan(&newID)

	if err != nil {
		return 0, err
	}

	return newID, nil
}

// ResetPassword sets the password of the user
func (u *User) ResetPassword(password string) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return err
	}

	sql := `
	UPDATE users
	SET
		password = $1
	WHERE
		id = $2
	`
	_, err = client.ExecContext(ctx, sql, hashedPassword, u.ID)
	if err != nil {
		return err
	}

	return nil
}

// PasswordMatches uses Go's bcrypt package to compare a user's supplied
// password with the hash stored in the database. Returns true if the password
// and hash match
func (u *User) PasswordMatches(plainText string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(plainText))
	if err != nil {
		switch {
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			// invalid password
			return false, nil
		default:
			return false, err
		}
	}

	return true, nil
}
