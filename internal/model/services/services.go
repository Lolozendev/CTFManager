package services

type Services struct {
	Wireguard Wireguard `yaml:"wireguard"`
	Dnsmasq   Dnsmasq   `yaml:"dnsmasq"`
	Challenge []Challenge
}

/*
services:
	{wireguard ...}
	{dnsmasq ...}
	{challenges ...}
		[ challenge ... ]
*/
