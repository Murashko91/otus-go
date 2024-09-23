package storage

import (
	"fmt"
)

type ErrDateBusy struct {
	event Event
}

func (e ErrDateBusy) Error() string {
	return fmt.Sprintf("%v - %v date busy with another event", e.event.StartDate, e.event.EndDate)
}

func NewErrDateBusy(event Event) ErrDateBusy {
	return ErrDateBusy{event: event}
}

type ErrMissedReqiredData struct {
	fields []string
}

func (e ErrMissedReqiredData) Error() string {
	return fmt.Sprintf("missed reqired fields: %v", e.fields)
}

func NewErrMissedReqiredData(fields []string) ErrMissedReqiredData {
	return ErrMissedReqiredData{fields: fields}
}

type WrongEventDatesError struct{}

func (e WrongEventDatesError) Error() string {
	return "end date should not be before start date"
}

func NewWrongEventDatesError() WrongEventDatesError {
	return WrongEventDatesError{}
}
