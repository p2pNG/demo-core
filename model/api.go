package model

import "fmt"

type StandardError struct {
	Code     int
	Message  string
	Internal string
}

func (s *StandardError) Error() string {
	return fmt.Sprintf("[%d] %s; %s", s.Code, s.Message, s.Internal)
}
