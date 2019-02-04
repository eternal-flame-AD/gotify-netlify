package main

// DefaultConfig implements plugin.Configurer
func (c *Plugin) DefaultConfig() interface{} {
	return new(Conf)
}

// ValidateAndSetConfig implements plugin.Configurer
func (c *Plugin) ValidateAndSetConfig(conf interface{}) error {
	c.conf = conf.(*Conf)
	return nil
}
