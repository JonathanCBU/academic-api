-- Academic Data Collection Tool - Database Schema
-- SQLite3 Compatible Version
-- Created: 2024
-- ============================================================================
-- DROP EXISTING TABLES (for clean reinstall)
-- ============================================================================
DROP TABLE IF EXISTS raw_data;
DROP TABLE IF EXISTS school_report;
DROP TABLE IF EXISTS school;
-- ============================================================================
-- SCHOOLS TABLE
-- ============================================================================
CREATE TABLE school (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    school_name TEXT NOT NULL,
    state_code TEXT NOT NULL,
    district_name TEXT NOT NULL,
    is_deleted BOOLEAN,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME,
    deleted_at DATETIME
);
-- ============================================================================
-- RAW DATA TABLE
-- ============================================================================
CREATE TABLE raw_data (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    scope TEXT NOT NULL CHECK (
        scope IN (
            'school',
            'district',
            'state'
        )
    ),
    source TEXT NOT NULL,
    structure TEXT NOT NULL CHECK (
        structure in (
            'json',
            'csv',
            'xml',
            'html',
            'xlsx',
            'xls'
        )
    ),
    actual TEXT NOT NULL,
    is_deleted BOOLEAN,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME,
    deleted_at DATETIME
);
-- ============================================================================
-- SINGLE SCHOOL DATA TABLE (Main Fact Table)
-- ============================================================================
CREATE TABLE school_report (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    school_id INTEGER NOT NULL,
    data_id INTEGER NOT NULL,
    academic_year INTEGER NOT NULL,
    subject TEXT NOT NULL CHECK (subject IN ('ela', 'math')),
    grade_level TEXT NOT NULL,
    demographic_group TEXT NOT NULL CHECK (
        demographic_group IN (
            'all',
            'black',
            'hispanic',
            'economically_disadvantaged'
        )
    ),
    n_tested INTEGER CHECK (n_tested >= 0),
    n_proficient INTEGER CHECK (n_proficient >= 0),
    pct_proficient REAL CHECK (
        pct_proficient >= 0
        AND pct_proficient <= 100
    ),
    is_deleted BOOLEAN,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME,
    deleted_at DATETIME,
    FOREIGN KEY (school_id) REFERENCES schools(id) ON DELETE CASCADE,
    FOREIGN KEY (data_id) REFERENCES raw_data(id) ON DELETE CASCADE,
    -- Ensure we don't have duplicate records for same school/year/subject/grade/demographic
    UNIQUE (
        school_id,
        academic_year,
        subject,
        grade_level,
        demographic_group
    ),
    -- Ensure n_proficient doesn't exceed n_tested
    CHECK (
        n_proficient IS NULL
        OR n_tested IS NULL
        OR n_proficient <= n_tested
    )
);