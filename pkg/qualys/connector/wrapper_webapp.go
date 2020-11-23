package connector

import (
	"context"
	"fmt"
	"github.com/nortonlifelock/aegis/pkg/domain"
)

type WebAppWrapper struct {
	// the webApp ID
	sourceID string

	// the finding UID
	findingID string
	name      string
	url       string
}

// ID is the ID of the device as reported by the backend database of Aegis
func (w *WebAppWrapper) ID() string {
	return ""
}

func (w *WebAppWrapper) SourceID() *string {
	return &w.findingID
}

func (w *WebAppWrapper) OS() string {
	return ""
}

func (w *WebAppWrapper) HostName() string {
	return w.url
}

func (w *WebAppWrapper) MAC() string {
	return w.name
}

func (w *WebAppWrapper) IP() string {
	return ""
}

// Vulnerabilities not implemented as the interface method is not yet used
func (w *WebAppWrapper) Vulnerabilities(ctx context.Context) (param <-chan domain.Detection, err error) {
	return nil, fmt.Errorf("not implemented")
}

func (w *WebAppWrapper) Region() *string {
	return nil
}

func (w *WebAppWrapper) InstanceID() *string {
	return nil
}

func (w *WebAppWrapper) GroupID() *string {
	return nil
}

func (w *WebAppWrapper) TrackingMethod() *string {
	return nil
}
