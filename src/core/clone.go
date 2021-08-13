package core

import (
	"fw/src/core/log"

	"go.mongodb.org/mongo-driver/bson"
)

// ============================================================================

func CloneBsonObject(obj interface{}) interface{} {
	buf, err := bson.Marshal(obj)
	if err != nil {
		log.Error("marshal object data failed:", err)
		return nil
	}

	out := make(map[string]interface{})
	err = bson.Unmarshal(buf, &out)
	if err != nil {
		log.Error("unmarshal object data failed:", err)
		return nil
	}

	return out
}

func CloneBsonArray(arr interface{}) interface{} {
	type obj_t struct {
		A interface{}
	}
	obj := &obj_t{arr}

	buf, err := bson.Marshal(obj)
	if err != nil {
		log.Error("marshal object data failed:", err)
		return nil
	}

	var out obj_t
	err = bson.Unmarshal(buf, &out)
	if err != nil {
		log.Error("unmarshal object data failed:", err)
		return nil
	}

	return out.A
}
