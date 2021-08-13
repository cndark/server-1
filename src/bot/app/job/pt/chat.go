package pt

import (
	"fw/src/bot/app"
	"fw/src/bot/msg"
	"fw/src/core"
	"fw/src/core/evtmgr"
	"io/ioutil"
	"math/rand"
	"regexp"
	"strings"
)

// ============================================================================

var (
	v_chat_gen *words_chain_t
)

// ============================================================================

type words_chain_t struct {
	chain    map[string][]string
	prefix_n int
}

type prefix_t []string

// ============================================================================

func (self prefix_t) Shift(v string) {
	copy(self, self[1:])
	self[len(self)-1] = v
}

func (self prefix_t) String() string {
	return strings.Join(self, " ")
}

// ============================================================================

func new_words_chain(n int) *words_chain_t {
	return &words_chain_t{
		chain:    make(map[string][]string),
		prefix_n: n,
	}
}

func (self *words_chain_t) Build(fn string) *words_chain_t {
	d, err := ioutil.ReadFile(fn)
	if err != nil {
		core.Panic("building chat-gen failed:", err)
	}

	re := regexp.MustCompile(`[\r\n;.!?；。！？]+`)
	for _, v := range re.Split(string(d), -1) {
		self.build_one(v)
	}

	return self
}

func (self *words_chain_t) build_one(v string) {
	p := make(prefix_t, self.prefix_n)
	for _, w := range v {
		k := p.String()
		self.chain[k] = append(self.chain[k], string(w))
		p.Shift(string(w))
	}
}

func (self *words_chain_t) Gen(n int) string {
	var ret strings.Builder

	p := make(prefix_t, self.prefix_n)
	for {
		k := p.String()
		arr := self.chain[k]
		L := len(arr)
		if L == 0 {
			break
		}

		next := arr[rand.Intn(L)]
		p.Shift(next)

		ret.WriteString(next)

		n--
		if n <= 0 {
			break
		}
	}

	return ret.String()
}

// ============================================================================

func init() {
	v_chat_gen = new_words_chain(3).Build("chat.txt")

	evtmgr.On("userinfo", func(args ...interface{}) {
		bot := args[0].(*app.Bot)

		bot.JobAdd("pt_chat_send", pt_chat_send)
	})
}

// ============================================================================

func pt_chat_send(bot *app.Bot) {
	var tp int32
	p := rand.Float32()
	if p < 0.7 { // world
		tp = 1
	} else if p < 0.95 { // guild
		tp = 3
	} else { // cross
		tp = 2
	}

	bot.SendMsg(&msg.C_ChatSend{
		Tp:      tp,
		Content: v_chat_gen.Gen(30),
	})
}
