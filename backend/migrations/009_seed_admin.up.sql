INSERT INTO users (id, username, password_hash, role)
VALUES ('admin1-uuid-placeholder-1234567890ab', 'admin1', '$2a$10$6Wqf8hI7vR9Z6J1uXJgYGea7a8dG09R8a9A5S3g4M6iY7a4r8d3aK', 'Admin')
ON CONFLICT (username) DO NOTHING;
