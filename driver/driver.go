package driver


type Driver interface {
	Init()
	Create()
	List()
	Close()
}

type Controller struct {
    driver Driver
}

func New(uri string) *Controller{
    c := &Controller{driver: New_kvm(uri)}
    c.driver.Init()

    return  c
}


func (c *Controller) Start() {
	c.driver.Create()
	c.driver.List()
	c.driver.Close()
}
