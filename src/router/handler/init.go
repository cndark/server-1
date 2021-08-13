package handler

import (
    "fw/src/router/handler/gs_rt"
    "fw/src/router/msg"
)

func Init() {
    msg.Handler(1200, gs_rt.GS_RegisterGame)
}
