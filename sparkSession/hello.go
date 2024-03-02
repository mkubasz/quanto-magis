package main

import (
	"fmt"
	"strings"

	"github.com/google/uuid"
)

const (
	// Local mode
	Local = "local"
	// Cluster mode
	Cluster = "cluster"
)

type QuantumSession struct {
	Id      string
	AppName string
	Mode    string
}

func NewQuantumSession() *QuantumSession {
	return &QuantumSession{}
}

func (s *QuantumSession) SetAppName(appName string) *QuantumSession {
	s.AppName = appName
	return s
}

func (s *QuantumSession) master(mode string) *QuantumSession {
	switch mode {
	case "local":
		fmt.Println("Running in local mode")
		s.Mode = Local
	case "cluster":
		fmt.Println("Running in cluster mode")
		s.Mode = Cluster
	}
	return s
}

func (s *QuantumSession) getOrCreate() *QuantumSession {
	s.Id = uuid.NewString()
	return s
}

func (s *QuantumSession) String() string {
	separationIndex := strings.Index(s.Id, "-")
	id := s.Id[:separationIndex]
	return fmt.Sprintf("AppName: %s @ %s, Mode: %s", s.AppName, id, s.Mode)
}

func main() {
	session := NewQuantumSession().
		SetAppName("Quantum Session").
		master("local").
		getOrCreate()
	fmt.Println(session)
}
