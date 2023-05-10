package cache

import (
	"strings"

	"github.com/seal-io/seal/pkg/dao/types"
	"github.com/seal-io/seal/utils/json"
)

const subjectPermissionKeyPrefix = "subject_permission:"

type SubjectPermission struct {
	Roles    types.SubjectRoles `json:"roles"`
	Policies types.RolePolicies `json:"policies"`
}

// StoreSubjectPermission stores the subject permission with the given subject.
func StoreSubjectPermission(subject string, permission SubjectPermission) {
	if subject == "" {
		return
	}
	bs, err := json.Marshal(permission)
	if err == nil {
		_ = cacher.Set(subjectPermissionKeyPrefix+subject, bs)
	}
}

// LoadSubjectPermission loads the subject permission via the given subject,
// if the subject is cached, returns the SubjectPermission,
// if the subject is not cached, returns a nil SubjectPermission.
func LoadSubjectPermission(subject string) (*SubjectPermission, bool) {
	if subject == "" {
		return nil, false
	}
	bs, _ := cacher.Get(subjectPermissionKeyPrefix + subject)
	if len(bs) > 0 {
		var permission SubjectPermission
		err := json.Unmarshal(bs, &permission)
		if err == nil {
			return &permission, true
		}
	}
	return nil, false
}

// CleanSubjectPermission cleans the subject permission of the given subject.
func CleanSubjectPermission(subject string) {
	if subject == "" {
		return
	}
	_ = cacher.Delete(subjectPermissionKeyPrefix + subject)
}

// CleanSubjectPermissions cleans all subject permissions.
func CleanSubjectPermissions() {
	it := cacher.Iterator()
	for it.SetNext() {
		e, err := it.Value()
		if err != nil {
			break
		}
		key := e.Key()
		if strings.HasPrefix(key, subjectPermissionKeyPrefix) {
			_ = cacher.Delete(key)
		}
	}
}
