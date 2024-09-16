package entity

type (
	UserID string
)
type User struct {
	ID        UserID
	UserName  string
	FirstName string
	LastName  string
	CreatedAt string
	UpdatedAt string
}
