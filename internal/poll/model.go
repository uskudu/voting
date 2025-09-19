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
	Votes  []Vote `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"votes,omitempty"`
}

type Vote struct {
	ID       string `gorm:"primaryKey"`
	UserID   string `gorm:"index"` // индекс для быстрого поиска
	OptionID string `gorm:"index"`
}
