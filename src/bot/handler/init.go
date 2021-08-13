package handler

import (
    "fw/src/bot/handler/c_gs"
    "fw/src/bot/handler/c_gw"
    "fw/src/bot/msg"
)

func Init() {
    msg.Handler(1001, c_gw.GW_Auth_R)
    msg.Handler(1003, c_gw.GW_Login_R)
    msg.Handler(1005, c_gw.GW_TokenGet_R)
    msg.Handler(1007, c_gw.GW_TokenAuth_R)
    msg.Handler(5000, c_gs.GS_LoginError)
    msg.Handler(5001, c_gs.GS_UserInfo)
}
