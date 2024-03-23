package main

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
	Read    Read
	Id      string
	AppName string
	Mode    Mode
}

func NewQuantoSession() *QuantoSession {
	return &QuantoSession{
		Id:   uuid.NewString(),
		Read: Read{},
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
	separationIndex := strings.Index(s.Id, "-")
	if separationIndex == -1 {
		return fmt.Sprintf("AppName: %s, Mode: %s", s.AppName, s.Mode)
	}
	id := s.Id[:separationIndex]
	return fmt.Sprintf("AppName: %s @ %s, Mode: %s", s.AppName, id, s.Mode)
}
