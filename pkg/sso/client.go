package sso

import (
	"context"
	"fmt"
	"os"

	"github.com/rs/zerolog/log"
	sssov1 "go.pilab.hu/pilab-cloud/ga-pi/gen/pilab/ssso/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Transport string

const (
	// ProtocolGRPC is a standard gRPC protocol.
	TransportGRPC Transport = "grpc"
	// ProtocolWS is a gRPC over WebSocket protocol.
	TransportWS Transport = "ws"
	// ProtocolREST is not implemented yet.
	TransportREST Transport = "rest"
)

var (
	serviceHost = "sso.pilab.hu"
	servicePort = "443"
	secure      = true
	transport   = TransportGRPC
)

type Client struct {
	sc sssov1.AuthServiceClient

	conn *grpc.ClientConn
}

func loadEnv() {
	if h := os.Getenv("SSSO_HOST"); h != "" {
		serviceHost = h
	}

	if p := os.Getenv("SSSO_PORT"); p != "" {
		servicePort = p
	}

	if s := os.Getenv("SSSO_SECURE"); s != "" {
		secure = s == "true"
	}

	if t := os.Getenv("SSSO_TRANSPORT"); t != "" {
		switch t {
		case "grpc":
			transport = TransportGRPC
		case "ws":
			transport = TransportWS
		case "rest":
			transport = TransportREST
		default:
			panic("invalid SSSO_TRANSPORT value")
		}
	}
}

type ClientOption func(*Client)

// WithConnection sets the connection for the client.
func WithConnection(conn *grpc.ClientConn) ClientOption {
	return func(o *Client) {
		o.conn = conn
	}
}

// NewClient creates a new client.
func NewClient(opts ...ClientOption) *Client {
	// Create a new client
	c := new(Client)

	// Load environment variables
	loadEnv()

	for _, o := range opts {
		o(c)
	}

	// if the connection is nil create a new one with default values
	if c.conn == nil {
		// Create a new connection
		conn, err := grpc.Dial(
			fmt.Sprintf("%s:%s", serviceHost, servicePort),
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		)
		if err != nil {
			panic(err)
		}

		c.conn = conn
	}

	// Create the client with the connection
	c.sc = sssov1.NewAuthServiceClient(c.conn)

	return c
}

type PiError struct {
	Code    string
	Message string
	TraceID string
}

func (e *PiError) Error() string {
	return "login error"
}

// Login tries to login with the given credentials.
func (c *Client) Login(ctx context.Context, realmID, clientID, username, password string) error {
	message, err := c.sc.Token(ctx, &sssov1.TokenRequest{
		Client: &sssov1.TenantClientMessage{
			Tenant:   realmID,
			ClientId: clientID,
		},
		Credentials: &sssov1.TokenRequest_PasswordAuth{
			PasswordAuth: &sssov1.PasswordAuthMessage{
				Username: username,
				Password: password,
			},
		},
	})
	if err != nil {
		return err
	}

	switch r := message.GetResponse().(type) {
	case *sssov1.TokenResponse_Tokens:
		tokens := r.Tokens
		fmt.Printf("Login successful: %s\n", tokens.AccessToken)
	case *sssov1.TokenResponse_Error:
		log.Printf("TraceID extraction not implemented yet")

		err := &PiError{
			Code:    fmt.Sprintf("%d", r.Error.Code),
			Message: r.Error.Message,
			TraceID: "",
		}

		return fmt.Errorf("%w: %s", err, r.Error.Message)
	case *sssov1.TokenResponse_ActionRequired:
		panic("Action Required case not implemented yet")
	}

	return nil
}
