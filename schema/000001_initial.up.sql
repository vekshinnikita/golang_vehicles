CREATE TABLE users (
    id serial NOT NULL UNIQUE,
    name VARCHAR(255) not null,
    surname VARCHAR(255) not null,
    patronymic VARCHAR(255)
);

CREATE TABLE vehicles (
    id serial NOT NULL UNIQUE,
    model VARCHAR(255) not null,
    mark VARCHAR(255) not null,
    year int not null,
    reg_num VARCHAR(255) not null,
    owner_id INTEGER not null REFERENCES users(id) ON DELETE CASCADE
);
