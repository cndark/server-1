package msg

var MsgCreators = map[uint32]func() Message{
    1200: func() Message {
        return &GS_RegisterGame{}
    },
    1201: func() Message {
        return &RT_RegisterGame_R{}
    },
}

func (self *GS_RegisterGame) MsgId() uint32 {
    return 1200
}

func (self *RT_RegisterGame_R) MsgId() uint32 {
    return 1201
}
