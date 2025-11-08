// Package session provides session management for Quanto applications.
package session

import (
	"fmt"
	"log"
	"strings"

	"github.com/google/uuid"
)

// Mode represents the execution mode for a Quanto session.
type Mode string

const (
	// Local mode runs Quanto operations on the local machine.
	Local Mode = "local"
	// Cluster mode runs Quanto operations in a distributed cluster.
	Cluster Mode = "cluster"
)

// QuantoSession represents a session for running Quanto operations.
type QuantoSession struct {
	ID      string
	AppName string
	Mode    Mode
}

// New creates a new QuantoSession with a unique identifier.
func New() *QuantoSession {
	return &QuantoSession{
		ID: uuid.NewString(),
	}
}

// SetAppName sets the application name for the session and returns the session for chaining.
func (s *QuantoSession) SetAppName(appName string) *QuantoSession {
	s.AppName = appName
	return s
}

// SetMode sets the execution mode for the session and returns the session for chaining.
func (s *QuantoSession) SetMode(mode string) *QuantoSession {
	s.Mode = map[string]Mode{"local": Local, "cluster": Cluster}[mode]
	log.Printf("Running in %s mode\n", mode)
	return s
}

// GetOrCreate returns the existing session or creates a new one if needed.
func (s *QuantoSession) GetOrCreate() *QuantoSession {
	return s
}

func (s *QuantoSession) String() string {
	separationIndex := strings.Index(s.ID, "-")
	if separationIndex == -1 {
		return fmt.Sprintf("AppName: %s, Mode: %s", s.AppName, s.Mode)
	}
	id := s.ID[:separationIndex]
	return fmt.Sprintf("AppName: %s @ %s, Mode: %s", s.AppName, id, s.Mode)
}
