package model

import (
	"github.com/Lolozendev/CTFManager/internal/model/network"
	"github.com/Lolozendev/CTFManager/internal/model/services"
)

type Member struct {
	Username string
}

type Team struct {
	Name     string
	Number int
	Members  []Member
	Network  network.Network   `yaml:"networks"`
	Services services.Services `yaml:"services"`
}

/*
{ networks ... }
{ services ... }
*/
