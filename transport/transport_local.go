package transport

var ConnectionLocal = Connection{
	Protocol:             "local",
	Destination:          "",
	EstabilishConnection: func() error { return nil }, // no need to estabilish connection

}
