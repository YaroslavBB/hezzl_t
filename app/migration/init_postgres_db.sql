\connect postgres postgres

CREATE TABLE projects (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE goods (
    id SERIAL PRIMARY KEY,
    project_id INTEGER REFERENCES projects (id),
    name VARCHAR(255) NOT NULL,
    description VARCHAR(255),
    priority INTEGER DEFAULT 1,
    removed BOOL NOT NULL DEFAULT FALSE,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX good_id_idx ON goods (id);
CREATE INDEX good_project_id_idx ON goods (project_id);
CREATE INDEX good_name_idx ON goods (name);

CREATE INDEX project_id_index ON projects (id);

INSERT INTO projects (id, name) VALUES (1, 'Первая запись');

CREATE FUNCTION set_priority_on_insert() RETURNS trigger
AS $$
BEGIN
  IF NEW.priority IS NULL THEN
    NEW.priority = 1;
  ELSE
    NEW.priority = (SELECT MAX(priority) + 1 FROM goods);
  END IF;
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER set_priority_on_insert
BEFORE INSERT ON goods
FOR EACH ROW
EXECUTE PROCEDURE set_priority_on_insert();