package types

const (
	Deployment     = "Deployment"
	StatefulSet    = "StatefulSet"
	DaemonSet      = "DaemonSet"
	TimedOutReason = "ProgressDeadlineExceeded"
)

type Operator string

const (
	OPERATOR_ADD    Operator = "ADD"
	OPERATOR_DELETE Operator = "DELETE"
)
