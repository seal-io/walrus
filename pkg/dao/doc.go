// Package dao collects the tools for accessing with Database.
//
// # TL;DR
//
//   - leverage validation syntax, like `NotEmpty` or `Validate` in basic types,
//     if we want to reduce the validation from the DAO util functions.
//   - leverage `Nillable` syntax in basic types, if we want to reduce the Go zero value validation.
//   - leverage `Default` syntax in complex types and always check in DAO util functions,
//     if we don't want a DB JSON search errors.
//
// # Style of Schema Definition
//
//	Ent provides some syntactic sugars to assist scheming, let's declare the style of DAO field schema.
//
//	```
//	field declaration.
//	  generation syntax.
//	  ddl syntax.
//	  creation syntax.
//	  modification syntax.
//	  other syntax
//	```
//
//	- field name always be camel case starting with a lower letter.
//	- field generation, composited with `Comment`, `Nillable`, `StructTag`, `GoType`.
//	- field validation, select from `NotEmpty`, `Validate`.
//	- field creation, select from `Default`.
//	- field modification, select from `UpdateDefault`, `Immutable`.
//	- field generation configuration, `Annotations`.
//	- ddl creation, composited with `StorageKey`, `SchemaType`, `Optional`.
//
//	it's worth noticing that reduce the use of `Optional` if we not sure whether to field can be a query condition. `Optional` allows **null** DB value,  which could cause a DB query accident, and it's not easy to migrate with an automation-generated DDL.
//
// # Logic for Creator
//
//	XXXCreate(s)(mc model.ClientSet, input <Entity or Entity Slice>)
//
//	```
//	1. validate required fields from the input.
//	2. get DAO creator.
//	3. call `SetXXX` directly for required fields.
//	  3.1. if we scheme any type with `NotEmpty` explicitly, we can use `SetXXX` directly,
//	       ent errors if it is empty.
//	  3.2. if we cannot scheme any type with `NotEmpty` explicitly, like JSON/Strings field,
//	       we need to check before `SetXXX` here.
//	4. call `SetXXX` or `SetNillableXXX` if optional fields satisfy something.
//	  4.1. if we scheme basic type without `Nillable`, `Default` or with `Optional`, like Bool/String field,
//	       we can use `SetXXX` directly.
//	  4.2. if we scheme any type with `Nillable` explicitly,
//	       we can use `SetNillableXXX` directly, ent defaults if it is nil.
//	  4.3. if we scheme any type with `Default` explicitly,
//	       since ent holds the pointer of field and only defaults field if holds a nil pointer,
//	       we also need to check Go zero value before `SetXXX`.
//	  4.4. if we cannot scheme any type with `Nillable` explicitly, like JSON/Strings field,
//	       we need to check before `SetXXX` here.
//	```
//
// # Logic for Updater
//
//	XXXUpdate(s)(mc model.ClientSet, input <Entity or Entity Slice>)
//
//	```
//	1. validate prediction from the input.
//	2. get DAO updater.
//	3. call `SetXXX` or `SetNillableXXX` if the fields should be updated.
//	  3.1. if we want a field to configure as zero, like blank string,
//	       it's better to scheme the field with `Nillable`.
//	```
package dao
