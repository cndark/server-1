package dbmgr

import (
	"fmt"
	"fw/src/core/db"
	"fw/src/core/log"
	"fw/src/shared/config"
)

// ============================================================================

type seqid_t struct {
	Id          int   `bson:"_id"`
	Seq_UserId  int64 `bson:"seq_uid"`
	Seq_GuildId int64 `bson:"seq_gid"`
}

// ============================================================================

func Center_GenUserId(dbname string) string {
	var obj seqid_t

	err := DBCenter.FindAndUpdate(
		C_tabname_seqid,
		db.M{"_id": 1},
		db.M{
			"$inc": db.M{"seq_uid": 1},
		},
		db.M{"seq_uid": 1},
		&obj,
	)
	if err != nil {
		log.Error("dbmgr.Center_GenUserId() failed:", err)
	}

	return fmt.Sprintf("%s-%d-%d", dbname, config.Common.Area.Id, obj.Seq_UserId)
}
