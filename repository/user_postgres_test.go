package repository

import (
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/tacomea/worldLetter/domain"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"regexp"
	"testing"
)

func getDBMock() (*gorm.DB, sqlmock.Sqlmock, error) {
	sqlDB, mock, err := sqlmock.New()
	if err != nil {
		return nil, nil, err
	}

	gdb, err := gorm.Open(postgres.New(postgres.Config{
		Conn: sqlDB,
	}), &gorm.Config{})
	if err != nil {
		return nil, nil, err
	}
	return gdb, mock, nil
}

const CreateQuery = `INSERT INTO "users" ("email","password") VALUES ($1,$2) RETURNING "users"."email"`
//const CreateQuery = `INSERT INTO "users" ("email","password") VALUES ($1,$2)"`

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
				rows := sqlmock.NewRows([]string{"email", "password"}).AddRow(want[0].Email, want[0].Password)
				mock.ExpectQuery(regexp.QuoteMeta(CreateQuery)).WithArgs(want[0].Email, want[0].Password).WillReturnRows(rows)
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
			db, mock, err := getDBMock()
			if err != nil {
				log.Fatalln("failed to init db mock:", err)
			}
			tt.mockClosure(mock)
			cr := &userRepositoryPG{db: db}

			err = cr.Create(tt.args.user)
			tt.assertion(t, err)

			mock.ExpectClose()
			//if err = mock.ExpectationsWereMet(); err != nil {
			//	t.Errorf("there were unfullfilled expectations: %s", err)
			//}
		})
	}
}
