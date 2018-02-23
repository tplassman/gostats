package models

import "time"

type Post struct {
	Id           uint
	Publish_Date uint
	Name         string
	Url          string
	SocialShares map[string]int
}

func (p Post) FormattedDate() time.Time {
	return time.Unix(int64(p.Publish_Date/1000), 0)
}
