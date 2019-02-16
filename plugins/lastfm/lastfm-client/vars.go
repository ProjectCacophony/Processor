package lastfmclient

// Period is a type for periods used for Last.FM requests
type Period string

// defines possible LastFM periods
const (
	PeriodOverall Period = "overall"
	Period7day    Period = "7day"
	Period1month  Period = "1month"
	Period3month  Period = "3month"
	Period6month  Period = "6month"
	Period12month Period = "12month"
)

const (
	lastFmTargetImageSize = "extralarge"
)
