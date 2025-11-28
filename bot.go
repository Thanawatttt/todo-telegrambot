package main

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/mem"
)

// Language constants
const (
	LangEN = "en"
	LangTH = "th"
)

// Translation structure
type Translation struct {
	Welcome          string
	MainMenu         string
	MyTasks          string
	Statistics       string
	AddTask          string
	Reminders        string
	ServerStats      string
	Help             string
	Settings         string
	Snooze           string
	YourTodos        string
	TaskCompleted    string
	TaskDeleted      string
	ReminderSet      string
	InvalidTaskID    string
	TaskNotFound     string
	Language         string
	CurrentLanguage  string
	LanguageChanged  string
	HelpText         string
	NoTasks          string
	ReminderOptions  string
	TimeFormats      string
	Examples         string
	ProTips          string
	CurrentReminders string
}

// Translations map
var translations = map[string]Translation{
	LangEN: {
		Welcome:         "👋 Welcome to Todo Bot!\n\nI'll help you manage your tasks efficiently.",
		MainMenu:        "🏠 <b>Main Menu</b>\n\nChoose an option below:",
		MyTasks:         "📋 My Tasks",
		Statistics:       "📊 Statistics",
		AddTask:         "➕ Add Task",
		Reminders:        "⏰ Reminders",
		ServerStats:      "🖥️ Server Stats",
		Help:            "❓ Help",
		Settings:        "⚙️ Settings",
		Snooze:          "😴 Snooze",
		YourTodos:       "📋 <b>Your Todos:</b>",
		TaskCompleted:   "✅ Task completed successfully!\n\n<b>%s</b>",
		TaskDeleted:     "🗑️ Task deleted successfully!",
		ReminderSet:     "⏰ Reminder set successfully!\n\nI'll remind you in %s\n\n📅 %s",
		InvalidTaskID:    "Invalid task ID. Please use a number like 1, 2, 3...",
		TaskNotFound:     "Task not found. Please use a number between 1 and %d",
		Language:         "🌐 Language",
		CurrentLanguage:  "Current language: %s",
		LanguageChanged:  "Language changed to %s!",
		HelpText:         `🤖 <b>Todo Bot Help</b>

📝 <b>Task Management:</b>
• /add &lt;title&gt; [description] - Create a new task
• /list - View all your tasks
• /stats - View your task statistics

🔧 <b>Task Actions:</b>
• /complete &lt;id&gt; - Mark a task as completed
• /delete &lt;id&gt; - Delete a task

⏰ <b>Reminders:</b>
• /remind &lt;id&gt; &lt;time&gt; - Set a reminder for a task
• /snooze &lt;id&gt; &lt;time&gt; - Snooze a reminder
• /reminders - View all reminder options

📊 <b>Examples:</b>
• /add Buy groceries
• /add Meeting with John at 3pm
• /complete 1
• /remind 1 2h
• /remind 1 1d (every day)
• /remind 1 1h (every hour)
• /snooze 1 30m

⚙️ <b>Settings:</b>
• /start - Main menu
• /help - Show this help message`,
		NoTasks:         "You don't have any todos yet. Use /add to create one!",
		ReminderOptions:  `⏰ <b>Reminder Options</b>

<b>Time Formats:</b>
• 1h - 1 hour
• 2h - 2 hours
• 30m - 30 minutes
• 1d - 1 day (repeats daily)
• 1w - 1 week (repeats weekly)
• 1h - 1 hour (repeats hourly)

<b>Examples:</b>
• /remind 1 30m - Remind in 30 minutes
• /remind 2 1h - Remind every hour
• /remind 3 1d - Remind every day
• /remind 4 1w - Remind every week

<b>Pro Tips:</b>
• Use daily reminders for habits
• Use hourly reminders for urgent tasks
• Use weekly reminders for goals

<b>Your Current Reminders:</b>
Use /list to see all tasks and their reminders`,
		TimeFormats:      `<b>Time Formats:</b>`,
		Examples:         `<b>Examples:</b>`,
		ProTips:          `<b>Pro Tips:</b>`,
		CurrentReminders: `<b>Your Current Reminders:</b>`,
	},
	LangTH: {
		Welcome:         "👋 ยินดีต้อนรับสู่ Todo Bot!\n\nฉันจะช่วยคุณจัดการงานของคุณอย่างมีประสิทธิภาพ",
		MainMenu:        "🏠 <b>เมนูหลัก</b>\n\nเลือกตัวเลือกด้านล่าง:",
		MyTasks:         "📋 งานของฉัน",
		Statistics:       "📊 สถิติ",
		AddTask:         "➕ เพิ่มงาน",
		Reminders:        "⏰ การแจ้งเตือน",
		ServerStats:      "🖥️ สถิติเซิร์ฟเวอร์",
		Help:            "❓ ความช่วยเหลือ",
		Settings:        "⚙️ การตั้งค่า",
		Snooze:          "😴 พักการแจ้งเตือน",
		YourTodos:       "📋 <b>งานของคุณ:</b>",
		TaskCompleted:   "✅ ทำงานเสร็จสิ้นแล้ว!\n\n<b>%s</b>",
		TaskDeleted:     "🗑️ ลบงานเรียบร้อยแล้ว!",
		ReminderSet:     "⏰ ตั้งการแจ้งเตือนเรียบร้อยแล้ว!\n\nฉันจะแจ้งเตือนในอีก %s\n\n📅 %s",
		InvalidTaskID:    "รหัสงานไม่ถูกต้อง กรุณาใช้ตัวเลขเช่น 1, 2, 3...",
		TaskNotFound:     "ไม่พบงาน กรุณาใช้ตัวเลขระหว่าง 1 ถึง %d",
		Language:         "🌐 ภาษา",
		CurrentLanguage:  "ภาษาปัจจุบัน: %s",
		LanguageChanged:  "เปลี่ยนภาษาเป็น %s เรียบร้อยแล้ว!",
		HelpText:         `🤖 <b>ความช่วยเหลือ Todo Bot</b>

📝 <b>การจัดการงาน:</b>
• /add &lt;ชื่องาน&gt; [คำอธิบาย] - สร้างงานใหม่
• /list - ดูงานทั้งหมดของคุณ
• /stats - ดูสถิติงานของคุณ

🔧 <b>การกระทำงาน:</b>
• /complete &lt;id&gt; - ทำเครื่องหมายว่างานเสร็จสิ้น
• /delete &lt;id&gt; - ลบงาน

⏰ <b>การแจ้งเตือน:</b>
• /remind &lt;id&gt; &lt;เวลา&gt; - ตั้งการแจ้งเตือนสำหรับงาน
• /snooze &lt;id&gt; &lt;เวลา&gt; - พักการแจ้งเตือน
• /reminders - ดูตัวเลือกการแจ้งเตือนทั้งหมด

📊 <b>ตัวอย่าง:</b>
• /add ซื้อของ
• /add นัดกับจอห์น 3โมงเย็น
• /complete 1
• /remind 1 2h
• /remind 1 1d (ทุกวัน)
• /remind 1 1h (ทุกชั่วโมง)
• /snooze 1 30m

⚙️ <b>การตั้งค่า:</b>
• /start - เมนูหลัก
• /help - แสดงข้อความช่วยเหลือนี้`,
		NoTasks:         "คุณยังไม่มีงานใดๆ เลย ใช้ /add เพื่อสร้างงานแรกของคุณ!",
		ReminderOptions:  `⏰ <b>ตัวเลือกการแจ้งเตือน</b>

<b>รูปแบบเวลา:</b>
• 1h - 1 ชั่วโมง
• 2h - 2 ชั่วโมง
• 30m - 30 นาที
• 1d - 1 วัน (ทำซ้ำทุกวัน)
• 1w - 1 สัปดาห์ (ทำซ้ำทุกสัปดาห์)
• 1h - 1 ชั่วโมง (ทำซ้ำทุกชั่วโมง)

<b>ตัวอย่าง:</b>
• /remind 1 30m - แจ้งเตือนใน 30 นาที
• /remind 2 1h - แจ้งเตือนทุกชั่วโมง
• /remind 3 1d - แจ้งเตือนทุกวัน
• /remind 4 1w - แจ้งเตือนทุกสัปดาห์

<b>เคล็ดลับ:</b>
• ใช้การแจ้งเตือนรายวันสำหรับนิสัย
• ใช้การแจ้งเตือนทุกชั่วโมงสำหรับงานเร่งด่วน
• ใช้การแจ้งเตือนรายสัปดาห์สำหรับเป้าหมาย

<b>การแจ้งเตือนปัจจุบันของคุณ:</b>
ใช้ /list เพื่อดูงานทั้งหมดและการแจ้งเตือน`,
		TimeFormats:      `<b>รูปแบบเวลา:</b>`,
		Examples:         `<b>ตัวอย่าง:</b>`,
		ProTips:          `<b>เคล็ดลับ:</b>`,
		CurrentReminders: `<b>การแจ้งเตือนปัจจุบันของคุณ:</b>`,
	},
}

// getTranslation gets the translation for a user
func (b *Bot) getTranslation(userID int64) Translation {
	// Get user from database to check language preference
	user, err := b.db.GetUserByTelegramID(userID)
	if err != nil || user == nil || user.Language == "" {
		return translations[LangEN] // Default to English
	}
	
	if trans, exists := translations[user.Language]; exists {
		return trans
	}
	return translations[LangEN] // Default to English
}

// setUserLanguage sets the user's language preference
func (b *Bot) setUserLanguage(userID int64, language string) error {
	user, err := b.db.GetUserByTelegramID(userID)
	if err != nil {
		return err
	}
	
	return b.db.UpdateUserLanguage(user.ID, language)
}

// Bot represents the Telegram bot
type Bot struct {
	api      *tgbotapi.BotAPI
	db       *Database
	commands map[string]func(*tgbotapi.Message) error
}

// NewBot creates a new bot instance
func NewBot(token string, db *Database) (*Bot, error) {
	api, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, fmt.Errorf("failed to create bot API: %w", err)
	}

	api.Debug = false
	log.Printf("Authorized on account %s", api.Self.UserName)

	bot := &Bot{
		api: api,
		db:  db,
	}

	bot.setupCommands()
	return bot, nil
}

// setupCommands sets up the command handlers
func (b *Bot) setupCommands() {
	b.commands = map[string]func(*tgbotapi.Message) error{
		"start":       b.handleStart,
		"add":         b.handleAdd,
		"list":        b.handleList,
		"help":        b.handleHelp,
		"stats":       b.handleStats,
		"reminders":   b.handleReminders,
		"serverstats": b.handleServerStats,
		"delete":      b.handleDelete,
		"complete":    b.handleComplete,
		"remind":      b.handleRemind,
		"snooze":      b.handleSnooze,
	}
}

// Start starts the bot
func (b *Bot) Start() error {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := b.api.GetUpdatesChan(u)

	// Start reminder checker in background
	go b.reminderChecker()

	for update := range updates {
		if update.Message != nil {
			log.Printf("Received message: %s", update.Message.Text)
			if err := b.handleMessage(update.Message); err != nil {
				log.Printf("Error handling message: %v", err)
			}
		} else if update.CallbackQuery != nil {
			log.Printf("Received callback query: %s", update.CallbackQuery.Data)
			if err := b.handleCallbackQuery(update.CallbackQuery); err != nil {
				log.Printf("Error handling callback query: %v", err)
			}
		}
	}

	return nil
}

// handleMessage handles incoming messages
func (b *Bot) handleMessage(message *tgbotapi.Message) error {
	if message.IsCommand() {
		command := message.Command()
		if handler, exists := b.commands[command]; exists {
			return handler(message)
		}
		return b.handleUnknownCommand(message)
	}

	// Handle non-command messages
	return b.handleTextMessage(message)
}

// handleCallbackQuery handles callback queries from inline keyboards
func (b *Bot) handleCallbackQuery(callback *tgbotapi.CallbackQuery) error {
	// Parse callback data
	data := callback.Data
	if data == "" {
		_, err := b.api.Request(tgbotapi.CallbackConfig{
			CallbackQueryID: callback.ID,
		})
		return err
	}

	// Handle different callback actions
	parts := strings.Split(data, ":")
	if len(parts) < 2 && data != "main_menu" && data != "list" && data != "stats" && data != "help" && data != "add" && data != "settings" && data != "reminders" && data != "serverstats" && data != "lang_en" && data != "lang_th" {
		_, err := b.api.Request(tgbotapi.CallbackConfig{
			CallbackQueryID: callback.ID,
		})
		return err
	}

	var action, id string
	if len(parts) >= 2 {
		action = parts[0]
		id = parts[1]
	} else {
		action = data
	}

	switch action {
	case "complete":
		return b.handleCompleteCallback(callback, id)
	case "delete":
		return b.handleDeleteCallback(callback, id)
	case "snooze":
		return b.handleSnoozeCallback(callback, id)
	case "main_menu":
		return b.handleMainMenu(callback)
	case "list":
		return b.handleListFromCallback(callback)
	case "stats":
		return b.handleStatsFromCallback(callback)
	case "help":
		return b.handleHelpFromCallback(callback)
	case "add":
		msg := tgbotapi.NewMessage(callback.Message.Chat.ID, "Please use: /add &lt;task title&gt; to create a new task")
		msg.ParseMode = "HTML"
		_, err := b.api.Send(msg)
		return err
	case "reminders":
		return b.handleRemindersFromCallback(callback)
	case "settings":
		return b.handleSettings(callback)
	case "serverstats":
		return b.handleServerStatsFromCallback(callback)
	case "lang_en":
		return b.handleLanguageChange(callback, LangEN)
	case "lang_th":
		return b.handleLanguageChange(callback, LangTH)
	default:
		_, err := b.api.Request(tgbotapi.CallbackConfig{
			CallbackQueryID: callback.ID,
		})
		return err
	}
}

// handleListFromCallback handles the list command from a callback
func (b *Bot) handleListFromCallback(callback *tgbotapi.CallbackQuery) error {
	// Get user
	user, err := b.db.GetUserByTelegramID(callback.From.ID)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}

	// Get user's todos
	todos, err := b.db.GetUserTodos(user.ID)
	if err != nil {
		return fmt.Errorf("failed to get todos: %w", err)
	}

	trans := b.getTranslation(callback.From.ID)

	if len(todos) == 0 {
		noTasksText := trans.NoTasks
		msg := tgbotapi.NewMessage(callback.Message.Chat.ID, noTasksText)
		_, err := b.api.Send(msg)
		return err
	}

	// Build todo list
	var listText strings.Builder
	listText.WriteString(fmt.Sprintf("%s\n\n", trans.YourTodos))

	for i, todo := range todos {
		status := "🔴"
		if todo.Status == "completed" {
			status = "✅"
		}
		
		priority := ""
		if todo.Priority == "high" {
			priority = "� "
		} else if todo.Priority == "medium" {
			priority = "� "
		}
		
		listText.WriteString(fmt.Sprintf("%d. %s %s%s\n", i+1, status, priority, todo.Title))
		
		if todo.DueTime != nil {
			listText.WriteString(fmt.Sprintf("   📅 Due: %s\n", todo.DueTime.Format("2006-01-02 15:04")))
		}
		
		if todo.Description != nil && *todo.Description != "" {
			listText.WriteString(fmt.Sprintf("   📝 %s\n", *todo.Description))
		}
		
		listText.WriteString("\n")
	}

	// Create inline keyboard for each todo
	var keyboardRows [][]tgbotapi.InlineKeyboardButton
	
	for _, todo := range todos {
		if todo.Status == "pending" {
			row := tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("✅", fmt.Sprintf("complete:%s", todo.ID)),
				tgbotapi.NewInlineKeyboardButtonData("🗑️", fmt.Sprintf("delete:%s", todo.ID)),
				tgbotapi.NewInlineKeyboardButtonData("⏰", fmt.Sprintf("remind:%s", todo.ID)),
			)
			keyboardRows = append(keyboardRows, row)
		}
	}
	
	// Add navigation buttons
	navRow := tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("🏠 Main Menu", "main_menu"),
	)
	keyboardRows = append(keyboardRows, navRow)

	keyboard := tgbotapi.NewInlineKeyboardMarkup(keyboardRows...)

	msg := tgbotapi.NewMessage(callback.Message.Chat.ID, listText.String())
	msg.ParseMode = "HTML"
	msg.ReplyMarkup = keyboard

	_, err = b.api.Send(msg)
	return err
}

// handleStatsFromCallback handles the stats command from a callback
func (b *Bot) handleStatsFromCallback(callback *tgbotapi.CallbackQuery) error {
	// Get user
	user, err := b.db.GetUserByTelegramID(callback.From.ID)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}
	if user == nil {
		msg := tgbotapi.NewMessage(callback.Message.Chat.ID, "Please start with /start first")
		_, err := b.api.Send(msg)
		return err
	}

	// Get user's todo stats
	stats, err := b.db.GetTodoStats(user.ID)
	if err != nil {
		return fmt.Errorf("failed to get todo stats: %w", err)
	}

	statsText := fmt.Sprintf(`📊 <b>Your Todo Statistics</b>

📈 <b>Overview:</b>
• Total tasks: %d
• Completed: %d
• Pending: %d
• Overdue: %d

🎯 <b>Priority Breakdown:</b>
• High priority: %d
• Medium priority: %d
• Low priority: %d

📈 <b>Completion Rate:</b>
• %.1f%% completed`,
		stats.Total,
		stats.Completed,
		stats.Pending,
		stats.Overdue,
		stats.HighPriority,
		stats.MediumPriority,
		stats.LowPriority,
		float64(stats.Completed)/float64(stats.Total)*100,
	)

	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🏠 Main Menu", "main_menu"),
			tgbotapi.NewInlineKeyboardButtonData("📋 My Tasks", "list"),
		),
	)

	msg := tgbotapi.NewMessage(callback.Message.Chat.ID, statsText)
	msg.ParseMode = "HTML"
	msg.ReplyMarkup = keyboard

	_, err = b.api.Send(msg)
	return err
}

// handleHelpFromCallback handles the help command from a callback
func (b *Bot) handleHelpFromCallback(callback *tgbotapi.CallbackQuery) error {
	trans := b.getTranslation(callback.From.ID)

	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🏠 Main Menu", "main_menu"),
		),
	)

	msg := tgbotapi.NewMessage(callback.Message.Chat.ID, trans.HelpText)
	msg.ParseMode = "HTML"
	msg.ReplyMarkup = keyboard

	_, err := b.api.Send(msg)
	return err
}

// handleReminders handles the /reminders command
func (b *Bot) handleReminders(message *tgbotapi.Message) error {
	trans := b.getTranslation(message.From.ID)

	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🏠 Main Menu", "main_menu"),
		),
	)

	msg := tgbotapi.NewMessage(message.Chat.ID, trans.ReminderOptions)
	msg.ParseMode = "HTML"
	msg.ReplyMarkup = keyboard

	_, err := b.api.Send(msg)
	return err
}

// handleRemindersFromCallback handles the reminders command from a callback
func (b *Bot) handleRemindersFromCallback(callback *tgbotapi.CallbackQuery) error {
	trans := b.getTranslation(callback.From.ID)

	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(trans.MyTasks, "list"),
			tgbotapi.NewInlineKeyboardButtonData(trans.AddTask, "add"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🏠 Main Menu", "main_menu"),
			tgbotapi.NewInlineKeyboardButtonData(trans.Help, "help"),
		),
	)

	msg := tgbotapi.NewMessage(callback.Message.Chat.ID, trans.ReminderOptions)
	msg.ParseMode = "HTML"
	msg.ReplyMarkup = keyboard

	_, err := b.api.Send(msg)
	return err
}

// handleServerStats handles the /serverstats command
func (b *Bot) handleServerStats(message *tgbotapi.Message) error {
	// Get system information
	hostInfo, _ := host.Info()

	// Get CPU info
	cpuPercent, _ := cpu.Percent(time.Second, false)

	// Get memory info
	memInfo, _ := mem.VirtualMemory()

	// Get disk info
	diskInfo, _ := disk.Usage("/")

	// Get uptime
	var uptime string
	if hostInfo != nil {
		uptime = formatUptime(hostInfo.Uptime)
	}

	// Get bot version (hardcoded for now)
	botVersion := "1.0.0"

	// Build stats message
	statsText := fmt.Sprintf(`🖥️ <b>Server Statistics</b>

📊 <b>System Info:</b>
• <b>OS:</b> %s %s
• <b>Platform:</b> %s
• <b>Architecture:</b> %s
• <b>Hostname:</b> %s
• <b>Uptime:</b> %s

💻 <b>Hardware:</b>
• <b>CPU Usage:</b> %.1f%%
• <b>CPU Cores:</b> %d
• <b>Memory:</b> %s / %s (%.1f%%)
• <b>Disk:</b> %s / %s (%.1f%%)

🤖 <b>Bot Info:</b>
• <b>Version:</b> %s
• <b>Go Version:</b> %s
• <b>Process ID:</b> %d`,
		func() string {
			if hostInfo != nil {
				return hostInfo.OS + " " + hostInfo.Platform
			}
			return "Unknown"
		}(),
		func() string {
			if hostInfo != nil {
				return hostInfo.PlatformFamily
			}
			return "Unknown"
		}(),
		func() string {
			if hostInfo != nil {
				return hostInfo.PlatformVersion
			}
			return "Unknown"
		}(),
		func() string {
			if hostInfo != nil {
				return hostInfo.KernelArch
			}
			return runtime.GOARCH
		}(),
		func() string {
			if hostInfo != nil {
				return hostInfo.Hostname
			}
			return "Unknown"
		}(),
		uptime,
		func() float64 {
			if len(cpuPercent) > 0 {
				return cpuPercent[0]
			}
			return 0
		}(),
		runtime.NumCPU(),
		func() string {
			if memInfo != nil {
				return formatBytes(memInfo.Used)
			}
			return "Unknown"
		}(),
		func() string {
			if memInfo != nil {
				return formatBytes(memInfo.Total)
			}
			return "Unknown"
		}(),
		func() float64 {
			if memInfo != nil {
				return memInfo.UsedPercent
			}
			return 0
		}(),
		func() string {
			if diskInfo != nil {
				return formatBytes(diskInfo.Used)
			}
			return "Unknown"
		}(),
		func() string {
			if diskInfo != nil {
				return formatBytes(diskInfo.Total)
			}
			return "Unknown"
		}(),
		func() float64 {
			if diskInfo != nil {
				return diskInfo.UsedPercent
			}
			return 0
		}(),
		botVersion,
		runtime.Version(),
		os.Getpid())

	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🔄 Refresh", "serverstats"),
			tgbotapi.NewInlineKeyboardButtonData("🏠 Main Menu", "main_menu"),
		),
	)

	msg := tgbotapi.NewMessage(message.Chat.ID, statsText)
	msg.ParseMode = "HTML"
	msg.ReplyMarkup = keyboard

	_, err := b.api.Send(msg)
	return err
}

// formatBytes formats bytes into human readable format
func formatBytes(bytes uint64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}

// formatUptime formats uptime in seconds to human readable format
func formatUptime(seconds uint64) string {
	days := seconds / 86400
	hours := (seconds % 86400) / 3600
	minutes := (seconds % 3600) / 60
	
	if days > 0 {
		return fmt.Sprintf("%dd %dh %dm", days, hours, minutes)
	} else if hours > 0 {
		return fmt.Sprintf("%dh %dm", hours, minutes)
	} else {
		return fmt.Sprintf("%dm", minutes)
	}
}

// handleServerStatsFromCallback handles the serverstats command from a callback
func (b *Bot) handleServerStatsFromCallback(callback *tgbotapi.CallbackQuery) error {
	// Get system information
	hostInfo, _ := host.Info()

	// Get CPU info
	cpuPercent, _ := cpu.Percent(time.Second, false)

	// Get memory info
	memInfo, _ := mem.VirtualMemory()

	// Get disk info
	diskInfo, _ := disk.Usage("/")

	// Get uptime
	var uptime string
	if hostInfo != nil {
		uptime = formatUptime(hostInfo.Uptime)
	}

	// Get bot version (hardcoded for now)
	botVersion := "1.0.0"

	// Build stats message
	statsText := fmt.Sprintf(`🖥️ <b>Server Statistics</b>

📊 <b>System Info:</b>
• <b>OS:</b> %s %s
• <b>Platform:</b> %s
• <b>Architecture:</b> %s
• <b>Hostname:</b> %s
• <b>Uptime:</b> %s

💻 <b>Hardware:</b>
• <b>CPU Usage:</b> %.1f%%
• <b>CPU Cores:</b> %d
• <b>Memory:</b> %s / %s (%.1f%%)
• <b>Disk:</b> %s / %s (%.1f%%)

🤖 <b>Bot Info:</b>
• <b>Version:</b> %s
• <b>Go Version:</b> %s
• <b>Process ID:</b> %d`,
		func() string {
			if hostInfo != nil {
				return hostInfo.OS + " " + hostInfo.Platform
			}
			return "Unknown"
		}(),
		func() string {
			if hostInfo != nil {
				return hostInfo.PlatformFamily
			}
			return "Unknown"
		}(),
		func() string {
			if hostInfo != nil {
				return hostInfo.PlatformVersion
			}
			return "Unknown"
		}(),
		func() string {
			if hostInfo != nil {
				return hostInfo.KernelArch
			}
			return runtime.GOARCH
		}(),
		func() string {
			if hostInfo != nil {
				return hostInfo.Hostname
			}
			return "Unknown"
		}(),
		uptime,
		func() float64 {
			if len(cpuPercent) > 0 {
				return cpuPercent[0]
			}
			return 0
		}(),
		runtime.NumCPU(),
		func() string {
			if memInfo != nil {
				return formatBytes(memInfo.Used)
			}
			return "Unknown"
		}(),
		func() string {
			if memInfo != nil {
				return formatBytes(memInfo.Total)
			}
			return "Unknown"
		}(),
		func() float64 {
			if memInfo != nil {
				return memInfo.UsedPercent
			}
			return 0
		}(),
		func() string {
			if diskInfo != nil {
				return formatBytes(diskInfo.Used)
			}
			return "Unknown"
		}(),
		func() string {
			if diskInfo != nil {
				return formatBytes(diskInfo.Total)
			}
			return "Unknown"
		}(),
		func() float64 {
			if diskInfo != nil {
				return diskInfo.UsedPercent
			}
			return 0
		}(),
		botVersion,
		runtime.Version(),
		os.Getpid())

	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🔄 Refresh", "serverstats"),
			tgbotapi.NewInlineKeyboardButtonData("🏠 Main Menu", "main_menu"),
		),
	)

	msg := tgbotapi.NewMessage(callback.Message.Chat.ID, statsText)
	msg.ParseMode = "HTML"
	msg.ReplyMarkup = keyboard

	_, err := b.api.Send(msg)
	return err
}

// handleStart handles the /start command
func (b *Bot) handleStart(message *tgbotapi.Message) error {
	userID := message.From.ID
	userName := message.From.FirstName + " " + message.From.LastName

	// Check if user exists
	user, err := b.db.GetUserByTelegramID(userID)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}

	if user == nil {
		// Create new user
		newUser := NewUser{
			TelegramID: userID,
			Name:       userName,
			Timezone:   "UTC",
			Language:   "en",
		}
		user, err = b.db.CreateUser(newUser)
		if err != nil {
			return fmt.Errorf("failed to create user: %w", err)
		}
	}

	// Show main menu directly
	return b.handleMainMenuFromMessage(message, user)
}

// handleMainMenuFromMessage shows main menu from a regular message
func (b *Bot) handleMainMenuFromMessage(message *tgbotapi.Message, user *User) error {
	trans := b.getTranslation(message.From.ID)
	
	// Get user statistics
	todos, err := b.db.GetUserTodos(user.ID)
	if err != nil {
		return err
	}
	
	totalTasks := len(todos)
	completedTasks := 0
	pendingTasks := 0
	
	for _, todo := range todos {
		if todo.Status == "completed" {
			completedTasks++
		} else {
			pendingTasks++
		}
	}
	
	completionRate := 0.0
	if totalTasks > 0 {
		completionRate = float64(completedTasks) / float64(totalTasks) * 100
	}
	
	menuText := fmt.Sprintf(`🏠 <b>Main Menu</b>

👋 Welcome back, <b>%s</b>!

📊 <b>Your Statistics:</b>
• Total Tasks: <b>%d</b>
• Completed: <b>%d</b>
• Pending: <b>%d</b>
• Success Rate: <b>%.1f%%</b>`, 
		user.Name, totalTasks, completedTasks, pendingTasks, completionRate)

	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(trans.MyTasks, "list"),
			tgbotapi.NewInlineKeyboardButtonData(trans.Statistics, "stats"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(trans.AddTask, "add"),
			tgbotapi.NewInlineKeyboardButtonData(trans.Reminders, "reminders"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(trans.ServerStats, "serverstats"),
			tgbotapi.NewInlineKeyboardButtonData(trans.Settings, "settings"),
		),
	)

	msg := tgbotapi.NewMessage(message.Chat.ID, menuText)
	msg.ParseMode = "HTML"
	msg.ReplyMarkup = keyboard

	_, err = b.api.Send(msg)
	return err
}

// handleAdd handles the /add command
func (b *Bot) handleAdd(message *tgbotapi.Message) error {
	args := message.CommandArguments()
	if args == "" {
		msg := tgbotapi.NewMessage(message.Chat.ID, "Please provide a task title. Example: /add Buy groceries")
		_, err := b.api.Send(msg)
		return err
	}

	// Get user
	user, err := b.db.GetUserByTelegramID(message.From.ID)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}
	if user == nil {
		msg := tgbotapi.NewMessage(message.Chat.ID, "Please start with /start first")
		_, err := b.api.Send(msg)
		return err
	}

	// Parse the task (simple implementation for now)
	parts := strings.SplitN(args, " ", 2)
	if len(parts) < 1 {
		msg := tgbotapi.NewMessage(message.Chat.ID, "Please provide a task title")
		_, err := b.api.Send(msg)
		return err
	}

	title := parts[0]
	var description *string
	if len(parts) > 1 {
		desc := strings.Join(parts[1:], " ")
		description = &desc
	}

	// Create todo
	newTodo := NewTodo{
		UserID:      user.ID,
		Title:       title,
		Description: description,
		Priority:    "medium",
	}

	todo, err := b.db.CreateTodo(newTodo)
	if err != nil {
		return fmt.Errorf("failed to create todo: %w", err)
	}

	msgText := fmt.Sprintf("✅ Task created successfully!\n\n<b>%s</b>", todo.Title)
	if description != nil {
		msgText += fmt.Sprintf("\n\n%s", *description)
	}

	msg := tgbotapi.NewMessage(message.Chat.ID, msgText)
	msg.ParseMode = "HTML"

	_, err = b.api.Send(msg)
	return err
}

// handleList handles the /list command
func (b *Bot) handleList(message *tgbotapi.Message) error {
	// Get user
	user, err := b.db.GetUserByTelegramID(message.From.ID)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}
	if user == nil {
		msg := tgbotapi.NewMessage(message.Chat.ID, "Please start with /start first")
		_, err := b.api.Send(msg)
		return err
	}

	// Get user's todos
	todos, err := b.db.GetUserTodos(user.ID)
	if err != nil {
		return fmt.Errorf("failed to get todos: %w", err)
	}

	if len(todos) == 0 {
		msg := tgbotapi.NewMessage(message.Chat.ID, "You don't have any todos yet. Use /add to create one!")
		_, err := b.api.Send(msg)
		return err
	}

	// Build message with todos
	var msgText strings.Builder
	msgText.WriteString("📋 <b>Your Todos:</b>\n\n")

	for i, todo := range todos {
		status := "⏳"
		if todo.Status == "completed" {
			status = "✅"
		}

		priority := ""
		switch todo.Priority {
		case "high":
			priority = "🔴"
		case "medium":
			priority = "🟡"
		case "low":
			priority = "🟢"
		}

		dueTime := ""
		if todo.DueTime != nil {
			dueTime = fmt.Sprintf(" 📅 %s", todo.DueTime.Format("2006-01-02 15:04"))
		}

		msgText.WriteString(fmt.Sprintf("%d\\. %s %s *%s*%s\n", i+1, status, priority, escapeMarkdown(todo.Title), escapeMarkdown(dueTime)))

		if todo.Description != nil {
			msgText.WriteString(fmt.Sprintf("   %s\n", escapeMarkdown(*todo.Description)))
		}
	}

	// Add action buttons
	keyboard := tgbotapi.NewInlineKeyboardMarkup()
	for _, todo := range todos {
		if todo.Status != "completed" {
			row := tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("✅ Complete", fmt.Sprintf("complete:%s", todo.ID)),
				tgbotapi.NewInlineKeyboardButtonData("🗑️ Delete", fmt.Sprintf("delete:%s", todo.ID)),
			)
			keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, row)
		}
	}

	// Add main menu button at the bottom
	menuRow := tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("🏠 Main Menu", "main_menu"),
	)
	keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, menuRow)

	msg := tgbotapi.NewMessage(message.Chat.ID, msgText.String())
	msg.ParseMode = "HTML"
	msg.ReplyMarkup = keyboard

	_, err = b.api.Send(msg)
	return err
}

// handleHelp handles the /help command
func (b *Bot) handleHelp(message *tgbotapi.Message) error {
	helpText := `🤖 <b>Todo Bot Help</b>

📝 <b>Task Management:</b>
• /add &lt;title&gt; [description] - Create a new task
• /list - View all your tasks
• /stats - View your task statistics

🔧 <b>Task Actions:</b>
• /complete &lt;id&gt; - Mark a task as completed
• /delete &lt;id&gt; - Delete a task

⏰ <b>Reminders:</b>
• /remind &lt;id&gt; &lt;time&gt; - Set a reminder for a task
• /snooze &lt;id&gt; &lt;time&gt; - Snooze a reminder

📊 <b>Examples:</b>
• /add Buy groceries
• /add Meeting with John at 3pm
• /complete 1
• /remind 1 2h
• /snooze 1 30m

⚙️ <b>Settings:</b>
• /start - Register or welcome message
• /help - Show this help message`

	msg := tgbotapi.NewMessage(message.Chat.ID, helpText)
	msg.ParseMode = "HTML"

	_, err := b.api.Send(msg)
	return err
}

// handleStats handles the /stats command
func (b *Bot) handleStats(message *tgbotapi.Message) error {
	// Get user
	user, err := b.db.GetUserByTelegramID(message.From.ID)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}
	if user == nil {
		msg := tgbotapi.NewMessage(message.Chat.ID, "Please start with /start first")
		_, err := b.api.Send(msg)
		return err
	}

	// Get user's todo stats
	stats, err := b.db.GetTodoStats(user.ID)
	if err != nil {
		return fmt.Errorf("failed to get todo stats: %w", err)
	}

	statsText := fmt.Sprintf(`📊 <b>Your Todo Statistics</b>

📈 <b>Overview:</b>
• Total tasks: %d
• Completed: %d
• Pending: %d
• Overdue: %d

🎯 <b>Priority Breakdown:</b>
• High priority: %d
• Medium priority: %d
• Low priority: %d

📈 <b>Completion Rate:</b>
• %.1f%% completed`,
		stats.Total,
		stats.Completed,
		stats.Pending,
		stats.Overdue,
		stats.HighPriority,
		stats.MediumPriority,
		stats.LowPriority,
		float64(stats.Completed)/float64(stats.Total)*100,
	)

	msg := tgbotapi.NewMessage(message.Chat.ID, statsText)
	msg.ParseMode = "HTML"

	_, err = b.api.Send(msg)
	return err
}

// handleDelete handles the /delete command
func (b *Bot) handleDelete(message *tgbotapi.Message) error {
	args := message.CommandArguments()
	if args == "" {
		msg := tgbotapi.NewMessage(message.Chat.ID, "Please provide a task ID. Example: /delete 1")
		_, err := b.api.Send(msg)
		return err
	}

	// Convert args to UUID (simplified - in real app you'd use task numbers)
	todoID, err := uuid.Parse(args)
	if err != nil {
		msg := tgbotapi.NewMessage(message.Chat.ID, "Invalid task ID")
		_, err := b.api.Send(msg)
		return err
	}

	// Delete todo
	err = b.db.DeleteTodo(todoID)
	if err != nil {
		msg := tgbotapi.NewMessage(message.Chat.ID, "Failed to delete task")
		_, err2 := b.api.Send(msg)
		return err2
	}

	msg := tgbotapi.NewMessage(message.Chat.ID, "🗑️ Task deleted successfully!")
	_, err = b.api.Send(msg)
	return err
}

// handleComplete handles the /complete command
func (b *Bot) handleComplete(message *tgbotapi.Message) error {
	args := message.CommandArguments()
	if args == "" {
		msg := tgbotapi.NewMessage(message.Chat.ID, "Please provide a task ID. Example: /complete 1")
		_, err := b.api.Send(msg)
		return err
	}

	// Parse task number
	taskNum, err := strconv.Atoi(args)
	if err != nil {
		msg := tgbotapi.NewMessage(message.Chat.ID, "Invalid task ID. Please use a number like 1, 2, 3...")
		_, err := b.api.Send(msg)
		return err
	}

	// Get user
	user, err := b.db.GetUserByTelegramID(message.From.ID)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}
	if user == nil {
		msg := tgbotapi.NewMessage(message.Chat.ID, "Please start with /start first")
		_, err := b.api.Send(msg)
		return err
	}

	// Get user's todos to find the task by number
	todos, err := b.db.GetUserTodos(user.ID)
	if err != nil {
		return fmt.Errorf("failed to get todos: %w", err)
	}

	if taskNum < 1 || taskNum > len(todos) {
		msg := tgbotapi.NewMessage(message.Chat.ID, fmt.Sprintf("Task not found. Please use a number between 1 and %d", len(todos)))
		_, err := b.api.Send(msg)
		return err
	}

	// Get the task by index
	todo := todos[taskNum-1]

	// Update todo status
	updatedTodo, err := b.db.UpdateTodoStatus(todo.ID, "completed")
	if err != nil {
		msg := tgbotapi.NewMessage(message.Chat.ID, "Failed to complete task")
		_, err2 := b.api.Send(msg)
		return err2
	}

	msgText := fmt.Sprintf("✅ Task completed successfully!\n\n<b>%s</b>", updatedTodo.Title)
	msg := tgbotapi.NewMessage(message.Chat.ID, msgText)
	msg.ParseMode = "HTML"

	_, err = b.api.Send(msg)
	return err
}

// handleRemind handles the /remind command
func (b *Bot) handleRemind(message *tgbotapi.Message) error {
	args := message.CommandArguments()
	if args == "" {
		msg := tgbotapi.NewMessage(message.Chat.ID, "Please provide task ID and time. Example: /remind 1 2h")
		_, err := b.api.Send(msg)
		return err
	}

	parts := strings.SplitN(args, " ", 2)
	if len(parts) != 2 {
		msg := tgbotapi.NewMessage(message.Chat.ID, "Please provide task ID and time. Example: /remind 1 2h")
		_, err := b.api.Send(msg)
		return err
	}

	// Parse task number (not UUID)
	taskNumStr := parts[0]
	timeStr := parts[1]

	taskNum, err := strconv.Atoi(taskNumStr)
	if err != nil {
		msg := tgbotapi.NewMessage(message.Chat.ID, "Invalid task ID. Please use a number like 1, 2, 3...")
		_, err := b.api.Send(msg)
		return err
	}

	// Get user
	user, err := b.db.GetUserByTelegramID(message.From.ID)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}
	if user == nil {
		msg := tgbotapi.NewMessage(message.Chat.ID, "Please start with /start first")
		_, err := b.api.Send(msg)
		return err
	}

	// Get user's todos to find the task by number
	todos, err := b.db.GetUserTodos(user.ID)
	if err != nil {
		return fmt.Errorf("failed to get todos: %w", err)
	}

	if taskNum < 1 || taskNum > len(todos) {
		msg := tgbotapi.NewMessage(message.Chat.ID, fmt.Sprintf("Task not found. Please use a number between 1 and %d", len(todos)))
		_, err := b.api.Send(msg)
		return err
	}

	// Get the task by index
	todo := todos[taskNum-1]

	// Parse time duration
	duration, err := parseDuration(timeStr)
	if err != nil {
		msg := tgbotapi.NewMessage(message.Chat.ID, "Invalid time format. Use '2h' for 2 hours or '30m' for 30 minutes")
		_, err := b.api.Send(msg)
		return err
	}

	// Calculate next notification time
	nextTime := time.Now().Add(duration)

	// Create reminder
	newReminder := NewReminder{
		TodoID:                 todo.ID,
		RepeatCount:            1,
		RepeatIntervalHours:    int(duration.Hours()),
		NextNotifyTime:         nextTime,
	}

	_, err = b.db.CreateReminder(newReminder)
	if err != nil {
		msg := tgbotapi.NewMessage(message.Chat.ID, "Failed to create reminder")
		_, err2 := b.api.Send(msg)
		return err2
	}

	msgText := fmt.Sprintf("⏰ Reminder set successfully!\n\nI'll remind you in %s\n\n📅 %s",
		duration.String(), nextTime.Format("2006-01-02 15:04"))
	msg := tgbotapi.NewMessage(message.Chat.ID, msgText)
	msg.ParseMode = "HTML"

	_, err = b.api.Send(msg)
	return err
}

// handleSnooze handles the /snooze command
func (b *Bot) handleSnooze(message *tgbotapi.Message) error {
	args := message.CommandArguments()
	if args == "" {
		msg := tgbotapi.NewMessage(message.Chat.ID, "Please provide reminder ID and time. Example: /snooze 1 30m")
		_, err := b.api.Send(msg)
		return err
	}

	// Parse args (simplified)
	parts := strings.Split(args, " ")
	if len(parts) < 2 {
		msg := tgbotapi.NewMessage(message.Chat.ID, "Please provide reminder ID and time. Example: /snooze 1 30m")
		_, err := b.api.Send(msg)
		return err
	}

	reminderID, err := uuid.Parse(parts[0])
	if err != nil {
		msg := tgbotapi.NewMessage(message.Chat.ID, "Invalid reminder ID")
		_, err := b.api.Send(msg)
		return err
	}

	// Parse time (simplified)
	timeStr := parts[1]
	var duration time.Duration
	if strings.HasSuffix(timeStr, "h") {
		hours := strings.TrimSuffix(timeStr, "h")
		hoursInt := 1 // default
		fmt.Sscanf(hours, "%d", &hoursInt)
		duration = time.Duration(hoursInt) * time.Hour
	} else if strings.HasSuffix(timeStr, "m") {
		minutes := strings.TrimSuffix(timeStr, "m")
		minutesInt := 30 // default
		fmt.Sscanf(minutes, "%d", &minutesInt)
		duration = time.Duration(minutesInt) * time.Minute
	} else {
		msg := tgbotapi.NewMessage(message.Chat.ID, "Invalid time format. Use '2h' for 2 hours or '30m' for 30 minutes")
		_, err := b.api.Send(msg)
		return err
	}

	snoozeUntil := time.Now().Add(duration)

	// Snooze reminder
	_, err = b.db.SnoozeReminder(reminderID, snoozeUntil)
	if err != nil {
		msg := tgbotapi.NewMessage(message.Chat.ID, "Failed to snooze reminder")
		_, err2 := b.api.Send(msg)
		return err2
	}

	msgText := fmt.Sprintf("😴 Reminder snoozed successfully!\n\nI'll remind you again in %s\n\n📅 %s",
		duration.String(), snoozeUntil.Format("2006-01-02 15:04"))
	msg := tgbotapi.NewMessage(message.Chat.ID, msgText)
	msg.ParseMode = "HTML"

	_, err = b.api.Send(msg)
	return err
}

// handleUnknownCommand handles unknown commands
func (b *Bot) handleUnknownCommand(message *tgbotapi.Message) error {
	msg := tgbotapi.NewMessage(message.Chat.ID, "Unknown command. Use /help to see available commands.")
	_, err := b.api.Send(msg)
	return err
}

// handleTextMessage handles non-command text messages
func (b *Bot) handleTextMessage(message *tgbotapi.Message) error {
	msg := tgbotapi.NewMessage(message.Chat.ID, "I can help you manage your todos! Use /help to see available commands.")
	_, err := b.api.Send(msg)
	return err
}

// handleCompleteCallback handles the complete callback
func (b *Bot) handleCompleteCallback(callback *tgbotapi.CallbackQuery, todoIDStr string) error {
	todoID, err := uuid.Parse(todoIDStr)
	if err != nil {
		_, err := b.api.Request(tgbotapi.CallbackConfig{
			CallbackQueryID: callback.ID,
		})
		return err
	}

	// Update todo status
	_, err = b.db.UpdateTodoStatus(todoID, "completed")
	if err != nil {
		_, err := b.api.Request(tgbotapi.CallbackConfig{
			CallbackQueryID: callback.ID,
			Text:            "Failed to complete task",
		})
		return err
	}

	// Send updated list
	return b.handleListFromCallback(callback)
}

// handleDeleteCallback handles the delete callback
func (b *Bot) handleDeleteCallback(callback *tgbotapi.CallbackQuery, todoIDStr string) error {
	todoID, err := uuid.Parse(todoIDStr)
	if err != nil {
		_, err := b.api.Request(tgbotapi.CallbackConfig{
			CallbackQueryID: callback.ID,
		})
		return err
	}

	// Delete todo
	err = b.db.DeleteTodo(todoID)
	if err != nil {
		_, err := b.api.Request(tgbotapi.CallbackConfig{
			CallbackQueryID: callback.ID,
			Text:            "Failed to delete task",
		})
		return err
	}

	// Send updated list
	return b.handleListFromCallback(callback)
}

// handleSnoozeCallback handles the snooze callback
func (b *Bot) handleSnoozeCallback(callback *tgbotapi.CallbackQuery, reminderIDStr string) error {
	reminderID, err := uuid.Parse(reminderIDStr)
	if err != nil {
		_, err := b.api.Request(tgbotapi.CallbackConfig{
			CallbackQueryID: callback.ID,
		})
		return err
	}

	// Snooze for 30 minutes by default
	snoozeUntil := time.Now().Add(30 * time.Minute)

	// Snooze reminder
	_, err = b.db.SnoozeReminder(reminderID, snoozeUntil)
	if err != nil {
		_, err := b.api.Request(tgbotapi.CallbackConfig{
			CallbackQueryID: callback.ID,
			Text:            "Failed to snooze reminder",
		})
		return err
	}

	// Send callback response
	callbackText := fmt.Sprintf("😴 Snoozed until %s", snoozeUntil.Format("15:04"))
	_, err = b.api.Request(tgbotapi.CallbackConfig{
		CallbackQueryID: callback.ID,
		Text:            callbackText,
	})
	return err
}

// handleSettings handles the settings callback
func (b *Bot) handleSettings(callback *tgbotapi.CallbackQuery) error {
	// Get user info
	user, err := b.db.GetUserByTelegramID(callback.From.ID)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}

	trans := b.getTranslation(callback.From.ID)
	
	currentLang := "English"
	if user != nil && user.Language == LangTH {
		currentLang = "ไทย (Thai)"
	}

	settingsText := fmt.Sprintf(`⚙️ <b>Settings</b>

👤 <b>User Info:</b>
• Name: <b>%s</b>
• Timezone: <b>%s</b>

%s: <b>%s</b>

🌐 <b>Language Selection:</b>
Choose your preferred language:`, 
		user.Name, user.Timezone,
		trans.CurrentLanguage,
		currentLang)

	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🇺🇸 English", "lang_en"),
			tgbotapi.NewInlineKeyboardButtonData("🇹🇭 ไทย", "lang_th"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🏠 Main Menu", "main_menu"),
			tgbotapi.NewInlineKeyboardButtonData("❓ Help", "help"),
		),
	)

	msg := tgbotapi.NewMessage(callback.Message.Chat.ID, settingsText)
	msg.ParseMode = "HTML"
	msg.ReplyMarkup = keyboard

	_, err = b.api.Send(msg)
	return err
}

// handleLanguageChange handles the language change callback
func (b *Bot) handleLanguageChange(callback *tgbotapi.CallbackQuery, language string) error {
	// Update user's language preference
	err := b.setUserLanguage(callback.From.ID, language)
	if err != nil {
		return err
	}

	trans := b.getTranslation(callback.From.ID)
	
	langName := "English"
	if language == LangTH {
		langName = "ไทย (Thai)"
	}

	confirmationText := fmt.Sprintf(trans.LanguageChanged, langName)

	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🏠 Main Menu", "main_menu"),
		),
	)

	msg := tgbotapi.NewMessage(callback.Message.Chat.ID, confirmationText)
	msg.ParseMode = "HTML"
	msg.ReplyMarkup = keyboard

	_, err = b.api.Send(msg)
	return err
}

// handleMainMenu handles the main menu callback
func (b *Bot) handleMainMenu(callback *tgbotapi.CallbackQuery) error {
	trans := b.getTranslation(callback.From.ID)
	
	// Get user info for statistics
	user, err := b.db.GetUserByTelegramID(callback.From.ID)
	if err != nil {
		return err
	}
	
	// Get user statistics
	todos, err := b.db.GetUserTodos(user.ID)
	if err != nil {
		return err
	}
	
	totalTasks := len(todos)
	completedTasks := 0
	pendingTasks := 0
	
	for _, todo := range todos {
		if todo.Status == "completed" {
			completedTasks++
		} else {
			pendingTasks++
		}
	}
	
	completionRate := 0.0
	if totalTasks > 0 {
		completionRate = float64(completedTasks) / float64(totalTasks) * 100
	}
	
	menuText := fmt.Sprintf(`🏠 <b>Main Menu</b>

👋 Welcome back, <b>%s</b>!

📊 <b>Your Statistics:</b>
• Total Tasks: <b>%d</b>
• Completed: <b>%d</b>
• Pending: <b>%d</b>
• Success Rate: <b>%.1f%%</b>`, 
		user.Name, totalTasks, completedTasks, pendingTasks, completionRate)

	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(trans.MyTasks, "list"),
			tgbotapi.NewInlineKeyboardButtonData(trans.Statistics, "stats"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(trans.AddTask, "add"),
			tgbotapi.NewInlineKeyboardButtonData(trans.Reminders, "reminders"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(trans.ServerStats, "serverstats"),
			tgbotapi.NewInlineKeyboardButtonData(trans.Settings, "settings"),
		),
	)

	msg := tgbotapi.NewMessage(callback.Message.Chat.ID, menuText)
	msg.ParseMode = "HTML"
	msg.ReplyMarkup = keyboard

	_, err = b.api.Send(msg)
	return err
}

// reminderChecker runs in background to check and send due reminders
func (b *Bot) reminderChecker() {
	ticker := time.NewTicker(30 * time.Second) // Check every 30 seconds
	defer ticker.Stop()

	for range ticker.C {
		b.checkAndSendReminders()
	}
}

// checkAndSendReminders checks for due reminders and sends notifications
func (b *Bot) checkAndSendReminders() {
	reminders, err := b.db.GetDueReminders()
	if err != nil {
		return
	}

	for _, reminder := range reminders {
		// Get the todo details
		todo, err := b.db.GetTodoByID(reminder.TodoID)
		if err != nil {
			continue
		}

		// Get the user
		user, err := b.db.GetUserByID(todo.UserID)
		if err != nil {
			continue
		}

		// Send reminder notification
		reminderText := fmt.Sprintf(`⏰ <b>Reminder!</b>

📝 <b>%s</b>

%s

Don't forget to complete this task! 💪

Use /complete %d to mark it done`, 
			todo.Title, 
			func() string {
				if todo.Description != nil && *todo.Description != "" {
					return *todo.Description
				}
				return "No description"
			}(),
			// We need to find the task number for this user
			b.getTaskNumber(user.ID, todo.ID))

		msg := tgbotapi.NewMessage(user.TelegramID, reminderText)
		msg.ParseMode = "HTML"

		_, err = b.api.Send(msg)
		if err != nil {
			continue
		}

		// Update the next reminder time
		if reminder.RepeatCount > 1 {
			// This is a repeating reminder
			nextTime := reminder.NextNotifyTime.Add(time.Duration(reminder.RepeatIntervalHours) * time.Hour)
			err = b.db.UpdateReminderTime(reminder.ID, nextTime, reminder.RepeatCount-1)
			if err != nil {
				// Continue even if update fails
			}
		} else {
			// One-time reminder, delete it
			err = b.db.DeleteReminder(reminder.ID)
			if err != nil {
				// Continue even if delete fails
			}
		}
	}
}

// getTaskNumber finds the task number for a given todo ID
func (b *Bot) getTaskNumber(userID uuid.UUID, todoID uuid.UUID) int {
	todos, err := b.db.GetUserTodos(userID)
	if err != nil {
		return 1 // fallback
	}

	for i, todo := range todos {
		if todo.ID == todoID {
			return i + 1
		}
	}
	return 1 // fallback
}

// parseDuration parses time strings like "2h", "30m", "1d"
func parseDuration(timeStr string) (time.Duration, error) {
	if strings.HasSuffix(timeStr, "h") {
		hours, err := strconv.Atoi(timeStr[:len(timeStr)-1])
		if err != nil {
			return 0, err
		}
		return time.Duration(hours) * time.Hour, nil
	} else if strings.HasSuffix(timeStr, "m") {
		minutes, err := strconv.Atoi(timeStr[:len(timeStr)-1])
		if err != nil {
			return 0, err
		}
		return time.Duration(minutes) * time.Minute, nil
	} else if strings.HasSuffix(timeStr, "d") {
		days, err := strconv.Atoi(timeStr[:len(timeStr)-1])
		if err != nil {
			return 0, err
		}
		return time.Duration(days) * 24 * time.Hour, nil
	} else if strings.HasSuffix(timeStr, "w") {
		weeks, err := strconv.Atoi(timeStr[:len(timeStr)-1])
		if err != nil {
			return 0, err
		}
		return time.Duration(weeks) * 7 * 24 * time.Hour, nil
	}
	
	return 0, fmt.Errorf("invalid time format: %s", timeStr)
}

// escapeMarkdown escapes special characters for Telegram MarkdownV2
func escapeMarkdown(text string) string {
	// Telegram MarkdownV2 requires escaping: _ * [ ] ( ) ~ ` > # + - = | { } . !
	specialChars := []string{"_", "*", "[", "]", "(", ")", "~", "`", ">", "#", "+", "-", "=", "|", "{", "}", ".", "!"}
	result := text
	for _, char := range specialChars {
		result = strings.ReplaceAll(result, char, "\\"+char)
	}
	return result
}
