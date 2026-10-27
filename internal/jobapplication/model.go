type Job struct {
	ID			uuid.UUID
	CompanyName	string
	RoleTitle	string
}

type JobApplication struct {
	ID				uuid.UUID
	UserID			uuid.UUID
	Job				Job
	Stage			string
	NextFollowUpOn	*time.Time
}