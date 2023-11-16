package valueobject

type CreateSubjectSetRelationTupleReq struct {
	Namespace           string
	Object              string
	Relation            string
	SubjectSetNamespace string
	SubjectSetObject    string
	SubjectSetRelation  string
	CreatedBy           string
}
