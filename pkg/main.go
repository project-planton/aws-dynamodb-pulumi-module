package pkg

import (
	"github.com/pkg/errors"
	"github.com/plantoncloud/planton-cloud-apis/zzgo/cloud/planton/apis/code2cloud/v1/aws/awsdynamodb"
	"github.com/pulumi/pulumi-aws/sdk/v6/go/aws"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type ResourceStack struct {
	StackInput *awsdynamodb.AwsDynamodbStackInput
	Labels     map[string]string
}

func (s *ResourceStack) Resources(ctx *pulumi.Context) error {
	locals := initializeLocals(ctx, s.StackInput)

	awsCredential := s.StackInput.AwsCredential

	//create aws provider using the credentials from the input
	awsProvider, err := aws.NewProvider(ctx,
		"classic-provider",
		&aws.ProviderArgs{
			AccessKey: pulumi.String(awsCredential.Spec.AccessKeyId),
			SecretKey: pulumi.String(awsCredential.Spec.SecretAccessKey),
			Region:    pulumi.String(awsCredential.Spec.Region),
		})

	createdDynamodbTable, err := table(ctx, locals, awsProvider)
	if err != nil {
		return errors.Wrap(err, "failed to create dynamo table resources")
	}

	if err = autoScale(ctx, locals, awsProvider, createdDynamodbTable); err != nil {
		return errors.Wrap(err, "failed to create dynamo db auto scaling resources")
	}
	return nil
}
