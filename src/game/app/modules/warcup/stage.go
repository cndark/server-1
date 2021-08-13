package warcup

import (
	"fw/src/game/app/modules/utils"
	"fw/src/game/msg"
	"time"
)

// ============================================================================
// 阶段

const (
	c_Stage_Close = iota
	c_Stage_Open
	c_Stage_Audition
	c_Stage_Top64
	c_Stage_Top8
	c_Stage_End
)

// ============================================================================

func broadcast_stage() {
	res := &msg.GS_WarCupStageUpdate{
		Stage: int32(g_stage),
		Ts2:   g_t2.Unix(),
	}

	res.Round, res.RoundPiece, res.PieceEndTs = round_piece()

	// notify
	utils.BroadcastPlayers(res)
}

func on_stage(stg string, t0, t1, t2 time.Time, f func(bool)) {
	if stg == "open" {
		force_close = false
	}

	g_t0 = t0
	g_t1 = t1
	g_t2 = t2

	f(!force_close)
}

func stage_open() {
	warcup_plrs = make(map[string]*warcup_plr_t)
	warcup_vsseq = 0
	warcup_vsdata = make(map[int32]*vs_data_t)
	warcup_plrguess = make(map[string]*plr_guess_data_t)

	g_stage = c_Stage_Open
	broadcast_stage()

	// notify mail
	send_open_mail()
}

func stage_audition() {
	g_stage = c_Stage_Audition
	new_audition()
}

func stage_top64() {
	g_stage = c_Stage_Top64
	new_top64()
}

func stage_top8() {
	g_stage = c_Stage_Top8
	new_top8()
}

func stage_end() {
	g_stage = c_Stage_End
	broadcast_stage()

	end_award()
}

func stage_close() {
	force_close = true

	g_stage = c_Stage_Close
	broadcast_stage()
}
