package main

import (
	"github.com/murashko91/otus-go/hw12_13_14_15_calendar/internal/config"
	"github.com/murashko91/otus-go/hw12_13_14_15_calendar/internal/storage"
	senderdb "github.com/murashko91/otus-go/hw12_13_14_15_calendar/internal/storage/sender_storage"
)

func getStorage(config config.DBConfig) storage.SenderStorage {
	return senderdb.New(storage.Info{
		Host:     config.Host,
		Port:     config.Port,
		User:     config.User,
		Password: config.Password,
		DBName:   config.DBName,
	})
}
