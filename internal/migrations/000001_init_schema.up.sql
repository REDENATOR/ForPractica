-- Создаём таблицу студентов (соответствует модели)
CREATE TABLE IF NOT EXISTS students (
    -- gorm.Model поля
    id SERIAL PRIMARY KEY,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,

    fio TEXT NOT NULL,
    group_of_students TEXT,
    phone_number TEXT
);

-- Индексы для частых запросов
CREATE INDEX idx_students_group_of_students ON students(group_of_students);
CREATE INDEX idx_students_fio ON students(fio);
CREATE INDEX idx_students_deleted_at ON students(deleted_at);

-- Триггер для обновления updated_at
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER update_students_updated_at
    BEFORE UPDATE ON students
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();