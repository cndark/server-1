package bill

import (
	"errors"
	"net/http"
)

// ============================================================================

var (
	ErrNoKey    = errors.New("invalid key")
	ErrNoPlayer = errors.New("player not found")
)

// ============================================================================

var handlers = map[string]func(*http.Request) (r string, err error){
	"give_items": handle_give_items,
}
