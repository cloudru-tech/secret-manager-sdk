package api

import (
	"crypto/tls"
	"errors"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/cloudru-tech/secret-manager-sdk/api/v1"
	"github.com/cloudru-tech/secret-manager-sdk/api/v2"
)

// Config is the connection configuration.
type Config struct {
	// Address is the address string to csm service, e.g: "secretmanager.api.cloud.ru:443"
	Host string
	// Insecure disables the SSL verification
	Insecure bool
}

// Client is the gRPC client to csm service.
type Client struct {
	conn *grpc.ClientConn

	// V2 is version 2 of the CSM client.
	V2 V2

	// SecretService is the client to the secret manager service.
	// Deprecated: use the SecretServiceV2 instead.
	SecretService v1.SecretManagerServiceClient
}

// V2 is version 2 of the CSM client.
type V2 struct {
	FolderService v2.FolderServiceClient
	SecretService v2.SecretServiceClient
}

// New возвращает новый клиент к gRPC серверу
func New(conf *Config, opts ...grpc.DialOption) (*Client, error) {
	if conf == nil {
		return nil, errors.New("empty config")
	}

	if conf.Insecure {
		opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	} else {
		opts = append(opts,
			grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{MinVersion: tls.VersionTLS13})),
		)
	}

	conn, err := grpc.NewClient(conf.Host, opts...)
	if err != nil {
		return nil, err
	}

	return &Client{
		conn: conn,

		V2: V2{
			FolderService: v2.NewFolderServiceClient(conn),
			SecretService: v2.NewSecretServiceClient(conn),
		},
		SecretService: v1.NewSecretManagerServiceClient(conn),
	}, nil
}

// Close закрывает соединение с сервером
func (c *Client) Close() error { return c.conn.Close() }
