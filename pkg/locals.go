package pkg

import (
	"github.com/plantoncloud/planton-cloud-apis/zzgo/cloud/planton/apis/code2cloud/v1/aws/awsdynamodb"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type Locals struct {
	AwsDynamodb *awsdynamodb.AwsDynamodb
	Labels      map[string]string
}

func initializeLocals(ctx *pulumi.Context, stackInput *awsdynamodb.AwsDynamodbStackInput) *Locals {
	locals := &Locals{}

	//assign value for the locals variable to make it available across the project
	locals.AwsDynamodb = stackInput.ApiResource

	return locals
}
