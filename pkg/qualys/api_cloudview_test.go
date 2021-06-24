package qualys

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetAccountType(t *testing.T) {
	t.Run("Should return as an AWS account", func(t *testing.T) {
		mockAWS := "123456789012"
		assert.Equal(t,AWS_CLOUD_ACCOUNT,getAccountType(mockAWS))
	})
	t.Run("Should return as an Azure account", func(t *testing.T) {
		mockAzure := "20ff7fc3-e762-44dd-bd96-b71116dcdc23"
		assert.Equal(t,AZURE_CLOUD_ACCOUNT,getAccountType(mockAzure))
	})
	t.Run("Should return as an GCP account", func(t *testing.T) {
		mockGCP := "12we"
		assert.Equal(t,GOOGLE_CLOUD_ACCOUNT,getAccountType(mockGCP))
	})

}

