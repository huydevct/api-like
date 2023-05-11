package device

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
)

// GetDeviceID ..
func (r Repo) GetDeviceID() (macAddress []string, err error) {

	ctx := context.Background()
	collection := r.Session.GetCollection(r.Collection)
	// set filter
	filter := bson.M{}
	values, err := collection.Distinct(ctx, "MacAddress", filter)
	if err != nil {
		return
	}

	for _, value := range values {
		if temp, ok := value.(string); ok {
			macAddress = append(macAddress, temp)
		} else {
			err = fmt.Errorf("Get macAddress parse to string")
			return
		}
	}

	return
}
