package main

import (
	"github.com/murashko91/otus-go/hw12_13_14_15_calendar/internal/config"
	"github.com/murashko91/otus-go/hw12_13_14_15_calendar/internal/storage"
	memorystorage "github.com/murashko91/otus-go/hw12_13_14_15_calendar/internal/storage/memory"
	sqlstorage "github.com/murashko91/otus-go/hw12_13_14_15_calendar/internal/storage/sql"
)

func getStorage(config config.DBConfig) storage.Storage {
	if config.InMemory {
		return memorystorage.New()
	}

	return sqlstorage.New(sqlstorage.StorageInfo{
		Host:     config.Host,
		Port:     config.Port,
		User:     config.User,
		Password: config.Password,
		DBName:   config.DBName,
	})
}
