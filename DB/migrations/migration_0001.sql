-- NOTE: "SERIAL" is the postgres version of auto incrementing a value. Has type integer with max value (2^31) - 1
CREATE TABLE users (
  id BIGSERIAL PRIMARY KEY,
  email TEXT UNIQUE NOT NULL,
  password TEXT NOT NULL
);
