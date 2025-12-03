package configuration

import (
	"os"

	"github.com/xendit/xendit-go/v7"
)

func NewXenditClient() *xendit.APIClient {
	apiKey := os.Getenv("XENDIT_API_KEY");
	xdt := xendit.NewClient(apiKey);
	return xdt;
}
