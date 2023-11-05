CREATE TABLE templates (
	id          serial  primary key,
	updatedAt   date    not null,
	createdAt   date    not null,
	name        text    default ''
	
);

INSERT INTO templates (updatedAt, createdAt, name) VALUES 
('2023-11-03', '2023-11-03', 'Vorlage 1'),
('2023-11-02', '2023-11-02', 'Vorlage 2');