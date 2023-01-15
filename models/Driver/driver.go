package driver

import (
	"sync"
	logger "github.com/AdityaSubrahmanyaBhat/golang/dashDB/models/Logger"
)

type Driver struct {
	Dir string
	Mutex   sync.Mutex
	Mutexes map[string]*sync.Mutex
	Log     logger.Logger
}