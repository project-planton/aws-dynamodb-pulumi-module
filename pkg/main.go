package pkg

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/plantoncloud/planton-cloud-apis/zzgo/cloud/planton/apis/code2cloud/v1/aws/awsdynamodb"
	"github.com/plantoncloud/pulumi-module-golang-commons/pkg/provider/aws/pulumiawsprovider"
	"github.com/pulumi/pulumi-aws/sdk/v6/go/aws/appautoscaling"
	"github.com/pulumi/pulumi-aws/sdk/v6/go/aws/dynamodb"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type ResourceStack struct {
	Input  *awsdynamodb.AwsDynamodbStackInput
	Labels map[string]string
}

func (s *ResourceStack) Resources(ctx *pulumi.Context) error {
	//create aws provider using the credentials from the input
	awsProvider, err := pulumiawsprovider.GetNative(ctx, s.Input.AwsCredential)
	if err != nil {
		return errors.Wrap(err, "failed to create aws provider")
	}

	awsDynamodb := s.Input.ApiResource
	streamEnabled := awsDynamodb.Spec.Table.EnableStreams
	streamViewType := ""
	if len(awsDynamodb.Spec.Table.ReplicaRegionNames) > 0 {
		streamEnabled = true
	}

	if len(awsDynamodb.Spec.Table.ReplicaRegionNames) > 0 || awsDynamodb.Spec.Table.EnableStreams {
		streamViewType = awsDynamodb.Spec.Table.StreamViewType
	}

	var replicaArray = dynamodb.TableReplicaTypeArray{}
	for _, regionName := range awsDynamodb.Spec.Table.ReplicaRegionNames {
		replicaArray = append(replicaArray, &dynamodb.TableReplicaTypeArgs{
			KmsKeyArn:           nil,
			PointInTimeRecovery: pulumi.Bool(false),
			PropagateTags:       pulumi.Bool(false),
			RegionName:          pulumi.String(regionName),
		})
	}

	var attributeArray = dynamodb.TableAttributeArray{}
	for _, attribute := range awsDynamodb.Spec.Table.Attributes {
		attributeArray = append(attributeArray, &dynamodb.TableAttributeArgs{
			Name: pulumi.String(attribute.Name),
			Type: pulumi.String(attribute.Type),
		})
	}

	var globalSecondaryIndexArray = dynamodb.TableGlobalSecondaryIndexArray{}
	for _, globalSecondaryIndex := range awsDynamodb.Spec.Table.GlobalSecondaryIndexes {
		globalSecondaryIndexArray = append(globalSecondaryIndexArray, &dynamodb.TableGlobalSecondaryIndexArgs{
			Name:             pulumi.String(globalSecondaryIndex.Name),
			HashKey:          pulumi.String(globalSecondaryIndex.HashKey),
			RangeKey:         pulumi.String(globalSecondaryIndex.RangeKey),
			ReadCapacity:     pulumi.Int(globalSecondaryIndex.ReadCapacity),
			WriteCapacity:    pulumi.Int(globalSecondaryIndex.WriteCapacity),
			ProjectionType:   pulumi.String(globalSecondaryIndex.ProjectionType),
			NonKeyAttributes: pulumi.ToStringArray(globalSecondaryIndex.NonKeyAttributes),
		})
	}

	var localSecondaryIndexArray = dynamodb.TableLocalSecondaryIndexArray{}
	for _, localSecondaryIndex := range awsDynamodb.Spec.Table.LocalSecondaryIndexes {
		localSecondaryIndexArray = append(localSecondaryIndexArray, &dynamodb.TableLocalSecondaryIndexArgs{
			Name:             pulumi.String(localSecondaryIndex.Name),
			RangeKey:         pulumi.String(localSecondaryIndex.RangeKey),
			ProjectionType:   pulumi.String(localSecondaryIndex.ProjectionType),
			NonKeyAttributes: pulumi.ToStringArray(localSecondaryIndex.NonKeyAttributes),
		})
	}

	var serverSideEncryption *dynamodb.TableServerSideEncryptionArgs
	if awsDynamodb.Spec.Table.ServerSideEncryption != nil {
		serverSideEncryption = &dynamodb.TableServerSideEncryptionArgs{
			Enabled:   pulumi.Bool(awsDynamodb.Spec.Table.ServerSideEncryption.IsEnabled),
			KmsKeyArn: pulumi.StringPtr(awsDynamodb.Spec.Table.ServerSideEncryption.KmsKeyArn),
		}
	}

	var pointInTimeRecovery *dynamodb.TablePointInTimeRecoveryArgs
	if awsDynamodb.Spec.Table.PointInTimeRecovery != nil {
		pointInTimeRecovery = &dynamodb.TablePointInTimeRecoveryArgs{
			Enabled: pulumi.Bool(awsDynamodb.Spec.Table.PointInTimeRecovery.IsEnabled),
		}
	}

	var ttl *dynamodb.TableTtlArgs
	if awsDynamodb.Spec.Table.Ttl != nil {
		ttl = &dynamodb.TableTtlArgs{
			Enabled:       pulumi.Bool(awsDynamodb.Spec.Table.Ttl.IsEnabled),
			AttributeName: pulumi.String(awsDynamodb.Spec.Table.Ttl.AttributeName),
		}
	}

	var importTable *dynamodb.TableImportTableArgs
	if awsDynamodb.Spec.Table.ImportTable != nil {
		inputFormatOptions := &dynamodb.TableImportTableInputFormatOptionsArgs{
			Csv: dynamodb.TableImportTableInputFormatOptionsCsvArgs{
				Delimiter:   pulumi.String(","),
				HeaderLists: pulumi.ToStringArray([]string{}),
			},
		}
		if awsDynamodb.Spec.Table.ImportTable.InputFormatOptions != nil && awsDynamodb.Spec.Table.ImportTable.InputFormatOptions.Csv != nil {
			inputFormatOptions = &dynamodb.TableImportTableInputFormatOptionsArgs{
				Csv: dynamodb.TableImportTableInputFormatOptionsCsvArgs{
					Delimiter:   pulumi.String(awsDynamodb.Spec.Table.ImportTable.InputFormatOptions.Csv.Delimiter),
					HeaderLists: pulumi.ToStringArray(awsDynamodb.Spec.Table.ImportTable.InputFormatOptions.Csv.HeaderLists),
				},
			}
		}

		s3BucketSource := &dynamodb.TableImportTableS3BucketSourceArgs{}
		if awsDynamodb.Spec.Table.ImportTable.S3BucketSource != nil {
			s3BucketSource = &dynamodb.TableImportTableS3BucketSourceArgs{
				Bucket:      pulumi.String(awsDynamodb.Spec.Table.ImportTable.S3BucketSource.Bucket),
				BucketOwner: pulumi.String(awsDynamodb.Spec.Table.ImportTable.S3BucketSource.BucketOwner),
				KeyPrefix:   pulumi.String(awsDynamodb.Spec.Table.ImportTable.S3BucketSource.KeyPrefix),
			}
		}
		importTable = &dynamodb.TableImportTableArgs{
			InputCompressionType: pulumi.String(awsDynamodb.Spec.Table.ImportTable.InputCompressionType),
			InputFormat:          pulumi.String(awsDynamodb.Spec.Table.ImportTable.InputFormat),
			InputFormatOptions:   inputFormatOptions,
			S3BucketSource:       s3BucketSource,
		}
	}

	createdDynamodbTable, err := dynamodb.NewTable(ctx, awsDynamodb.Metadata.Name, &dynamodb.TableArgs{
		Name:                      pulumi.String(awsDynamodb.Metadata.Name),
		BillingMode:               pulumi.String(awsDynamodb.Spec.Table.BillingMode),
		ReadCapacity:              pulumi.Int(awsDynamodb.Spec.Table.ReadCapacity),
		WriteCapacity:             pulumi.Int(awsDynamodb.Spec.Table.WriteCapacity),
		HashKey:                   pulumi.String(awsDynamodb.Spec.Table.HashKey),
		RangeKey:                  pulumi.String(awsDynamodb.Spec.Table.RangeKey),
		StreamEnabled:             pulumi.Bool(streamEnabled),
		StreamViewType:            pulumi.String(streamViewType),
		TableClass:                pulumi.String(awsDynamodb.Spec.Table.TableClass),
		DeletionProtectionEnabled: pulumi.Bool(awsDynamodb.Spec.Table.DeletionProtectionEnabled),
		ServerSideEncryption:      serverSideEncryption,
		PointInTimeRecovery:       pointInTimeRecovery,
		Ttl:                       ttl,
		Tags:                      pulumi.ToStringMap(s.Labels),
		Attributes:                attributeArray,
		GlobalSecondaryIndexes:    globalSecondaryIndexArray,
		LocalSecondaryIndexes:     localSecondaryIndexArray,
		Replicas:                  replicaArray,
		ImportTable:               importTable,
	}, pulumi.Provider(awsProvider))
	if err != nil {
		return errors.Wrap(err, "failed to create dynamo table resources")
	}

	enableAutoScale := true

	//read
	autoScaleMinReadCapacity := 5
	autoScaleMaxReadCapacity := 20
	autoScaleReadTarget := 50
	autoScaleReadTargetIndex := 50
	autoScaleMinReadCapacityIndex := 5
	autoScaleMaxReadCapacityIndex := 20

	//write
	autoScaleMinWriteCapacity := 5
	autoScaleMaxWriteCapacity := 20
	autoScaleWriteTarget := 50
	autoScaleWriteTargetIndex := 50
	autoScaleMinWriteCapacityIndex := 5
	autoScaleMaxWriteCapacityIndex := 20

	autoScaleScaleInCoolDown := 60
	autoScaleScaleOutCoolDown := 60

	if enableAutoScale && awsDynamodb.Spec.Table.BillingMode == "PROVISIONED" {

		readTarget, err := appautoscaling.NewTarget(ctx, "readTarget", &appautoscaling.TargetArgs{
			MaxCapacity:       pulumi.Int(autoScaleMaxReadCapacity),
			MinCapacity:       pulumi.Int(autoScaleMinReadCapacity),
			ResourceId:        pulumi.String("table/" + awsDynamodb.Metadata.Name),
			ScalableDimension: pulumi.String("dynamodb:table:ReadCapacityUnits"),
			ServiceNamespace:  pulumi.String("dynamodb"),
			Tags:              pulumi.ToStringMap(s.Labels),
		}, pulumi.Provider(awsProvider),
			pulumi.Parent(createdDynamodbTable),
			pulumi.DependsOn([]pulumi.Resource{createdDynamodbTable}))
		if err != nil {
			return errors.Wrap(err, "failed to create read target auto scaling resources")
		}

		_, err = appautoscaling.NewPolicy(ctx, "readPolicy", &appautoscaling.PolicyArgs{
			Name:              pulumi.Sprintf("DynamoDBReadCapacityUtilization:%s", readTarget.ID().ElementType()),
			PolicyType:        pulumi.String("TargetTrackingScaling"),
			ResourceId:        readTarget.ResourceId,
			ScalableDimension: readTarget.ScalableDimension,
			ServiceNamespace:  readTarget.ServiceNamespace,
			TargetTrackingScalingPolicyConfiguration: &appautoscaling.PolicyTargetTrackingScalingPolicyConfigurationArgs{
				PredefinedMetricSpecification: &appautoscaling.PolicyTargetTrackingScalingPolicyConfigurationPredefinedMetricSpecificationArgs{
					PredefinedMetricType: pulumi.String("DynamoDBReadCapacityUtilization"),
				},
				TargetValue:      pulumi.Float64(autoScaleReadTarget),
				ScaleInCooldown:  pulumi.Int(autoScaleScaleInCoolDown),
				ScaleOutCooldown: pulumi.Int(autoScaleScaleOutCoolDown),
			},
		}, pulumi.Provider(awsProvider),
			pulumi.Parent(readTarget),
			pulumi.DependsOn([]pulumi.Resource{createdDynamodbTable, readTarget}))
		if err != nil {
			return errors.Wrap(err, "failed to create read policy")
		}

		for _, index := range awsDynamodb.Spec.Table.GlobalSecondaryIndexes {
			indexTarget, err := appautoscaling.NewTarget(ctx, fmt.Sprintf("readTargetIndex-%s", index.Name), &appautoscaling.TargetArgs{
				MaxCapacity:       pulumi.Int(autoScaleMaxReadCapacityIndex),
				MinCapacity:       pulumi.Int(autoScaleMinReadCapacityIndex),
				ResourceId:        pulumi.String(fmt.Sprintf("table/%s/index/%s", awsDynamodb.Metadata.Name, index.Name)),
				ScalableDimension: pulumi.String("dynamodb:index:ReadCapacityUnits"),
				ServiceNamespace:  pulumi.String("dynamodb"),
				Tags:              pulumi.ToStringMap(s.Labels),
			}, pulumi.Provider(awsProvider),
				pulumi.Parent(readTarget),
				pulumi.DependsOn([]pulumi.Resource{createdDynamodbTable, readTarget}))
			if err != nil {
				return errors.Wrap(err, "failed to create read target index auto scaling resources")
			}

			// Create a Scaling Policy
			_, err = appautoscaling.NewPolicy(ctx, fmt.Sprintf("readPolicyIndex-%s", index.Name), &appautoscaling.PolicyArgs{
				PolicyType:        pulumi.String("TargetTrackingScaling"),
				ResourceId:        indexTarget.ResourceId,
				ScalableDimension: indexTarget.ScalableDimension,
				ServiceNamespace:  indexTarget.ServiceNamespace,
				TargetTrackingScalingPolicyConfiguration: &appautoscaling.PolicyTargetTrackingScalingPolicyConfigurationArgs{
					PredefinedMetricSpecification: &appautoscaling.PolicyTargetTrackingScalingPolicyConfigurationPredefinedMetricSpecificationArgs{
						PredefinedMetricType: pulumi.String("DynamoDBReadCapacityUtilization"),
					},
					TargetValue:      pulumi.Float64(autoScaleReadTargetIndex),
					ScaleInCooldown:  pulumi.Int(autoScaleScaleInCoolDown),
					ScaleOutCooldown: pulumi.Int(autoScaleScaleOutCoolDown),
				},
			}, pulumi.Provider(awsProvider),
				pulumi.Parent(indexTarget),
				pulumi.DependsOn([]pulumi.Resource{createdDynamodbTable, indexTarget}))
			if err != nil {
				return errors.Wrap(err, "failed to create read policy index auto scaling resources")
			}
		}

		writeTarget, err := appautoscaling.NewTarget(ctx, "writeTarget", &appautoscaling.TargetArgs{
			MaxCapacity:       pulumi.Int(autoScaleMaxWriteCapacity),
			MinCapacity:       pulumi.Int(autoScaleMinWriteCapacity),
			ResourceId:        pulumi.String("table/" + awsDynamodb.Metadata.Name),
			ScalableDimension: pulumi.String("dynamodb:table:WriteCapacityUnits"),
			ServiceNamespace:  pulumi.String("dynamodb"),
			Tags:              pulumi.ToStringMap(s.Labels),
		}, pulumi.Provider(awsProvider),
			pulumi.Parent(createdDynamodbTable),
			pulumi.DependsOn([]pulumi.Resource{createdDynamodbTable}))
		if err != nil {
			return errors.Wrap(err, "failed to create write target auto scaling resources")
		}

		_, err = appautoscaling.NewPolicy(ctx, "writePolicy", &appautoscaling.PolicyArgs{
			Name:              pulumi.Sprintf("DynamoDBWriteCapacityUtilization:%s", writeTarget.ID().ElementType()),
			PolicyType:        pulumi.String("TargetTrackingScaling"),
			ResourceId:        writeTarget.ResourceId,
			ScalableDimension: writeTarget.ScalableDimension,
			ServiceNamespace:  writeTarget.ServiceNamespace,
			TargetTrackingScalingPolicyConfiguration: &appautoscaling.PolicyTargetTrackingScalingPolicyConfigurationArgs{
				PredefinedMetricSpecification: &appautoscaling.PolicyTargetTrackingScalingPolicyConfigurationPredefinedMetricSpecificationArgs{
					PredefinedMetricType: pulumi.String("DynamoDBWriteCapacityUtilization"),
				},
				TargetValue:      pulumi.Float64(autoScaleWriteTarget),
				ScaleInCooldown:  pulumi.Int(autoScaleScaleInCoolDown),
				ScaleOutCooldown: pulumi.Int(autoScaleScaleOutCoolDown),
			},
		}, pulumi.Provider(awsProvider),
			pulumi.Parent(writeTarget),
			pulumi.DependsOn([]pulumi.Resource{createdDynamodbTable, writeTarget}))
		if err != nil {
			return errors.Wrap(err, "failed to create write policy")
		}

		for _, index := range awsDynamodb.Spec.Table.GlobalSecondaryIndexes {
			indexTarget, err := appautoscaling.NewTarget(ctx, fmt.Sprintf("writeTargetIndex-%s", index.Name), &appautoscaling.TargetArgs{
				MaxCapacity:       pulumi.Int(autoScaleMaxWriteCapacityIndex),
				MinCapacity:       pulumi.Int(autoScaleMinWriteCapacityIndex),
				ResourceId:        pulumi.String(fmt.Sprintf("table/%s/index/%s", awsDynamodb.Metadata.Name, index.Name)),
				ScalableDimension: pulumi.String("dynamodb:index:WriteCapacityUnits"),
				ServiceNamespace:  pulumi.String("dynamodb"),
				Tags:              pulumi.ToStringMap(s.Labels),
			}, pulumi.Provider(awsProvider),
				pulumi.Parent(writeTarget),
				pulumi.DependsOn([]pulumi.Resource{createdDynamodbTable, writeTarget}))
			if err != nil {
				return errors.Wrap(err, "failed to create write target index auto scaling resources")
			}

			// Create a Scaling Policy
			_, err = appautoscaling.NewPolicy(ctx, fmt.Sprintf("writePolicyIndex-%s", index.Name), &appautoscaling.PolicyArgs{
				PolicyType:        pulumi.String("TargetTrackingScaling"),
				ResourceId:        indexTarget.ResourceId,
				ScalableDimension: indexTarget.ScalableDimension,
				ServiceNamespace:  indexTarget.ServiceNamespace,
				TargetTrackingScalingPolicyConfiguration: &appautoscaling.PolicyTargetTrackingScalingPolicyConfigurationArgs{
					PredefinedMetricSpecification: &appautoscaling.PolicyTargetTrackingScalingPolicyConfigurationPredefinedMetricSpecificationArgs{
						PredefinedMetricType: pulumi.String("DynamoDBWriteCapacityUtilization"),
					},
					TargetValue:      pulumi.Float64(autoScaleWriteTargetIndex),
					ScaleInCooldown:  pulumi.Int(autoScaleScaleInCoolDown),
					ScaleOutCooldown: pulumi.Int(autoScaleScaleOutCoolDown),
				},
			}, pulumi.Provider(awsProvider),
				pulumi.Parent(indexTarget),
				pulumi.DependsOn([]pulumi.Resource{createdDynamodbTable, indexTarget}))
			if err != nil {
				return errors.Wrap(err, "failed to create write policy index auto scaling resources")
			}
		}
	}
	return nil
}
