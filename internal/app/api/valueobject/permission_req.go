package valueobject

type CheckPermissionReq struct {
	Namespace string
	Object    string
	Relation  string
	SubjectID string
	MaxDepth  uint8
}
