package domain

type Session struct {
	ID               string `gorm:"id; not null; primary_key"`
	RefreshTokenHash string `gorm:"refresh_token_hash"`
	IP               string `gorm:"ip"`
}
