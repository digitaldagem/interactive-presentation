package storage

import (
	"database/sql"
	"fmt"
	"log"
	"reflect"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"interactive-presentation/src/config"
)

func connectToDatabase() (*sql.DB, error) {
	configuration, err := config.New()
	if err != nil {
		return nil, fmt.Errorf("error initializing config: %v", err)
	}

	db, err := sql.Open("postgres", configuration.DatabaseURL)
	if err != nil {
		return nil, fmt.Errorf("error opening db connection: %v", err)
	}
	return db, nil
}

func scanArgs(elem reflect.Value) []interface{} {
	var args []interface{}
	for i := 0; i < elem.NumField(); i++ {
		args = append(args, elem.Field(i).Addr().Interface())
	}
	return args
}

func InsertIntoDatabase(table string, args ...interface{}) error {
	db, err := connectToDatabase()
	if err != nil {
		return fmt.Errorf("error connecting to database: %v", err)
	}
	defer func() {
		if err = db.Close(); err != nil {
			log.Printf("warning: error closing database connection: %v", err)
		}
	}()

	if err != nil {
		return fmt.Errorf("error beginning database transaction: %v", err)
	}

	var insertStatement string
	switch table {
	case "presentation":
		insertStatement = fmt.Sprintf("INSERT INTO %s (presentation_id, current_poll_index) VALUES ($1, $2)", table)
	case "poll":
		insertStatement = fmt.Sprintf("INSERT INTO %s (poll_id, question, presentation_id, index) VALUES ($1, $2, $3, $4)", table)
	case "option":
		insertStatement = fmt.Sprintf("INSERT INTO %s (key, value, poll_id, index) VALUES ($1, $2, $3, $4)", table)
	case "vote":
		insertStatement = fmt.Sprintf("INSERT INTO %s (key, client_id, poll_id) VALUES ($1, $2, $3)", table)
	default:
		return fmt.Errorf("unknown table: %s", table)
	}

	var argsList []interface{}
	for _, arg := range args {
		if reflect.TypeOf(arg).Kind() == reflect.Struct {
			elem := reflect.ValueOf(arg)
			for i := 0; i < elem.NumField(); i++ {
				argsList = append(argsList, elem.Field(i).Interface())
			}
		} else {
			argsList = append(argsList, arg)
		}
	}

	_, err = db.Exec(insertStatement, argsList...)
	if err != nil {
		return fmt.Errorf("error executing insert statement for table %s: %v", table, err)
	}

	return nil
}

func SelectFromTable[T any](table string, conditionID interface{}, dest *[]T) error {
	db, err := connectToDatabase()
	if err != nil {
		return fmt.Errorf("error connecting to database: %v", err)
	}
	defer func() {
		if err = db.Close(); err != nil {
			log.Printf("warning: error closing database connection: %v", err)
		}
	}()

	var selectStatement string
	switch table {
	case "presentation":
		selectStatement = fmt.Sprintf("SELECT * FROM %s WHERE presentation_id=$1", table)
	case "poll":
		selectStatement = fmt.Sprintf("SELECT * FROM %s WHERE presentation_id=$1", table)
	case "option":
		selectStatement = fmt.Sprintf("SELECT * FROM %s WHERE poll_id=$1", table)
	case "vote":
		selectStatement = fmt.Sprintf("SELECT * FROM %s WHERE poll_id=$1", table)
	default:
		return fmt.Errorf("unknown table: %s", table)
	}

	rows, err := db.Query(selectStatement, conditionID)
	if err != nil {
		return fmt.Errorf("error running select query for table %s: %v", table, err)
	}
	defer func(rows *sql.Rows) {
		err = rows.Close()
		if err != nil {
			log.Printf("warning: error closing rows: %v", err)
		}
	}(rows)

	elemType := reflect.TypeOf(*dest).Elem()
	for rows.Next() {
		elem := reflect.New(elemType).Elem()
		if err = rows.Scan(scanArgs(elem)...); err != nil {
			return fmt.Errorf("error scanning row from select query for table %s: %v", table, err)
		}
		*dest = append(*dest, elem.Interface().(T))
	}

	if err = rows.Err(); err != nil {
		return fmt.Errorf("error iterating over rows: %v", err)
	}

	return nil
}

func UpdatePresentation(presentationID uuid.UUID, currentPollIndex int) error {
	db, err := connectToDatabase()
	if err != nil {
		return err
	}
	defer func() {
		if err = db.Close(); err != nil {
			log.Printf("warning: error closing database connection: %v", err)
		}
	}()

	query := fmt.Sprintf("UPDATE %s SET current_poll_index = $1 WHERE presentation_id = $2", "presentation")
	_, err = db.Exec(query, currentPollIndex, presentationID)
	if err != nil {
		return err
	}

	return nil
}
