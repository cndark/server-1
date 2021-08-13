package utils

import (
	"fw/src/core"
	"fw/src/core/sched/loop"
	"fw/src/game/msg"
	Err "fw/src/proto/errorcode"
	"fw/src/shared/config"
	"reflect"
	"sync/atomic"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

// ============================================================================

const (
	c_GP_Timeout = 10 // in seconds
)

// ============================================================================

var (
	gp_seqid     int64
	gp_callbacks = make(map[int64]*gp_callback_t)
)

// ============================================================================

type gp_callback_t struct {
	id          int64
	f           func(int32, interface{})
	ret_proto   interface{}
	tid_timeout *core.Timer
}

// ============================================================================

func GsPush(svrid int32, evt_name string, sarg []string, oarg interface{}) {
	// marshal object arg
	var oarg_data []byte
	var err error
	if oarg != nil {
		oarg_data, err = MarshalArg(oarg)
		if err != nil {
			return
		}
	}

	// send msg
	Send2Game(svrid, &msg.GS_Push{
		EvtName: evt_name,
		SArg:    sarg,
		OArg:    oarg_data,
	})
}

func GsPushAll(evt_name string, sarg []string, oarg interface{}, except_self ...bool) {
	// marshal object arg
	var oarg_data []byte
	var err error
	if oarg != nil {
		oarg_data, err = MarshalArg(oarg)
		if err != nil {
			return
		}
	}

	// send msg
	BroadcastGames(&msg.GS_Push{
		EvtName: evt_name,
		SArg:    sarg,
		OArg:    oarg_data,
	}, except_self...)
}

func GsPull(svrid int32, evt_name string, sarg []string, oarg interface{}, ret_proto interface{}, f func(ec int32, r interface{})) {
	// marshal object arg
	var oarg_data []byte
	var err error
	if oarg != nil {
		oarg_data, err = MarshalArg(oarg)
		if err != nil {
			f(Err.Failed, nil)
			return
		}
	}

	// create callback object
	cb := &gp_callback_t{
		id:        atomic.AddInt64(&gp_seqid, 1),
		f:         f,
		ret_proto: ret_proto,
	}

	// add
	gp_callbacks[cb.id] = cb

	// set timeout
	cb.tid_timeout = loop.SetTimeout(time.Now().Add(time.Second*c_GP_Timeout), func() {
		delete(gp_callbacks, cb.id)
		cb.f(Err.Common_Timeout, nil)
	})

	// send msg
	Send2Game(svrid, &msg.GS_Pull{
		Tk:       cb.id,
		SrcSvrId: config.CurGame.Id,
		EvtName:  evt_name,
		SArg:     sarg,
		OArg:     oarg_data,
	})
}

func GSPull_R(tk int64, ec int32, result []byte) {
	// find callback
	cb := gp_callbacks[tk]
	if cb == nil {
		return
	}

	// remove
	delete(gp_callbacks, cb.id)

	// clear timeout
	loop.CancelTimer(cb.tid_timeout)

	// check error
	if ec != Err.OK {
		cb.f(ec, nil)
		return
	}

	// result
	if result != nil {
		// unmarshal result
		err := UnmarshalArg(result, cb.ret_proto)
		if err != nil {
			cb.f(Err.Failed, nil)
			return
		}

		// keep ptr form only if ret-proto is a struct
		e := reflect.ValueOf(cb.ret_proto).Elem()
		if e.Kind() != reflect.Struct {
			cb.ret_proto = e.Interface()
		}
	}

	// ok
	cb.f(Err.OK, cb.ret_proto)
}

// ============================================================================

func MarshalArg(v interface{}) ([]byte, error) {
	return bson.Marshal(&struct{ V interface{} }{v})
}

func UnmarshalArg(data []byte, v interface{}) error {
	rv, err := bson.Raw(data).LookupErr("v")
	if err != nil {
		return err
	}

	return rv.Unmarshal(v)
}
