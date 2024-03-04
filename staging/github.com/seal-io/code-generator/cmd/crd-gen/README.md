# CRD Gen

This directory contains the code generators for the CRD (Custom Resource Definition) API.

# Marker Usage

- Package Markers
    - `+groupName=example.domain.io`: Specify the group name.
    - `+versionName=v1`: Specify the version.
- Type Markers
    - `+k8s:crd-gen:resource:scope="Namespaced",categories=["..."],shortName=["..."],plural="...",subResources=["status","scale"]`:
      Specify the resource.
        - `scope`: Select from `Namespaced` or `Cluster`.
        - `categories`, `shortName` and `subResources`: Specify in the JSON array format.
        - `subResources` only supports `status` and `scale`.
    - `+k8s:crd-gen:printcolumn:name="Name",type="string",jsonPath=".metadata.name",format="string",priority=0`:
      Specify the print column, `name`, `type`, `jsonPath` are required.
        - `jsonPath`: Specify the JSON path.
        - `priority`: Specify in the integer format.
- Field Markers
    - `+k8s:validation:default={"name":"abc","age":10}`: Specify the default value of the field, in JSON format.
      - do not compatible with `openapi-gen`, please use `+default=` replacing `+k8s:validation:default=` on top of the Field for `openapi-gen`.
    - `+k8s:validation:example={"name":"abc","age":10}`: Specify the example value of the field, in JSON format.
    - `+k8s:validation:enum=["a","b","c"]`: Specify the enum value of the field, in JSON array format.
      - do not compatible with `openapi-gen`, please use `+enum` on top of the enumeration Type for `openapi-gen`.
    - `+k8s:validation:maximum=10`: Specify the maximum value of the number field.
    - `+k8s:validation:minimum=1`: Specify the minimum value of the number field.
    - `+k8s:validation:exclusiveMaximum`: Specify whether be exclusive the minimum value of the number field.
    - `+k8s:validation:exclusiveMinimum`: Specify whether be exclusive the minimum value of the number field.
    - `+k8s:validation:multipleOf=5`: Specify the basic multiple of the number field.
    - `+k8s:validation:maxItems=10`: Specify the maximum items of the array field.
    - `+k8s:validation:minItems=1`: Specify the minimum items of the array field.
    - `+k8s:validation:uniqueItems=true`: Specify the items of the array field are unique or not.
    - `+k8s:validation:maxProperties=true`: Specify the maximum items of the struct/map field.
    - `+k8s:validation:minProperties=true`: Specify the minimum items of the struct/map field.
    - `+k8s:validation:maxLength=10`: Specify the maximum length of the string field.
    - `+k8s:validation:minLength=1`: Specify the minimum length of the string field.
    - `+k8s:validation:format="date"`: Specify the format of the string field.
    - `+k8s:validation:pattern="^a.*$"`: Specify the pattern of the field, in Go Regexp format.
    - `+k8s:validation:preserveUnknownFields"`: Specify whether to preserve the unknown (sub) fields.
    - `+k8s:validation:embeddedResource`: Specify whether to embed the resource to the field, the resource that has
      apiVersion, kind and metadata fields.
    - `+k8s:validation:cel[?]:rule=""`, `+k8s:validation:cel[?]:rule> [block]\n`: Specify the CEL (Common Expression
      Language) rule of the field, it is required to have at least one rule.
    - `+k8s:validation:cel[?]:message=""`, `+k8s:validation:cel[?]:message> [block]\n`: Specify the message if the given
      data violates the CEL rule.
    - `+k8s:validation:cel[?]:messageExpression=""`, `+k8s:validation:cel[?]:messageExpression> [block]\n`: Specify the
      evaluative message if the given data violates the CEL rule.
    - `+k8s:validation:cel[?]:reason="FieldValueInvalid"`: Specify the violation reason, select
      from `FieldValueInvalid`, `FieldValueForbidden`, `FieldValueRequired` or `FieldValueDuplicate`.
    - `+k8s:validation:cel[?]:fieldPath=""`: Specify the field path for violation.
    - `+k8s:validation:cel[?]:optionalOldSelf=""`: Inject `oldSelf` to avoid not validating when creating or introducing
      new value updates.
    - `+nullable`, `+nullable=false`: Specify whether the field is nullable or not.
    - `+optional`: Specify whether the field is optional, by default all fields are required.
    - `+listType="atomic"`: Specify the list type of the field, select from `map`, `set` or `atomic`.
        - `atomic`: Default, all the fields are treated as one unit, any changes have to replace the entire list.
        - `map`: It needs to have a key field (`+listMapKey=`), which will be used to build an associative list.
        - `set`: Fields need to be "scalar", and there can be only one occurrence of each.
    - `+listMapKey="name"`: Specify the (sub) field of the list type to be the map key, it is required when the list
      type
      is `map`.
    - `+mapType="granular"`: Specify the map type of the field, select from `granular` or `atomic`.
        - `atomic`: All fields are treated as one unit, any changes have to replace the entire map.
        - `granular`: Default, items in the map are independent of each other, and can be manipulated by different
          actors.

# License

Copyright (c) 2024 [Seal, Inc.](https://seal.io)

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at [LICENSE](../../LICENSE) file for details.

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.