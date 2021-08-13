package msg

var MsgCreators = map[uint32]func() Message{
    1000: func() Message {
        return &C_Auth{}
    },
    1001: func() Message {
        return &GW_Auth_R{}
    },
    1002: func() Message {
        return &C_Login{}
    },
    1003: func() Message {
        return &GW_Login_R{}
    },
    1004: func() Message {
        return &C_TokenGet{}
    },
    1005: func() Message {
        return &GW_TokenGet_R{}
    },
    1006: func() Message {
        return &C_TokenAuth{}
    },
    1007: func() Message {
        return &GW_TokenAuth_R{}
    },
    1100: func() Message {
        return &GW_RegisterGate{}
    },
    1101: func() Message {
        return &GS_RegisterGate_R{}
    },
    1102: func() Message {
        return &GW_UserOnline{}
    },
    1103: func() Message {
        return &GW_LogoutPlayer{}
    },
    1104: func() Message {
        return &GS_Kick{}
    },
}

func (self *C_Auth) MsgId() uint32 {
    return 1000
}

func (self *GW_Auth_R) MsgId() uint32 {
    return 1001
}

func (self *C_Login) MsgId() uint32 {
    return 1002
}

func (self *GW_Login_R) MsgId() uint32 {
    return 1003
}

func (self *C_TokenGet) MsgId() uint32 {
    return 1004
}

func (self *GW_TokenGet_R) MsgId() uint32 {
    return 1005
}

func (self *C_TokenAuth) MsgId() uint32 {
    return 1006
}

func (self *GW_TokenAuth_R) MsgId() uint32 {
    return 1007
}

func (self *GW_RegisterGate) MsgId() uint32 {
    return 1100
}

func (self *GS_RegisterGate_R) MsgId() uint32 {
    return 1101
}

func (self *GW_UserOnline) MsgId() uint32 {
    return 1102
}

func (self *GW_LogoutPlayer) MsgId() uint32 {
    return 1103
}

func (self *GS_Kick) MsgId() uint32 {
    return 1104
}
