package consul

import (
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

func pathConfig() *framework.Path {
	return &framework.Path{
		Pattern: "config",
		Fields: map[string]*framework.FieldSchema{
			"address": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "Consul server address",
			},

			"scheme": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "URI scheme for the Consul address",

				// https would be a better default but Consul on its own
				// defaults to HTTP access, and when HTTPS is enabled it
				// disables HTTP, so there isn't really any harm done here.
				Default: "http",
			},

			"token": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "Token for API calls",
			},
		},

		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.WriteOperation: pathConfigWrite,
		},
	}
}

func pathConfigWrite(
	req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	entry, err := logical.StorageEntryJSON("config", config{
		Address: data.Get("address").(string),
		Scheme:  data.Get("scheme").(string),
		Token:   data.Get("token").(string),
	})
	if err != nil {
		return nil, err
	}

	if err := req.Storage.Put(entry); err != nil {
		return nil, err
	}

	return nil, nil
}

type config struct {
	Address string `json:"address"`
	Scheme  string `json:"scheme"`
	Token   string `json:"token"`
}