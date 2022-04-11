package helpers

import (
	"context"
)

func ExtractClientID(ctx context.Context) string {
	return ctx.Value("client_id").(string)
}