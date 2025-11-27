package database

import (
    "database/sql"
    "fmt"
    "log"
    "os"
    "github.com/joho/godotenv"
    _ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

func InitDB() {
    // Charger le fichier .env
    if err := godotenv.Load(); err != nil {
        log.Println("Warning: .env file not found, using environment variables")
    }

    host := os.Getenv("DB_HOST")
    port := os.Getenv("DB_PORT")
    user := os.Getenv("DB_USER")
    password := os.Getenv("DB_PASSWORD")
    dbname := os.Getenv("DB_NAME")

    // Vérifier que les variables sont définies
    if host == "" || port == "" || user == "" || password == "" || dbname == "" {
        log.Fatal("❌ Missing database configuration in .env file")
    }

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