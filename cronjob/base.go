package cronjob

import (
	"app/common/config"
	"hash/fnv"
	"math"
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
