INSERT INTO status (id, name) VALUES
   (0, 'Pending'),
   (1, 'In Progress'),
   (2, 'Completed'),
   (3, 'Cancelled');

INSERT INTO task (title, body, status_id) VALUES 
   ('Complete project documentation', 'Write comprehensive documentation covering all API endpoints, authentication methods, and usage examples. Include code samples in multiple languages and ensure all examples are tested and working. The documentation should be organized into clear sections: Getting Started, Authentication, Core API Reference, Advanced Features, and Troubleshooting. Each endpoint should include request/response examples, parameter descriptions, and common error codes. Add interactive examples where possible and ensure the documentation is accessible to both beginner and advanced developers.', 2),
   ('Review code changes', NULL, 1),
   ('Fix bug in authentication', NULL, 0),
   ('Deploy to production', NULL, 2),
   ('Write unit tests', NULL, 0),
   ('Update dependencies', NULL, 3);
