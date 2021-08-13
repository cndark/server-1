package mail

import (
	"fw/src/game/app/comp"
)

type IPlayer interface {
	comp.IPlayer

	GetMailBox() *MailBox
	CheckGmailDeliverCond(cond string) bool
}
