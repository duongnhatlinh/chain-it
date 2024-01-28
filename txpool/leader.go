package txpool

//Aggregate root must implement aggregate interface
type Leader struct {
	Id string
}

type LeaderRepository interface {
	Set(leader Leader)
	Get() Leader
}
