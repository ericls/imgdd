package domainmodels

type StorageDefinition struct {
	Id          string
	Identifier  string
	StorageType string
	Config      string
	IsEnabled   bool
	Priority    int32
}
