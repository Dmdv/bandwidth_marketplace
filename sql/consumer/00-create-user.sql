CREATE extension ltree;
CREATE DATABASE consumer_meta;
\connect consumer_meta;
CREATE USER consumer_user WITH ENCRYPTED PASSWORD 'consumer';
GRANT ALL PRIVILEGES ON DATABASE consumer_meta TO consumer_user;