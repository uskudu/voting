package poll

type Poll struct {
	ID      string   `gorm:"primaryKey" json:"id"`
	Title   string   `json:"title"`
	Options []Option `gorm:"foreignKey:PollID" json:"options"`
}

type Option struct {
	ID     string `gorm:"primaryKey" json:"id"`
	Text   string `json:"text"`
	PollID string `json:"-"`
	Votes  []Vote `gorm:"foreignKey:OptionID" json:"votes,omitempty"`
}

type Vote struct {
	ID       string `gorm:"primaryKey"`
	UserID   string
	OptionID string
}
