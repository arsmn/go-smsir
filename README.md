# go-smsir #

go-smsir is a Go client library for accessing the [sms.ir](https://sms.ir).

## Install ##

```sh
go get -u github.com/arsmn/go-smsir
```

## Usage ##

```go
import "github.com/arsmn/go-smsir/smsir"
```

Construct a new smsir client, then use the various services on the client to
access different parts of the smsir API. For example:

```go
client := github.NewClient().WithAuthentication("Your API Key", "Your Secret Key")

// send sms with template
req := &smsir.UltraFastSendRequest{
	Mobile:     "xxx",
	TemplateID: "xxx",
	Parameters: []smsir.UltraFastParameter{
		{Key: "xxx", Value: "xxx"},
		{Key: "xxx", Value: "xxx"},
	},
}
_, err := client.Verification.UltraFastSend(context.Background(), req)
```
