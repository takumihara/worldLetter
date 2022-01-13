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
	}), &gorm.Config{
		SkipDefaultTransaction: true,
	})
	if err != nil {
		return nil, nil, err
	}
	return gdb, mock, nil
}

//const CreateQuery = `INSERT INTO "users" ("email","password") VALUES ($1,$2) RETURNING "users"."email"`
const CreateQuery = `INSERT INTO "users" ("email","password") VALUES ($1,$2)`

func TestCreate(t *testing.T) {
	type args struct {
		user domain.User
	}
	//want := []domain.User{
	//	{
	//		Email:    "test@domain.com",
	//		Password: []byte("password"),
	//	},
	//}

	//tests := []struct {
	//	name        string
	//	mockClosure func(sqlmock2 sqlmock.Sqlmock)
	//	args        args
	//	want        []domain.User
	//	assertion   assert.ErrorAssertionFunc
	//}{
	//	{
	//		name: "Success",
	//		mockClosure: func(mock sqlmock.Sqlmock) {
	//			//mock.ExpectBegin()
	//			//rows := sqlmock.NewRows([]string{"email", "password"}).AddRow(want[0].Email, want[0].Password)
	//			//mock.ExpectExec(CreateQuery).WithArgs(want[0].Email, want[0].Password).WillReturnResult(sqlmock.NewResult(1, 1))
	//			mock.ExpectQuery(CreateQuery).WithArgs(want[0].Email, want[0].Password)
	//		},
	//		args: args{
	//			user: want[0],
	//		},
	//		want:      want,
	//		assertion: assert.NoError,
	//	},
	//}
	//
	//t.Parallel()

	//for _, tt := range tests {
	//	t.Run(tt.name, func(t *testing.T) {
	//		db, mock, err := getDBMock()
	//		if err != nil {
	//			log.Fatalln("failed to init db mock:", err)
	//		}
	//		tt.mockClosure(mock)
	//		cr := &userRepositoryPG{db: db}
	//
	//		err = cr.Create(tt.args.user)
	//		tt.assertion(t, err)
	//
	//		//mock.ExpectClose()
	//		//if err = mock.ExpectationsWereMet(); err != nil {
	//		//	t.Errorf("there were unfullfilled expectations: %s", err)
	//		//}
	//	})
	//}

	db, mock, err := getDBMock()
	if err != nil {
		log.Fatalln("failed to init db mock:", err)
	}
	user := domain.User{
		Email:    "test@domain.com",
		Password: []byte("password"),
	}
	mock.ExpectBegin()
	mock.ExpectQuery(CreateQuery).WithArgs(user.Email, user.Password)
	mock.ExpectCommit()

	cr := &userRepositoryPG{db: db}

	err = cr.Create(user)
	assert.NoError(t, err)

}

func TestRead(t *testing.T) {
	db, mock, err := getDBMock()
	if err != nil {
		t.Errorf("Failed to initialize mock DB: %v", err)
		return
	}

	email := "test@domain.com"
	password := []byte("password")
	ur := &userRepositoryPG{db: db}
	rows := sqlmock.NewRows([]string{"email", "password"}).AddRow(email, password)
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users" WHERE email = $1 ORDER BY "users"."email" LIMIT 1`)).WithArgs(email).WillReturnRows(rows)

	res, err := ur.Read("test@domain.com")

	assert.Equal(t, err, nil)
	assert.Equal(t, res.Email, email)
	assert.Equal(t, res.Password, password)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Read: %v", err)
	}
}
