package mapper

import (
	"context"
	"errors"

	"github.com/jinzhu/copier"
	"github.com/mohae/deepcopy"

	"github.com/banbridge/common/pkg/ptr"
)

func Convert[S, T any](ctx context.Context, src S, dst T) error {
	if ptr.IsNil(src) || ptr.IsNil(dst) {
		return errors.New("src or dst is nil")
	}

	if !ptr.IsPtr(dst) {
		return errors.New("dst is not a pointer")
	}

	newS := deepcopy.Copy(&src)
	if err := copier.Copy(dst, newS); err != nil {
		return err
	}
	return nil
}
