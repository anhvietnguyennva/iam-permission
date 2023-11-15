package valueobject

type CreateSubjectSetReq struct {
	Namespace   string
	Object      string
	Relation    string
	Description string
	CreatedBy   string
}
