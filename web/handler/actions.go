package handler

import (
	"app/constants"
	"app/model"
	cacheRepo "app/repo/cache/service"
	cloneRepo "app/repo/mongo/clone"
	codeRepo "app/repo/mongo/code"
	"app/utils"
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ActionRes : Chứa info của action
type ActionRes struct {
	MsgErr      string `json:"msg_err,omitempty"`
	ServiceCode string `json:"service_code,omitempty"`
	Action      string `json:"action"`
	Count       int    `json:"count"`
	// Autofarmer data
	Number      int    `json:"number,omitempty"`
	ReRun       int    `json:"re_run,omitempty"`
	Page        string `json:"page,omitempty"`
	Group       string `json:"group,omitempty"`
	OldPassword string `json:"old_password,omitempty"`
	// Autolike data
	FbID           string                `json:"fb_id,omitempty"`
	PhotoID        string                `json:"photo_id,omitempty"`
	PostID         string                `json:"post_id,omitempty"`
	URLService     string                `json:"url_service,omitempty"`
	LinkService    string                `json:"link_service,omitempty"`
	SharePostID    *primitive.ObjectID   `json:"share_post_id,omitempty"`
	ShareComment   string                `json:"share_comment,omitempty"`
	ShareCommentID *primitive.ObjectID   `json:"share_comment_id,omitempty"`
	Comment        string                `json:"comment,omitempty"`
	Max            int                   `json:"max,omitempty"`
	Data           interface{}           `json:"data,omitempty"`
	ViplikeData    *model.ViplikeItem    `json:"viplike_data,omitempty"`
	Kind           constants.ServiceKind `json:"kind,omitempty"`
	ViewTime       int                   `json:"view_time,omitempty"`
	InstagramId    string                `json:"insta_id,omitempty"`
	ProductSearch  string                `json:"product_search,omitempty"`
	ChannelId      string                `json:"channel_id,omitempty"`

	// Instagram
	UID       string `json:"uid,omitempty"`
	TargetUID string `json:"target_uid,omitempty"`
	Keyword   string `json:"keyword,omitempty"`
}

// ActionsHandler : Chứa các api cho mobile gọi: getDoAction, doResult
type ActionsHandler struct{}

//NewActionsHandler : Tạo đối tượng ActionsHandler
func NewActionsHandler() *ActionsHandler {
	return &ActionsHandler{}
}

func (h ActionsHandler) GetDoActions(c echo.Context) (err error) {
	// GetDoActionPublish ..
	type GetDoActionPublish struct {
		CloneID    primitive.ObjectID `json:"clone_id" query:"clone_id" validate:"required"`
		Token      string             `json:"token" query:"token" validate:"required"`
		AppName    string             `json:"app_name" query:"app_name" validate:"required"`
		DeviceInfo model.DeviceInfo   `json:"device_info" query:"device_info"`
	}
	httpCtx := context.Background()
	request := new(GetDoActionPublish)
	if err = c.Bind(request); err != nil {
		return
	}
	if err = c.Validate(request); err != nil {
		return
	}
	// Lấy thông tin clone theo cloneID
	clone, err := getCloneByID(request.CloneID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("Get clone by cloneID %s %s", request.CloneID, err))
	}
	if !clone.IsExists() {
		return fmt.Errorf("Not found clone %s", request.CloneID)
	}
	// Kiểm tra clone này có thuộc về token này hay không ?
	if clone.Token != request.Token {
		return echo.NewHTTPError(http.StatusForbidden, fmt.Errorf("Clone is not belong to you"))
	}

	// TODO: Lấy thông tin action profile theo lazyload
	actionProfile, _ := getActionProfileByID(clone.ActionProfileID)
	if !actionProfile.IsExists() {

		// Kiểm tra lại nếu không có actionProfile, trả lõi
		if !actionProfile.IsExists() {
			return fmt.Errorf("Clone not yet register action profile")
		}
	}

	// Chạy action profile default
	todayActions := actionProfile.ActionDefault
	now := time.Now()
	if len(actionProfile.Actions) > 0 {
		aliveDays := utils.GetDaysDuration(*clone.CreatedDate, now)

		if aliveDays < (len(actionProfile.Actions) - 1) {
			// cập nhật todayActions thành theo action profile ngày
			todayActions = actionProfile.Actions[aliveDays]
		}
	}

	// Tạo ra channel chứa response có len = len(nextActions) + 5 action autolike
	response := make([]ActionRes, 0)
	responseCh := make(chan ActionRes, 100)
	cloneActivitiesCh := make(chan model.CloneActivityInfo, 100)
	wg := new(sync.WaitGroup)

	// 1. Tạo ra 100 codes
	codeCh := make(chan string, 100)

	codes, _ := h.genAutofarmerCode(100)
	for _, code := range codes {
		codeCh <- code
	}

	// TODO: Chạy các action trong action_profile
	for _, action := range todayActions {
		// mỗi action tạo 1 goroutine để chạy độc lâp
		wg.Add(1)
		go func(wg *sync.WaitGroup, action model.Action) {
			defer wg.Done()

			//  Tuỳ theo "action" sẽ lấy data tương ứng
			switch action.Action {

			// Bộ action trả về kết quả liền, không cần query thêm
			case constants.FeedlikeAction,
				constants.WatchAction,
				constants.FeedAction:

				serviceCode := <-codeCh

				responseCh <- ActionRes{
					ServiceCode: serviceCode,
					Action:      action.Action.String(),
					Count:       action.Quantity,
					ReRun:       action.ReRun,
				}

			default:
				// case chưa được xử lý
			}
		}(wg, action)
	}
	// TODO: Kiểm tra clone này có chạy các action của: "autolike" hay không ?
	// actionProfile có type: "fblike"
	// user info có isSubLike = true
	if actionProfile.Type == "fblike" {
		settingPrices, errSettingPrices := getSettingPrices(clone.AppName)
		if errSettingPrices != nil {
			err = errSettingPrices
			return
		}

		for _, item := range settingPrices.Settings {
			// mỗi action tạo 1 goroutine để chạy độc lâp
			wg.Add(1)
			go func(wg *sync.WaitGroup, likeSubAction model.SettingPriceItem) {
				defer wg.Done()
				serviceType := likeSubAction.Key
				quantity := likeSubAction.Limit
				for i := 0; i < quantity; i++ {
					service := h.GetRandom(string(serviceType))
					actionRes := ActionRes{
						Action:      serviceType,
						ServiceCode: service["service_code"].(string),
						FbID:        service["fanpage_id"].(string),
						LinkService: service["link_service"].(string),
						Number:      service["number"].(int),
					}
					responseCh <- actionRes
				}
			}(wg, item)
		}
	}

	wg.Wait()
	close(cloneActivitiesCh)
	close(responseCh)
	close(codeCh)

	// TODO: Cập nhật lại clone
	err = cloneRepo.New(httpCtx).Update(clone)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("Update info clone %s", err))
	}

	// update response
	for item := range responseCh {
		response = append(response, item)
	}

	return c.JSON(success(response))
}

// GetRandom : Kiểm tra tồn tại
func (h ActionsHandler) GetRandom(Type string) (service map[string]interface{}) {
	if Type == "follow" {
		keyFollow := "FOLLOW"
		follow := h.GetServiceFollow()
		if follow != "" {
			resArr := strings.Split(follow, keyFollow)
			if len(resArr) > 4 {
				service = bson.M{
					"service_code": resArr[0],
					"fanpage_id":   resArr[1],
					"type":         resArr[2],
					"link_service": resArr[3],
					"number":       resArr[4],
				}
			}
		}
	}
	if Type == "likepage" {
		likepage := h.GetServiceLikepage()
		if likepage != "" {
			resArr := strings.Split(likepage, "LIKEPAGE")
			if len(resArr) >= 5 {
				service = bson.M{
					"service_code": resArr[0],
					"fanpage_id":   resArr[1],
					"type":         resArr[2],
					"link_service": resArr[3],
					"number":       resArr[4],
				}
			}
		}
	}

	if Type == "bufflike" {
		bufflike := h.GetServiceBufflike()
		if bufflike != "" {
			resArr := strings.Split(bufflike, "BUFFLIKE")
			if len(resArr) >= 5 {
				service = bson.M{
					"service_code": resArr[0],
					"fanpage_id":   resArr[1],
					"type":         resArr[2],
					"link_service": resArr[3],
					"number":       resArr[4],
				}
			}
		}
	}
	return service
}

// GetServiceFollow : lấy 1 follow
func (h ActionsHandler) GetServiceFollow() (data string) {
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
func (h ActionsHandler) GetServiceLikepage() (data string) {
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
func (h ActionsHandler) GetServiceBufflike() (data string) {
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

// genAutofarmerCode : Tạo mã autofarmer theo số lượng
func (ActionsHandler) genAutofarmerCode(quantity int) (codes []string, err error) {

	codes, err = codeRepo.New(context.Background()).Generate("autofarmer_code", quantity)
	if err != nil || len(codes) == 0 {
		err = fmt.Errorf("Gen autofarmer code %s", err)
	}

	return
}
