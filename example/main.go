package main

import (
	"context"
	"log"
	"os"
	"time"

	agi "github.com/Shubham-Thakur06/go-asterisk-agi"
)

func main() {
	// Check if we're running as FastAGI server
	if len(os.Args) > 1 && os.Args[1] == "server" {
		runFastAGIServer()
		return
	}

	// Regular AGI script
	runAGIScript()
}

func runAGIScript() {
	// Create new AGI session
	session, err := agi.NewAgiSession()
	if err != nil {
		log.Fatal(err)
	}
	defer session.Close()

	// Enable debug mode to see AGI commands and responses
	session.SetDebug(true)

	// Answer the call
	if err := session.Answer(); err != nil {
		log.Fatal(err)
	}

	// Get channel variables
	uniqueid, err := session.GetVariable("UNIQUEID")
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Call UNIQUEID: %s", uniqueid)

	// Play welcome message
	if err := session.StreamFile("welcome", "0123456789*#"); err != nil {
		log.Fatal(err)
	}

	// Wait for digit
	digit, err := session.WaitForDigit(5000)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("User pressed: %s", digit)

	// Get more digits
	digits, err := session.GetData("enter-account", 5000, 4)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("User entered: %s", digits)

	// Set a variable
	if err := session.SetVariable("ACCOUNT", digits); err != nil {
		log.Fatal(err)
	}

	// Hangup
	if err := session.Hangup(); err != nil {
		log.Fatal(err)
	}
}

func runFastAGIServer() {
	// Create a handler function
	handler := agi.HandlerFunc(func(ctx context.Context, s *agi.AgiSession) error {
		// Enable debug mode
		s.SetDebug(true)

		// Set operation timeout
		s.SetTimeout(30 * time.Second)

		// Answer the call
		if err := s.Answer(); err != nil {
			return err
		}

		// Get some AGI environment variables
		log.Printf("Network Script: %s", s.GetEnv("agi_network_script"))
		log.Printf("Caller ID: %s", s.GetEnv("agi_callerid"))

		// Play welcome message
		if err := s.StreamFile("welcome", "0123456789*#"); err != nil {
			return err
		}

		// Get user input
		digits, err := s.GetData("enter-account", 5000, 4)
		if err != nil {
			return err
		}
		log.Printf("User entered: %s", digits)

		// Hangup
		return s.Hangup()
	})

	// Create and start FastAGI server
	server, err := agi.NewFastAGIServer(":4573", handler)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("FastAGI server listening on :4573")
	if err := server.Serve(); err != nil {
		log.Fatal(err)
	}
}
