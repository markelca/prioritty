CREATE TABLE status (
   id INTEGER PRIMARY KEY,
   name TEXT NOT NULL UNIQUE
);

CREATE TABLE task (
   id INTEGER PRIMARY KEY,
   title TEXT NOT NULL,
   status_id INTEGER NOT NULL,
   FOREIGN KEY (status_id) REFERENCES status(id)
);

INSERT INTO status (id, name) VALUES
   (0, 'Pending'),
   (1, 'In Progress'),
   (2, 'Completed'),
   (3, 'Cancelled');

INSERT INTO task (title, status_id) VALUES 
   ('Complete project documentation', 0),
   ('Review code changes', 1),
   ('Fix bug in authentication', 0),
   ('Deploy to production', 2),
   ('Write unit tests', 0),
   ('Update dependencies', 3);
