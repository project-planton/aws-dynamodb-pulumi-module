package main

import (
	"github.com/pkg/errors"
	"github.com/project-planton/aws-dynamodb-pulumi-module/pkg"
	awsdynamodbv1 "github.com/project-planton/project-planton/apis/go/project/planton/provider/aws/awsdynamodb/v1"
	"github.com/project-planton/project-planton/pkg/pulmod/stackinput"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		stackInput := &awsdynamodbv1.AwsDynamodbStackInput{}

		if err := stackinput.LoadStackInput(ctx, stackInput); err != nil {
			return errors.Wrap(err, "failed to load stack-input")
		}

		return pkg.Resources(ctx, stackInput)
	})
}
