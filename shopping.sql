CREATE TABLE products(
                         id SERIAL primary key ,
                         name varchar(50) not null ,
                         price numeric(8,2) not null ,
                         tax   numeric(8,2) not null
);

CREATE TABLE cart(
                     id SERIAL primary key ,
                     product_id int references products(id) ,
                     quantity  int not null,
                     total_price numeric(8,2) not null ,
                     total_discount  numeric(8,2) not null
);

CREATE TABLE customer(
                         id SERIAL primary key ,
                         name varchar(50) not null,
                         surname varchar(50) not null
);

CREATE TABLE payment(
                        id SERIAL primary key ,
                        customer_id int references customer(id) ,
                        total_price numeric(8,2) not null
);