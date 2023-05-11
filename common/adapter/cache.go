package adapter

import (
	"log"
	"sync"
	"time"

	"github.com/akrylysov/pogreb"
)

// Caches ..
type Caches map[string]*Cache

// Get Cache
func (adapters Caches) Get(name string) (result *Cache) {
	if adapter, ok := adapters[name]; ok {
		result = adapter
	} else {
		panic("Not found config Cache " + name)
	}
	return
}

// Cache struct
type Cache struct {
	Name       string        `mapstructure:"name"`
	Address    []string      `mapstructure:"address"`
	RepSetName *string       `mapstructure:"repset_name"`
	DBAuthen   string        `mapstructure:"dbauthen"`
	User       string        `mapstructure:"user"`
	Pass       string        `mapstructure:"pass"`
	Timeout    time.Duration `mapstructure:"timeout"`
	IsSSL      bool          `mapstructure:"is_ssl"`
	DBName     string        `mapstructure:"dbname"`
	PoolLimit  *uint64       `mapstructure:"pool_limit"`
	ReadPref   string        `mapstructure:"read_pref"`
	ConClient  *pogreb.DB
}

var (
	onceCache      map[string]*sync.Once = make(map[string]*sync.Once)
	onceCacheMutex                       = sync.RWMutex{}
)

// Init func
func (config *Cache) Init() {
	onceCacheMutex.Lock()

	if onceCache[config.Name] == nil {
		onceCache[config.Name] = &sync.Once{}
	}

	var connectError error
	onceCache[config.Name].Do(func() {
		log.Printf("[%s][%s] Cache V2 [connecting]\n", config.Name, config.Address)
		db, err := pogreb.Open("cache", nil)
		if err != nil {
			log.Fatal(err)
			return
		}
		config.ConClient = db
		// defer db.Close()
	})

	onceCacheMutex.Unlock()

	if connectError != nil {
		log.Printf("[%s][%s] Cache V2[error]: %v \n", config.Name, config.Address, connectError)
		time.Sleep(1 * time.Second)
		onceCache[config.Name] = &sync.Once{}
		config.Init()
		return
	}
}
