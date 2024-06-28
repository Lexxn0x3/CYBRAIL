package checkstrategies

// CheckStrategy interface for different check strategies
type CheckStrategy interface {
	Execute(logPath string) (string, error)
}

// CheckModule struct to hold the strategy and its name
type CheckModule struct {
	Name     string
	Strategy CheckStrategy
}
