package group

type Group struct {
	GroupID    *uint
	GroupName  string
	Users      []*uint
	Msgs       []*uint
	Owner      *uint
	TotalUsers uint
}
