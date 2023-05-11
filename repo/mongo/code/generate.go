package code

import (
	"context"
	"fmt"

	"app/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Generate func
func (r Repo) Generate(id string, qtt int) (codes []string, err error) {

	ctx := context.Background()

	collection := r.Session.GetCollection(r.Collection)
	result := collection.FindOneAndUpdate(ctx,
		bson.M{"_id": id},
		bson.M{"$inc": bson.M{"num": qtt}},
		options.FindOneAndUpdate().SetUpsert(true))
	if result.Err() != nil {
		err = result.Err()
		return
	}

	doc := model.Code{}
	err = result.Decode(&doc)

	if err != nil {
		return
	}

	cache := make(map[int]string)

	for i := 0; i < qtt; i++ {
		code := generate(doc.Num+qtt-i, cache)
		codes = append(codes, code)
	}

	return
}

func generate(num int, cache map[int]string) (result string) {
	if str, ok := cache[num]; ok {
		result = str
	} else {

		char := "SKWGHXA07DU2Z51C9V68IQLFNTPBO34YJRME"
		if num > 0 {
			mod := num % len(char)
			newNum := num / len(char)
			c := string(char[mod])
			result = fmt.Sprintf("%v%v%v", generate(newNum, cache), c, result)
		}
		cache[num] = result
	}

	return result
}
