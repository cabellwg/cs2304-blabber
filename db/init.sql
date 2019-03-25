-- Give Docker permission to run this script

CREATE USER docker;
CREATE DATABASE blabdb;
GRANT ALL PRIVILEGES ON DATABASE blabdb TO docker;
\connect blabdb


-- Create blab user

CREATE USER blabclient PASSWORD 'r$J89ka&36';


-- Build Blab tables

CREATE TABLE users (
  id serial UNIQUE,
  name varchar(255),
  email varchar(255)
);

CREATE TABLE blabs (
  id serial UNIQUE,
  postTime timestamp,
  author serial REFERENCES users(id),
  message text
);


-- Grant permissions to web api client

GRANT ALL ON users, blabs TO blabclient;
