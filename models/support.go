package models

type SupportTicket struct {
    ID        uint `json:"id" gorm:"primaryKey"`

    UserID    uint
    User      User `gorm:"foreignKey:UserID"`

    Subject   string
    Message   string
    Status    string // open, in_progress, closed
}
