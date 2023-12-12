package aoi

const (
	OPER_AGENT_JOIN            = 1
	OPER_AGENT_LEAVE           = 2
	OPER_AGENT_MOVE            = 3
	OPER_NOAGENT_JOIN          = 4
	OPER_NOAGENT_LEAVE         = 5
	OPER_NOAGENT_MOVE          = 6
	OPER_ACTION_MOVE_REUSE     = 16
	OPER_ACTION_AOI            = 17
	OPER_ACTION_INTERSTING_ALL = 18
	OPER_ACTION_AGENT_ALL      = 19
	OPER_QUIT                  = 100
)

type OperEvent struct {
	Entity Entity
	Op     uint8
	Handle func(ming, other Entity)
}
