INSERT INTO status (id, name) VALUES
   (0, 'Pending'),
   (1, 'In Progress'),
   (2, 'Completed'),
   (3, 'Cancelled');

INSERT INTO tag (id, name) VALUES 
   (1, 'coding'),
   (2, 'docs');

INSERT INTO task (title, body, status_id, tag_id) VALUES 
   ('Complete project documentation', 'Write comprehensive documentation covering all API endpoints, authentication methods, and usage examples. Include code samples in multiple languages and ensure all examples are tested and working. The documentation should be organized into clear sections: Getting Started, Authentication, Core API Reference, Advanced Features, and Troubleshooting. Each endpoint should include request/response examples, parameter descriptions, and common error codes. Add interactive examples where possible and ensure the documentation is accessible to both beginner and advanced developers.', 2, 2),
   ('Review code changes', NULL, 1, NULL),
   ('Fix bug in authentication', NULL, 0, 1),
   ('Deploy to production', NULL, 2, NULL),
   ('Write unit tests', NULL, 0, 1),
   ('Update dependencies', NULL, 3, 1);

INSERT INTO note (title) VALUES 
   ('Some note');

