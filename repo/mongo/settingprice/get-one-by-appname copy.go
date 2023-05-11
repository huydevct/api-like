package settingprice

import (
	"app/model"
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// GetOneByAppname ..
func (r Repo) GetOneByAppname(appname string) (result model.SettingPrice, err error) {

	ctx := context.Background()

	collection := r.Session.GetCollection(r.Collection)

	condition := bson.M{}
	condition["appname"] = appname

	err = collection.FindOne(ctx, condition).Decode(&result)
	if err == mongo.ErrNoDocuments {
		err = nil
	}

	return
}
