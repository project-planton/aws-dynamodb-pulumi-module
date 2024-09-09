package main

import (
	"github.com/plantoncloud/aws-dynamodb-pulumi-module/pkg"
	"github.com/plantoncloud/planton-cloud-apis/zzgo/cloud/planton/apis/code2cloud/v1/aws/awsdynamodb"
	"github.com/plantoncloud/planton-cloud-apis/zzgo/cloud/planton/apis/commons/apiresource"
	"github.com/plantoncloud/planton-cloud-apis/zzgo/cloud/planton/apis/connect/v1/awscredential"
	"github.com/plantoncloud/planton-cloud-apis/zzgo/cloud/planton/apis/connect/v1/pulumibackendcredential"
	"github.com/plantoncloud/planton-cloud-apis/zzgo/cloud/planton/apis/connect/v1/pulumibackendcredential/enums/pulumibackendtype"
	"github.com/plantoncloud/planton-cloud-apis/zzgo/cloud/planton/apis/iac/v1/stackjob"
	"github.com/plantoncloud/planton-cloud-apis/zzgo/cloud/planton/apis/iac/v1/stackjob/progress/progressstatus"
	"github.com/plantoncloud/planton-cloud-apis/zzgo/cloud/planton/apis/resourcemanager/v1/environment"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"os"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		s := pkg.ResourceStack{
			StackInput: &awsdynamodb.AwsDynamodbStackInput{
				ApiResource: &awsdynamodb.AwsDynamodb{
					ApiVersion: "code2cloud.planton.cloud/v1",
					Kind:       "AwsDynamodb",
					Metadata: &apiresource.ApiResourceMetadata{
						Name: "orders",
						Id:   "awsdyn-planton-cloud-aws-module-test-orders",
					},
					Spec: &awsdynamodb.AwsDynamodbSpec{
						EnvironmentInfo: &environment.ApiResourceEnvironmentInfo{
							EnvId: os.Getenv("ENV_ID"),
						},
						StackJobSettings: &stackjob.StackJobSettings{
							PulumiBackendCredentialId: os.Getenv("PULUMI_BACKEND_CREDENTIAL_ID"),
							StackJobRunnerId:          os.Getenv("STACK_JOB_RUNNER_ID"),
						},
						Table: &awsdynamodb.AwsDynamodbTable{
							TableName:   "orders",
							BillingMode: "PAY_PER_REQUEST",
							HashKey: &awsdynamodb.AwsDynamodbTableAttribute{
								Name: "HashKey",
								Type: "S",
							},
							RangeKey: &awsdynamodb.AwsDynamodbTableAttribute{
								Name: "RangeKey",
								Type: "N",
							},
							AutoScale: &awsdynamodb.AwsDynamodbAutoScaleCapacity{
								IsEnabled: true,
							},
						},
					},
				},
				AwsCredential: &awscredential.AwsCredential{
					Spec: &awscredential.AwsCredentialSpec{
						AccessKeyId:     os.Getenv("AWS_ACCESS_KEY_ID"),
						SecretAccessKey: os.Getenv("AWS_SECRET_ACCESS_KEY"),
						Region:          os.Getenv("AWS_REGION"),
					},
				},
				PulumiBackendCredential: &pulumibackendcredential.PulumiBackendCredential{
					Spec: &pulumibackendcredential.PulumiBackendCredentialSpec{
						HttpBackend: &pulumibackendcredential.PulumiBackendCredentialHttpBackendSpec{
							AccessToken: os.Getenv("PULUMI_ACCESS_TOKEN"),
							ApiUrl:      os.Getenv("PULUMI_API_URL"),
						},
						PulumiBackendType:  pulumibackendtype.PulumiBackendType_http,
						PulumiOrganization: os.Getenv("PULUMI_ORGANIZATION"),
					},
				},
				StackJob: &stackjob.StackJob{
					Metadata: &apiresource.ApiResourceMetadata{
						Id: "awsdyn-stack-job",
					},
					Spec: &stackjob.StackJobSpec{
						EnvId:           "planton-cloud-aws-module-test",
						ResourceId:      "awsdyn-planton-cloud-aws-module-test-orders",
						PulumiStackName: "awsdyn-planton-cloud-aws-module-test-orders",
					},
					Status: &stackjob.StackJobStatus{
						PulumiOperations: &stackjob.StackJobStatusPulumiOperationsStatus{
							Apply: &progressstatus.StackJobProgressPulumiOperationStatus{
								IsRequired: true,
							},
							ApplyPreview: &progressstatus.StackJobProgressPulumiOperationStatus{
								IsRequired: false,
							},
							Destroy: &progressstatus.StackJobProgressPulumiOperationStatus{
								IsRequired: false,
							},
							DestroyPreview: &progressstatus.StackJobProgressPulumiOperationStatus{
								IsRequired: false,
							},
						},
					},
				},
			},
		}
		err := s.Resources(ctx)
		if err != nil {
			return err
		}
		return nil
	})
}
