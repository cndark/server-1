package handler

import (
    "fw/src/gate/handler/c_gw"
    "fw/src/gate/handler/gw_gs"
    "fw/src/gate/msg"
)

func Init() {
    msg.Handler(1000, c_gw.C_Auth)
    msg.Handler(1002, c_gw.C_Login)
    msg.Handler(1004, c_gw.C_TokenGet)
    msg.Handler(1006, c_gw.C_TokenAuth)
    msg.Handler(1101, gw_gs.GS_RegisterGate_R)
    msg.Handler(1104, gw_gs.GS_Kick)
}
