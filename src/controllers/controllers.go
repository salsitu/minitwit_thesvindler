package controllers

import (
	"database/sql"
	"io/ioutil"
	"log"
	"reflect"

	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Id       int64  `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Pw_hash  string `json:"pw_hash"`
}

type Follower struct {
	Follower_id int `json:"follower_id"`
	Followed_id int `json:"followed_id"`
}

type Message struct {
	Message_id int64  `json:"message_id"`
	Author_id  int64  `json:"author_id"`
	Text       string `json:"text"`
	Pub_date   int64  `json:"pub_date"`
	Flagged    int64  `json:"flagged"`
}

const (
	DBPath       = "/tmp/minitwit.db"
	InitDBSchema = "../../db_init.sql"
)

func CheckError(err error) bool {
	if err != nil {
		log.Printf("Error: %s\n", err)
	}

	return err != nil
}

func Init_db(schemaDest, dbDest string) {
	query, err := ioutil.ReadFile(schemaDest)

	if CheckError(err) {
		panic(err)
	}

	db := Connect_db(dbDest)

	if _, err := db.Exec(string(query)); err != nil {
		panic(err)
	}
	db.Close()
	log.Println("Initialised database")
}

func Connect_db(dbDest string) *sql.DB {
	db, err := sql.Open("sqlite3", dbDest)
	CheckError(err)
	return db
}

func HandleQuery(rows *sql.Rows, err error) []map[string]interface{} {
	if CheckError(err) {
		return nil
	} else {
		defer rows.Close()
	}

	cols, err := rows.Columns()
	if CheckError(err) {
		return nil
	}

	values := make([]interface{}, len(cols))
	for i := range cols {
		values[i] = new(interface{})
	}

	//var dicts []map[string]interface{}
	dicts := make([]map[string]interface{}, 0)
	//dictIdx := 0

	rowsCount := 0

	for rows.Next() {
		rowsCount++
		err = rows.Scan(values...)
		if CheckError(err) {
			continue
		}

		m := make(map[string]interface{})

		for i, v := range values {
			val := reflect.Indirect(reflect.ValueOf(v)).Interface()
			m[cols[i]] = val
		}

		dicts = append(dicts, m)
		//dicts[dictIdx] = m
		//dictIdx++
	}

	log.Printf("	Columns %v returned dictionaries: %v", cols, dicts)

	log.Printf("Length of dicts: %d", len(dicts))
	return dicts
	/* if rowsCount == 0 {
		var noData []map[string]interface{}
		return noData
	} else {
		return dicts
	} */
}

// The function below has been copied from: https://gowebexamples.com/password-hashing/
func Generate_password_hash(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 8)
	return string(bytes), err
}