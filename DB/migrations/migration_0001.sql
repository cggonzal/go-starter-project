-- NOTE: "SERIAL" is the postgres version of auto incrementing a value
CREATE TABLE users (
  id SERIAL PRIMARY KEY,
  username TEXT UNIQUE NOT NULL,
  password TEXT NOT NULL
);
