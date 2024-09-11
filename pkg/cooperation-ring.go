package pkg

type CooperationRing struct {
	ID          string
	MemberCount uint
	Weight      float64
	Next        string
	Prev        string
	Investor    string
	Rounds      uint
}
