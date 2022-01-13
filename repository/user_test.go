package repository

import (
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/tacomea/worldLetter/domain"
	"log"
	"os"
	"regexp"
	"testing"
)

func openTestConnectionUser() (domain.UserRepository, sqlmock.Sqlmock, error) {
	switch os.Getenv("DB_DIALECT") {
	case "postgres":
		log.Println("testing postgres...")
		gdb, mock, err := getDBMock()
		if err != nil {
			return nil, nil, err
		}
		return NewUserRepositoryPG(gdb), mock, nil
	case "map":
		log.Println("testing sync.Map...")
		return NewSyncMapUserRepository(), nil, nil
	default:
		return nil, nil, errors.New("please set DB_DIALECT as an environment variable")
	}

}

//const CreateQuery = `INSERT INTO "users" ("email","password") VALUES ($1,$2) RETURNING "users"."email"`
const CreateQuery = `INSERT INTO "users" ("email","password") VALUES ($1,$2)`

func TestCreate(t *testing.T) {
	type args struct {
		user domain.User
	}
	want := []domain.User{
		{
			Email:    "test@domain.com",
			Password: []byte("password"),
		},
	}

	tests := []struct {
		name        string
		mockClosure func(sqlmock2 sqlmock.Sqlmock)
		args        args
		want        []domain.User
		assertion   assert.ErrorAssertionFunc
	}{
		{
			name: "Success",
			mockClosure: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(regexp.QuoteMeta(CreateQuery)).WithArgs(want[0].Email, want[0].Password).WillReturnResult(sqlmock.NewResult(1, 1))
			},
			args: args{
				user: want[0],
			},
			want:      want,
			assertion: assert.NoError,
		},
	}

	t.Parallel()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ur, mock, err := openTestConnectionUser()
			if err != nil {
				t.Fatalf("failed to init UserRepository: %v", err)
			}
			if mock != nil {
				tt.mockClosure(mock)
			}

			err = ur.Create(tt.args.user)
			tt.assertion(t, err)

			if mock != nil {
				if err = mock.ExpectationsWereMet(); err != nil {
					t.Errorf("there were unfullfilled expectations: %s", err)
				}
			}
		})
	}
}

func TestRead(t *testing.T) {
	ur, mock, err := openTestConnectionUser()
	if err != nil {
		t.Fatalf("failed to init UserRepository: %v", err)
	}

	email := "test@domain.com"
	password := []byte("password")
	if mock != nil {
		rows := sqlmock.NewRows([]string{"email", "password"}).AddRow(email, password)
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users" WHERE email = $1 ORDER BY "users"."email" LIMIT 1`)).WithArgs(email).WillReturnRows(rows)
	} else {
		err = ur.Create(domain.User{Email: email, Password: password})
		if err != nil {
			t.Fatalf("Create: %v", err)
		}
	}

	res, err := ur.Read("test@domain.com")

	if err != nil {
		t.Errorf("Read: %v", err)
	}
	if res.Email != email {
		t.Errorf("res.Email != email")
	}
	if string(res.Password) != string(password) {
		t.Errorf("res.Password != password")
	}

	if mock != nil {
		err = mock.ExpectationsWereMet()
		if err != nil {
			t.Errorf("Read: %v", err)
		}
	}
}
