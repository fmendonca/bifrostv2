CREATE DATABASE libvirtdb;

\c libvirtdb

CREATE TABLE hosts (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255),
    address VARCHAR(255),
    port INT,
    user VARCHAR(255),
    auth_method VARCHAR(50),
    password VARCHAR(255),
    ssh_key_path VARCHAR(255),
    created_at TIMESTAMP DEFAULT now(),
    updated_at TIMESTAMP DEFAULT now()
);

CREATE TABLE vms (
    id SERIAL PRIMARY KEY,
    host_id INT REFERENCES hosts(id),
    name VARCHAR(255),
    cpu INT,
    memory INT,
    disk INT,
    network VARCHAR(255),
    created_at TIMESTAMP DEFAULT now(),
    updated_at TIMESTAMP DEFAULT now()
);
