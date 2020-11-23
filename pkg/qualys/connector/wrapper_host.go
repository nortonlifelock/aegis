package connector

import (
	"context"
	"fmt"
	"strconv"

	"github.com/nortonlifelock/aegis/pkg/domain"
	"github.com/nortonlifelock/aegis/pkg/qualys"
)

type host struct {
	h qualys.QHost
}

func (h *host) SourceID() *string {
	id := strconv.Itoa(h.h.HostID)
	return &id
}

func (h *host) OS() string {
	return h.h.OperatingSystem.Text
}

func (h *host) HostName() string {
	return h.h.DNS.Text
}

func (h *host) MAC() string {
	return ""
}

func (h *host) IP() string {
	return h.h.IPAddress
}

func (h *host) Region() *string {
	return nil
}

func (h *host) InstanceID() *string {
	if len(h.h.EC2Id) > 0 {
		return &h.h.EC2Id
	} else {
		return nil
	}
}

func (h *host) ID() string {
	return ""
}

func (h *host) Vulnerabilities(ctx context.Context) (param <-chan domain.Detection, err error) {
	var out = make(chan domain.Detection)
	defer close(out)
	err = fmt.Errorf("not implemented")
	return out, err
}

func (h *host) GroupID() *string {
	return nil
}

func (h *host) TrackingMethod() *string {
	return &h.h.TrackingMethod
}
