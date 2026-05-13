// Package service provides business logic.
// See https://www.alexedwards.net/blog/the-fat-service-pattern
package service

import (
	"context"

	"maragu.dev/glue/email/postmark"

	"app/model"
	"app/sqlite"
)

type Fat struct {
	db     *sqlite.Database
	sender *postmark.Sender
}

type NewFatOptions struct {
	Database *sqlite.Database
	Sender *postmark.Sender
}

func NewFat(opts NewFatOptions) *Fat {
	return &Fat{
		db:     opts.Database,
		sender: opts.Sender,
	}
}

func (f *Fat) GetUser(ctx context.Context, id model.UserID) (model.User, error) {
	return f.db.GetUser(ctx, id)
}
