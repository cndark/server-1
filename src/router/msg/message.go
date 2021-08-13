package msg

type Message interface {
    MsgId() uint32
    Marshal() ([]byte, error)
    Unmarshal([]byte) error
}

func Marshal(m Message) ([]byte, error) {
    return m.Marshal()
}

func Unmarshal(b []byte, m Message) error {
    return m.Unmarshal(b)
}
