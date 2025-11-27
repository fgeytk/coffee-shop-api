package database

import (
    "database/sql"
    "fmt"
    "log"
    "os"

    _ "github.com/go-sql-driver/mysql"  // <- Ajouter le underscore
)

var DB *sql.DB

func InitDB() {
    host := getEnv("DB_HOST", "localhost")
    port := getEnv("DB_PORT", "3306")
    user := getEnv("DB_USER", "coffee_user")
    password := getEnv("DB_PASSWORD", "coffee123")
    dbname := getEnv("DB_NAME", "coffee_shop")

    dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", user, password, host, port, dbname)
    
    var err error
    DB, err = sql.Open("mysql", dsn)
    if err != nil {
        log.Fatal("Error opening database: ", err)
    }

    if err = DB.Ping(); err != nil {
        log.Fatal("Error connecting to database: ", err)
    }

    fmt.Println("✅ Connected to MySQL database!")
}

func getEnv(key, defaultValue string) string {
    if value := os.Getenv(key); value != "" {
        return value
    }
    return defaultValue
}