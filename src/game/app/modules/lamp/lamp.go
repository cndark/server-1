package lamp

import (
	"fw/src/core/evtmgr"
	"fw/src/core/log"
	"fw/src/game/app/gconst"
	"fw/src/game/app/modules/utils"
	"fw/src/game/msg"
	"time"
)

// ============================================================================
const (
	C_Lamp_Type_Local = 0 // 本服
	C_Lamp_Type_Cross = 1 // 跨服

)

// ============================================================================

type lamp_one_t struct {
	Id    int32
	Param map[string]string
	Ts    time.Time
}

// ============================================================================

func init() {
	//remote
	evtmgr.On(gconst.Evt_GsPush_LampCross, func(args ...interface{}) {
		oarg := args[1].([]byte)

		// unmarshal object arg
		var res *lamp_one_t
		err := utils.UnmarshalArg(oarg, &res)
		if err != nil {
			log.Error("unmarshal lamp failed:", err)
			return
		}

		add_list(res)
	})
}

// ============================================================================

func Init() {
	handler_evt_on()
}

func ClearAll() {
}

func ClearSys() {
}

func ClearUser() {
}

func add_list(one *lamp_one_t) {
	utils.BroadcastPlayers(&msg.GS_LampMsg{
		One: &msg.LampOne{
			Id:    one.Id,
			Param: one.Param,
			Ts:    one.Ts.Unix(),
		},
	})
}

func Add(tp int32, id int32, param map[string]string) {
	one := &lamp_one_t{
		Id:    id,
		Param: param,
		Ts:    time.Now(),
	}

	switch int(tp) {
	case C_Lamp_Type_Cross:
		if id != 1 {
			utils.GsPushAll(gconst.Evt_GsPush_LampCross, nil, one)
		}

	case C_Lamp_Type_Local:
		add_list(one)

	default:
	}
}

// ============================================================================

func Lamp_ToMsg() *msg.LampData {
	return nil
}
