package crypto

import (
	"reflect"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go/service/kms"
)

func TestCreateKMSClientWithProfile(t *testing.T) {
	type args struct {
		keyID   string
		profile string
		region  string
	}
	tests := []struct {
		name       string
		args       args
		wantClient *KMSClient
		wantErr    bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotClient, err := CreateKMSClientWithProfile(tt.args.keyID, tt.args.profile, tt.args.region)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateKMSClientWithProfile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotClient, tt.wantClient) {
				t.Errorf("CreateKMSClientWithProfile() = %v, want %v", gotClient, tt.wantClient)
			}
		})
	}
}

func TestCreateKMSClient(t *testing.T) {
	type args struct {
		keyID  string
		region string
	}
	tests := []struct {
		name       string
		args       args
		wantClient *KMSClient
		wantErr    bool
	}{
		{
			name: "Create KMS Client",
			args: args{
				keyID:  "abc123",
				region: "us-west-2",
			},
			wantClient: &KMSClient{
				keyID:   "abc123",
				KeySpec: encryptionType256,
				Client:  &kms.KMS{},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotClient, err := CreateKMSClient(tt.args.keyID, tt.args.region)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateKMSClient() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if strings.Compare(gotClient.keyID, tt.wantClient.keyID) != 0 {
				t.Errorf("CreateKMSClient() = %v, want %v", gotClient.keyID, tt.wantClient.keyID)
			}
			if strings.Compare(gotClient.KeySpec, tt.wantClient.KeySpec) != 0 {
				t.Errorf("CreateKMSClient() = %v, want %v", gotClient.KeySpec, tt.wantClient.KeySpec)
			}
		})
	}
}

func TestKMSClient_Encrypt(t *testing.T) {
	type fields struct {
		Client  *kms.KMS
		keyID   string
		KeySpec string
	}
	type args struct {
		message string
	}
	tests := []struct {
		name                string
		fields              fields
		args                args
		wantEncryptedString string
		wantErr             bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			kmsClient := &KMSClient{
				Client:  tt.fields.Client,
				keyID:   tt.fields.keyID,
				KeySpec: tt.fields.KeySpec,
			}
			gotEncryptedString, err := kmsClient.Encrypt(tt.args.message)
			if (err != nil) != tt.wantErr {
				t.Errorf("KMSClient.Encrypt() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotEncryptedString != tt.wantEncryptedString {
				t.Errorf("KMSClient.Encrypt() = %v, want %v", gotEncryptedString, tt.wantEncryptedString)
			}
		})
	}
}

func TestKMSClient_Decrypt(t *testing.T) {
	type fields struct {
		Client  *kms.KMS
		keyID   string
		KeySpec string
	}
	type args struct {
		encryptedText string
	}
	tests := []struct {
		name        string
		fields      fields
		args        args
		wantMessage string
		wantErr     bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			kmsClient := &KMSClient{
				Client:  tt.fields.Client,
				keyID:   tt.fields.keyID,
				KeySpec: tt.fields.KeySpec,
			}
			gotMessage, err := kmsClient.Decrypt(tt.args.encryptedText)
			if (err != nil) != tt.wantErr {
				t.Errorf("KMSClient.Decrypt() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotMessage != tt.wantMessage {
				t.Errorf("KMSClient.Decrypt() = %v, want %v", gotMessage, tt.wantMessage)
			}
		})
	}
}
