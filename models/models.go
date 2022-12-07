package models

type MessageEvent struct {
	Provider string
	Message  string
	OrgId    int
	Region   string
}
