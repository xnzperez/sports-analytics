package database

import (
	"embed"
	"log"
)

//go:embed migrations/*.sql
var migrationFiles embed.FS // Ahora busca en la subcarpeta "migrations" (Correcto)

// Migrate ejecuta el script SQL inicial
func Migrate() {
	log.Println("🔄 Iniciando migración de base de datos...")

	// Leemos el archivo desde la memoria incrustada
	script, err := migrationFiles.ReadFile("migrations/001_initial_schema.sql")
	if err != nil {
		log.Fatalf("❌ Error leyendo archivo de migración: %v", err)
	}

	if err := Instance.Exec(string(script)).Error; err != nil {
		log.Printf("⚠️  Advertencia al migrar: %v", err)
	} else {
		log.Println("✅ Tablas y Schema creados exitosamente.")
	}
}
