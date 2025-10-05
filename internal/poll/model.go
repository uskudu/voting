package poll

type Poll struct {
	ID      string   `gorm:"primaryKey" json:"id"`
	Title   string   `json:"title"`
	UserID  string   `gorm:"index" json:"-"`
	Options []Option `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"options"`
}

type CreateOrPatchPollRequest struct {
	Title   string                `json:"title"`
	Options []CreateOptionRequest `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"options"`
}

type Option struct {
	PollID string `gorm:"index" json:"-"`
	ID     int    `gorm:"primaryKey;autoIncrement" json:"id"`
	Text   string `json:"text"`
	Votes  int    `json:"votes"`
}

type CreateOptionRequest struct {
	Text string `json:"text"`
}
