package pkg

const (
	RoundsCount     = 10
	FractalMin      = 500
	FractalMax      = 2000
	VerificationMin = 20
	VerificationMax = 50
)

type CooperationTable struct {
	ID          string   `json:"id"`
	MemberCount uint     `json:"member_count"`
	Weight      float64  `json:"weight"`
	Next        string   `json:"next"`
	Prev        string   `json:"prev"`
	Investor    string   `json:"investor"`
	Rounds      uint     `json:"rounds"`
	Members     []string `json:"-"`
}
