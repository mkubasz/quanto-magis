package main

import (
	"fmt"
	"github.com/google/uuid"
	"strings"
)

type Mode string

const (
	Local   Mode = "local"
	Cluster Mode = "cluster"
)

type QuantoSession struct {
	Id      string
	AppName string
	Mode    Mode
}

func NewQuantoSession() *QuantoSession {
	return &QuantoSession{
		Id: uuid.NewString(),
	}
}

func (s *QuantoSession) SetAppName(appName string) *QuantoSession {
	s.AppName = appName
	return s
}

func (s *QuantoSession) master(mode string) *QuantoSession {
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

func (s *QuantoSession) getOrCreate() *QuantoSession {
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
