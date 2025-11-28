package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	_ "github.com/jackc/pgx/v5/stdlib"
)

// Database handles all database operations
type Database struct {
	db *sql.DB
}

// NewDatabase creates a new database connection
func NewDatabase(databaseURL string) (*Database, error) {
	db, err := sql.Open("pgx/v5", databaseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Test the connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	log.Println("Successfully connected to NeonDB")

	// Create tables if they don't exist
	if err := createTables(db); err != nil {
		return nil, fmt.Errorf("failed to create tables: %w", err)
	}

	return &Database{db: db}, nil
}

// Close closes the database connection
func (d *Database) Close() error {
	return d.db.Close()
}

// createTables creates the necessary tables if they don't exist
func createTables(db *sql.DB) error {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS users (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			telegram_id BIGINT UNIQUE NOT NULL,
			name VARCHAR(255) NOT NULL,
			timezone VARCHAR(50) DEFAULT 'Asia/Bangkok',
			language VARCHAR(10) DEFAULT 'en',
			default_reminder_interval INTEGER DEFAULT 24,
			notification_style VARCHAR(20) DEFAULT 'detailed',
			created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
			updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
		)`,
		`CREATE TABLE IF NOT EXISTS todos (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
			title VARCHAR(500) NOT NULL,
			description TEXT,
			due_time TIMESTAMP WITH TIME ZONE,
			priority VARCHAR(20) DEFAULT 'medium',
			status VARCHAR(20) DEFAULT 'pending',
			tags TEXT,
			created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
			updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
		)`,
		`CREATE TABLE IF NOT EXISTS reminders (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			todo_id UUID NOT NULL REFERENCES todos(id) ON DELETE CASCADE,
			repeat_count INTEGER DEFAULT 1,
			repeat_interval_hours INTEGER DEFAULT 24,
			next_notify_time TIMESTAMP WITH TIME ZONE NOT NULL,
			snoozed_until TIMESTAMP WITH TIME ZONE,
			is_active BOOLEAN DEFAULT true,
			created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
			updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
		)`,
		`CREATE INDEX IF NOT EXISTS idx_users_telegram_id ON users(telegram_id)`,
		`CREATE INDEX IF NOT EXISTS idx_todos_user_id ON todos(user_id)`,
		`CREATE INDEX IF NOT EXISTS idx_todos_status ON todos(status)`,
		`CREATE INDEX IF NOT EXISTS idx_reminders_todo_id ON reminders(todo_id)`,
		`CREATE INDEX IF NOT EXISTS idx_reminders_next_notify ON reminders(next_notify_time) WHERE is_active = true`,
	}

	for _, query := range queries {
		if _, err := db.Exec(query); err != nil {
			return fmt.Errorf("failed to execute table creation query: %w", err)
		}
	}

	log.Println("Database tables created/verified successfully")
	return nil
}

// CreateUser creates a new user
func (d *Database) CreateUser(user NewUser) (*User, error) {
	ctx := context.Background()
	now := time.Now()

	query := `
		INSERT INTO users (telegram_id, name, timezone, language, default_reminder_interval, notification_style, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id, telegram_id, name, timezone, language, default_reminder_interval, notification_style, created_at, updated_at
	`

	var result User
	err := d.db.QueryRowContext(ctx, query,
		user.TelegramID, user.Name, user.Timezone, user.Language,
		24, "detailed", now, now,
	).Scan(
		&result.ID, &result.TelegramID, &result.Name, &result.Timezone,
		&result.Language, &result.DefaultReminderInterval, &result.NotificationStyle,
		&result.CreatedAt, &result.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return &result, nil
}

// UpdateUserLanguage updates the user's language preference
func (d *Database) UpdateUserLanguage(userID uuid.UUID, language string) error {
	ctx := context.Background()
	now := time.Now()

	query := `
		UPDATE users 
		SET language = $1, updated_at = $2
		WHERE id = $3
	`

	_, err := d.db.ExecContext(ctx, query, language, now, userID)
	if err != nil {
		return fmt.Errorf("failed to update user language: %w", err)
	}

	return nil
}

// GetUserByTelegramID gets a user by their Telegram ID
func (d *Database) GetUserByTelegramID(telegramID int64) (*User, error) {
	ctx := context.Background()

	query := `
		SELECT id, telegram_id, name, timezone, language, default_reminder_interval, notification_style, created_at, updated_at
		FROM users
		WHERE telegram_id = $1
	`

	var user User
	err := d.db.QueryRowContext(ctx, query, telegramID).Scan(
		&user.ID, &user.TelegramID, &user.Name, &user.Timezone,
		&user.Language, &user.DefaultReminderInterval, &user.NotificationStyle,
		&user.CreatedAt, &user.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get user by telegram ID: %w", err)
	}

	return &user, nil
}

// GetUserByID gets a user by their UUID
func (d *Database) GetUserByID(userID uuid.UUID) (*User, error) {
	ctx := context.Background()

	query := `
		SELECT id, telegram_id, name, timezone, language, default_reminder_interval, notification_style, created_at, updated_at
		FROM users
		WHERE id = $1
	`

	var user User
	err := d.db.QueryRowContext(ctx, query, userID).Scan(
		&user.ID, &user.TelegramID, &user.Name, &user.Timezone,
		&user.Language, &user.DefaultReminderInterval, &user.NotificationStyle,
		&user.CreatedAt, &user.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get user by ID: %w", err)
	}

	return &user, nil
}

// UpdateUserSettings updates user settings
func (d *Database) UpdateUserSettings(userID uuid.UUID, timezone, language *string, defaultInterval *int) (*User, error) {
	ctx := context.Background()
	now := time.Now()

	// Build dynamic update query
	query := "UPDATE users SET updated_at = $1"
	args := []interface{}{now}
	argIndex := 2

	if timezone != nil {
		query += fmt.Sprintf(", timezone = $%d", argIndex)
		args = append(args, *timezone)
		argIndex++
	}

	if language != nil {
		query += fmt.Sprintf(", language = $%d", argIndex)
		args = append(args, *language)
		argIndex++
	}

	if defaultInterval != nil {
		query += fmt.Sprintf(", default_reminder_interval = $%d", argIndex)
		args = append(args, *defaultInterval)
		argIndex++
	}

	query += fmt.Sprintf(" WHERE id = $%d", argIndex)
	args = append(args, userID)

	if _, err := d.db.ExecContext(ctx, query, args...); err != nil {
		return nil, fmt.Errorf("failed to update user settings: %w", err)
	}

	return d.GetUserByID(userID)
}

// CreateTodo creates a new todo
func (d *Database) CreateTodo(todo NewTodo) (*Todo, error) {
	ctx := context.Background()
	now := time.Now()

	query := `
		INSERT INTO todos (user_id, title, description, due_time, priority, status, tags, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id, user_id, title, description, due_time, priority, status, tags, created_at, updated_at
	`

	var result Todo
	err := d.db.QueryRowContext(ctx, query,
		todo.UserID, todo.Title, todo.Description, todo.DueTime,
		todo.Priority, "pending", todo.Tags, now, now,
	).Scan(
		&result.ID, &result.UserID, &result.Title, &result.Description,
		&result.DueTime, &result.Priority, &result.Status, &result.Tags,
		&result.CreatedAt, &result.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create todo: %w", err)
	}

	return &result, nil
}

// GetUserTodos gets all todos for a user
func (d *Database) GetUserTodos(userID uuid.UUID) ([]Todo, error) {
	ctx := context.Background()

	query := `
		SELECT id, user_id, title, description, due_time, priority, status, tags, created_at, updated_at
		FROM todos
		WHERE user_id = $1
		ORDER BY created_at DESC
	`

	rows, err := d.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user todos: %w", err)
	}
	defer rows.Close()

	var todos []Todo
	for rows.Next() {
		var todo Todo
		err := rows.Scan(
			&todo.ID, &todo.UserID, &todo.Title, &todo.Description,
			&todo.DueTime, &todo.Priority, &todo.Status, &todo.Tags,
			&todo.CreatedAt, &todo.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan todo: %w", err)
		}
		todos = append(todos, todo)
	}

	return todos, nil
}

// UpdateTodoStatus updates the status of a todo
func (d *Database) UpdateTodoStatus(todoID uuid.UUID, status string) (*Todo, error) {
	ctx := context.Background()
	now := time.Now()

	query := `
		UPDATE todos 
		SET status = $1, updated_at = $2
		WHERE id = $3
		RETURNING id, user_id, title, description, due_time, priority, status, tags, created_at, updated_at
	`

	var todo Todo
	err := d.db.QueryRowContext(ctx, query, status, now, todoID).Scan(
		&todo.ID, &todo.UserID, &todo.Title, &todo.Description,
		&todo.DueTime, &todo.Priority, &todo.Status, &todo.Tags,
		&todo.CreatedAt, &todo.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to update todo status: %w", err)
	}

	return &todo, nil
}

// DeleteTodo deletes a todo
func (d *Database) DeleteTodo(todoID uuid.UUID) error {
	ctx := context.Background()

	query := "DELETE FROM todos WHERE id = $1"
	if _, err := d.db.ExecContext(ctx, query, todoID); err != nil {
		return fmt.Errorf("failed to delete todo: %w", err)
	}

	return nil
}

// CreateReminder creates a new reminder
func (d *Database) CreateReminder(reminder NewReminder) (*Reminder, error) {
	ctx := context.Background()
	now := time.Now()

	query := `
		INSERT INTO reminders (todo_id, repeat_count, repeat_interval_hours, next_notify_time, is_active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, todo_id, repeat_count, repeat_interval_hours, next_notify_time, snoozed_until, is_active, created_at, updated_at
	`

	var result Reminder
	err := d.db.QueryRowContext(ctx, query,
		reminder.TodoID, reminder.RepeatCount, reminder.RepeatIntervalHours,
		reminder.NextNotifyTime, true, now, now,
	).Scan(
		&result.ID, &result.TodoID, &result.RepeatCount, &result.RepeatIntervalHours,
		&result.NextNotifyTime, &result.SnoozedUntil, &result.IsActive,
		&result.CreatedAt, &result.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create reminder: %w", err)
	}

	return &result, nil
}

// GetRemindersForTodo gets all reminders for a todo
func (d *Database) GetRemindersForTodo(todoID uuid.UUID) ([]Reminder, error) {
	ctx := context.Background()

	query := `
		SELECT id, todo_id, repeat_count, repeat_interval_hours, next_notify_time, snoozed_until, is_active, created_at, updated_at
		FROM reminders
		WHERE todo_id = $1
		ORDER BY created_at DESC
	`

	rows, err := d.db.QueryContext(ctx, query, todoID)
	if err != nil {
		return nil, fmt.Errorf("failed to get reminders for todo: %w", err)
	}
	defer rows.Close()

	var reminders []Reminder
	for rows.Next() {
		var reminder Reminder
		err := rows.Scan(
			&reminder.ID, &reminder.TodoID, &reminder.RepeatCount, &reminder.RepeatIntervalHours,
			&reminder.NextNotifyTime, &reminder.SnoozedUntil, &reminder.IsActive,
			&reminder.CreatedAt, &reminder.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan reminder: %w", err)
		}
		reminders = append(reminders, reminder)
	}

	return reminders, nil
}

// UpdateReminderNextTime updates the next notification time for a reminder
func (d *Database) UpdateReminderNextTime(reminderID uuid.UUID, nextTime time.Time) (*Reminder, error) {
	ctx := context.Background()
	now := time.Now()

	query := `
		UPDATE reminders 
		SET next_notify_time = $1, snoozed_until = NULL, updated_at = $2
		WHERE id = $3
		RETURNING id, todo_id, repeat_count, repeat_interval_hours, next_notify_time, snoozed_until, is_active, created_at, updated_at
	`

	var reminder Reminder
	err := d.db.QueryRowContext(ctx, query, nextTime, now, reminderID).Scan(
		&reminder.ID, &reminder.TodoID, &reminder.RepeatCount, &reminder.RepeatIntervalHours,
		&reminder.NextNotifyTime, &reminder.SnoozedUntil, &reminder.IsActive,
		&reminder.CreatedAt, &reminder.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to update reminder next time: %w", err)
	}

	return &reminder, nil
}

// SnoozeReminder snoozes a reminder
func (d *Database) SnoozeReminder(reminderID uuid.UUID, snoozeUntil time.Time) (*Reminder, error) {
	ctx := context.Background()
	now := time.Now()

	query := `
		UPDATE reminders 
		SET snoozed_until = $1, updated_at = $2
		WHERE id = $3
		RETURNING id, todo_id, repeat_count, repeat_interval_hours, next_notify_time, snoozed_until, is_active, created_at, updated_at
	`

	var reminder Reminder
	err := d.db.QueryRowContext(ctx, query, snoozeUntil, now, reminderID).Scan(
		&reminder.ID, &reminder.TodoID, &reminder.RepeatCount, &reminder.RepeatIntervalHours,
		&reminder.NextNotifyTime, &reminder.SnoozedUntil, &reminder.IsActive,
		&reminder.CreatedAt, &reminder.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to snooze reminder: %w", err)
	}

	return &reminder, nil
}

// GetDueReminders gets all reminders that are due to be sent
func (d *Database) GetDueReminders() ([]Reminder, error) {
	ctx := context.Background()
	now := time.Now()

	query := `
		SELECT id, todo_id, repeat_count, repeat_interval_hours, next_notify_time, snoozed_until, is_active, created_at, updated_at
		FROM reminders
		WHERE is_active = true 
		AND next_notify_time <= $1
		AND (snoozed_until IS NULL OR snoozed_until <= $1)
		ORDER BY next_notify_time ASC
	`

	rows, err := d.db.QueryContext(ctx, query, now)
	if err != nil {
		return nil, fmt.Errorf("failed to get due reminders: %w", err)
	}
	defer rows.Close()

	var reminders []Reminder
	for rows.Next() {
		var reminder Reminder
		err := rows.Scan(
			&reminder.ID, &reminder.TodoID, &reminder.RepeatCount, &reminder.RepeatIntervalHours,
			&reminder.NextNotifyTime, &reminder.SnoozedUntil, &reminder.IsActive,
			&reminder.CreatedAt, &reminder.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan reminder: %w", err)
		}
		reminders = append(reminders, reminder)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating reminders: %w", err)
	}

	return reminders, nil
}

// UpdateReminderTime updates the next notification time and repeat count for a reminder
func (d *Database) UpdateReminderTime(reminderID uuid.UUID, nextTime time.Time, repeatCount int) error {
	ctx := context.Background()
	now := time.Now()

	query := `
		UPDATE reminders 
		SET next_notify_time = $1, repeat_count = $2, updated_at = $3
		WHERE id = $4
	`

	_, err := d.db.ExecContext(ctx, query, nextTime, repeatCount, now, reminderID)
	if err != nil {
		return fmt.Errorf("failed to update reminder time: %w", err)
	}

	return nil
}

// DeleteReminder deletes a reminder
func (d *Database) DeleteReminder(reminderID uuid.UUID) error {
	ctx := context.Background()

	query := `DELETE FROM reminders WHERE id = $1`

	_, err := d.db.ExecContext(ctx, query, reminderID)
	if err != nil {
		return fmt.Errorf("failed to delete reminder: %w", err)
	}

	return nil
}

// GetAllOverdueTodos gets all overdue todos with their users
func (d *Database) GetAllOverdueTodos() ([]struct {
	Todo Todo
	User User
}, error) {
	ctx := context.Background()
	now := time.Now()

	query := `
		SELECT t.id, t.user_id, t.title, t.description, t.due_time, t.priority, t.status, t.tags, t.created_at, t.updated_at,
			   u.id, u.telegram_id, u.name, u.timezone, u.language, u.default_reminder_interval, u.notification_style, u.created_at, u.updated_at
		FROM todos t
		JOIN users u ON t.user_id = u.id
		WHERE t.status = 'pending' AND t.due_time IS NOT NULL AND t.due_time < $1
		ORDER BY t.due_time ASC
	`

	rows, err := d.db.QueryContext(ctx, query, now)
	if err != nil {
		return nil, fmt.Errorf("failed to get overdue todos: %w", err)
	}
	defer rows.Close()

	var results []struct {
		Todo Todo
		User User
	}

	for rows.Next() {
		var result struct {
			Todo Todo
			User User
		}
		err := rows.Scan(
			&result.Todo.ID, &result.Todo.UserID, &result.Todo.Title, &result.Todo.Description,
			&result.Todo.DueTime, &result.Todo.Priority, &result.Todo.Status, &result.Todo.Tags,
			&result.Todo.CreatedAt, &result.Todo.UpdatedAt,
			&result.User.ID, &result.User.TelegramID, &result.User.Name, &result.User.Timezone,
			&result.User.Language, &result.User.DefaultReminderInterval, &result.User.NotificationStyle,
			&result.User.CreatedAt, &result.User.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan overdue todo: %w", err)
		}
		results = append(results, result)
	}

	return results, nil
}

// GetTodoStats gets statistics for a user's todos
func (d *Database) GetTodoStats(userID uuid.UUID) (*TodoStats, error) {
	ctx := context.Background()
	now := time.Now()

	query := `
		SELECT 
			COUNT(*) as total,
			COUNT(*) FILTER (WHERE status = 'completed') as completed,
			COUNT(*) FILTER (WHERE status = 'pending') as pending,
			COUNT(*) FILTER (WHERE status = 'pending' AND due_time IS NOT NULL AND due_time < $1) as overdue,
			COUNT(*) FILTER (WHERE priority = 'high') as high_priority,
			COUNT(*) FILTER (WHERE priority = 'medium') as medium_priority,
			COUNT(*) FILTER (WHERE priority = 'low') as low_priority
		FROM todos
		WHERE user_id = $2
	`

	var stats TodoStats
	err := d.db.QueryRowContext(ctx, query, now, userID).Scan(
		&stats.Total, &stats.Completed, &stats.Pending, &stats.Overdue,
		&stats.HighPriority, &stats.MediumPriority, &stats.LowPriority,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to get todo stats: %w", err)
	}

	return &stats, nil
}

// GetTodoByID gets a todo by its ID
func (d *Database) GetTodoByID(todoID uuid.UUID) (*Todo, error) {
	ctx := context.Background()

	query := `
		SELECT id, user_id, title, description, due_time, priority, status, tags, created_at, updated_at
		FROM todos
		WHERE id = $1
	`

	var todo Todo
	err := d.db.QueryRowContext(ctx, query, todoID).Scan(
		&todo.ID, &todo.UserID, &todo.Title, &todo.Description,
		&todo.DueTime, &todo.Priority, &todo.Status, &todo.Tags,
		&todo.CreatedAt, &todo.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to get todo: %w", err)
	}

	return &todo, nil
}
