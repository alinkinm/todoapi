\c todoapi;

CREATE TABLE Task (
    id SERIAL PRIMARY KEY,
    header VARCHAR NOT NULL,
    descr VARCHAR NOT NULL,
    task_date DATE NOT NULL,
    done boolean NOT NULL
);

INSERT INTO Task (header, descr, task_date, done) VALUES 
('task struct', 'create struct with fields...', '2024-05-31', false),
('crudl', 'create update delete list...', '2024-05-31', false),
('pagination', 'implement pagination on status...', '2024-05-31', false),
('date filter', 'implement tasks date filter...', '2024-06-01', false),
('swagger', 'create api documentaion...', '2024-06-01', false),
('docker', 'dockerize api...', '2024-06-01', false),
('unit-test', 'cover api with unit-tests...', '2024-06-01', false);