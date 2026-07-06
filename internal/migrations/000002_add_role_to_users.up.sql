ALTER TABLE students ADD COLUMN role VARCHAR(20) NOT NULL DEFAULT 'student';
CREATE INDEX idx_students_role ON students(role);
