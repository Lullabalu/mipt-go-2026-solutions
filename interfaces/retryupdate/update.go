//go:build !solution

package retryupdate

import (
	"errors"

	"github.com/gofrs/uuid"
	"gitlab.com/slon/shad-go/retryupdate/kvapi"
)

func UpdateValue(c kvapi.Client, key string, updateFn func(oldValue *string) (string, error)) error {
	for {
		req := kvapi.GetRequest{Key: key}
		resp, err := c.Get(&req)

		var old *string
		var nversion uuid.UUID
		var oversion uuid.UUID

		if err != nil {
			var autherr *kvapi.AuthError
			if errors.As(err, &autherr) {
				return err
			}

			if errors.Is(err, kvapi.ErrKeyNotFound) {
				old = nil
				oversion = uuid.UUID{}
			} else {
				var apierr *kvapi.APIError
				if errors.As(err, &apierr) {
					continue
				}
				return err
			}
		} else {
			oldv := resp.Value
			old = &oldv
			oversion = resp.Version
		}

		newv, err := updateFn(old)
		if err != nil {
			return err
		}

		nversion = uuid.Must(uuid.NewV4())

		for {
			setreq := kvapi.SetRequest{
				Key:        key,
				Value:      newv,
				OldVersion: oversion,
				NewVersion: nversion,
			}

			_, err = c.Set(&setreq)
			if err == nil {
				return nil
			}

			var autherr *kvapi.AuthError
			if errors.As(err, &autherr) {
				return err
			}

			if errors.Is(err, kvapi.ErrKeyNotFound) {
				old = nil
				oversion = uuid.UUID{}
				newv, err = updateFn(old)
				if err != nil {
					return err
				}
				nversion = uuid.Must(uuid.NewV4())
				continue
			}

			var apierr *kvapi.APIError
			if errors.As(err, &apierr) {
				var conerr *kvapi.ConflictError
				if errors.As(err, &conerr) {
					if conerr.ExpectedVersion == nversion {
						return nil
					}
					break
				}
				continue
			}

			var conerr *kvapi.ConflictError
			if errors.As(err, &conerr) {
				if conerr.ExpectedVersion == nversion {
					return nil
				}
				break
			}

			return err
		}
	}
}
