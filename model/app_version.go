package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AppVersion struct {
	ID              primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	AppName         string             `json:"app_name" bson:"app_name, omitempty"`
	VersionCode     string             `json:"version_code" bson:"version_code,omitempty"`
	VersionName     string             `json:"version_name" bson:"version_name,omitempty"`
	AppCount        int                `json:"app_count" bson:"app_count, omitempty"`
	AppCollectionID primitive.ObjectID `json:"app_collection_id" bson:"app_collection_id,omitempty"`
	CreatedDate     *time.Time         `json:"created_date,omitempty" bson:"created_date,omitempty"`
	UpdatedDate     *time.Time         `json:"updated_date,omitempty" bson:"updated_date,omitempty"`
	Order           int                `json:"order" bson:"order"`
}

type AppVersionCollection struct {
	ID                primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	AppCollectionName string             `json:"app_collection_name" bson:"app_collection_name, omitempty"`
	CreatedDate       *time.Time         `json:"created_date,omitempty" bson:"created_date,omitempty"`
	UpdatedDate       *time.Time         `json:"updated_date,omitempty" bson:"updated_date,omitempty"`
	Order             int                `json:"order" bson:"order"`

	Apps []AppVersion `json:"apps" bson:"apps, omitempty"` // ko lưu db, lookup to appversion
}
