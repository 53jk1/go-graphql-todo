package config

import (
	"log"
	"os"
	"strconv"
)

const (
	HOST     = "db"
	PORT     = 5432
	USER     = "postgres"
	PASSWORD = "postgres"
	DBNAME   = "postgres"
)

func SetEnvForDatabaseConnection() error {
	var err error
	err = os.Setenv("DB_HOST", HOST)
	if err != nil {
		log.Fatal(err)
		return err
	}
	err = os.Setenv("DB_PORT", strconv.Itoa(PORT))
	if err != nil {
		log.Fatal(err)
		return err
	}
	err = os.Setenv("DB_USER", USER)
	if err != nil {
		log.Fatal(err)
		return err
	}
	err = os.Setenv("DB_PASSWORD", PASSWORD)
	if err != nil {
		log.Fatal(err)
		return err
	}
	err = os.Setenv("DB_NAME", DBNAME)
	if err != nil {
		log.Fatal(err)
		return err
	}
	return nil
}
