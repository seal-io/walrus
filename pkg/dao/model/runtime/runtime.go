// SPDX-FileCopyrightText: 2023 Seal, Inc
// SPDX-License-Identifier: Apache-2.0

// GENERATED, DO NOT EDIT.

package runtime

import (
	"time"

	"github.com/seal-io/seal/pkg/dao/model/role"
	"github.com/seal-io/seal/pkg/dao/model/setting"
	"github.com/seal-io/seal/pkg/dao/model/subject"
	"github.com/seal-io/seal/pkg/dao/model/token"
	"github.com/seal-io/seal/pkg/dao/schema"
)

// The init function reads all schema descriptors with runtime code
// (default values, validators, hooks and policies) and stitches it
// to their package variables.
func init() {
	roleMixin := schema.Role{}.Mixin()
	roleMixinHooks0 := roleMixin[0].Hooks()
	role.Hooks[0] = roleMixinHooks0[0]
	roleMixinFields1 := roleMixin[1].Fields()
	_ = roleMixinFields1
	roleFields := schema.Role{}.Fields()
	_ = roleFields
	// roleDescCreateTime is the schema descriptor for createTime field.
	roleDescCreateTime := roleMixinFields1[0].Descriptor()
	// role.DefaultCreateTime holds the default value on creation for the createTime field.
	role.DefaultCreateTime = roleDescCreateTime.Default.(func() time.Time)
	// roleDescUpdateTime is the schema descriptor for updateTime field.
	roleDescUpdateTime := roleMixinFields1[1].Descriptor()
	// role.DefaultUpdateTime holds the default value on creation for the updateTime field.
	role.DefaultUpdateTime = roleDescUpdateTime.Default.(func() time.Time)
	// role.UpdateDefaultUpdateTime holds the default value on update for the updateTime field.
	role.UpdateDefaultUpdateTime = roleDescUpdateTime.UpdateDefault.(func() time.Time)
	// roleDescDomain is the schema descriptor for domain field.
	roleDescDomain := roleFields[0].Descriptor()
	// role.DefaultDomain holds the default value on creation for the domain field.
	role.DefaultDomain = roleDescDomain.Default.(string)
	// roleDescName is the schema descriptor for name field.
	roleDescName := roleFields[1].Descriptor()
	// role.NameValidator is a validator for the "name" field. It is called by the builders before save.
	role.NameValidator = roleDescName.Validators[0].(func(string) error)
	// roleDescDescription is the schema descriptor for description field.
	roleDescDescription := roleFields[2].Descriptor()
	// role.DefaultDescription holds the default value on creation for the description field.
	role.DefaultDescription = roleDescDescription.Default.(string)
	// roleDescPolicies is the schema descriptor for policies field.
	roleDescPolicies := roleFields[3].Descriptor()
	// role.DefaultPolicies holds the default value on creation for the policies field.
	role.DefaultPolicies = roleDescPolicies.Default.(schema.RolePolicies)
	// roleDescBuiltin is the schema descriptor for builtin field.
	roleDescBuiltin := roleFields[4].Descriptor()
	// role.DefaultBuiltin holds the default value on creation for the builtin field.
	role.DefaultBuiltin = roleDescBuiltin.Default.(bool)
	// roleDescSession is the schema descriptor for session field.
	roleDescSession := roleFields[5].Descriptor()
	// role.DefaultSession holds the default value on creation for the session field.
	role.DefaultSession = roleDescSession.Default.(bool)
	settingMixin := schema.Setting{}.Mixin()
	settingMixinHooks0 := settingMixin[0].Hooks()
	setting.Hooks[0] = settingMixinHooks0[0]
	settingMixinFields1 := settingMixin[1].Fields()
	_ = settingMixinFields1
	settingFields := schema.Setting{}.Fields()
	_ = settingFields
	// settingDescCreateTime is the schema descriptor for createTime field.
	settingDescCreateTime := settingMixinFields1[0].Descriptor()
	// setting.DefaultCreateTime holds the default value on creation for the createTime field.
	setting.DefaultCreateTime = settingDescCreateTime.Default.(func() time.Time)
	// settingDescUpdateTime is the schema descriptor for updateTime field.
	settingDescUpdateTime := settingMixinFields1[1].Descriptor()
	// setting.DefaultUpdateTime holds the default value on creation for the updateTime field.
	setting.DefaultUpdateTime = settingDescUpdateTime.Default.(func() time.Time)
	// setting.UpdateDefaultUpdateTime holds the default value on update for the updateTime field.
	setting.UpdateDefaultUpdateTime = settingDescUpdateTime.UpdateDefault.(func() time.Time)
	// settingDescValue is the schema descriptor for value field.
	settingDescValue := settingFields[1].Descriptor()
	// setting.DefaultValue holds the default value on creation for the value field.
	setting.DefaultValue = settingDescValue.Default.(string)
	// settingDescHidden is the schema descriptor for hidden field.
	settingDescHidden := settingFields[2].Descriptor()
	// setting.DefaultHidden holds the default value on creation for the hidden field.
	setting.DefaultHidden = settingDescHidden.Default.(bool)
	// settingDescEditable is the schema descriptor for editable field.
	settingDescEditable := settingFields[3].Descriptor()
	// setting.DefaultEditable holds the default value on creation for the editable field.
	setting.DefaultEditable = settingDescEditable.Default.(bool)
	// settingDescPrivate is the schema descriptor for private field.
	settingDescPrivate := settingFields[4].Descriptor()
	// setting.DefaultPrivate holds the default value on creation for the private field.
	setting.DefaultPrivate = settingDescPrivate.Default.(bool)
	subjectMixin := schema.Subject{}.Mixin()
	subjectMixinHooks0 := subjectMixin[0].Hooks()
	subject.Hooks[0] = subjectMixinHooks0[0]
	subjectMixinFields1 := subjectMixin[1].Fields()
	_ = subjectMixinFields1
	subjectFields := schema.Subject{}.Fields()
	_ = subjectFields
	// subjectDescCreateTime is the schema descriptor for createTime field.
	subjectDescCreateTime := subjectMixinFields1[0].Descriptor()
	// subject.DefaultCreateTime holds the default value on creation for the createTime field.
	subject.DefaultCreateTime = subjectDescCreateTime.Default.(func() time.Time)
	// subjectDescUpdateTime is the schema descriptor for updateTime field.
	subjectDescUpdateTime := subjectMixinFields1[1].Descriptor()
	// subject.DefaultUpdateTime holds the default value on creation for the updateTime field.
	subject.DefaultUpdateTime = subjectDescUpdateTime.Default.(func() time.Time)
	// subject.UpdateDefaultUpdateTime holds the default value on update for the updateTime field.
	subject.UpdateDefaultUpdateTime = subjectDescUpdateTime.UpdateDefault.(func() time.Time)
	// subjectDescKind is the schema descriptor for kind field.
	subjectDescKind := subjectFields[0].Descriptor()
	// subject.DefaultKind holds the default value on creation for the kind field.
	subject.DefaultKind = subjectDescKind.Default.(string)
	// subjectDescGroup is the schema descriptor for group field.
	subjectDescGroup := subjectFields[1].Descriptor()
	// subject.DefaultGroup holds the default value on creation for the group field.
	subject.DefaultGroup = subjectDescGroup.Default.(string)
	// subjectDescName is the schema descriptor for name field.
	subjectDescName := subjectFields[2].Descriptor()
	// subject.NameValidator is a validator for the "name" field. It is called by the builders before save.
	subject.NameValidator = subjectDescName.Validators[0].(func(string) error)
	// subjectDescDescription is the schema descriptor for description field.
	subjectDescDescription := subjectFields[3].Descriptor()
	// subject.DefaultDescription holds the default value on creation for the description field.
	subject.DefaultDescription = subjectDescDescription.Default.(string)
	// subjectDescMountTo is the schema descriptor for mountTo field.
	subjectDescMountTo := subjectFields[4].Descriptor()
	// subject.DefaultMountTo holds the default value on creation for the mountTo field.
	subject.DefaultMountTo = subjectDescMountTo.Default.(bool)
	// subjectDescLoginTo is the schema descriptor for loginTo field.
	subjectDescLoginTo := subjectFields[5].Descriptor()
	// subject.DefaultLoginTo holds the default value on creation for the loginTo field.
	subject.DefaultLoginTo = subjectDescLoginTo.Default.(bool)
	// subjectDescRoles is the schema descriptor for roles field.
	subjectDescRoles := subjectFields[6].Descriptor()
	// subject.DefaultRoles holds the default value on creation for the roles field.
	subject.DefaultRoles = subjectDescRoles.Default.(schema.SubjectRoles)
	// subjectDescPaths is the schema descriptor for paths field.
	subjectDescPaths := subjectFields[7].Descriptor()
	// subject.DefaultPaths holds the default value on creation for the paths field.
	subject.DefaultPaths = subjectDescPaths.Default.([]string)
	// subjectDescBuiltin is the schema descriptor for builtin field.
	subjectDescBuiltin := subjectFields[8].Descriptor()
	// subject.DefaultBuiltin holds the default value on creation for the builtin field.
	subject.DefaultBuiltin = subjectDescBuiltin.Default.(bool)
	tokenMixin := schema.Token{}.Mixin()
	tokenMixinHooks0 := tokenMixin[0].Hooks()
	token.Hooks[0] = tokenMixinHooks0[0]
	tokenMixinFields1 := tokenMixin[1].Fields()
	_ = tokenMixinFields1
	tokenFields := schema.Token{}.Fields()
	_ = tokenFields
	// tokenDescCreateTime is the schema descriptor for createTime field.
	tokenDescCreateTime := tokenMixinFields1[0].Descriptor()
	// token.DefaultCreateTime holds the default value on creation for the createTime field.
	token.DefaultCreateTime = tokenDescCreateTime.Default.(func() time.Time)
	// tokenDescUpdateTime is the schema descriptor for updateTime field.
	tokenDescUpdateTime := tokenMixinFields1[1].Descriptor()
	// token.DefaultUpdateTime holds the default value on creation for the updateTime field.
	token.DefaultUpdateTime = tokenDescUpdateTime.Default.(func() time.Time)
	// token.UpdateDefaultUpdateTime holds the default value on update for the updateTime field.
	token.UpdateDefaultUpdateTime = tokenDescUpdateTime.UpdateDefault.(func() time.Time)
}

const (
	Version = "v0.11.7"                                         // Version of ent codegen.
	Sum     = "h1:V+wKFh0jhAbY/FoU+PPbdMOf2Ma5vh07R/IdF+N/nFg=" // Sum of ent codegen.
)
