package services

import (
	"app/common/config"
	"app/constants"
	"app/model"
	settingPriceRepo "app/repo/mongo/settingprice"
	"context"
	"fmt"
	"hash/fnv"
	"math"
	"sort"
)

var (
	cfg = config.GetConfig()
)

func hashToInt(text string) (result int) {
	h := fnv.New64a()
	h.Write([]byte(text))
	result = int(h.Sum64())
	result = int(math.Abs(float64(result)))

	if result > 1000 {
		result = result / 1000
	}

	return
}

func updateNumberRestByType(service model.Service, numberRestReal int, appname constants.AppName) (result int) {
	result = numberRestReal
	var numberRests []model.UpdateNumberRest

	settingPrices, err := settingPriceRepo.New(context.Background()).GetOneByAppname(appname.String())
	if err != nil {
		err = fmt.Errorf("Get setting prices %s", err)
		return
	}

	// get update number rest by type
	{
		for _, item := range settingPrices.Settings {
			if item.Key == service.Type.String() {
				numberRests = item.UpdateNumberRest
			}
		}
	}

	// set number rest by number
	{
		sort.Sort(model.ByNumber(numberRests))

		for _, numberRest := range numberRests {
			if service.Number < numberRest.Number {
				if numberRestReal >= numberRest.NumberRest {
					result = numberRest.NumberRest
				}
				break
			}
		}
	}

	return
}
