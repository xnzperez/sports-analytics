package database

import (
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Instance mantendrá la conexión activa a la DB
var Instance *gorm.DB

// Connect inicializa la conexión a PostgreSQL
func Connect() {
	// 1. Construimos el DSN
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_SSLMODE"),
	)

	// 2. Abrimos la conexión con el Logger forzado en INFO
	// Esto nos permitirá ver la consulta SQL exacta en la terminal
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		log.Fatal("❌ No se pudo conectar a la base de datos: ", err)
	}

	// 3. Configuración del Connection Pool
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal("❌ Error obteniendo la instancia genérica de DB")
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	Instance = db
	log.Println("⚡ [STAKEWISE-CLOUD] Conexión establecida")
	log.Println("🔍 Logs SQL activados para depuración")
}
