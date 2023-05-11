package adapter

import "sync"

// Others info
type Others map[string]string

var (
	onceOtherMutex = sync.RWMutex{}
)

//Get ..
func (other Others) Get(name string) (result string) {
	onceOtherMutex.Lock()
	if value, ok := other[name]; ok {
		result = value
	} else {
		onceOtherMutex.Unlock()
		panic("Không tìm thấy config Other " + name)
	}
	onceOtherMutex.Unlock()
	return
}

//Set ..
func (other Others) Set(name string, value string) (result bool) {
	onceOtherMutex.Lock()
	other[name] = value
	onceOtherMutex.Unlock()
	return true
}
