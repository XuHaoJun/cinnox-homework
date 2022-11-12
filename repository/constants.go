package repository

type CollectionName = string

const (
	CollNameCinnoxMessage CollectionName = "cinnoxMessage"
	CollNameLineEvent     CollectionName = "lineEvent"
	CollNameLineUser      CollectionName = "lineUser"
)

type DatabaseName = string

const (
	DbNameMain DatabaseName = "cinnox"
)
