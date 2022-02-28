package main

import (
	"database/sql"
	//"fmt"
	"errors"

	_ "github.com/go-sql-driver/mysql"
)

type info struct {
	id    int
	name  string
	email string
	role  string
}

// func Create(db *sql.DB) {
// 	_, err := db.Exec("CREATE TABLE IF NOT EXISTS employee (id INT NOT NULL, name VARCHAR(100), email VARCHAR(100), role VARCHAR(100),PRIMARY KEY(id))")
// 	if err != nil {
// 		fmt.Printf("Error %v while creating table", err)
// 	}
// }
func Insert(db *sql.DB, a info) error {
	query := "INSERT INTO employee(id, name, email, role) VALUES(?, ?, ?, ?)"

	res, err := db.Exec(query, a.id, a.name, a.email, a.role)
	if err != nil {
		return err
	}
	_, _ = res.LastInsertId()
	return nil
}

func UpdateById(db *sql.DB, id int, Name string, Email string, role string) error {

	res, err := db.Prepare("update Employee_Details set Name=?, Email=?, role=? where id=?")
	if err != nil {
		return errors.New("prepare error")
	}

	defer res.Close()

	_, err = res.Exec(Name, Email, role, id)

	return err

}
func Read(db *sql.DB, i int) (b info, err error) {
	var a info
	row, err := db.Query("SELECT * FROM employee WHERE id = ?", i)
	if err != nil {
		return a, err
	}
	for row.Next() {
		err = row.Scan(&(a.id), &(a.name), &(a.email), &(a.role))
		if err != nil {
			return a, err
		} 
	}
	return a, err
}

func Delete(db *sql.DB, i int) error {
	_, err := db.Exec("DELETE FROM employee WHERE id = ?", i)
	if err != nil {
		return err
	}
	return nil
}

// func main() {
// 	fmt.Println("GO MYSQL TUTORIAL")

// 	db, err := sql.Open("mysql", "root:password@tcp(127.0.0.1:3306)/Employee")

// 	if err != nil {
// 		panic(err.Error())
// 	}

// 	defer db.Close()
// 	x, err := Read(db, 4)
// 	if err != nil {
// 		fmt.Printf("Error %v while reading row", err)
// 	} else {
// 		fmt.Println(x.id, " ", x.name, " ", x.email, " ", x.role)
// 	}

// 	fmt.Println("Connection done Successfully")

// 	// ins, e := db.Query("INSERT INTO employee VALUES(1,'Ishan','ishankochar0099@gmail.com','SDE-INTERN')")

// 	// if e != nil {
// 	// 	panic(e.Error())
// 	// }

// 	// defer ins.Close()

// 	// ins1, e := db.Query("INSERT INTO employee VALUES(5,'sharif','shariif999@gmail.com','SDE-INTERN')")
// 	// if e != nil {
// 	// 	panic(e.Error())
// 	// }

// 	// defer ins1.Close()
// 	// fmt.Println("Inserted value in the database")
// 	//READ THE DATABASE WITH ID
// 	// x := info{
// 	// 	id:    6,
// 	// 	name:  "mayank",
// 	// 	email: "mayank@gmail.com",
// 	// 	role:  "tester",
// 	// }
// 	//Insert(db, x)
// 	//UpdateById(db, 4, "prateekGarg", "prateekG@gmail.com", "SDET")
// 	//DELETE THE DATABASE WITH ID
// 	//Delete(db, 3)

//}
