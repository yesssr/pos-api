CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE TYPE roles AS ENUM ('admin', 'user');

CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    username VARCHAR(50) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    role roles NOT NULL,
    image_url TEXT NOT NULL DEFAULT '',
    is_active BOOLEAN DEFAULT TRUE
);
