package email

import "context"

type EmailService interface {
	SendOTP(ctx context.Context, to, code string) error
}
