package implementations

import (
	"context"
	"fmt"
	"github.com/nortonlifelock/domain"
	"testing"
)

const testingID = "TEST"

type mockSQLDriver struct {
	domain.DatabaseConnection
}

func TestSomething(t *testing.T) {
	asj := &AssetSyncJob{}
	fmt.Println("here")

	asj.Process(context.Background(), testingID, nil, mockSQLDriver{}, nil, "", nil, nil, nil)
}
