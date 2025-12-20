-- Habilitamos la extensión para UUIDs (Identificadores únicos universales)
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- 1. TABLA DE USUARIOS
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    username VARCHAR(50) UNIQUE NOT NULL,
    
    -- Manejo de Bankroll (Dinero)
    -- Usamos DECIMAL para precisión financiera, nunca FLOAT
    bankroll_units DECIMAL(10, 2) DEFAULT 0.00,
    bankroll_currency DECIMAL(15, 2) DEFAULT 0.00, -- Saldo en USD/COP
    currency_code VARCHAR(3) DEFAULT 'USD', -- 'USD', 'COP'
    
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- 2. TABLA DE APUESTAS (BETS)
CREATE TABLE bets (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    
    -- Metadatos Generales
    title VARCHAR(150), -- Ej: "Parlay Fin de Semana"
    is_parlay BOOLEAN DEFAULT FALSE,
    sport_key VARCHAR(50) NOT NULL, -- 'nba', 'cs2', 'lol'
    
    -- Estado de la apuesta
    status VARCHAR(20) DEFAULT 'pending' CHECK (status IN ('pending', 'won', 'lost', 'void', 'cashout')),
    
    -- Datos Financieros
    stake_units DECIMAL(10, 2) NOT NULL CHECK (stake_units > 0),
    odds DECIMAL(10, 2) NOT NULL CHECK (odds >= 1.01),
    potential_payout DECIMAL(10, 2) GENERATED ALWAYS AS (stake_units * odds) STORED,
    
    -- El corazón flexible: JSONB para guardar los "legs" o detalles específicos
    -- Aquí guardaremos el JSON que diseñaste (equipos, mercado, selección)
    details JSONB NOT NULL,
    
    -- Análisis (Tus notas y la IA)
    user_notes TEXT,
    ai_prediction JSONB, -- { "probability": 0.85, "risk": "low" }
    
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    resulted_at TIMESTAMP WITH TIME ZONE -- Cuándo se definió si ganó/perdió
);

-- 3. TABLA DE TRANSACCIONES (Ledger)
-- Para tener un historial inmutable de cada movimiento de dinero (Académico y Profesional)
CREATE TABLE transactions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id),
    bet_id UUID REFERENCES bets(id), -- Puede ser NULL si es un depósito manual
    
    amount DECIMAL(10, 2) NOT NULL, -- Positivo (Ganancia/Depósito) o Negativo (Apuesta/Retiro)
    type VARCHAR(50) NOT NULL, -- 'bet_placement', 'bet_won', 'deposit', 'adjustment'
    
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- ÍNDICES PARA OPTIMIZACIÓN (Concepto del PDF)
CREATE INDEX idx_bets_user_status ON bets(user_id, status);
CREATE INDEX idx_bets_details ON bets USING GIN (details); -- Permite buscar DENTRO del JSON