package cond

import (
	"fw/src/core/evtmgr"
)

// ============================================================================
// obj值接口

type ICondObj interface {
	GetVal() float64
	SetVal(v float64)
	AddVal(v float64)

	// body:(plr, guild); confid:配置表id; isChange: 值是否有改动
	Done(body interface{}, confid int32, isChange bool)
}

// ============================================================================
// 管理需要监听的模块

type cond_mgr_t struct {
	Module string                          // 注册监听条件的模块
	ConfId int32                           // 表id: 透传给各模块,用于判断进度用
	Cond   int32                           // 条件id
	P1     int32                           // 条件参数类型
	GetObj func(body interface{}) ICondObj // 把各模块的obj对象传进来
}

var cond_mgr = map[int32][]*cond_mgr_t{}

// 各模块注册监听
func RegistCondObj(module string, confid int32, cond int32, p1 int32, f func(body interface{}) ICondObj) {
	if cond == 0 || f == nil {
		return
	}

	cond_mgr[cond] = append(cond_mgr[cond], &cond_mgr_t{
		Module: module,
		ConfId: confid,
		Cond:   cond,
		P1:     p1,
		GetObj: f,
	})
}

// 移除一个模块所有的监听对象
func RemoveModuleCondObjs(module string) {
	for key := range cond_mgr {
		for i := 0; i < len(cond_mgr[key]); i++ {
			if cond_mgr[key][i].Module == module {
				cond_mgr[key] = append(cond_mgr[key][:i], cond_mgr[key][i+1:]...)
				i--
			}
		}
	}
}

// ============================================================================
// event监听on

func init() {
	for cid, cb := range cond_cb {
		cid, cb := cid, cb

		evtmgr.On(cb.Evt, func(args ...interface{}) {
			body := args[0] // (body = plr, guild)

			for _, o := range cond_mgr[cid] {
				obj := o.GetObj(body)
				if obj == nil {
					continue
				}

				oldVal := obj.GetVal()
				cb.F(o.P1, obj, args)

				isChange := oldVal != obj.GetVal()
				obj.Done(body, o.ConfId, isChange)
			}
		})
	}
}
