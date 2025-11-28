package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: Could not load .env file: %v", err)
	}

	// Get bot token
	botToken := os.Getenv("TELEGRAM_BOT_TOKEN")
	if botToken == "" {
		log.Fatal("TELEGRAM_BOT_TOKEN environment variable is required")
	}

	// Get database URL
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL environment variable is required")
	}

	log.Printf("Starting Telegram Todo Bot...")
	log.Printf("Bot Token: %s", botToken[:10]+"...")
	log.Printf("Database URL: %s", dbURL[:30]+"...")

	// Initialize database connection
	db, err := NewDatabase(dbURL)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	// Initialize Telegram bot
	bot, err := NewBot(botToken, db)
	if err != nil {
		log.Fatalf("Failed to initialize bot: %v", err)
	}

	log.Println("Bot initialized successfully!")

	// Start bot in a goroutine
	go func() {
		if err := bot.Start(); err != nil {
			log.Printf("Bot error: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the bot
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c

	log.Println("Shutting down bot...")
}
