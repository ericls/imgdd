package domainmodels

type StorageDefinition struct {
	Id         string
	Identifier string
	Type       string
	Config     string
	IsEnabled  bool
	Priority   int32
}
