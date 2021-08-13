package dbmgr

import (
	"fw/src/core/db"
	"fw/src/core/log"
	"fw/src/shared/config"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// ============================================================================

type userload_t struct {
	DbName string `bson:"_id"`
	N      int64  `bson:"n"`
}

// ============================================================================

func Center_CreateUserLoad() {
	// create only if userload does NOT exist
	if DBCenter.HasCollection(C_tabname_userload) {
		return
	}

	// create
	var obj userload_t

	for k, _ := range config.Common.DBUser {
		obj.DbName = k
		obj.N = 0

		err := DBCenter.Insert(C_tabname_userload, &obj)
		if err != nil && !db.IsDup(err) {
			log.Error("dbmgr.Center_CreateUserLoad() failed:", err)
		}
	}
}

func Center_AllocUserDB() string {
	var obj []*userload_t

	DBCenter.Execute(func(thedb *mongo.Database) {
		opts := options.Find().SetSort(db.M{"n": 1}).SetLimit(1)
		cursor, err := thedb.Collection(C_tabname_userload).Find(nil, db.D{}, opts)
		if err == nil {
			err = cursor.All(nil, &obj)
		}

		if err != nil {
			log.Error("dbmgr.Center_AllocUserDB() failed:", err)
		}
	})

	if len(obj) > 0 {
		// ok
		return obj[0].DbName
	} else {
		// failed
		return ""
	}
}

func Center_IncUserLoad(dbname string) {
	err := DBCenter.Update(
		C_tabname_userload,
		dbname,
		db.M{"$inc": db.M{"n": 1}},
	)
	if err != nil {
		log.Error("dbmgr.Center_IncUserLoad() failed:", err)
	}
}
