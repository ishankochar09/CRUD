package main

import (
	"database/sql"
	"errors"
	"reflect"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	_ "github.com/go-sql-driver/mysql"
)

func TestGetById(t *testing.T) {

	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Errorf("an error '%s' was not expected when opening a stub database connection", err)
	}

	defer db.Close()
	rows := sqlmock.NewRows([]string{"id", "name", "email", "role"}).
		AddRow(1, "kochar", "kochar@gmail.com", "sde")

	testCases := []struct {
		id          int
		index       info
		mockQuery   interface{}
		expectError error
	}{
		// Success case
		{
			id:          1,
			index:       info{1, "kochar", "kochar@gmail.com", "sde"},
			mockQuery:   mock.ExpectQuery("SELECT * FROM employee WHERE id = ?").WithArgs(1).WillReturnRows(rows),
			expectError: nil},
		// Error case
		{
			id:          3,
			index:       info{},
			mockQuery:   mock.ExpectQuery("SELECT * FROM employee WHERE id = ?").WithArgs(3).WillReturnError(sql.ErrNoRows),
			expectError: sql.ErrNoRows},
	}

	for i, testCase := range testCases {
		t.Run("", func(t *testing.T) {
			user, err := Read(db, testCase.id)
			if err != nil && err.Error() != testCase.expectError.Error() {
				t.Errorf("expected error: %v, got: %v test %v", testCase.expectError, err, i+1)
			}
			if err == nil && !reflect.DeepEqual(user, testCase.index) {
				t.Errorf("expected user: %v, got: %v test %v", testCase.index, user, i+1)
			}
		})
	}
}
func TestInsertValues(t *testing.T) {

	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Errorf("an error '%s' was not expected when opening a stub database connection", err)
	}

	defer db.Close()

	testCases := []struct {
		employee    info
		mockQuery   interface{}
		expectError error
	}{
		// Success case
		{
			employee:    info{1, "kochar", "kochar@gmail.com", "sde"},
			mockQuery:   mock.ExpectExec("INSERT INTO employee(id, name, email, role) VALUES(?, ?, ?, ?)").WithArgs(1, "kochar", "kochar@gmail.com", "sde").WillReturnResult(sqlmock.NewResult(1, 1)),
			expectError: nil},
		//failure case
		{
			employee:       info{1, "kochar", "kochar@gmail.com", "sde"},
			mockQuery:   mock.ExpectExec("INSERT INTO employee(id, name, email, role) VALUES(?, ?, ?, ?)").WithArgs(1, "kochar", "kochar@gmail.com", "sde").WillReturnError((errors.New("error while inserting"))),
			expectError: errors.New("error while inserting")},
	}

	for i, testCase := range testCases {
		t.Run("", func(t *testing.T) {
			err := Insert(db, testCase.employee)
			if err != nil && !reflect.DeepEqual(err, testCase.expectError) {
				t.Errorf("expected error: %v, got: %v test %v", testCase.expectError, err, i+1)
			}
		})
	}
}
func TestDeleteById(t *testing.T) {

    db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Errorf("an error '%s' was not expected when opening a stub database connection", err)
	}

	defer db.Close()

	testCases := []struct {
		id          int
		employee    info
		mockQuery   interface{}
		expectError error
	}{
		// Success case
		{
			id : 1,
			employee:    info{1, "kochar", "kochar@gmail.com", "sde"},
			mockQuery:   mock.ExpectExec("DELETE FROM employee WHERE id = ?").WithArgs(1).WillReturnResult(sqlmock.NewResult(1, 1)),
			expectError: nil},
		//failure case
		{
			id : 2,
			employee:       info{1, "kochar", "kochar@gmail.com", "sde"},
			mockQuery:   mock.ExpectExec("DELETE FROM employee WHERE id = ?").WithArgs(2).WillReturnError((errors.New("error while inserting"))),
			expectError: errors.New("error while inserting")},
	}

    for i, testCase := range testCases {
		t.Run("", func(t *testing.T) {
			err := Delete(db, testCase.id)
			if err != nil && !reflect.DeepEqual(err, testCase.expectError) {
				t.Errorf("expected error: %v, got: %v test %v", testCase.expectError, err, i+1)
			}
		})
	}
}
func TestUpdateById(t *testing.T) {
    db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
    if err != nil {
        t.Errorf(err.Error())
    }
    defer db.Close()

    query := "update Employee_Details set Name=?, Email=?, role=? where id=?"
    kquery := "update Employee_Details1 set Name=?, Email=?, role=? where id=?"

    tc := []struct {
        id    int
        Name  string
        Email string
        role  string
        err   error
        mockQ interface{}
    }{
        {
            id:    1,
            Name:  "Ishan",
            Email: "ishankochar@gmail.com",
            role:  "SDE-I",
            err:   nil,
            mockQ: mock.ExpectPrepare(query).ExpectExec().WithArgs("Ishan", "ishankochar@gmail.com", "SDE-I", 1).WillReturnResult(sqlmock.NewResult(1, 0)),
        },
        {
            id:    2,
            Name:  "prateek",
            Email: "prTEEK@k.com",
            role:  "SDE-I",
            err:   sql.ErrNoRows,
            mockQ: mock.ExpectPrepare(query).ExpectExec().WithArgs("prateek", "prTEEK@k.com", "SDE-I", 2).WillReturnError(sql.ErrNoRows),
        },
        {
            id:    2,
            Name:  "KAKA",
            Email: "k@k.com",
            role:  "SDE-I",
            err:   errors.New("prepare error"),
            mockQ: mock.ExpectPrepare(kquery).ExpectExec().WithArgs("KAKA", "k@k.com", "SDE-I", 3).WillReturnError(errors.New("Prepare Table")),
        },
    }

    for _, tt := range tc {
        t.Run("", func(t *testing.T) {

            err := UpdateById(db, tt.id, tt.Name, tt.Email, tt.role)
            if err != nil && err.Error() != tt.err.Error() {
                t.Errorf("expected error:%v, got:%v", tt.err, err)
            }

        })
    }

}

// func TestUpdateById(t *testing.T) {

// 	db, mock := NewMock()

// 	query := "update Employee_Details set Name=?, Email=?, role=? where id=?"

// 	prep := mock.ExpectPrepare(query)
// 	prep.ExpectExec().WithArgs("sukant", "sukant@zopsmart.com", "sde", 1).WillReturnResult(sqlmock.NewResult(0, 1))

// 	err := UpdateById(db, emp.id, emp.Name, emp.Email, emp.role)

// 	assert.NoError(t, err)
// }

// func TestUpdateByIdError(t *testing.T) {

// 	db, mock := NewMock()

// 	query := "update Employee_Details set Name=?, Email=?, role=? where id=?"

// 	prep := mock.ExpectPrepare(query)
// 	prep.ExpectExec().WithArgs("sukant", "sukant@zopsmart.com", "sde", 1).WillReturnResult(sqlmock.NewResult(0, 0))

// 	err := UpdateById(db, emp.id, emp.Name, emp.Email, emp.role)

// 	assert.Error(t, err)
// }

// func TestDeleteById(t *testing.T) {

// 	db, mock := NewMock()
// 	query := "delete from Employee_Details where id=?"

// 	prep := mock.ExpectPrepare(query)
// 	prep.ExpectExec().WithArgs(1).WillReturnResult(sqlmock.NewResult(0, 1))

// 	err := DeleteById(db, 1)

// 	assert.NoError(t, err)

// }
