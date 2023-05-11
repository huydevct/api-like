package handler

import (
	cacheRepo "app/repo/cache/service"
	serviceRepo "app/repo/mongo/backup/service"
	"app/utils"
	"context"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
)

var indexFollow = 0
var indexLikepage = 0
var indexBufflike = 0

// ServiceHandler : struct init Service
type ServiceHandler struct{}

// NewServiceHandler : Tạo mới 1 đối tượng Service handler
func NewServiceHandler() *ServiceHandler {
	return &ServiceHandler{}
}

// GetRandom : Kiểm tra tồn tại
func (h ServiceHandler) GetRandom(c echo.Context) (err error) {
	//
	services := make([]map[string]interface{}, 0)
	keyFollow := "FOLLOW"
	follow := h.GetServiceFollow()
	if follow != "" {
		resArr := strings.Split(follow, keyFollow)
		if len(resArr) > 4 {
			if err == nil {
				modelService := bson.M{
					"service_code": resArr[0],
					"fanpage_id":   resArr[1],
					"type":         resArr[2],
					"link_service": resArr[3],
					"number":       resArr[4],
				}
				services = append(services, modelService)
			}
		}
	}
	likepage := h.GetServiceLikepage()
	if likepage != "" {
		resArr := strings.Split(likepage, "LIKEPAGE")
		if len(resArr) >= 5 {
			modelService := bson.M{
				"service_code": resArr[0],
				"fanpage_id":   resArr[1],
				"type":         resArr[2],
				"link_service": resArr[3],
				"number":       resArr[4],
			}
			services = append(services, modelService)
		}
	}

	bufflike := h.GetServiceBufflike()
	if bufflike != "" {
		resArr := strings.Split(bufflike, "BUFFLIKE")
		if len(resArr) >= 5 {
			modelService := bson.M{
				"service_code": resArr[0],
				"fanpage_id":   resArr[1],
				"type":         resArr[2],
				"link_service": resArr[3],
				"number":       resArr[4],
			}
			services = append(services, modelService)
		}
	}

	return c.JSON(success(services))
}

// GetServiceFollow : lấy 1 follow
func (h ServiceHandler) GetServiceFollow() (data string) {
	httpCtx := context.Background()
	key := "follow" + strconv.Itoa(indexFollow)
	res, _ := cacheRepo.New(httpCtx).GetCacheBy(key)
	total, _ := cacheRepo.New(httpCtx).GetCacheBy("total_follow")
	indexFollow = indexFollow + 1

	if utils.ConvertToInt(total) <= indexFollow {
		indexFollow = 0
	}
	if res != "" {
		resArr := strings.Split(res, "FOLLOW")
		if len(resArr) >= 5 {
			NumberRest := utils.ConvertToInt(resArr[4])
			if NumberRest > 0 {
				NumberRest = NumberRest - 1
				data = resArr[0] + "FOLLOW" + resArr[1] + "FOLLOW" + resArr[2] + "FOLLOW" + resArr[3] + "FOLLOW" + strconv.Itoa(NumberRest)
				cacheRepo.New(httpCtx).SetCache(key, data)
				return
			}
		}
	}
	return
}

// GetServiceLikepage : lấy 1 likepage
func (h ServiceHandler) GetServiceLikepage() (data string) {
	httpCtx := context.Background()
	key := "likepage" + strconv.Itoa(indexLikepage)
	res, _ := cacheRepo.New(httpCtx).GetCacheBy(key)
	total, _ := cacheRepo.New(httpCtx).GetCacheBy("total_likepage")
	indexLikepage = indexLikepage + 1

	if utils.ConvertToInt(total) <= indexLikepage {
		indexLikepage = 0
	}
	if res != "" {
		resArr := strings.Split(res, "LIKEPAGE")
		if len(resArr) >= 5 {
			NumberRest := utils.ConvertToInt(resArr[4])
			if NumberRest > 0 {
				NumberRest = NumberRest - 1
				data = resArr[0] + "LIKEPAGE" + resArr[1] + "LIKEPAGE" + resArr[2] + "LIKEPAGE" + resArr[3] + "LIKEPAGE" + strconv.Itoa(NumberRest)
				cacheRepo.New(httpCtx).SetCache(key, data)
				return
			}
		}
	}
	return
}

// GetServiceBufflike : lấy 1 bufflike
func (h ServiceHandler) GetServiceBufflike() (data string) {
	httpCtx := context.Background()
	key := "bufflike" + strconv.Itoa(indexBufflike)
	res, _ := cacheRepo.New(httpCtx).GetCacheBy(key)
	total, _ := cacheRepo.New(httpCtx).GetCacheBy("total_bufflike")
	indexBufflike = indexBufflike + 1

	if utils.ConvertToInt(total) <= indexBufflike {
		indexBufflike = 0
	}
	if res != "" {
		resArr := strings.Split(res, "BUFFLIKE")
		if len(resArr) >= 5 {
			NumberRest := utils.ConvertToInt(resArr[4])
			if NumberRest > 0 {
				NumberRest = NumberRest - 1
				data = resArr[0] + "BUFFLIKE" + resArr[1] + "BUFFLIKE" + resArr[2] + "BUFFLIKE" + resArr[3] + "BUFFLIKE" + strconv.Itoa(NumberRest)
				cacheRepo.New(httpCtx).SetCache(key, data)
				return
			}
		}
	}
	return
}

// Update : cập nhật
func (h ServiceHandler) Update(c echo.Context) (err error) {
	//
	type myRequest struct {
		ServiceCode string `json:"service_code" query:"service_code" validate:"required"`
		Type        string `json:"type" query:"type" validate:"required"`
	}
	request := new(myRequest)
	if err = c.Bind(request); err != nil {
		return
	}
	if err = c.Validate(request); err != nil {
		return
	}
	httpCtx := c.Request().Context()
	serviceRepo.New(httpCtx).IncNumberSuccess(request.ServiceCode)
	return c.JSON(success(true))
}
