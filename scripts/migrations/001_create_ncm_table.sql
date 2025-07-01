-- Ativar extensão para UUID (se não estiver ativa)
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Criação da tabela NCM
CREATE TABLE IF NOT EXISTS ncm (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    code VARCHAR(20) NOT NULL,
    code_no_symbols VARCHAR(20) NOT NULL,
    description TEXT NOT NULL,
    initial_date VARCHAR(20),
    final_date VARCHAR(20),
    type_year_ini VARCHAR(50),
    number_ato_ini VARCHAR(20),
    year_ato_ini VARCHAR(10),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Índices para melhorar performance
CREATE INDEX IF NOT EXISTS idx_ncm_code ON ncm(code);
CREATE INDEX IF NOT EXISTS idx_ncm_description ON ncm(description);
CREATE INDEX IF NOT EXISTS idx_ncm_code_no_symbols ON ncm(code_no_symbols);

-- Trigger para atualizar updated_at automaticamente
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_ncm_updated_at 
    BEFORE UPDATE ON ncm 
    FOR EACH ROW 
    EXECUTE FUNCTION update_updated_at_column();
