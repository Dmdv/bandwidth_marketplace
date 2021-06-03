CREATE extension ltree;
CREATE DATABASE provider_meta;
\connect provider_meta;
CREATE USER provider_user WITH ENCRYPTED PASSWORD 'provider';
GRANT ALL PRIVILEGES ON DATABASE provider_meta TO provider_user;