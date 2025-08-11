CREATE TABLE books (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    title STRING NOT NULL,
    author STRING NOT NULL,
    edition INT NOT NULL,
    isbn STRING NOT NULL,
    created_at TIMESTAMPTZ DEFAULT now(),
    updated_at TIMESTAMPTZ DEFAULT now()
);