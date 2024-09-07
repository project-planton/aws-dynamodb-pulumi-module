package pkg

import (
	"github.com/pkg/errors"
	"github.com/plantoncloud/planton-cloud-apis/zzgo/cloud/planton/apis/code2cloud/v1/aws/awsdynamodb"
	"github.com/pulumi/pulumi-aws-native/sdk/go/aws"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type ResourceStack struct {
	StackInput *awsdynamodb.AwsDynamodbStackInput
	Labels     map[string]string
}

func (s *ResourceStack) Resources(ctx *pulumi.Context) error {
	locals := initializeLocals(ctx, s.StackInput)

	//create aws provider using the credentials from the input
	awsProvider, err := aws.NewProvider(ctx,
		"aws-provider",
		&aws.ProviderArgs{
			AccessKey:                 pulumi.String(s.StackInput.AwsCredential.Spec.AccessKeyId),
			SecretKey:                 pulumi.String(s.StackInput.AwsCredential.Spec.SecretAccessKey),
			Region:                    pulumi.String(s.StackInput.AwsCredential.Spec.Region),
			SkipCredentialsValidation: pulumi.Bool(true),
		})
	if err != nil {
		return errors.Wrap(err, "failed to create aws provider")
	}

	createdDynamodbTable, err := table(ctx, locals, awsProvider)
	if err != nil {
		return errors.Wrap(err, "failed to create dynamo table resources")
	}

	if err = autoScale(ctx, locals, awsProvider, createdDynamodbTable); err != nil {
		return errors.Wrap(err, "failed to create dynamo db auto scaling resources")
	}
	return nil
}
