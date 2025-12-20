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
	// 1. Construimos el DSN (Data Source Name) usando las variables de entorno
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_SSLMODE"),
	)

	// 2. Configuración Avanzada (Senior Tip)
	// Usamos un logger silencioso en prod para no saturar los logs, pero info en dev.
	logLevel := logger.Silent
	if os.Getenv("ENV") == "development" {
		logLevel = logger.Info
	}

	// 3. Abrimos la conexión
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logLevel),
	})

	if err != nil {
		log.Fatal("❌ No se pudo conectar a la base de datos: ", err)
	}

	// 4. Configuración del Connection Pool (Optimizaciones de rendimiento)
	// Esto es vital para APIs de alto tráfico.
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal("❌ Error obteniendo la instancia genérica de DB")
	}

	// SetMaxIdleConns: Cuantas conexiones mantener "dormidas" listas para usar.
	sqlDB.SetMaxIdleConns(10)

	// SetMaxOpenConns: Máximo de conexiones simultáneas (Evita tumbar la base de datos).
	sqlDB.SetMaxOpenConns(100)

	// SetConnMaxLifetime: Cuánto tiempo puede vivir una conexión antes de ser reciclada.
	sqlDB.SetConnMaxLifetime(time.Hour)

	Instance = db
	log.Println("✅ Conexión a PostgreSQL exitosa y Pool configurado.")
}