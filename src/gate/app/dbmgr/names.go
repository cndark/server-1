package dbmgr

import (
	"fw/src/core/db"
)

// ============================================================================

// _id:  objectid
// name: string (uk)

// ============================================================================

func Center_InsertName(name string) bool {
	var rec struct {
		Name string `bson:"name"`
	}

	rec.Name = name

	err := DBCenter.Insert(C_tabname_names, &rec)
	return err == nil
}

func Center_ChangeName(oldname, newname string) bool {
	err := DBCenter.UpdateByCond(
		C_tabname_names,
		db.M{"name": oldname},
		db.M{"$set": db.M{"name": newname}},
	)
	return err == nil
}
