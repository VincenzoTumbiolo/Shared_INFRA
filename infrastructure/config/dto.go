package config

type SubnetVisibility string

const (
	Public   SubnetVisibility = "Public"
	Private  SubnetVisibility = "Private"
	Isolated SubnetVisibility = "Isolated"
)
