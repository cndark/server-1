package ladder

import (
	"fw/src/game/app/modules/calendar"
)

// ============================================================================

const (
	NAME = "ladder"
)

func init() {
	calendar.Register(&calendar.Reg{
		Name: NAME,

		OnStage: on_stage,

		StageFunc: map[string]func(){
			"robot1":  stage_robot1,
			"robot2":  stage_robot2,
			"prepare": stage_prepare,
			"start":   stage_start,
			"reward":  stage_reward,
			"close":   stage_close,
		},
	})
}
