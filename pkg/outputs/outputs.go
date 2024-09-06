package outputs

import (
	"github.com/plantoncloud/planton-cloud-apis/zzgo/cloud/planton/apis/code2cloud/v1/aws/awsdynamodb"
	"github.com/plantoncloud/stack-job-runner-golang-sdk/pkg/automationapi/autoapistackoutput"
	"github.com/pulumi/pulumi/sdk/v3/go/auto"
)

const (
	TableName                          = "table-name"
	TableArn                           = "table-arn"
	TableStreamArn                     = "table-stream-arn"
	AutoscalingReadPolicyArn           = "autoscaling-read-policy-arn"
	AutoscalingWritePolicyArn          = "autoscaling-write-policy-arn"
	AutoscalingIndexReadPolicyArnList  = "autoscaling-index-read-policy-arn-list"
	AutoscalingIndexWritePolicyArnList = "autoscaling-index-write-policy-arn-list"
)

func PulumiOutputsToStackOutputsConverter(pulumiOutputs auto.OutputMap,
	input *awsdynamodb.AwsDynamodbStackInput) *awsdynamodb.AwsDynamodbStackOutputs {
	return &awsdynamodb.AwsDynamodbStackOutputs{
		TableName:                          autoapistackoutput.GetVal(pulumiOutputs, TableName),
		TableArn:                           autoapistackoutput.GetVal(pulumiOutputs, TableArn),
		TableStreamArn:                     autoapistackoutput.GetVal(pulumiOutputs, TableStreamArn),
		AutoscalingReadPolicyArn:           autoapistackoutput.GetVal(pulumiOutputs, AutoscalingReadPolicyArn),
		AutoscalingWritePolicyArn:          autoapistackoutput.GetVal(pulumiOutputs, AutoscalingWritePolicyArn),
		AutoscalingIndexReadPolicyArnList:  autoapistackoutput.GetVal(pulumiOutputs, AutoscalingIndexReadPolicyArnList),
		AutoscalingIndexWritePolicyArnList: autoapistackoutput.GetVal(pulumiOutputs, AutoscalingIndexWritePolicyArnList),
	}
}
