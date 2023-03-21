package cache

import (
	"context"

	"github.com/seal-io/seal/pkg/dao/types"
	"github.com/seal-io/seal/utils/cache"
	"github.com/seal-io/seal/utils/json"
)

const subjectPermissionKeyPrefix = "subject_permission:"

type SubjectPermission struct {
	Roles    types.SubjectRoles `json:"roles"`
	Policies types.RolePolicies `json:"policies"`
}

// StoreSubjectPermission stores the subject permission with the given subject.
func StoreSubjectPermission(ctx context.Context, subject string, permission SubjectPermission) {
	if subject == "" {
		return
	}
	var bs, err = json.Marshal(permission)
	if err == nil {
		_ = cacher.Set(ctx, subjectPermissionKeyPrefix+subject, bs)
	}
}

// LoadSubjectPermission loads the subject permission via the given subject,
// if the subject is cached, returns the SubjectPermission,
// if the subject is not cached, returns a nil SubjectPermission.
func LoadSubjectPermission(ctx context.Context, subject string) (*SubjectPermission, bool) {
	if subject == "" {
		return nil, false
	}
	var bs, _ = cacher.Get(ctx, subjectPermissionKeyPrefix+subject)
	if len(bs) > 0 {
		var permission SubjectPermission
		var err = json.Unmarshal(bs, &permission)
		if err == nil {
			return &permission, true
		}
	}
	return nil, false
}

// CleanSubjectPermission cleans the subject permission of the given subject.
func CleanSubjectPermission(ctx context.Context, subject string) {
	if subject == "" {
		return
	}
	_ = cacher.Delete(ctx, subjectPermissionKeyPrefix+subject)
}

// CleanSubjectPermissions cleans all subject permissions.
func CleanSubjectPermissions(ctx context.Context) {
	_ = cacher.Iterate(ctx, cache.HasPrefix(subjectPermissionKeyPrefix),
		func(ctx context.Context, e cache.Entry) (bool, error) {
			_ = cacher.Delete(ctx, e.Key())
			return true, nil
		})
}
