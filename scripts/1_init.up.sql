CREATE TABLE IF NOT EXISTS reserved
(
    reserv_id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    day DATE NOT NULL,
    fk_sup_id INT NOT NULL
);

CREATE TABLE IF NOT EXISTS sups (
    sup_id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    model_name VARCHAR(100) NOT NULL,
    price INT NOT NULL DEFAULT 1000
);

CREATE TABLE IF NOT EXISTS approve (
    approve_id INT GENERATED ALWAYS AS IDENTITY,
    client_phone	VARCHAR(20) NOT NULL,
    client_name	VARCHAR(100) NOT NULL,
    price INT NOT NULL,
    order_info TEXT NOT NULL,
    status SMALLINT DEFAULT 1
);

ALTER TABLE approve
ADD CONSTRAINT pk_approve_id_phone PRIMARY KEY (approve_id, client_phone);

ALTER TABLE reserved 
ADD CONSTRAINT fk_reserved_sup_id FOREIGN KEY (fk_sup_id) REFERENCES sups (sup_id) ON DELETE CASCADE;

CREATE INDEX idx_reserved_day ON reserved(day);

CREATE INDEX idx_approve_status ON approve(status);

INSERT INTO sups (model_name, price) VALUES
('GLADIATOR', 1000),
('BOMBITTO', 1000),
('PANORAMA', 1000);
