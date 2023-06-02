package stmts

var (
	GetEmailsForUserByID = Statement[any]{
		Query: `SELECT * FROM user_emails WHERE user_id = ?`,
		Args:  []string{"user_id"},
	}
)
