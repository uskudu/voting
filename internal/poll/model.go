package poll

type Poll struct {
	ID      int      `gorm:"primaryKey;autoIncrement" json:"id"`
	Title   string   `json:"title"`
	UserID  int      `gorm:"index" json:"-"`
	Options []Option `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"options"`
}

type CreateOrPatchPollRequest struct {
	Title   string                `json:"title"`
	Options []CreateOptionRequest `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"options"`
}

type Option struct {
	ID     int    `gorm:"primaryKey;autoIncrement" json:"id"`
	Text   string `json:"text"`
	PollID int    `gorm:"index" json:"-"`
	Votes  int    `json:"votes"`
}

type CreateOptionRequest struct {
	Text string `json:"text"`
}
