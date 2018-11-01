package network


func (network *nbsNetwork) NewHost(options ...HostOption) Host{

	instance := &NbsHost{

	}

	return instance
}

func (network *nbsNetwork) ListenAddrString(address string) HostOption{
	return func() error{
		return nil
	}
}