package main

import (
	"database/sql"
	"flag"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/sirupsen/logrus"
)

const (
	defaultDBPath         = "./data/academic_data.db"
	defaultMigrationsDir  = "./migrations"
	migrationsTableSchema = `
		CREATE TABLE IF NOT EXISTS schema_migrations (
			version INTEGER PRIMARY KEY,
			filename TEXT NOT NULL,
			applied_at TEXT NOT NULL,
			checksum TEXT NOT NULL
		);
	`
)

// MigrationFile represents a migration file
type MigrationFile struct {
	Version  int
	Filename string
	Path     string
	Type     string // "up" or "down"
}

// Config holds migration configuration
type Config struct {
	DBPath        string
	MigrationsDir string
	Command       string
	Steps         int
	Verbose       bool
	Force         bool
	DryRun        bool
}

func init() {
	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05",
	})
	logrus.SetLevel(logrus.InfoLevel)
}

func main() {
	config := parseFlags()

	if config.Verbose {
		logrus.SetLevel(logrus.DebugLevel)
	}

	log := logrus.WithFields(logrus.Fields{
		"db_path":        config.DBPath,
		"migrations_dir": config.MigrationsDir,
		"command":        config.Command,
	})

	// Validate command
	validCommands := []string{"up", "down", "down-all", "status", "version", "create"}
	if !contains(validCommands, config.Command) {
		log.Fatalf("Invalid command: %s. Valid commands: %s", config.Command, strings.Join(validCommands, ", "))
	}

	// Execute command
	if err := run(config, log); err != nil {
		log.WithError(err).Fatal("Migration failed")
	}

	log.Info("Migration completed successfully")
}

func run(config Config, log *logrus.Entry) error {
	switch config.Command {
	case "up":
		return migrateUp(config, log)
	case "down":
		return migrateDown(config, log)
	case "down-all":
		return migrateDownAll(config, log)
	case "status":
		return showStatus(config, log)
	case "version":
		return showVersion(config, log)
	case "create":
		return createMigration(config, log)
	default:
		return fmt.Errorf("unknown command: %s", config.Command)
	}
}

// migrateUp applies all pending migrations
func migrateUp(config Config, log *logrus.Entry) error {
	db, err := openDB(config.DBPath)
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}
	defer db.Close()

	// Ensure migrations table exists
	if err := createMigrationsTable(db); err != nil {
		return fmt.Errorf("failed to create migrations table: %w", err)
	}

	// Get applied migrations
	appliedVersions, err := getAppliedMigrations(db)
	if err != nil {
		return fmt.Errorf("failed to get applied migrations: %w", err)
	}

	// Get pending migrations
	allMigrations, err := getMigrationFiles(config.MigrationsDir, "up")
	if err != nil {
		return fmt.Errorf("failed to read migration files: %w", err)
	}

	pendingMigrations := filterPendingMigrations(allMigrations, appliedVersions)

	if len(pendingMigrations) == 0 {
		log.Info("No pending migrations")
		return nil
	}

	log.Infof("Found %d pending migration(s)", len(pendingMigrations))

	// Apply migrations
	for _, migration := range pendingMigrations {
		if config.Steps > 0 && len(appliedVersions) >= config.Steps {
			log.Infof("Reached step limit (%d), stopping", config.Steps)
			break
		}

		log.Infof("Applying migration: %s", migration.Filename)

		if config.DryRun {
			log.Info("[DRY RUN] Would apply migration")
			continue
		}

		if err := applyMigration(db, migration, log); err != nil {
			return fmt.Errorf("failed to apply migration %s: %w", migration.Filename, err)
		}

		log.Infof("✓ Applied: %s", migration.Filename)
		appliedVersions = append(appliedVersions, migration.Version)
	}

	return nil
}

// migrateDown rolls back the last migration
func migrateDown(config Config, log *logrus.Entry) error {
	db, err := openDB(config.DBPath)
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}
	defer db.Close()

	// Get applied migrations
	appliedVersions, err := getAppliedMigrations(db)
	if err != nil {
		return fmt.Errorf("failed to get applied migrations: %w", err)
	}

	if len(appliedVersions) == 0 {
		log.Info("No migrations to rollback")
		return nil
	}

	// Get last applied version
	lastVersion := appliedVersions[len(appliedVersions)-1]

	// Find corresponding down migration
	downMigrations, err := getMigrationFiles(config.MigrationsDir, "down")
	if err != nil {
		return fmt.Errorf("failed to read migration files: %w", err)
	}

	var targetMigration *MigrationFile
	for _, m := range downMigrations {
		if m.Version == lastVersion {
			targetMigration = &m
			break
		}
	}

	if targetMigration == nil {
		return fmt.Errorf("rollback migration not found for version %d", lastVersion)
	}

	log.Infof("Rolling back migration: %s", targetMigration.Filename)

	if config.DryRun {
		log.Info("[DRY RUN] Would rollback migration")
		return nil
	}

	if err := rollbackMigration(db, *targetMigration, log); err != nil {
		return fmt.Errorf("failed to rollback migration %s: %w", targetMigration.Filename, err)
	}

	log.Infof("✓ Rolled back: %s", targetMigration.Filename)
	return nil
}

// migrateDownAll rolls back all migrations
func migrateDownAll(config Config, log *logrus.Entry) error {
	db, err := openDB(config.DBPath)
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}
	defer db.Close()

	appliedVersions, err := getAppliedMigrations(db)
	if err != nil {
		return fmt.Errorf("failed to get applied migrations: %w", err)
	}

	if len(appliedVersions) == 0 {
		log.Info("No migrations to rollback")
		return nil
	}

	if !config.Force {
		log.Warn("This will rollback ALL migrations. Use --force to confirm.")
		return nil
	}

	log.Infof("Rolling back %d migration(s)", len(appliedVersions))

	// Rollback in reverse order
	for i := len(appliedVersions) - 1; i >= 0; i-- {
		version := appliedVersions[i]

		downMigrations, err := getMigrationFiles(config.MigrationsDir, "down")
		if err != nil {
			return fmt.Errorf("failed to read migration files: %w", err)
		}

		var targetMigration *MigrationFile
		for _, m := range downMigrations {
			if m.Version == version {
				targetMigration = &m
				break
			}
		}

		if targetMigration == nil {
			log.Warnf("Rollback migration not found for version %d, skipping", version)
			continue
		}

		log.Infof("Rolling back: %s", targetMigration.Filename)

		if config.DryRun {
			log.Info("[DRY RUN] Would rollback migration")
			continue
		}

		if err := rollbackMigration(db, *targetMigration, log); err != nil {
			return fmt.Errorf("failed to rollback migration %s: %w", targetMigration.Filename, err)
		}

		log.Infof("✓ Rolled back: %s", targetMigration.Filename)
	}

	return nil
}

// showStatus displays migration status
func showStatus(config Config, log *logrus.Entry) error {
	db, err := openDB(config.DBPath)
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}
	defer db.Close()

	if err := createMigrationsTable(db); err != nil {
		return fmt.Errorf("failed to create migrations table: %w", err)
	}

	appliedVersions, err := getAppliedMigrations(db)
	if err != nil {
		return fmt.Errorf("failed to get applied migrations: %w", err)
	}

	allMigrations, err := getMigrationFiles(config.MigrationsDir, "up")
	if err != nil {
		return fmt.Errorf("failed to read migration files: %w", err)
	}

	fmt.Println("\n=== Migration Status ===")
	fmt.Printf("Database: %s\n", config.DBPath)
	fmt.Printf("Migrations Directory: %s\n", config.MigrationsDir)
	fmt.Printf("Applied Migrations: %d\n", len(appliedVersions))
	fmt.Printf("Total Migrations: %d\n", len(allMigrations))
	fmt.Println()

	if len(allMigrations) == 0 {
		fmt.Println("No migration files found")
		return nil
	}

	fmt.Println("Version | Status   | Filename")
	fmt.Println("--------|----------|----------------------------------")

	for _, migration := range allMigrations {
		status := "Pending"
		if contains(appliedVersions, migration.Version) {
			status = "Applied"
		}
		fmt.Printf("%-7d | %-8s | %s\n", migration.Version, status, migration.Filename)
	}
	fmt.Println()

	return nil
}

// showVersion displays current schema version
func showVersion(config Config, log *logrus.Entry) error {
	db, err := openDB(config.DBPath)
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}
	defer db.Close()

	if err := createMigrationsTable(db); err != nil {
		return fmt.Errorf("failed to create migrations table: %w", err)
	}

	appliedVersions, err := getAppliedMigrations(db)
	if err != nil {
		return fmt.Errorf("failed to get applied migrations: %w", err)
	}

	if len(appliedVersions) == 0 {
		fmt.Println("No migrations applied (version: 0)")
		return nil
	}

	currentVersion := appliedVersions[len(appliedVersions)-1]
	fmt.Printf("Current schema version: %d\n", currentVersion)
	fmt.Printf("Total migrations applied: %d\n", len(appliedVersions))

	return nil
}

// createMigration creates a new migration file pair
func createMigration(config Config, log *logrus.Entry) error {
	if config.Command != "create" {
		return fmt.Errorf("invalid command for creating migration")
	}

	// Migration name from flag or generate default
	name := os.Getenv("NAME")
	if name == "" {
		return fmt.Errorf("migration name required. Usage: NAME=migration_name migrate create")
	}

	timestamp := time.Now().Unix()
	upFile := fmt.Sprintf("%s/%d_%s.up.sql", config.MigrationsDir, timestamp, name)
	downFile := fmt.Sprintf("%s/%d_%s.down.sql", config.MigrationsDir, timestamp, name)

	// Create up migration
	upContent := fmt.Sprintf("-- Migration: %s\n-- Created: %s\n\n-- Add your UP migration here\n\n", name, time.Now().Format(time.RFC3339))
	if err := os.WriteFile(upFile, []byte(upContent), 0644); err != nil {
		return fmt.Errorf("failed to create up migration: %w", err)
	}

	// Create down migration
	downContent := fmt.Sprintf("-- Migration: %s\n-- Created: %s\n\n-- Add your DOWN migration here\n\n", name, time.Now().Format(time.RFC3339))
	if err := os.WriteFile(downFile, []byte(downContent), 0644); err != nil {
		return fmt.Errorf("failed to create down migration: %w", err)
	}

	log.Infof("Created migration files:")
	log.Infof("  %s", upFile)
	log.Infof("  %s", downFile)

	return nil
}

// Database helper functions

func openDB(dbPath string) (*sql.DB, error) {
	// Ensure directory exists
	dir := filepath.Dir(dbPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create database directory: %w", err)
	}

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}

func createMigrationsTable(db *sql.DB) error {
	_, err := db.Exec(migrationsTableSchema)
	return err
}

func getAppliedMigrations(db *sql.DB) ([]int, error) {
	rows, err := db.Query("SELECT version FROM schema_migrations ORDER BY version ASC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var versions []int
	for rows.Next() {
		var version int
		if err := rows.Scan(&version); err != nil {
			return nil, err
		}
		versions = append(versions, version)
	}

	return versions, rows.Err()
}

func getMigrationFiles(dir string, migType string) ([]MigrationFile, error) {
	var migrations []MigrationFile

	err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}

		filename := d.Name()
		if !strings.HasSuffix(filename, fmt.Sprintf(".%s.sql", migType)) {
			return nil
		}

		var version int
		var name string
		_, err = fmt.Sscanf(filename, "%d_%s", &version, &name)
		if err != nil {
			return nil
		}

		migrations = append(migrations, MigrationFile{
			Version:  version,
			Filename: filename,
			Path:     path,
			Type:     migType,
		})

		return nil
	})

	if err != nil {
		return nil, err
	}

	// Sort by version
	sort.Slice(migrations, func(i, j int) bool {
		return migrations[i].Version < migrations[j].Version
	})

	return migrations, nil
}

func filterPendingMigrations(all []MigrationFile, applied []int) []MigrationFile {
	var pending []MigrationFile
	for _, m := range all {
		if !contains(applied, m.Version) {
			pending = append(pending, m)
		}
	}
	return pending
}

func applyMigration(db *sql.DB, migration MigrationFile, log *logrus.Entry) error {
	// Read migration file
	content, err := os.ReadFile(migration.Path)
	if err != nil {
		return fmt.Errorf("failed to read migration file: %w", err)
	}

	// Start transaction
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Execute migration
	if _, err := tx.Exec(string(content)); err != nil {
		return fmt.Errorf("failed to execute migration: %w", err)
	}

	// Record migration
	checksum := fmt.Sprintf("%x", len(content)) // Simple checksum
	_, err = tx.Exec(
		"INSERT INTO schema_migrations (version, filename, applied_at, checksum) VALUES (?, ?, ?, ?)",
		migration.Version, migration.Filename, time.Now().Format(time.RFC3339), checksum,
	)
	if err != nil {
		return fmt.Errorf("failed to record migration: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func rollbackMigration(db *sql.DB, migration MigrationFile, log *logrus.Entry) error {
	// Read migration file
	content, err := os.ReadFile(migration.Path)
	if err != nil {
		return fmt.Errorf("failed to read migration file: %w", err)
	}

	// Start transaction
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Execute migration
	if _, err := tx.Exec(string(content)); err != nil {
		return fmt.Errorf("failed to execute migration: %w", err)
	}

	// Remove migration record
	_, err = tx.Exec("DELETE FROM schema_migrations WHERE version = ?", migration.Version)
	if err != nil {
		return fmt.Errorf("failed to remove migration record: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// Utility functions

func parseFlags() Config {
	config := Config{}

	flag.StringVar(&config.DBPath, "db", getEnv("DB_PATH", defaultDBPath), "Database file path")
	flag.StringVar(&config.MigrationsDir, "dir", getEnv("MIGRATIONS_DIR", defaultMigrationsDir), "Migrations directory")
	flag.IntVar(&config.Steps, "steps", 0, "Number of migrations to apply (0 = all)")
	flag.BoolVar(&config.Verbose, "verbose", false, "Enable verbose output")
	flag.BoolVar(&config.Force, "force", false, "Force operation without confirmation")
	flag.BoolVar(&config.DryRun, "dry-run", false, "Show what would be done without executing")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [options] <command>\n\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Commands:\n")
		fmt.Fprintf(os.Stderr, "  up         Apply all pending migrations\n")
		fmt.Fprintf(os.Stderr, "  down       Rollback the last migration\n")
		fmt.Fprintf(os.Stderr, "  down-all   Rollback all migrations (requires --force)\n")
		fmt.Fprintf(os.Stderr, "  status     Show migration status\n")
		fmt.Fprintf(os.Stderr, "  version    Show current schema version\n")
		fmt.Fprintf(os.Stderr, "  create     Create new migration files (requires NAME env var)\n")
		fmt.Fprintf(os.Stderr, "\nOptions:\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\nExamples:\n")
		fmt.Fprintf(os.Stderr, "  %s up\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s down\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s status\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  NAME=add_users %s create\n", os.Args[0])
	}

	flag.Parse()

	// Get command from remaining args
	if flag.NArg() > 0 {
		config.Command = flag.Arg(0)
	} else {
		flag.Usage()
		os.Exit(1)
	}

	return config
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func contains[T comparable](slice []T, item T) bool {
	for _, v := range slice {
		if v == item {
			return true
		}
	}
	return false
}
