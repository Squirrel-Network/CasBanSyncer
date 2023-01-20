package main

import (
	"CASBanSyncer/src/consts"
	"github.com/joho/godotenv"
	"os"
)

func InitEnvironment() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}
	consts.DbHost = os.Getenv("DB_HOST")
	consts.DbPort = os.Getenv("DB_PORT")
	consts.DbName = os.Getenv("DB_NAME")
	consts.DbCharset = os.Getenv("DB_CHARSET")
	consts.DBUser = os.Getenv("DB_USER")
	consts.DbPassword = os.Getenv("DB_PASSWORD")
}
