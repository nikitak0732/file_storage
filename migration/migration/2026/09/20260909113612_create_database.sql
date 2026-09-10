-- migrate:up

-- Даёт функцию gen_random_uuid().
CREATE EXTENSION IF NOT EXISTS pgcrypto;

-- Пользователи.
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    login VARCHAR(64) NOT NULL UNIQUE,
    password_hash TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- Авторизованные сессии.
CREATE TABLE sessions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,

    -- SHA-256 от исходного токена в hex-виде.
    token_hash CHAR(64) NOT NULL UNIQUE,

    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    expires_at TIMESTAMPTZ
);

-- JSON- и файловые документы.
CREATE TABLE documents (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    owner_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,

    name VARCHAR(255) NOT NULL,
    mime VARCHAR(255),
    is_file BOOLEAN NOT NULL DEFAULT false,
    is_public BOOLEAN NOT NULL DEFAULT false,

    json_data JSONB,
    file_data BYTEA,

    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- Выданные права на непубличные документы.
CREATE TABLE document_grants (
    document_id UUID NOT NULL REFERENCES documents(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),

    PRIMARY KEY (document_id, user_id)
);

-- Индексы добавим после появления нагрузки и понимания реальных запросов.

-- CREATE INDEX sessions_user_id_idx ON sessions(user_id);

-- CREATE INDEX documents_owner_name_created_idx
--     ON documents(owner_id, name, created_at);

-- CREATE INDEX documents_public_name_created_idx
--     ON documents(name, created_at)
--     WHERE is_public = true;

-- CREATE INDEX document_grants_user_document_idx
--     ON document_grants(user_id, document_id);

-- CREATE INDEX documents_json_data_gin_idx
--     ON documents USING GIN(json_data);

-- migrate:down

