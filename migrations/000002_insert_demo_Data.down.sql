-- Migration Rollback: Remove sample schools and proficiency data
-- Created: 2024-11-16
-- =============================================================================
-- Delete in reverse order to respect foreign key constraints
-- =============================================================================
-- Delete collection logs for sample schools
DELETE FROM data_collection_log
WHERE collected_by = 'seed_migration';
-- Delete state benchmarks for CA and TX (2024 data only)
DELETE FROM district_state_benchmarks
WHERE state_code IN ('CA', 'TX')
    AND academic_year = 2024;
-- Delete proficiency data for sample schools
DELETE FROM proficiency_data
WHERE school_id IN (
        SELECT school_id
        FROM schools
        WHERE nces_id IN (
                '062961004587',
                '062961004588',
                '062961004589',
                -- CA schools
                '484329004587',
                '484329004588',
                '484329004589' -- TX schools
            )
    );
-- Delete sample schools
DELETE FROM schools
WHERE nces_id IN (
        '062961004587',
        -- Lincoln Elementary School (CA)
        '062961004588',
        -- Washington Middle School (CA)
        '062961004589',
        -- Roosevelt High School (CA)
        '484329004587',
        -- Jefferson Elementary (TX)
        '484329004588',
        -- Kennedy Middle School (TX)
        '484329004589' -- Adams High School (TX)
    );