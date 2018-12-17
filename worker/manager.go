package worker

import (
	"sync"
	"time"
)

var _managerIsEnable = false
var _managerRWMutex = new(sync.RWMutex)
var _managerWorkers = []func(){test, test2, test2}

func test() {
	println("HI!")
	time.Sleep(time.Second)
}

func test2() {
	println("YO!")
	time.Sleep(2 * time.Second)
}

func _managerRunWorkers(worker func()) {
	for {
		_managerRWMutex.RLock()
		isEnable := _managerIsEnable
		_managerRWMutex.RUnlock()
		if isEnable {
			worker()
		} else {
			break
		}
	}
}

// Start start workers
func Start() {

	println("worker starter!")

	// set flag to true
	_managerRWMutex.RLock()
	isEnable := _managerIsEnable
	_managerRWMutex.RUnlock()
	if !isEnable {
		_managerRWMutex.RLock()
		_managerIsEnable = true
		_managerRWMutex.RUnlock()
	}

	// start go routin
	go func() {
		for _, worker := range _managerWorkers {
			go _managerRunWorkers(worker)
		}
	}()
}

// Stop stop workers
func Stop() {

	println("worker stopper!")

	// set flag to false
	_managerRWMutex.RLock()
	isEnable := _managerIsEnable
	_managerRWMutex.RUnlock()
	if isEnable {
		_managerRWMutex.RLock()
		_managerIsEnable = false
		_managerRWMutex.RUnlock()
	}
}
