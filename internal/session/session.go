package session

import (
	"fmt"
	"strings"

	"github.com/google/uuid"
)

type Mode string

const (
	Local   Mode = "local"
	Cluster Mode = "cluster"
)

type QuantoSession struct {
	ID      string
	AppName string
	Mode    Mode
}

func New() *QuantoSession {
	return &QuantoSession{
		ID: uuid.NewString(),
	}
}

func (s *QuantoSession) SetAppName(appName string) *QuantoSession {
	s.AppName = appName
	return s
}

func (s *QuantoSession) SetMode(mode string) *QuantoSession {
	s.Mode = map[string]Mode{"local": Local, "cluster": Cluster}[mode]
	fmt.Printf("Running in %s mode\n", mode)
	return s
}

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
