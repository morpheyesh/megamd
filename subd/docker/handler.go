package docker

import (
	"github.com/megamsys/megamd/carton"
)

type Handler struct {
	Provider string
	D        *Config
}

// NewHandler returns a new instance of handler with routes.
func NewHandler(c *Config) *Handler {
	return &Handler{D: c}
}

func (h *Handler) serveAMQP(r *carton.Requests) error {
	p, err := carton.ParseRequest(r.CatId, r.Category, r.Action)
	if err != nil {
		return err
	}

	if rp := carton.NewReqOperator(r.CatId); rp != nil {
		return rp.Accept(&p) //error is swalled here.
	}
	return nil
}
