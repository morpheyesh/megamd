package carton

import (
	log "github.com/Sirupsen/logrus"
	"github.com/megamsys/megamd/carton/bind"
	"github.com/megamsys/megamd/provision"
	"gopkg.in/yaml.v2"
)

// A carton represents a real world assembly.
// This struct provides and easy way to manage information about an assembly, instead passing it around
type Carton struct {
	Id           string //assemblyid
	Name         string
	CartonsId    string
	Tosca        string
	ImageVersion string
	Compute      provision.BoxCompute
	DomainName   string
	Provider     string
	PublicIp     string
	Boxes        *[]provision.Box
}

func (a *Carton) String() string {
	if d, err := yaml.Marshal(a); err != nil {
		return err.Error()
	} else {
		return string(d)
	}
}

//Global provisioners set by the subd daemons.
var ProvisionerMap map[string]provision.Provisioner = make(map[string]provision.Provisioner)

//If there are boxes, then it set the enum BoxSome or its BoxZero
func (c *Carton) lvl() provision.BoxLevel {
	if len(*c.Boxes) > 0 {
		return provision.BoxSome
	} else {
		return provision.BoxNone
	}
}

//Converts a carton to a box, (applicable in torpedo case)
func (c *Carton) toBox() error { //assemblies id.
	switch c.lvl() {
	case provision.BoxNone:
		c.Boxes = &[]provision.Box{provision.Box{
			Id:           c.Id,        //should be the component id, but in case of BoxNone there is no component id.
			CartonId:     c.Id,        //We stick the assemlyid here.
			CartonsId:    c.CartonsId, //assembliesId,
			CartonName:   c.Name,
			Name:         c.Name,
			DomainName:   c.DomainName,
			Level:        c.lvl(), //based on the level, we decide to use the Box-Id as ComponentId or AssemblyId
			ImageVersion: c.ImageVersion,
			Compute:      c.Compute,
			Provider:     c.Provider,
			PublicIp:     c.PublicIp,
			Tosca:        c.Tosca,
		},
		}
	}
	return nil
}

// Deploy carton, which basically deploys the boxes.
func (c *Carton) Deploy() error {
	for _, box := range *c.Boxes {
		err := Deploy(&DeployOpts{B: &box})
		if err != nil {
			return err
		}
	}
	return nil
}

// Destroys a carton, which deletes its boxes.
func (c *Carton) Destroy() error {
	for _, box := range *c.Boxes {
		err := Destroy(&DestroyOpts{B: &box})
		if err != nil {
			return err
		}
	}
	return nil
}

// moves the state to the desired state
// changing the boxes state to StatusStateup.
func (c *Carton) Stateup() error {
	for _, box := range *c.Boxes {
		err := ChangeState(&StateChangeOpts{B: &box, Changed: provision.StatusStateup})
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *Carton) Bind(box *provision.Box) error {
	/*	bis := c.Group()
		for _, bi := range bis {
			err = bi.Bind(c, bi)
			if err != nil {
				log.Errorf("Error binding the box %s with the service %s: %s", box.Name, bi.Name, err)
			}
		}
	*/
	return nil
}

func (c *Carton) Unbind(box *provision.Box) error {
	/*	bis := c.Group()
		for _, bi := range bis {
			err = bi.UnBind(c, bi)
			if err != nil {
				log.Errorf("Error binding the box %s with the service %s: %s", box.Name, bi.Name, err)
			}
		}*/
	return nil
}

// Group the related_components into BindInstances.
func (c *Carton) Group() ([]*bind.YBoundBox, error) {
	return nil, nil
}

// Available returns true if at least one of N boxes which is started
func (c *Carton) Available() bool {
	for _, box := range *c.Boxes {
		if box.Available() {
			return true
		}
	}
	return false
}

// starts the box calling the provisioner.
// changing the boxes state to StatusStarted.
func (c *Carton) Start() error {
	for _, box := range *c.Boxes {
		err := ProvisionerMap[box.Provider].Start(&box, "", nil)
		if err != nil {
			log.Errorf("Unable to start the box  %s", err)
			return err
		}
	}
	return nil
}

// stops the box calling the provisioner.
// changing the boxes state to StatusStopped.
func (c *Carton) Stop() error {
	for _, box := range *c.Boxes {
		err := ProvisionerMap[box.Provider].Stop(&box, "", nil)
		if err != nil {
			log.Errorf("Unable to stop the box %s", err)
			return err
		}
	}
	return nil
}

// restarts the box calling the provisioner.
// changing the boxes state to StatusStarted.
func (c *Carton) Restart() error {
	for _, box := range *c.Boxes {
		err := ProvisionerMap[box.Provider].Restart(&box, "", nil)
		if err != nil {
			log.Errorf("[start] error on start the box %s - %s", box.Name, err)
			return err
		}
	}
	return nil
}
