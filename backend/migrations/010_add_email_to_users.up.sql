ALTER TABLE users ADD COLUMN email VARCHAR(255);
UPDATE users SET email = 'admin@HaberBot.com' WHERE username = 'admin1';
ALTER TABLE users ALTER COLUMN email SET NOT NULL;
ALTER TABLE users ADD CONSTRAINT users_email_key UNIQUE (email);
