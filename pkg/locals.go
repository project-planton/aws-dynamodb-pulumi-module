package pkg

import (
	awsdynamodbv1 "buf.build/gen/go/plantoncloud/project-planton/protocolbuffers/go/project/planton/provider/aws/awsdynamodb/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type Locals struct {
	AwsDynamodb *awsdynamodbv1.AwsDynamodb
	Labels      map[string]string
}

func initializeLocals(ctx *pulumi.Context, stackInput *awsdynamodbv1.AwsDynamodbStackInput) *Locals {
	locals := &Locals{}

	//assign value for the locals variable to make it available across the project
	locals.AwsDynamodb = stackInput.Target

	return locals
}
