CREATE TABLE status (
   id INTEGER PRIMARY KEY,
   name TEXT NOT NULL UNIQUE
);

CREATE TABLE task (
   id INTEGER PRIMARY KEY,
   title TEXT NOT NULL,
   body TEXT,
   status_id INTEGER NOT NULL,
   FOREIGN KEY (status_id) REFERENCES status(id)
);
