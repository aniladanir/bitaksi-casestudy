package importer

import (
	"context"
	"io"
)

type Importer interface {
	ImportCoordinates(ctx context.Context, reader io.Reader) error
}
