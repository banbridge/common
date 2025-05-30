package mapper

import (
	"context"
	"errors"

	"github.com/jinzhu/copier"
	"github.com/mohae/deepcopy"

	"github.com/banbridge/common/pkg/gptr"
)

func Convert[S, T any](ctx context.Context, src S, dst T) error {
	if gptr.IsNil(src) || gptr.IsNil(dst) {
		return errors.New("src or dst is nil")
	}

	if !gptr.IsPtr(dst) {
		return errors.New("dst is not a pointer")
	}

	newS := deepcopy.Copy(&src)
	if err := copier.Copy(dst, newS); err != nil {
		return err
	}
	return nil
}
