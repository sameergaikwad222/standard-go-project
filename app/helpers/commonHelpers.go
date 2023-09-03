package helpers

import (
	"encoding/json"
	"log"

	"go.mongodb.org/mongo-driver/bson"
)

func GetBsonMapFromMapDataType(mapType map[string]interface{}) (bson.M, error) {
	var bsonMap = make(bson.M)

	// Converting Map to Raw Bson First
	rawBsonData, err := bson.Marshal(mapType)
	if err != nil {
		log.Fatal("Error while marshalling json to raw bson")
		return nil, err
	}

	//Converting raw bson to bson Map
	err = bson.Unmarshal(rawBsonData, &bsonMap)
	if err != nil {
		log.Fatal("error while unmarshalling raw bson to bson map")
		return nil, err
	}
	return bsonMap, nil
}

func GetBsonMapFromJsonString(jsonString string) (bson.M, error) {
	var filterMap = make(map[string]interface{})

	//Converting Json string to Map data Type
	err := json.Unmarshal([]byte(jsonString), &filterMap)
	if err != nil {
		log.Fatal("error while unmarshalling json string to map")
		return nil, err
	}
	bsonMap, err := GetBsonMapFromMapDataType(filterMap)
	if err != nil {
		log.Fatal("error while unmarshalling json string to map")
		return nil, err
	}
	return bsonMap, nil
}
