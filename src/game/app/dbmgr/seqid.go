package dbmgr

import (
	"fmt"
	"fw/src/core/db"
	"fw/src/core/log"
)

// ============================================================================

type seqid_t struct {
	Id          int   `bson:"_id"`
	Seq_UserId  int64 `bson:"seq_uid"`
	Seq_GuildId int64 `bson:"seq_gid"`
}

// ============================================================================

func Center_CreateSeqId() {
	if DBCenter.HasCollection(C_tabname_seqid) {
		return
	}

	var obj seqid_t

	obj.Id = 1
	obj.Seq_UserId = 999999
	obj.Seq_GuildId = 999999

	err := DBCenter.Insert(C_tabname_seqid, &obj)
	if err != nil {
		log.Error("dbmgr.Center_CreateSeqId() failed:", err)
	}
}

func Center_GenGuildId(dbname string) string {
	var obj seqid_t

	err := DBCenter.FindAndUpdate(
		C_tabname_seqid,
		db.M{"_id": 1},
		db.M{
			"$inc": db.M{"seq_gid": 1},
		},
		db.M{"seq_gid": 1},
		&obj,
	)
	if err != nil {
		log.Error("dbmgr.Center_GenGuildId() failed:", err)
	}

	return fmt.Sprintf("g%s-%d", dbname, obj.Seq_GuildId)
}
