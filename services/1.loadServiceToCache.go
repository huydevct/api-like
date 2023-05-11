package services

import (
	"app/model"
	cacheRepo "app/repo/cache/service"
	serviceRepo "app/repo/mongo/backup/service"
	"context"
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/panjf2000/ants"
	"github.com/robfig/cron"
	"gopkg.in/mgo.v2/bson"
)

// loadServiceToQueue :
type loadServiceToQueue struct{}

// NewLoadServiceToQueue : tạo mới đối tượng loadServiceToQueue
func NewLoadServiceToQueue() *loadServiceToQueue {
	loadServiceToQueue := loadServiceToQueue{}

	return &loadServiceToQueue
}

func (s *loadServiceToQueue) Run() {
	cfg.Mongo.Get("autolike").Init()
	cfg.Cache.Get("core").Init()
	c := cron.New()
	c.AddFunc("@every 0h01m30s", func() {
		s.LoadDBToCache("follow", "FOLLOW")
	})
	c.AddFunc("@every 0h02m30s", func() {
		s.LoadDBToCache("likepage", "LIKEPAGE")
	})
	c.AddFunc("@every 0h01m30s", func() {
		s.LoadDBToCache("bufflike", "BUFFLIKE")
	})
	c.Start()
}

type dataPage struct {
	Page  int
	Limit int
	Type  string
}

func (s *loadServiceToQueue) LoadDBToCache(TypeService string, slug string) {
	ctx := context.Background()
	limit := 100
	serviceRepoInstance := serviceRepo.New(ctx)
	fmt.Println("------THOI GIAN BAT DAU aaa---------", time.Now())
	timeStart := time.Now()
	totalService, _ := serviceRepoInstance.TotalTypeAll(TypeService)

	totalPage := int(totalService / limit)
	wg := new(sync.WaitGroup)
	serviceLists := make([][]model.ActiveService, 0)
	pool, _ := ants.NewPoolWithFunc(1, func(data interface{}) {
		serviceList := s.GetListService(data)
		serviceLists = append(serviceLists, serviceList)
		wg.Done()
	})

	defer pool.Release()
	for page := 1; page <= totalPage+1; page++ {
		wg.Add(1)
		user := dataPage{Page: page, Limit: limit, Type: TypeService}
		pool.Invoke(user)
	}

	wg.Wait()
	totalAff := 0
	service_follow := bson.M{}
	for _, services := range serviceLists {
		for _, service := range services {
			data := service.ServiceCode + slug + service.FanpageID + slug + string(service.Type) + slug + service.LinkService + slug + strconv.Itoa(service.Number)
			service_follow[service.ServiceCode] = data
			key := TypeService + strconv.Itoa(totalAff)
			cacheRepo.New(ctx).SetCache(key, data)
			totalAff++
		}
	}
	cacheRepo.New(ctx).SetCache("total_"+TypeService, strconv.Itoa(totalAff))
	fmt.Println("total_"+TypeService, strconv.Itoa(totalAff))
	fmt.Println("totalService "+TypeService+" ", totalService)
	fmt.Println("totalService after ", totalAff)
	fmt.Println("------THOI GIAN KET THUC--------", timeStart, time.Now())
}

func (s *loadServiceToQueue) GetListService(data interface{}) (serviceList []model.ActiveService) {
	dataTh := data.(dataPage)
	page := dataTh.Page
	limit := dataTh.Limit
	Type := dataTh.Type
	ctx := context.Background()
	serviceRepoInstance := serviceRepo.New(ctx)
	services, _ := serviceRepoInstance.SearchAll(Type, page, limit)
	serviceList = services
	fmt.Println("số lượng service ", len(services), "limit ", limit, "page ", page)
	return
}
