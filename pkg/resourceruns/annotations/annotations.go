package annotations

import (
	"context"
	"fmt"

	"github.com/seal-io/walrus/pkg/auths/session"
	"github.com/seal-io/walrus/pkg/dao/model"
	"github.com/seal-io/walrus/pkg/dao/types"
	"github.com/seal-io/walrus/pkg/dao/types/object"
)

func GetSubjectID(entity *model.ResourceRun) (object.ID, error) {
	if entity == nil {
		return "", fmt.Errorf("resource is nil")
	}

	subjectIDStr := entity.Annotations[types.AnnotationSubjectID]

	return object.ID(subjectIDStr), nil
}

func SetSubjectID(ctx context.Context, resources ...*model.ResourceRun) error {
	sj, err := session.GetSubject(ctx)
	if err != nil {
		return err
	}

	for i := range resources {
		if resources[i].Annotations == nil {
			resources[i].Annotations = make(map[string]string)
		}
		resources[i].Annotations[types.AnnotationSubjectID] = string(sj.ID)
	}

	return nil
}
