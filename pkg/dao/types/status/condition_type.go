package status

import (
	"reflect"
	"time"
)

// ConditionStatus is the value of status
type ConditionStatus string

// These are valid condition statuses.
// "ConditionStatusTrue" means a resource is in the condition.
// "ConditionStatusFalse" means a resource is not in the condition.
// "ConditionStatusUnknown" means a resource is in the condition or not.
const (
	ConditionStatusTrue    ConditionStatus = "True"
	ConditionStatusFalse   ConditionStatus = "False"
	ConditionStatusUnknown ConditionStatus = "Unknown"
)

// ConditionType is the type of status
type ConditionType string

// True set status value to True for object field .Status.Conditions, object must be a pointer
func (ct ConditionType) True(obj interface{}, message string) {
	setCondStatusAndMessage(obj, ct, ConditionStatusTrue, message)
}

// False set status value to False for object field .Status.Conditions, object must be a pointer
func (ct ConditionType) False(obj interface{}, message string) {
	setCondStatusAndMessage(obj, ct, ConditionStatusFalse, message)
}

// Unknown set status value to Unknown for object field .Status.Conditions, object must be a pointer
func (ct ConditionType) Unknown(obj interface{}, message string) {
	setCondStatusAndMessage(obj, ct, ConditionStatusUnknown, message)
}

// Status set status value to custom value for object field .Status.Conditions, object must be a pointer
func (ct ConditionType) Status(obj interface{}, status ConditionStatus) {
	setCondStatus(obj, ct, status)
}

// Message set message to conditionType for object field .Status.Conditions, object must be a pointer
func (ct ConditionType) Message(obj interface{}, message string) {
	setCondMessage(obj, ct, message)
}

// IsTrue check status value for object, object must be a pointer
func (ct ConditionType) IsTrue(obj interface{}) bool {
	return getCondStatus(obj, ct) == "True"
}

// IsFalse check status value for object, object must be a pointer
func (ct ConditionType) IsFalse(obj interface{}) bool {
	return getCondStatus(obj, ct) == "False"
}

// IsUnknown check status value for object, object must be a pointer
func (ct ConditionType) IsUnknown(obj interface{}) bool {
	return getCondStatus(obj, ct) == "Unknown"
}

// GetMessage get message from conditionType for object field .Status.Conditions
func (ct ConditionType) GetMessage(obj interface{}) string {
	return getCondMessage(obj, ct)
}

func setCondStatusAndMessage(obj interface{}, condType ConditionType, status ConditionStatus, message string) {
	if obj == nil || reflect.TypeOf(obj).Kind() != reflect.Ptr {
		return
	}

	stField, st := getStatus(obj)
	if st == nil {
		return
	}

	applyCondStatusAndMessage(st, condType, status, message)
	stField.Set(reflect.ValueOf(*st))
}

func setCondStatus(obj interface{}, condType ConditionType, status ConditionStatus) {
	if obj == nil || reflect.TypeOf(obj).Kind() != reflect.Ptr {
		return
	}

	stField, st := getStatus(obj)
	if st == nil {
		return
	}

	applyCondStatus(st, condType, status)
	stField.Set(reflect.ValueOf(*st))
}

func getCondStatus(obj interface{}, condType ConditionType) ConditionStatus {
	if obj == nil || reflect.TypeOf(obj).Kind() != reflect.Ptr {
		return ""
	}

	_, st := getStatus(obj)
	if st == nil {
		return ""
	}

	cond := getCond(st, condType)
	if cond == nil {
		return ""
	}
	return cond.Status
}

func setCondMessage(obj interface{}, condType ConditionType, message string) {
	if obj == nil || reflect.TypeOf(obj).Kind() != reflect.Ptr {
		return
	}

	stField, st := getStatus(obj)
	if st == nil {
		return
	}

	applyCondMessage(st, condType, message)
	stField.Set(reflect.ValueOf(*st))
}

func getCondMessage(obj interface{}, condType ConditionType) string {
	if obj == nil || reflect.TypeOf(obj).Kind() != reflect.Ptr {
		return ""
	}

	_, st := getStatus(obj)
	if st == nil {
		return ""
	}

	cond := getCond(st, condType)
	if cond != nil {
		return ""
	}
	return cond.Message
}

func getStatus(obj interface{}) (reflect.Value, *Status) {
	v := reflect.ValueOf(obj)
	if v.Type().Kind() == reflect.Ptr {
		v = v.Elem()
	}

	if v.Kind() != reflect.Struct {
		return reflect.Value{}, nil
	}

	field := v.FieldByName("Status")
	st, ok := field.Interface().(Status)
	if !ok {
		return reflect.Value{}, nil
	}
	return field, &st
}

func getCond(st *Status, condType ConditionType) *Condition {
	for i, v := range st.Conditions {
		if v.Type == condType {
			return &st.Conditions[i]
		}
	}
	return nil
}

func applyCondStatusAndMessage(st *Status, condType ConditionType, status ConditionStatus, message string) {
	var (
		now       = time.Now().UTC()
		condSlice = st.Conditions
		newCond   = Condition{
			Type:           condType,
			Status:         status,
			LastUpdateTime: now,
			Message:        message,
		}
	)

	if len(condSlice) == 0 {
		st.Conditions = []Condition{newCond}
		st.conditionChanged = true
	}

	var existed bool
	for i, v := range condSlice {
		if v.Type == condType {
			if v.Status != status || v.Message != message {
				condSlice[i].Status = status
				condSlice[i].Message = message
				condSlice[i].LastUpdateTime = now
				st.conditionChanged = true
			}
			existed = true
		}
	}

	if !existed {
		condSlice = append(condSlice, newCond)
		st.conditionChanged = true
	}
	st.Conditions = condSlice
}

func applyCondStatus(st *Status, condType ConditionType, status ConditionStatus) {
	var (
		now       = time.Now().UTC()
		condSlice = st.Conditions
		newCond   = Condition{
			Type:           condType,
			Status:         status,
			LastUpdateTime: now,
		}
	)

	if len(condSlice) == 0 {
		st.Conditions = []Condition{newCond}
		st.conditionChanged = true
	}

	var existed bool
	for i, v := range condSlice {
		if v.Type == condType {
			if v.Status != status {
				condSlice[i].Status = status
				condSlice[i].LastUpdateTime = now
				st.conditionChanged = true
			}
			existed = true
		}
	}

	if !existed {
		condSlice = append(condSlice, newCond)
		st.conditionChanged = true
	}
	st.Conditions = condSlice
}

func applyCondMessage(st *Status, condType ConditionType, message string) {
	if len(st.Conditions) == 0 {
		return
	}

	for i, v := range st.Conditions {
		if v.Type == condType {
			if st.Conditions[i].Message != message {
				st.Conditions[i].Message = message
				st.Conditions[i].LastUpdateTime = time.Now().UTC()
				st.conditionChanged = true
			}
		}
	}
}
