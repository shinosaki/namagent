package signalcontext

type ContextKey string

const (
	SELF   ContextKey = "self"
	CONFIG ContextKey = "config"
)
