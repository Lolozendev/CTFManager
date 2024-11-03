package services

type Services struct {
	Wireguard Wireguard `yaml:"wireguard"`
	Dnsmasq   Dnsmasq   `yaml:"dnsmasq"`
}

/*
services:
	{wireguard ...}
	{dnsmasq ...}
	{challenges ...} // This is a list of challenges it is concatenated in the template that's why it's not here
		[ challenge ... ]
*/
