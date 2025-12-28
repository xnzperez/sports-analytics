# 🏆 Sports Analytics Platform (Imagine Cup 2025)

> Plataforma Full Stack de análisis predictivo y gestión de apuestas deportivas, impulsada por Inteligencia Artificial y Cloud Computing.

![Go](https://img.shields.io/badge/Backend-Go_%2F_Fiber-00ADD8?style=flat&logo=go)
![React](https://img.shields.io/badge/Frontend-React_%2F_TypeScript-20232A?style=flat&logo=react)
![Azure](https://img.shields.io/badge/Cloud-Microsoft_Azure-0078D4?style=flat&logo=microsoftazure)
![PostgreSQL](https://img.shields.io/badge/Database-PostgreSQL-316192?style=flat&logo=postgresql)

## 📖 Descripción del Proyecto

**Sports Analytics** es una solución Fintech aplicada al entretenimiento deportivo. Permite a los usuarios gestionar su capital (bankroll), realizar apuestas simuladas basadas en cuotas reales de mercados de E-Sports (LoL, Valorant, CS2) y recibir consejos de inversión basados en su rendimiento histórico.

El núcleo del sistema cuenta con un **Motor de Resolución Automática (Worker)** que opera en segundo plano para liquidar apuestas, procesar transacciones financieras y mantener el estado de la plataforma actualizado en tiempo real.

## 🚀 Arquitectura Técnica

El sistema sigue una **Clean Architecture** modularizada para garantizar escalabilidad y mantenibilidad.

### Backend (Go + Fiber)
- **API RESTful:** Implementada con Fiber v2 para alto rendimiento.
- **Worker Autónomo:** Gorutinas concurrentes que monitorean y resuelven apuestas automáticamente sin intervención humana.
- **Database:** PostgreSQL con GORM. Uso estricto de **Transacciones ACID** para garantizar la integridad financiera (row-locking `FOR UPDATE`).
- **External API:** Integración con Pinnacle Odds para obtención de mercados en tiempo real.

### Frontend (React + TypeScript)
- **Dashboard en Tiempo Real:** Implementación de Long-Polling para reflejar cambios de saldo y resultados al instante.
- **UI/UX:** Diseño moderno y responsivo utilizando **TailwindCSS** y componentes de **Lucide React**.
- **State Management:** Zustand para gestión de estado global ligero y persistencia de sesión.
- **Charts:** Visualización de datos financieros con Recharts.

## 🛠️ Stack Tecnológico

| Componente | Tecnología | Uso Principal |
|------------|------------|---------------|
| **Lenguaje** | Go (Golang) 1.23 | Lógica de negocio, Worker, API |
| **Framework** | Fiber v2 | HTTP Server & Routing |
| **Base de Datos**| PostgreSQL 16 | Persistencia relacional, JSONB |
| **ORM** | GORM | Manejo de datos y migraciones |
| **Frontend** | React + Vite | Interfaz de Usuario |
| **Estilos** | Tailwind CSS | Diseño responsivo Dark Mode |
| **Cloud** | **Azure Container Apps** | Despliegue del Backend |
| **Cloud DB** | **Azure Database for PostgreSQL** | Base de datos gestionada |

## ✨ Características Clave (Key Features)

1.  **Gestión de Bankroll (Ledger):** Sistema de contabilidad de doble entrada simplificado. Cada apuesta genera una transacción inmutable.
2.  **Auto-Settlement Worker:** Un proceso en segundo plano verifica periódicamente el estado de los partidos. Si un partido termina, el sistema determina automáticamente si la apuesta fue `WON` o `LOST` y acredita las ganancias.
3.  **Prevención de Fraude:** Validaciones de saldo atómicas a nivel de base de datos para evitar condiciones de carrera (Race Conditions).
4.  **Simulación de Mercados:** Algoritmo de simulación para demostraciones en vivo (Demo Mode) que permite visualizar el ciclo completo de la apuesta en segundos.

## 📦 Instalación y Despliegue Local

### Prerrequisitos
- Go 1.23+
- Node.js 18+
- PostgreSQL local o Docker

### 1. Backend Setup
```bash
cd backend
# Crear archivo .env basado en el ejemplo
cp .env.example .env
# Instalar dependencias
go mod tidy
# Ejecutar servidor
go run cmd/api/main.go
```

### 2. Frontend Setup
```bash
cd frontend
# Crear archivo .env
echo "VITE_API_URL=http://localhost:3000" > .env
# Instalar dependencias
pnpm install
# Iniciar interfaz
pnpm dev
```

## ☁️ Infraestructura Azure
El proyecto está diseñado para desplegarse utilizando Docker containers.

**Backend:** Empaquetado en Docker y desplegado en Azure Container Apps.

**Base de Datos:** Azure Database for PostgreSQL (Flexible Server).

Desarrollado por **Carlos Pérez** para la Microsoft Imagine Cup 2025.
