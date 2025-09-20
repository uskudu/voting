package poll

import "github.com/google/uuid"

type Poll struct {
	ID      uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	Title   string    `json:"title"`
	Options []Option  `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"options"`
}

type Option struct {
	ID     uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	Text   string    `json:"text"`
	PollID uuid.UUID `gorm:"index" json:"-"`
	Votes  int       `json:"votes"`
}
