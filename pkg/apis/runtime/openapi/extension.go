package openapi

// OpenAPI Extensions.
const (
	// ExtCliOperationName define the extension key to set the CLI operation name.
	ExtCliOperationName = "x-cli-operation-name"

	// ExtCliSchemaTypeName define the extension key to set the CLI operation params schema type.
	ExtCliSchemaTypeName = "x-cli-schema-type"

	// ExtCliIgnore define the extension key to ignore generate the api operation to the cli used api.json.
	ExtCliIgnore = "x-cli-ignore"

	// ExtCliCmdIgnore define the extension key to generate the operation to api.json but will not generate cli command.
	ExtCliCmdIgnore = "x-cli-cmd-ignore"

	// ExtCliOutputFormat define the output format set the CLI operation command.
	ExtCliOutputFormat = "x-cli-output-format"
)
