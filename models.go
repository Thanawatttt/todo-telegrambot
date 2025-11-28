package main

import (
	"time"

	"github.com/google/uuid"
)

// User represents a user in the database
type User struct {
	ID                       uuid.UUID  `json:"id"`
	TelegramID              int64      `json:"telegram_id"`
	Name                    string     `json:"name"`
	Timezone                string     `json:"timezone"`
	Language                string     `json:"language"`
	DefaultReminderInterval int        `json:"default_reminder_interval"`
	NotificationStyle       string     `json:"notification_style"`
	CreatedAt               time.Time  `json:"created_at"`
	UpdatedAt               time.Time  `json:"updated_at"`
}

// Todo represents a todo item in the database
type Todo struct {
	ID          uuid.UUID  `json:"id"`
	UserID      uuid.UUID  `json:"user_id"`
	Title       string     `json:"title"`
	Description *string    `json:"description,omitempty"`
	DueTime     *time.Time `json:"due_time,omitempty"`
	Priority    string     `json:"priority"`
	Status      string     `json:"status"`
	Tags        *string    `json:"tags,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

// Reminder represents a reminder in the database
type Reminder struct {
	ID                     uuid.UUID  `json:"id"`
	TodoID                 uuid.UUID  `json:"todo_id"`
	RepeatCount            int        `json:"repeat_count"`
	RepeatIntervalHours    int        `json:"repeat_interval_hours"`
	NextNotifyTime         time.Time  `json:"next_notify_time"`
	SnoozedUntil          *time.Time `json:"snoozed_until,omitempty"`
	IsActive               bool       `json:"is_active"`
	CreatedAt              time.Time  `json:"created_at"`
	UpdatedAt              time.Time  `json:"updated_at"`
}

// TodoStats represents statistics for todos
type TodoStats struct {
	Total         int `json:"total"`
	Completed     int `json:"completed"`
	Pending       int `json:"pending"`
	Overdue       int `json:"overdue"`
	HighPriority  int `json:"high_priority"`
	MediumPriority int `json:"medium_priority"`
	LowPriority   int `json:"low_priority"`
}

// NewUser represents a new user to be created
type NewUser struct {
	TelegramID int64  `json:"telegram_id"`
	Name       string `json:"name"`
	Timezone   string `json:"timezone"`
	Language   string `json:"language"`
}

// NewTodo represents a new todo to be created
type NewTodo struct {
	UserID      uuid.UUID  `json:"user_id"`
	Title       string     `json:"title"`
	Description *string    `json:"description,omitempty"`
	DueTime     *time.Time `json:"due_time,omitempty"`
	Priority    string     `json:"priority"`
	Tags        *string    `json:"tags,omitempty"`
}

// NewReminder represents a new reminder to be created
type NewReminder struct {
	TodoID                 uuid.UUID `json:"todo_id"`
	RepeatCount            int       `json:"repeat_count"`
	RepeatIntervalHours    int       `json:"repeat_interval_hours"`
	NextNotifyTime         time.Time `json:"next_notify_time"`
}
