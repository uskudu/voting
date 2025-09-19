package poll

type Poll struct {
	ID      string   `gorm:"primaryKey" json:"id"`
	Title   string   `json:"title"`
	Options []Option `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"options"`
}

type Option struct {
	ID     string `gorm:"primaryKey" json:"id"`
	Text   string `json:"text"`
	PollID string `gorm:"index" json:"-"`
	Votes  int    `json:"votes"`
}
