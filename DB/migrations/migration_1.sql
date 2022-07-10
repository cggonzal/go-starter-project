-- NOTE: "SERIAL" is the postgres version of auto incrementing a value. BIGSERIAL is a 64 bit integer.
CREATE TABLE users (
  id BIGSERIAL PRIMARY KEY,
  email TEXT UNIQUE NOT NULL,
  password TEXT NOT NULL
);
