package app

import (
	"fmt"
	"fw/src/bot/msg"
	"strings"
)

// ============================================================================

func MsgEvt(m msg.Message) string {
	return fmt.Sprintf("msg.%d", m.MsgId())
}

func parse_args(args string) (ret map[string]string) {
	ret = make(map[string]string)

	for _, kv := range strings.Split(args, ",") {
		arr := strings.Split(kv, "=")
		if len(arr) != 2 {
			continue
		}

		k := strings.Trim(arr[0], " ")
		v := strings.Trim(arr[1], " ")

		ret[k] = v
	}

	return
}
