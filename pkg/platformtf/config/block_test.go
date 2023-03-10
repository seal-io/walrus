package config

import (
	"reflect"
	"testing"
)

func TestBlock_Raw(t *testing.T) {
	testCases := []struct {
		block    *Block
		expected map[string]interface{}
	}{
		{
			block: &Block{
				Type:   "resource",
				Labels: []string{"aws_instance", "test"},
				Attributes: map[string]interface{}{
					"ami":           "ami-0c55b159cbfafe1f0",
					"instance_type": "t2.micro",
				},
			},
			expected: map[string]interface{}{
				"resource": map[string]interface{}{
					"aws_instance": map[string]interface{}{
						"test": map[string]interface{}{
							"ami":           "ami-0c55b159cbfafe1f0",
							"instance_type": "t2.micro",
						},
					},
				},
			},
		},
		{
			block: &Block{
				Type:   "provider",
				Labels: []string{"aws"},
				Attributes: map[string]interface{}{
					"region": "us-east-1",
				},
			},
			expected: map[string]interface{}{
				"provider": map[string]interface{}{
					"aws": map[string]interface{}{
						"region": "us-east-1",
					},
				},
			},
		},
		{
			block: &Block{
				Type:   "module",
				Labels: []string{"mysql"},
				Attributes: map[string]interface{}{
					"source":                  "github.com/terraform-aws-modules/terraform-aws-rds",
					"identifier":              "test",
					"engine":                  "mysql",
					"engine_version":          "5.7",
					"instance_class":          "db.t2.micro",
					"allocated_storage":       "20",
					"username":                "test",
					"password":                "test",
					"port":                    "3306",
					"backup_retention_period": "7",
					"backup_window":           "03:00-04:00",
					"maintenance_window":      "Mon:03:00-Mon:04:00",
				},
			},
			expected: map[string]interface{}{
				"module": map[string]interface{}{
					"mysql": map[string]interface{}{
						"source":                  "github.com/terraform-aws-modules/terraform-aws-rds",
						"identifier":              "test",
						"engine":                  "mysql",
						"engine_version":          "5.7",
						"instance_class":          "db.t2.micro",
						"allocated_storage":       "20",
						"username":                "test",
						"password":                "test",
						"port":                    "3306",
						"backup_retention_period": "7",
						"backup_window":           "03:00-04:00",
						"maintenance_window":      "Mon:03:00-Mon:04:00",
					},
				},
			},
		},
		{
			block: &Block{
				Type:   "data",
				Labels: []string{"aws_ami", "test"},
				Attributes: map[string]interface{}{
					"owners": []string{"amazon"},
					"filter": map[string]interface{}{
						"name":   "name",
						"values": []string{"amzn-ami-hvm-2018.03.0.20180409-x86_64-gp2"},
					},
				},
			},
			expected: map[string]interface{}{
				"data": map[string]interface{}{
					"aws_ami": map[string]interface{}{
						"test": map[string]interface{}{
							"owners": []string{"amazon"},
							"filter": map[string]interface{}{
								"name":   "name",
								"values": []string{"amzn-ami-hvm-2018.03.0.20180409-x86_64-gp2"},
							},
						},
					},
				},
			},
		},
	}

	for _, tc := range testCases {
		actual := tc.block.ToNestedMap()
		if !reflect.DeepEqual(actual, tc.expected) {
			t.Errorf("expected %#v, got %#v", tc.expected, actual)
		}
	}
}

func TestBlock_Print(t *testing.T) {
	testCases := []struct {
		block    *Block
		expected []string
	}{
		{
			block: &Block{
				Type:   "resource",
				Labels: []string{"aws_instance", "test"},
				Attributes: map[string]interface{}{
					"ami":           "ami-0c55b159cbfafe1f0",
					"instance_type": "t2.micro",
				},
			},
			expected: []string{
				`resource "aws_instance" "test" {
  ami           = "ami-0c55b159cbfafe1f0"
  instance_type = "t2.micro"
}
`,
				`resource "aws_instance" "test" {
  instance_type = "t2.micro"
  ami           = "ami-0c55b159cbfafe1f0"
},
`,
			},
		},
		{
			block: &Block{
				Type:   "provider",
				Labels: []string{"aws"},
				Attributes: map[string]interface{}{
					"region": "us-east-1",
				},
			},
			expected: []string{
				`provider "aws" {
  region = "us-east-1"
}
`,
			},
		},
		{
			block: &Block{
				Type:   "module",
				Labels: []string{"mysql"},
				Attributes: map[string]interface{}{
					"source":   "github.com/terraform-aws-modules/terraform-aws-rds",
					"username": "test",
					"password": "test",
				},
			},
			expected: []string{
				`module "mysql" {
  source   = "github.com/terraform-aws-modules/terraform-aws-rds"
  username = "test"
  password = "test"
}
`,
				`module "mysql" {
  source   = "github.com/terraform-aws-modules/terraform-aws-rds"
  password = "test"
  username = "test"
}
`,
				`module "mysql" {
  password = "test"
  source   = "github.com/terraform-aws-modules/terraform-aws-rds"
  username = "test"
}
`,
				`module "mysql" {
  password = "test"
  username = "test"
  source   = "github.com/terraform-aws-modules/terraform-aws-rds"
}
`,
				`module "mysql" {
  username = "test"
  source   = "github.com/terraform-aws-modules/terraform-aws-rds"
  password = "test"
}
`,
				`module "mysql" {
  username = "test"
  password = "test"
  source   = "github.com/terraform-aws-modules/terraform-aws-rds"
}
`,
			},
		},
		{
			block: &Block{
				Type:   "data",
				Labels: []string{"aws_ami", "test"},
				Attributes: map[string]interface{}{
					"owners": []string{"amazon"},
					"filter": map[string]interface{}{
						"name":   "name",
						"values": []string{"amzn-ami-hvm-2018.03.0.20180409-x86_64-gp2"},
					},
				},
			},
			expected: []string{
				`data "aws_ami" "test" {
  owners = ["amazon"]

  filter {
    name   = "name"
    values = ["amzn-ami-hvm-2018.03.0.20180409-x86_64-gp2"]
}
`,
				`data "aws_ami" "test" {
  filter {
    name   = "name"
    values = ["amzn-ami-hvm-2018.03.0.20180409-x86_64-gp2"]
  }

  owners = ["amazon"]
}
`,
			},
		},
	}

	for _, tc := range testCases {
		mapObjects := make(map[string]struct{}, 0)
		actual, err := tc.block.Print("hcl", mapObjects)
		if err != nil {
			t.Errorf("unexpected error: %s", err)
		}

		found := false
		for _, expected := range tc.expected {
			if string(actual) == expected {
				found = true
				break
			}
		}

		if !found {
			t.Errorf("expected %#v, got %#v", tc.expected, string(actual))
		}
	}
}
