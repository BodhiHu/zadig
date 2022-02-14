package mongodb

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/koderover/zadig/pkg/config"
	"github.com/koderover/zadig/pkg/microservice/aslan/core/common/repository/models"
	"github.com/koderover/zadig/pkg/setting"
	mongotool "github.com/koderover/zadig/pkg/tool/mongo"
)

type SystemSettingColl struct {
	*mongo.Collection

	coll string
}

func NewSystemSettingColl() *SystemSettingColl {
	name := models.SystemSetting{}.TableName()
	return &SystemSettingColl{Collection: mongotool.Database(config.MongoDatabase()).Collection(name), coll: name}
}

func (c *SystemSettingColl) GetCollectionName() string {
	return c.coll
}

func (c *SystemSettingColl) EnsureIndex(ctx context.Context) error {
	return nil
}

func (c *SystemSettingColl) Get() (*models.SystemSetting, error) {
	query := bson.M{}
	resp := &models.SystemSetting{}

	err := c.FindOne(context.TODO(), query).Decode(resp)
	return resp, err
}

func (c *SystemSettingColl) InitSystemSettings() error {
	_, err := c.Get()
	// if we didn't find anything
	if err != nil {
		return c.Create(setting.LocalClusterID, &models.SystemSetting{
			WorkflowConcurrency: 2,
			BuildConcurrency:    5,
		})
	}
	return nil
}

func (c *SystemSettingColl) Create(id string, args *models.SystemSetting) error {
	var objectID primitive.ObjectID
	if id != "" {
		objectID, _ = primitive.ObjectIDFromHex(id)
	} else {
		objectID = primitive.NewObjectID()
	}

	args.UpdateTime = time.Now().Unix()

	query := bson.M{"_id": objectID}
	change := bson.M{"$set": args}
	_, err := c.UpdateOne(context.TODO(), query, change, options.Update().SetUpsert(true))
	return err
}
