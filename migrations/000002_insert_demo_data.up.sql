-- Migration: Add sample schools and proficiency data
-- Created: 2024-11-16
-- Purpose: Seed database with realistic example data for testing and development
-- =============================================================================
-- Insert Sample Schools
-- =============================================================================
-- California Schools
INSERT INTO schools (
        school_name,
        state_code,
        district_name,
        nces_id,
        address
    )
VALUES (
        'Lincoln Elementary School',
        'CA',
        'Los Angeles Unified',
        '062961004587',
        '1234 Main St, Los Angeles, CA 90012'
    ),
    (
        'Washington Middle School',
        'CA',
        'Los Angeles Unified',
        '062961004588',
        '5678 Oak Ave, Los Angeles, CA 90013'
    ),
    (
        'Roosevelt High School',
        'CA',
        'San Diego Unified',
        '062961004589',
        '9012 Pine Blvd, San Diego, CA 92101'
    );
-- Texas Schools
INSERT INTO schools (
        school_name,
        state_code,
        district_name,
        nces_id,
        address
    )
VALUES (
        'Jefferson Elementary',
        'TX',
        'Houston ISD',
        '484329004587',
        '2345 Elm St, Houston, TX 77001'
    ),
    (
        'Kennedy Middle School',
        'TX',
        'Houston ISD',
        '484329004588',
        '6789 Maple Dr, Houston, TX 77002'
    ),
    (
        'Adams High School',
        'TX',
        'Dallas ISD',
        '484329004589',
        '3456 Cedar Ln, Dallas, TX 75201'
    );
-- =============================================================================
-- Insert Proficiency Data for California Schools
-- =============================================================================
-- Lincoln Elementary School (CA) - 2024 Data
-- All Students
INSERT INTO proficiency_data (
        school_id,
        academic_year,
        subject,
        grade_level,
        demographic_group,
        n_tested,
        n_proficient,
        pct_proficient,
        data_source,
        collection_method,
        collected_by
    )
VALUES -- ELA - All Students
    (
        (
            SELECT school_id
            FROM schools
            WHERE nces_id = '062961004587'
        ),
        2024,
        'ela',
        '3-8',
        'all',
        450,
        315,
        70.00,
        'https://caaspp-elpac.ets.org/',
        'manual',
        'seed_migration'
    ),
    -- Math - All Students
    (
        (
            SELECT school_id
            FROM schools
            WHERE nces_id = '062961004587'
        ),
        2024,
        'math',
        '3-8',
        'all',
        450,
        293,
        65.11,
        'https://caaspp-elpac.ets.org/',
        'manual',
        'seed_migration'
    ),
    -- ELA - Black Students
    (
        (
            SELECT school_id
            FROM schools
            WHERE nces_id = '062961004587'
        ),
        2024,
        'ela',
        '3-8',
        'black',
        85,
        51,
        60.00,
        'https://caaspp-elpac.ets.org/',
        'manual',
        'seed_migration'
    ),
    -- Math - Black Students
    (
        (
            SELECT school_id
            FROM schools
            WHERE nces_id = '062961004587'
        ),
        2024,
        'math',
        '3-8',
        'black',
        85,
        47,
        55.29,
        'https://caaspp-elpac.ets.org/',
        'manual',
        'seed_migration'
    ),
    -- ELA - Hispanic Students
    (
        (
            SELECT school_id
            FROM schools
            WHERE nces_id = '062961004587'
        ),
        2024,
        'ela',
        '3-8',
        'hispanic',
        280,
        182,
        65.00,
        'https://caaspp-elpac.ets.org/',
        'manual',
        'seed_migration'
    ),
    -- Math - Hispanic Students
    (
        (
            SELECT school_id
            FROM schools
            WHERE nces_id = '062961004587'
        ),
        2024,
        'math',
        '3-8',
        'hispanic',
        280,
        168,
        60.00,
        'https://caaspp-elpac.ets.org/',
        'manual',
        'seed_migration'
    ),
    -- ELA - FRL Students
    (
        (
            SELECT school_id
            FROM schools
            WHERE nces_id = '062961004587'
        ),
        2024,
        'ela',
        '3-8',
        'frl',
        320,
        208,
        65.00,
        'https://caaspp-elpac.ets.org/',
        'manual',
        'seed_migration'
    ),
    -- Math - FRL Students
    (
        (
            SELECT school_id
            FROM schools
            WHERE nces_id = '062961004587'
        ),
        2024,
        'math',
        '3-8',
        'frl',
        320,
        192,
        60.00,
        'https://caaspp-elpac.ets.org/',
        'manual',
        'seed_migration'
    );
-- Washington Middle School (CA) - 2024 Data
INSERT INTO proficiency_data (
        school_id,
        academic_year,
        subject,
        grade_level,
        demographic_group,
        n_tested,
        n_proficient,
        pct_proficient,
        data_source,
        collection_method,
        collected_by
    )
VALUES -- All Students
    (
        (
            SELECT school_id
            FROM schools
            WHERE nces_id = '062961004588'
        ),
        2024,
        'ela',
        '3-8',
        'all',
        520,
        364,
        70.00,
        'https://caaspp-elpac.ets.org/',
        'manual',
        'seed_migration'
    ),
    (
        (
            SELECT school_id
            FROM schools
            WHERE nces_id = '062961004588'
        ),
        2024,
        'math',
        '3-8',
        'all',
        520,
        327,
        62.88,
        'https://caaspp-elpac.ets.org/',
        'manual',
        'seed_migration'
    ),
    -- Black Students
    (
        (
            SELECT school_id
            FROM schools
            WHERE nces_id = '062961004588'
        ),
        2024,
        'ela',
        '3-8',
        'black',
        95,
        62,
        65.26,
        'https://caaspp-elpac.ets.org/',
        'manual',
        'seed_migration'
    ),
    (
        (
            SELECT school_id
            FROM schools
            WHERE nces_id = '062961004588'
        ),
        2024,
        'math',
        '3-8',
        'black',
        95,
        52,
        54.74,
        'https://caaspp-elpac.ets.org/',
        'manual',
        'seed_migration'
    ),
    -- Hispanic Students
    (
        (
            SELECT school_id
            FROM schools
            WHERE nces_id = '062961004588'
        ),
        2024,
        'ela',
        '3-8',
        'hispanic',
        310,
        201,
        64.84,
        'https://caaspp-elpac.ets.org/',
        'manual',
        'seed_migration'
    ),
    (
        (
            SELECT school_id
            FROM schools
            WHERE nces_id = '062961004588'
        ),
        2024,
        'math',
        '3-8',
        'hispanic',
        310,
        186,
        60.00,
        'https://caaspp-elpac.ets.org/',
        'manual',
        'seed_migration'
    ),
    -- FRL Students
    (
        (
            SELECT school_id
            FROM schools
            WHERE nces_id = '062961004588'
        ),
        2024,
        'ela',
        '3-8',
        'frl',
        380,
        247,
        65.00,
        'https://caaspp-elpac.ets.org/',
        'manual',
        'seed_migration'
    ),
    (
        (
            SELECT school_id
            FROM schools
            WHERE nces_id = '062961004588'
        ),
        2024,
        'math',
        '3-8',
        'frl',
        380,
        228,
        60.00,
        'https://caaspp-elpac.ets.org/',
        'manual',
        'seed_migration'
    );
-- Roosevelt High School (CA) - 2024 Data
INSERT INTO proficiency_data (
        school_id,
        academic_year,
        subject,
        grade_level,
        demographic_group,
        n_tested,
        n_proficient,
        pct_proficient,
        data_source,
        collection_method,
        collected_by
    )
VALUES -- All Students
    (
        (
            SELECT school_id
            FROM schools
            WHERE nces_id = '062961004589'
        ),
        2024,
        'ela',
        '3-8',
        'all',
        380,
        304,
        80.00,
        'https://caaspp-elpac.ets.org/',
        'manual',
        'seed_migration'
    ),
    (
        (
            SELECT school_id
            FROM schools
            WHERE nces_id = '062961004589'
        ),
        2024,
        'math',
        '3-8',
        'all',
        380,
        285,
        75.00,
        'https://caaspp-elpac.ets.org/',
        'manual',
        'seed_migration'
    ),
    -- Black Students
    (
        (
            SELECT school_id
            FROM schools
            WHERE nces_id = '062961004589'
        ),
        2024,
        'ela',
        '3-8',
        'black',
        65,
        49,
        75.38,
        'https://caaspp-elpac.ets.org/',
        'manual',
        'seed_migration'
    ),
    (
        (
            SELECT school_id
            FROM schools
            WHERE nces_id = '062961004589'
        ),
        2024,
        'math',
        '3-8',
        'black',
        65,
        45,
        69.23,
        'https://caaspp-elpac.ets.org/',
        'manual',
        'seed_migration'
    ),
    -- Hispanic Students
    (
        (
            SELECT school_id
            FROM schools
            WHERE nces_id = '062961004589'
        ),
        2024,
        'ela',
        '3-8',
        'hispanic',
        180,
        135,
        75.00,
        'https://caaspp-elpac.ets.org/',
        'manual',
        'seed_migration'
    ),
    (
        (
            SELECT school_id
            FROM schools
            WHERE nces_id = '062961004589'
        ),
        2024,
        'math',
        '3-8',
        'hispanic',
        180,
        126,
        70.00,
        'https://caaspp-elpac.ets.org/',
        'manual',
        'seed_migration'
    ),
    -- FRL Students
    (
        (
            SELECT school_id
            FROM schools
            WHERE nces_id = '062961004589'
        ),
        2024,
        'ela',
        '3-8',
        'frl',
        220,
        165,
        75.00,
        'https://caaspp-elpac.ets.org/',
        'manual',
        'seed_migration'
    ),
    (
        (
            SELECT school_id
            FROM schools
            WHERE nces_id = '062961004589'
        ),
        2024,
        'math',
        '3-8',
        'frl',
        220,
        154,
        70.00,
        'https://caaspp-elpac.ets.org/',
        'manual',
        'seed_migration'
    );
-- =============================================================================
-- Insert Proficiency Data for Texas Schools
-- =============================================================================
-- Jefferson Elementary (TX) - 2024 Data
INSERT INTO proficiency_data (
        school_id,
        academic_year,
        subject,
        grade_level,
        demographic_group,
        n_tested,
        n_proficient,
        pct_proficient,
        data_source,
        collection_method,
        collected_by
    )
VALUES -- All Students
    (
        (
            SELECT school_id
            FROM schools
            WHERE nces_id = '484329004587'
        ),
        2024,
        'ela',
        '3-8',
        'all',
        410,
        287,
        70.00,
        'https://tea.texas.gov/',
        'manual',
        'seed_migration'
    ),
    (
        (
            SELECT school_id
            FROM schools
            WHERE nces_id = '484329004587'
        ),
        2024,
        'math',
        '3-8',
        'all',
        410,
        266,
        64.88,
        'https://tea.texas.gov/',
        'manual',
        'seed_migration'
    ),
    -- Black Students
    (
        (
            SELECT school_id
            FROM schools
            WHERE nces_id = '484329004587'
        ),
        2024,
        'ela',
        '3-8',
        'black',
        125,
        75,
        60.00,
        'https://tea.texas.gov/',
        'manual',
        'seed_migration'
    ),
    (
        (
            SELECT school_id
            FROM schools
            WHERE nces_id = '484329004587'
        ),
        2024,
        'math',
        '3-8',
        'black',
        125,
        69,
        55.20,
        'https://tea.texas.gov/',
        'manual',
        'seed_migration'
    ),
    -- Hispanic Students
    (
        (
            SELECT school_id
            FROM schools
            WHERE nces_id = '484329004587'
        ),
        2024,
        'ela',
        '3-8',
        'hispanic',
        210,
        136,
        64.76,
        'https://tea.texas.gov/',
        'manual',
        'seed_migration'
    ),
    (
        (
            SELECT school_id
            FROM schools
            WHERE nces_id = '484329004587'
        ),
        2024,
        'math',
        '3-8',
        'hispanic',
        210,
        126,
        60.00,
        'https://tea.texas.gov/',
        'manual',
        'seed_migration'
    ),
    -- FRL Students
    (
        (
            SELECT school_id
            FROM schools
            WHERE nces_id = '484329004587'
        ),
        2024,
        'ela',
        '3-8',
        'frl',
        295,
        192,
        65.08,
        'https://tea.texas.gov/',
        'manual',
        'seed_migration'
    ),
    (
        (
            SELECT school_id
            FROM schools
            WHERE nces_id = '484329004587'
        ),
        2024,
        'math',
        '3-8',
        'frl',
        295,
        177,
        60.00,
        'https://tea.texas.gov/',
        'manual',
        'seed_migration'
    );
-- Kennedy Middle School (TX) - 2024 Data
INSERT INTO proficiency_data (
        school_id,
        academic_year,
        subject,
        grade_level,
        demographic_group,
        n_tested,
        n_proficient,
        pct_proficient,
        data_source,
        collection_method,
        collected_by
    )
VALUES -- All Students
    (
        (
            SELECT school_id
            FROM schools
            WHERE nces_id = '484329004588'
        ),
        2024,
        'ela',
        '3-8',
        'all',
        485,
        340,
        70.10,
        'https://tea.texas.gov/',
        'manual',
        'seed_migration'
    ),
    (
        (
            SELECT school_id
            FROM schools
            WHERE nces_id = '484329004588'
        ),
        2024,
        'math',
        '3-8',
        'all',
        485,
        306,
        63.09,
        'https://tea.texas.gov/',
        'manual',
        'seed_migration'
    ),
    -- Black Students
    (
        (
            SELECT school_id
            FROM schools
            WHERE nces_id = '484329004588'
        ),
        2024,
        'ela',
        '3-8',
        'black',
        110,
        71,
        64.55,
        'https://tea.texas.gov/',
        'manual',
        'seed_migration'
    ),
    (
        (
            SELECT school_id
            FROM schools
            WHERE nces_id = '484329004588'
        ),
        2024,
        'math',
        '3-8',
        'black',
        110,
        60,
        54.55,
        'https://tea.texas.gov/',
        'manual',
        'seed_migration'
    ),
    -- Hispanic Students
    (
        (
            SELECT school_id
            FROM schools
            WHERE nces_id = '484329004588'
        ),
        2024,
        'ela',
        '3-8',
        'hispanic',
        285,
        185,
        64.91,
        'https://tea.texas.gov/',
        'manual',
        'seed_migration'
    ),
    (
        (
            SELECT school_id
            FROM schools
            WHERE nces_id = '484329004588'
        ),
        2024,
        'math',
        '3-8',
        'hispanic',
        285,
        171,
        60.00,
        'https://tea.texas.gov/',
        'manual',
        'seed_migration'
    ),
    -- FRL Students
    (
        (
            SELECT school_id
            FROM schools
            WHERE nces_id = '484329004588'
        ),
        2024,
        'ela',
        '3-8',
        'frl',
        350,
        228,
        65.14,
        'https://tea.texas.gov/',
        'manual',
        'seed_migration'
    ),
    (
        (
            SELECT school_id
            FROM schools
            WHERE nces_id = '484329004588'
        ),
        2024,
        'math',
        '3-8',
        'frl',
        350,
        210,
        60.00,
        'https://tea.texas.gov/',
        'manual',
        'seed_migration'
    );
-- Adams High School (TX) - 2024 Data
INSERT INTO proficiency_data (
        school_id,
        academic_year,
        subject,
        grade_level,
        demographic_group,
        n_tested,
        n_proficient,
        pct_proficient,
        data_source,
        collection_method,
        collected_by
    )
VALUES -- All Students
    (
        (
            SELECT school_id
            FROM schools
            WHERE nces_id = '484329004589'
        ),
        2024,
        'ela',
        '3-8',
        'all',
        420,
        336,
        80.00,
        'https://tea.texas.gov/',
        'manual',
        'seed_migration'
    ),
    (
        (
            SELECT school_id
            FROM schools
            WHERE nces_id = '484329004589'
        ),
        2024,
        'math',
        '3-8',
        'all',
        420,
        315,
        75.00,
        'https://tea.texas.gov/',
        'manual',
        'seed_migration'
    ),
    -- Black Students
    (
        (
            SELECT school_id
            FROM schools
            WHERE nces_id = '484329004589'
        ),
        2024,
        'ela',
        '3-8',
        'black',
        75,
        56,
        74.67,
        'https://tea.texas.gov/',
        'manual',
        'seed_migration'
    ),
    (
        (
            SELECT school_id
            FROM schools
            WHERE nces_id = '484329004589'
        ),
        2024,
        'math',
        '3-8',
        'black',
        75,
        52,
        69.33,
        'https://tea.texas.gov/',
        'manual',
        'seed_migration'
    ),
    -- Hispanic Students
    (
        (
            SELECT school_id
            FROM schools
            WHERE nces_id = '484329004589'
        ),
        2024,
        'ela',
        '3-8',
        'hispanic',
        195,
        146,
        74.87,
        'https://tea.texas.gov/',
        'manual',
        'seed_migration'
    ),
    (
        (
            SELECT school_id
            FROM schools
            WHERE nces_id = '484329004589'
        ),
        2024,
        'math',
        '3-8',
        'hispanic',
        195,
        137,
        70.26,
        'https://tea.texas.gov/',
        'manual',
        'seed_migration'
    ),
    -- FRL Students
    (
        (
            SELECT school_id
            FROM schools
            WHERE nces_id = '484329004589'
        ),
        2024,
        'ela',
        '3-8',
        'frl',
        245,
        184,
        75.10,
        'https://tea.texas.gov/',
        'manual',
        'seed_migration'
    ),
    (
        (
            SELECT school_id
            FROM schools
            WHERE nces_id = '484329004589'
        ),
        2024,
        'math',
        '3-8',
        'frl',
        245,
        172,
        70.20,
        'https://tea.texas.gov/',
        'manual',
        'seed_migration'
    );
-- =============================================================================
-- Insert District Benchmarks
-- =============================================================================
-- Los Angeles Unified District (CA) - 2024
INSERT INTO district_state_benchmarks (
        state_code,
        district_name,
        academic_year,
        subject,
        grade_level,
        demographic_group,
        pct_proficient,
        data_source
    )
VALUES (
        'CA',
        'Los Angeles Unified',
        2024,
        'ela',
        '3-8',
        'all',
        65.50,
        'https://caaspp-elpac.ets.org/'
    ),
    (
        'CA',
        'Los Angeles Unified',
        2024,
        'math',
        '3-8',
        'all',
        58.70,
        'https://caaspp-elpac.ets.org/'
    ),
    (
        'CA',
        'Los Angeles Unified',
        2024,
        'ela',
        '3-8',
        'black',
        58.20,
        'https://caaspp-elpac.ets.org/'
    ),
    (
        'CA',
        'Los Angeles Unified',
        2024,
        'math',
        '3-8',
        'black',
        51.30,
        'https://caaspp-elpac.ets.org/'
    ),
    (
        'CA',
        'Los Angeles Unified',
        2024,
        'ela',
        '3-8',
        'hispanic',
        61.40,
        'https://caaspp-elpac.ets.org/'
    ),
    (
        'CA',
        'Los Angeles Unified',
        2024,
        'math',
        '3-8',
        'hispanic',
        55.20,
        'https://caaspp-elpac.ets.org/'
    ),
    (
        'CA',
        'Los Angeles Unified',
        2024,
        'ela',
        '3-8',
        'frl',
        60.10,
        'https://caaspp-elpac.ets.org/'
    ),
    (
        'CA',
        'Los Angeles Unified',
        2024,
        'math',
        '3-8',
        'frl',
        53.80,
        'https://caaspp-elpac.ets.org/'
    );
-- San Diego Unified District (CA) - 2024
INSERT INTO district_state_benchmarks (
        state_code,
        district_name,
        academic_year,
        subject,
        grade_level,
        demographic_group,
        pct_proficient,
        data_source
    )
VALUES (
        'CA',
        'San Diego Unified',
        2024,
        'ela',
        '3-8',
        'all',
        72.30,
        'https://caaspp-elpac.ets.org/'
    ),
    (
        'CA',
        'San Diego Unified',
        2024,
        'math',
        '3-8',
        'all',
        68.50,
        'https://caaspp-elpac.ets.org/'
    ),
    (
        'CA',
        'San Diego Unified',
        2024,
        'ela',
        '3-8',
        'black',
        68.70,
        'https://caaspp-elpac.ets.org/'
    ),
    (
        'CA',
        'San Diego Unified',
        2024,
        'math',
        '3-8',
        'black',
        63.20,
        'https://caaspp-elpac.ets.org/'
    ),
    (
        'CA',
        'San Diego Unified',
        2024,
        'ela',
        '3-8',
        'hispanic',
        69.50,
        'https://caaspp-elpac.ets.org/'
    ),
    (
        'CA',
        'San Diego Unified',
        2024,
        'math',
        '3-8',
        'hispanic',
        64.80,
        'https://caaspp-elpac.ets.org/'
    ),
    (
        'CA',
        'San Diego Unified',
        2024,
        'ela',
        '3-8',
        'frl',
        68.20,
        'https://caaspp-elpac.ets.org/'
    ),
    (
        'CA',
        'San Diego Unified',
        2024,
        'math',
        '3-8',
        'frl',
        63.10,
        'https://caaspp-elpac.ets.org/'
    );
-- Houston ISD (TX) - 2024
INSERT INTO district_state_benchmarks (
        state_code,
        district_name,
        academic_year,
        subject,
        grade_level,
        demographic_group,
        pct_proficient,
        data_source
    )
VALUES (
        'TX',
        'Houston ISD',
        2024,
        'ela',
        '3-8',
        'all',
        66.80,
        'https://tea.texas.gov/'
    ),
    (
        'TX',
        'Houston ISD',
        2024,
        'math',
        '3-8',
        'all',
        60.40,
        'https://tea.texas.gov/'
    ),
    (
        'TX',
        'Houston ISD',
        2024,
        'ela',
        '3-8',
        'black',
        60.20,
        'https://tea.texas.gov/'
    ),
    (
        'TX',
        'Houston ISD',
        2024,
        'math',
        '3-8',
        'black',
        52.80,
        'https://tea.texas.gov/'
    ),
    (
        'TX',
        'Houston ISD',
        2024,
        'ela',
        '3-8',
        'hispanic',
        62.50,
        'https://tea.texas.gov/'
    ),
    (
        'TX',
        'Houston ISD',
        2024,
        'math',
        '3-8',
        'hispanic',
        56.90,
        'https://tea.texas.gov/'
    ),
    (
        'TX',
        'Houston ISD',
        2024,
        'ela',
        '3-8',
        'frl',
        61.30,
        'https://tea.texas.gov/'
    ),
    (
        'TX',
        'Houston ISD',
        2024,
        'math',
        '3-8',
        'frl',
        55.70,
        'https://tea.texas.gov/'
    );
-- Dallas ISD (TX) - 2024
INSERT INTO district_state_benchmarks (
        state_code,
        district_name,
        academic_year,
        subject,
        grade_level,
        demographic_group,
        pct_proficient,
        data_source
    )
VALUES (
        'TX',
        'Dallas ISD',
        2024,
        'ela',
        '3-8',
        'all',
        73.20,
        'https://tea.texas.gov/'
    ),
    (
        'TX',
        'Dallas ISD',
        2024,
        'math',
        '3-8',
        'all',
        69.10,
        'https://tea.texas.gov/'
    ),
    (
        'TX',
        'Dallas ISD',
        2024,
        'ela',
        '3-8',
        'black',
        69.50,
        'https://tea.texas.gov/'
    ),
    (
        'TX',
        'Dallas ISD',
        2024,
        'math',
        '3-8',
        'black',
        64.20,
        'https://tea.texas.gov/'
    ),
    (
        'TX',
        'Dallas ISD',
        2024,
        'ela',
        '3-8',
        'hispanic',
        70.80,
        'https://tea.texas.gov/'
    ),
    (
        'TX',
        'Dallas ISD',
        2024,
        'math',
        '3-8',
        'hispanic',
        66.40,
        'https://tea.texas.gov/'
    ),
    (
        'TX',
        'Dallas ISD',
        2024,
        'ela',
        '3-8',
        'frl',
        69.90,
        'https://tea.texas.gov/'
    ),
    (
        'TX',
        'Dallas ISD',
        2024,
        'math',
        '3-8',
        'frl',
        65.30,
        'https://tea.texas.gov/'
    );
-- =============================================================================
-- Insert State Benchmarks
-- =============================================================================
-- California State - 2024
INSERT INTO district_state_benchmarks (
        state_code,
        district_name,
        academic_year,
        subject,
        grade_level,
        demographic_group,
        pct_proficient,
        data_source
    )
VALUES (
        'CA',
        NULL,
        2024,
        'ela',
        '3-8',
        'all',
        60.20,
        'https://caaspp-elpac.ets.org/'
    ),
    (
        'CA',
        NULL,
        2024,
        'math',
        '3-8',
        'all',
        55.10,
        'https://caaspp-elpac.ets.org/'
    ),
    (
        'CA',
        NULL,
        2024,
        'ela',
        '3-8',
        'black',
        52.40,
        'https://caaspp-elpac.ets.org/'
    ),
    (
        'CA',
        NULL,
        2024,
        'math',
        '3-8',
        'black',
        46.80,
        'https://caaspp-elpac.ets.org/'
    ),
    (
        'CA',
        NULL,
        2024,
        'ela',
        '3-8',
        'hispanic',
        55.70,
        'https://caaspp-elpac.ets.org/'
    ),
    (
        'CA',
        NULL,
        2024,
        'math',
        '3-8',
        'hispanic',
        50.30,
        'https://caaspp-elpac.ets.org/'
    ),
    (
        'CA',
        NULL,
        2024,
        'ela',
        '3-8',
        'frl',
        54.20,
        'https://caaspp-elpac.ets.org/'
    ),
    (
        'CA',
        NULL,
        2024,
        'math',
        '3-8',
        'frl',
        48.90,
        'https://caaspp-elpac.ets.org/'
    );
-- Texas State - 2024
INSERT INTO district_state_benchmarks (
        state_code,
        district_name,
        academic_year,
        subject,
        grade_level,
        demographic_group,
        pct_proficient,
        data_source
    )
VALUES (
        'TX',
        NULL,
        2024,
        'ela',
        '3-8',
        'all',
        62.50,
        'https://tea.texas.gov/'
    ),
    (
        'TX',
        NULL,
        2024,
        'math',
        '3-8',
        'all',
        58.20,
        'https://tea.texas.gov/'
    ),
    (
        'TX',
        NULL,
        2024,
        'ela',
        '3-8',
        'black',
        55.80,
        'https://tea.texas.gov/'
    ),
    (
        'TX',
        NULL,
        2024,
        'math',
        '3-8',
        'black',
        49.70,
        'https://tea.texas.gov/'
    ),
    (
        'TX',
        NULL,
        2024,
        'ela',
        '3-8',
        'hispanic',
        58.40,
        'https://tea.texas.gov/'
    ),
    (
        'TX',
        NULL,
        2024,
        'math',
        '3-8',
        'hispanic',
        53.60,
        'https://tea.texas.gov/'
    ),
    (
        'TX',
        NULL,
        2024,
        'ela',
        '3-8',
        'frl',
        57.10,
        'https://tea.texas.gov/'
    ),
    (
        'TX',
        NULL,
        2024,
        'math',
        '3-8',
        'frl',
        52.40,
        'https://tea.texas.gov/'
    );
-- =============================================================================
-- Insert Collection Logs
-- =============================================================================
INSERT INTO data_collection_log (
        school_id,
        state_code,
        status,
        records_collected,
        method_used,
        time_spent_minutes,
        collected_by
    )
VALUES (
        (
            SELECT school_id
            FROM schools
            WHERE nces_id = '062961004587'
        ),
        'CA',
        'success',
        8,
        'manual',
        15,
        'seed_migration'
    ),
    (
        (
            SELECT school_id
            FROM schools
            WHERE nces_id = '062961004588'
        ),
        'CA',
        'success',
        8,
        'manual',
        12,
        'seed_migration'
    ),
    (
        (
            SELECT school_id
            FROM schools
            WHERE nces_id = '062961004589'
        ),
        'CA',
        'success',
        8,
        'manual',
        10,
        'seed_migration'
    ),
    (
        (
            SELECT school_id
            FROM schools
            WHERE nces_id = '484329004587'
        ),
        'TX',
        'success',
        8,
        'manual',
        14,
        'seed_migration'
    ),
    (
        (
            SELECT school_id
            FROM schools
            WHERE nces_id = '484329004588'
        ),
        'TX',
        'success',
        8,
        'manual',
        11,
        'seed_migration'
    ),
    (
        (
            SELECT school_id
            FROM schools
            WHERE nces_id = '484329004589'
        ),
        'TX',
        'success',
        8,
        'manual',
        9,
        'seed_migration'
    );