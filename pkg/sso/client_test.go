package sso_test

import (
	"context"
	"log"
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
	sssov1 "go.pilab.hu/pilab-cloud/ga-pi/gen/pilab/ssso/v1"
	"go.pilab.hu/pilab-cloud/ga-pi/mocks"
	"go.pilab.hu/pilab-cloud/ga-pi/pkg/sso"
	"go.uber.org/mock/gomock"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
)

const bufSize = 1024 * 1024

var lis *bufconn.Listener

func bufDialer(context.Context, string) (net.Conn, error) {
	return lis.Dial()
}

func createMockClient(t *testing.T) (*sso.Client, *mocks.MockAuthServiceServer, error) {
	t.Helper()

	lis = bufconn.Listen(bufSize)
	s := grpc.NewServer()

	// Create a new mock server
	ctrl := gomock.NewController(t)
	srv := mocks.NewMockAuthServiceServer(ctrl)

	sssov1.RegisterAuthServiceServer(s, srv)

	go func() {
		if err := s.Serve(lis); err != nil {
			log.Fatalf("Server exited with error: %v", err)
		}
	}()

	conn, err := grpc.DialContext(context.Background(), "bufnet",
		grpc.WithContextDialer(bufDialer),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, nil, err
	}

	return sso.NewClient(
		sso.WithConnection(conn),
	), srv, nil
}

func TestClient_Login(t *testing.T) {
	// Create a new client
	client, srv, err := createMockClient(t)
	assert.NoError(t, err, "Failed to create gRPC client")

	offlineToken := new(string)
	*offlineToken = "test-offline-token"

	t.Run("TestSuccessfulLogin", func(t *testing.T) {
		srv.EXPECT().Login(gomock.Any(), gomock.Any()).Return(
			&sssov1.LoginResponse{
				Response: &sssov1.LoginResponse_TokenResponse{
					TokenResponse: &sssov1.AuthTokenResponse{
						AccessToken:  "test-token",
						RefreshToken: "test-refresh-token",
						OfflineToken: offlineToken,
					},
				},
			}, nil,
		).Times(1)

		// Test valid login credentials
		err := client.Login(context.Background(), "test-realm-id", "test-client-id", "test@example.com", "password")
		assert.NoError(t, err, "Login should succeed with valid credentials")
	})

	t.Run("TestWrongPassword", func(t *testing.T) { // Test invalid login credentials
		srv.EXPECT().Login(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, req *sssov1.LoginRequest) (*sssov1.LoginResponse, error) {
			creds := req.GetPasswordAuth()
			assert.NotNil(t, creds, "Credentials should not be nil")
			assert.Equal(t, creds.Username, "test@example.com")
			assert.Equal(t, creds.Password, "wrongpassword")

			return nil, sso.ErrInvalidCredentials
		})

		// srv.EXPECT().Login(gomock.Any(), gomock.Any()).Return(
		// 	sso.ErrInvalidCredentials.ToLoginResponse(),
		// 	nil,
		// ).Times(1)

		err := client.Login(context.Background(), "test-realm-id", "test-client-id", "test@example.com", "wrongpassword")
		assert.Error(t, err, "Login should fail with invalid credentials")
	})
}

func TestPiError_Error(t *testing.T) {
	// Create a new PiError
	err := &sso.PiError{
		Code:    "500",
		Message: "Internal Server Error",
		TraceID: "123456789",
	}

	// Test the Error() method
	expectedError := "500: Internal Server Error (TraceID: 123456789)"
	assert.Equal(t, expectedError, err.Error(), "Error message should be formatted correctly")
}

func TestIntegration(t *testing.T) {
	client, err := createIntegrationClient()
	assert.NoError(t, err, "Failed to create gRPC client")

	t.Run("TestSuccessfulLogin", func(t *testing.T) {
		err := client.Login(context.Background(), "test-realm-id", "test-client-id", "user1@asd.hu", "password")
		assert.NoError(t, err)
	})
}

func createIntegrationClient() (*sso.Client, error) {
	conn, err := grpc.Dial(
		"localhost:5000",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, err
	}

	return sso.NewClient(
		sso.WithConnection(conn),
	), nil
}
