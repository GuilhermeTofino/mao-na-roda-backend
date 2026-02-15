-- Habilitar extensão PostGIS (essencial para geolocalização)
CREATE EXTENSION IF NOT EXISTS postgis;

-- Tabela de Usuários
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(255) UNIQUE NOT NULL,
    name VARCHAR(255) NOT NULL,
    type VARCHAR(50) NOT NULL CHECK (type IN ('client', 'professional')),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Tabela de Profissionais
CREATE TABLE professionals (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    cpf VARCHAR(14) NOT NULL UNIQUE,
    bio TEXT,
    category VARCHAR(100) NOT NULL,
    price_hour NUMERIC(10, 2) NOT NULL,
    rating NUMERIC(3, 2) DEFAULT 0,
    review_count INTEGER DEFAULT 0,
    verified BOOLEAN DEFAULT FALSE,
    documents_url TEXT[], -- Array de URLs
    location GEOGRAPHY(POINT, 4326), -- PostGIS Point
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    CONSTRAINT unique_user_professional UNIQUE (user_id)
);

-- Índice Geoespacial (Essencial para busca por raio)
CREATE INDEX idx_professionals_location ON professionals USING GIST (location);
CREATE INDEX idx_professionals_category ON professionals (category);

-- Tabela de Solicitações de Serviço
CREATE TABLE requests (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    client_id UUID NOT NULL REFERENCES users(id),
    professional_id UUID NOT NULL REFERENCES professionals(id),
    description TEXT,
    status VARCHAR(50) NOT NULL DEFAULT 'Aberto', -- Aberto, Em andamento, Concluído, Cancelado
    scheduled_for TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Tabela de Avaliações
CREATE TABLE reviews (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    service_id UUID NOT NULL REFERENCES requests(id),
    client_id UUID NOT NULL REFERENCES users(id),
    professional_id UUID NOT NULL REFERENCES professionals(id),
    rating INTEGER NOT NULL CHECK (rating >= 1 AND rating <= 5),
    comment TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);
