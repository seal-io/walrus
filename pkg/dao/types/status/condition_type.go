package status

import (
	"reflect"
	"time"
)

// ConditionStatus is the value of status.
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

// ConditionType is the type of status.
type ConditionType string

func (ct ConditionType) String() string {
	return string(ct)
}

// True set status value to True for object field .Status.Conditions,
// object must be a pointer.
func (ct ConditionType) True(obj any, message string) {
	setCondStatusAndMessage(obj, ct, ConditionStatusTrue, message)
}

// False set status value to False for object field .Status.Conditions,
// object must be a pointer.
func (ct ConditionType) False(obj any, message string) {
	setCondStatusAndMessage(obj, ct, ConditionStatusFalse, message)
}

// Unknown set status value to Unknown for object field .Status.Conditions,
// object must be a pointer.
func (ct ConditionType) Unknown(obj any, message string) {
	setCondStatusAndMessage(obj, ct, ConditionStatusUnknown, message)
}

// Status set status value to custom value for object field .Status.Conditions,
// object must be a pointer.
func (ct ConditionType) Status(obj any, status ConditionStatus) {
	setCondStatus(obj, ct, status)
}

// Remove drop status from the object field .Status.Conditions,
// object must be a pointer.
func (ct ConditionType) Remove(obj any) {
	delCondStatus(obj, ct)
}

// Reset clean the object field .Status.Conditions,
// and set the status as Unknown type into the object field .Status.Conditions,
// object must be a pointer.
func (ct ConditionType) Reset(obj any, message string) {
	resetCondStatus(obj, ct, ConditionStatusUnknown, message)
}

// Message set message to conditionType for object field .Status.Conditions,
// object must be a pointer.
func (ct ConditionType) Message(obj any, message string) {
	setCondMessage(obj, ct, message)
}

// IsTrue check status value for object,
// object must be a pointer.
func (ct ConditionType) IsTrue(obj any) bool {
	s, _ := getCondStatus(obj, ct)
	return s == ConditionStatusTrue
}

// IsFalse check status value for object,
// object must be a pointer.
func (ct ConditionType) IsFalse(obj any) bool {
	s, _ := getCondStatus(obj, ct)
	return s == ConditionStatusFalse
}

// IsUnknown check status value for object,
// object must be a pointer.
func (ct ConditionType) IsUnknown(obj any) bool {
	s, _ := getCondStatus(obj, ct)
	return s == ConditionStatusUnknown
}

// Exist returns true if the status is existed,
// object must be a pointer.
func (ct ConditionType) Exist(obj any) bool {
	_, exist := getCondStatus(obj, ct)
	return exist
}

// GetMessage get message from conditionType for object field .Status.Conditions.
func (ct ConditionType) GetMessage(obj any) string {
	return getCondMessage(obj, ct)
}

func setCondStatusAndMessage(obj any, condType ConditionType, status ConditionStatus, message string) {
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

func setCondStatus(obj any, condType ConditionType, status ConditionStatus) {
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

func delCondStatus(obj any, condType ConditionType) {
	if obj == nil || reflect.TypeOf(obj).Kind() != reflect.Ptr {
		return
	}

	stField, st := getStatus(obj)
	if st == nil {
		return
	}

	if len(st.Conditions) == 0 {
		return
	}

	for i, v := range st.Conditions {
		if v.Type != condType {
			continue
		}

		switch s := len(st.Conditions); {
		case i == 0 && s == 1:
			st.Conditions = st.Conditions[:0]
		case i == 0:
			st.Conditions = st.Conditions[1:]
		case i == s-1:
			st.Conditions = st.Conditions[:i]
		default:
			st.Conditions = append(st.Conditions[:i],
				st.Conditions[i+1:]...)
		}
		// NB(thxCode): do not skip this loop until we check all condition items,
		// which can avoid duplicate items in condition list.
	}

	stField.Set(reflect.ValueOf(*st))
}

func resetCondStatus(obj any, condType ConditionType, status ConditionStatus, message string) {
	if obj == nil || reflect.TypeOf(obj).Kind() != reflect.Ptr {
		return
	}

	stField, st := getStatus(obj)
	if st == nil {
		return
	}

	if st.Conditions != nil {
		// NB(thxCode): clean conditions.
		st.Conditions = st.Conditions[:0]
	}

	applyCondStatusAndMessage(st, condType, status, message)
	stField.Set(reflect.ValueOf(*st))
}

func getCondStatus(obj any, condType ConditionType) (ConditionStatus, bool) {
	if obj == nil || reflect.TypeOf(obj).Kind() != reflect.Ptr {
		return "", false
	}

	_, st := getStatus(obj)
	if st == nil {
		return "", false
	}

	cond := getCond(st, condType)
	if cond == nil {
		return "", false
	}

	return cond.Status, true
}

func setCondMessage(obj any, condType ConditionType, message string) {
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

func getCondMessage(obj any, condType ConditionType) string {
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

	return cond.Message
}

func getStatus(obj any) (reflect.Value, *Status) {
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
