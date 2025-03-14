CREATE TABLE repositories (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    owner VARCHAR(255) NOT NULL,
    UNIQUE(owner, name)
)