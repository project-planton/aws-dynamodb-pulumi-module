package pkg

import (
	"github.com/pkg/errors"
	"github.com/plantoncloud/aws-dynamodb-pulumi-module/pkg/outputs"
	"github.com/pulumi/pulumi-aws/sdk/v6/go/aws"
	"github.com/pulumi/pulumi-aws/sdk/v6/go/aws/dynamodb"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func table(ctx *pulumi.Context, locals *Locals, awsProvider *aws.Provider) (*dynamodb.Table, error) {
	awsDynamodb := locals.AwsDynamodb

	// stream
	streamEnabled := awsDynamodb.Spec.Table.EnableStreams
	streamViewType := ""
	if len(awsDynamodb.Spec.Table.ReplicaRegionNames) > 0 {
		streamEnabled = true
	}

	if len(awsDynamodb.Spec.Table.ReplicaRegionNames) > 0 || awsDynamodb.Spec.Table.EnableStreams {
		streamViewType = awsDynamodb.Spec.Table.StreamViewType
	}

	// capacity
	readCapacity := 0
	writeCapacity := 0
	if awsDynamodb.Spec.Table.BillingMode == "PROVISIONED" {
		readCapacity = 5
		writeCapacity = 5
		if awsDynamodb.Spec.Table.AutoScale != nil &&
			awsDynamodb.Spec.Table.AutoScale.ReadCapacity != nil {
			readCapacity = int(awsDynamodb.Spec.Table.AutoScale.ReadCapacity.MinCapacity)
		}
		if awsDynamodb.Spec.Table.AutoScale != nil &&
			awsDynamodb.Spec.Table.AutoScale.WriteCapacity != nil {
			writeCapacity = int(awsDynamodb.Spec.Table.AutoScale.WriteCapacity.MinCapacity)
		}
	}

	// range key
	rangeKey := ""
	if awsDynamodb.Spec.Table.RangeKey != nil {
		rangeKey = awsDynamodb.Spec.Table.RangeKey.Name
	}

	// replicas
	var replicaArray = dynamodb.TableReplicaTypeArray{}
	for _, regionName := range awsDynamodb.Spec.Table.ReplicaRegionNames {
		replicaArray = append(replicaArray, &dynamodb.TableReplicaTypeArgs{
			KmsKeyArn:           nil,
			PointInTimeRecovery: pulumi.Bool(false),
			PropagateTags:       pulumi.Bool(false),
			RegionName:          pulumi.String(regionName),
		})
	}

	// attributes
	var attributeArray = dynamodb.TableAttributeArray{}
	var attributeMap = make(map[string]bool)

	attributeMap[awsDynamodb.Spec.Table.HashKey.Name] = true
	attributeArray = append(attributeArray, &dynamodb.TableAttributeArgs{
		Name: pulumi.String(awsDynamodb.Spec.Table.HashKey.Name),
		Type: pulumi.String(awsDynamodb.Spec.Table.HashKey.Type),
	})

	if awsDynamodb.Spec.Table.RangeKey.Name != "" {
		attributeMap[awsDynamodb.Spec.Table.RangeKey.Name] = true
		attributeArray = append(attributeArray, &dynamodb.TableAttributeArgs{
			Name: pulumi.String(awsDynamodb.Spec.Table.RangeKey.Name),
			Type: pulumi.String(awsDynamodb.Spec.Table.RangeKey.Type),
		})
	}

	for _, attribute := range awsDynamodb.Spec.Table.Attributes {
		_, exists := attributeMap[attribute.Name]
		if !exists {
			attributeMap[attribute.Name] = true
			attributeArray = append(attributeArray, &dynamodb.TableAttributeArgs{
				Name: pulumi.String(attribute.Name),
				Type: pulumi.String(attribute.Type),
			})
		}
	}

	// global secondary index
	var globalSecondaryIndexArray = dynamodb.TableGlobalSecondaryIndexArray{}
	for _, globalSecondaryIndex := range awsDynamodb.Spec.Table.GlobalSecondaryIndexes {
		globalIndexReadCapacity := globalSecondaryIndex.ReadCapacity
		globalIndexWriteCapacity := globalSecondaryIndex.WriteCapacity
		if awsDynamodb.Spec.Table.BillingMode == "PROVISIONED" && globalIndexReadCapacity == 0 {
			globalIndexReadCapacity = int32(readCapacity)
		}
		if awsDynamodb.Spec.Table.BillingMode == "PROVISIONED" && globalIndexWriteCapacity == 0 {
			globalIndexWriteCapacity = int32(writeCapacity)
		}
		globalSecondaryIndexArray = append(globalSecondaryIndexArray, &dynamodb.TableGlobalSecondaryIndexArgs{
			Name:             pulumi.String(globalSecondaryIndex.Name),
			HashKey:          pulumi.String(globalSecondaryIndex.HashKey),
			RangeKey:         pulumi.String(globalSecondaryIndex.RangeKey),
			ReadCapacity:     pulumi.Int(globalIndexReadCapacity),
			WriteCapacity:    pulumi.Int(globalIndexWriteCapacity),
			ProjectionType:   pulumi.String(globalSecondaryIndex.ProjectionType),
			NonKeyAttributes: pulumi.ToStringArray(globalSecondaryIndex.NonKeyAttributes),
		})
	}

	// local secondary index
	var localSecondaryIndexArray = dynamodb.TableLocalSecondaryIndexArray{}
	for _, localSecondaryIndex := range awsDynamodb.Spec.Table.LocalSecondaryIndexes {
		localSecondaryIndexArray = append(localSecondaryIndexArray, &dynamodb.TableLocalSecondaryIndexArgs{
			Name:             pulumi.String(localSecondaryIndex.Name),
			RangeKey:         pulumi.String(localSecondaryIndex.RangeKey),
			ProjectionType:   pulumi.String(localSecondaryIndex.ProjectionType),
			NonKeyAttributes: pulumi.ToStringArray(localSecondaryIndex.NonKeyAttributes),
		})
	}

	// server side encryption
	var serverSideEncryption *dynamodb.TableServerSideEncryptionArgs
	if awsDynamodb.Spec.Table.ServerSideEncryption != nil {
		serverSideEncryption = &dynamodb.TableServerSideEncryptionArgs{
			Enabled:   pulumi.Bool(awsDynamodb.Spec.Table.ServerSideEncryption.IsEnabled),
			KmsKeyArn: pulumi.StringPtr(awsDynamodb.Spec.Table.ServerSideEncryption.KmsKeyArn),
		}
	}

	// point in time recovery
	var pointInTimeRecovery *dynamodb.TablePointInTimeRecoveryArgs
	if awsDynamodb.Spec.Table.PointInTimeRecovery != nil {
		pointInTimeRecovery = &dynamodb.TablePointInTimeRecoveryArgs{
			Enabled: pulumi.Bool(awsDynamodb.Spec.Table.PointInTimeRecovery.IsEnabled),
		}
	}

	// ttl
	var ttl *dynamodb.TableTtlArgs
	if awsDynamodb.Spec.Table.Ttl != nil {
		ttl = &dynamodb.TableTtlArgs{
			Enabled:       pulumi.Bool(awsDynamodb.Spec.Table.Ttl.IsEnabled),
			AttributeName: pulumi.String(awsDynamodb.Spec.Table.Ttl.AttributeName),
		}
	}

	// import table
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
					HeaderLists: pulumi.ToStringArray(awsDynamodb.Spec.Table.ImportTable.InputFormatOptions.Csv.Headers),
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
		Name:                      pulumi.String(awsDynamodb.Spec.Table.TableName),
		BillingMode:               pulumi.String(awsDynamodb.Spec.Table.BillingMode),
		ReadCapacity:              pulumi.Int(readCapacity),
		WriteCapacity:             pulumi.Int(writeCapacity),
		HashKey:                   pulumi.String(awsDynamodb.Spec.Table.HashKey.Name),
		RangeKey:                  pulumi.String(rangeKey),
		StreamEnabled:             pulumi.Bool(streamEnabled),
		StreamViewType:            pulumi.String(streamViewType),
		TableClass:                pulumi.String("STANDARD"),
		DeletionProtectionEnabled: pulumi.Bool(false),
		ServerSideEncryption:      serverSideEncryption,
		PointInTimeRecovery:       pointInTimeRecovery,
		Ttl:                       ttl,
		Tags:                      pulumi.ToStringMap(locals.Labels),
		Attributes:                attributeArray,
		GlobalSecondaryIndexes:    globalSecondaryIndexArray,
		LocalSecondaryIndexes:     localSecondaryIndexArray,
		Replicas:                  replicaArray,
		ImportTable:               importTable,
	}, pulumi.Provider(awsProvider))
	if err != nil {
		return nil, errors.Wrap(err, "failed to create dynamo table resources")
	}

	ctx.Export(outputs.TableName, createdDynamodbTable.Name)
	ctx.Export(outputs.TableArn, createdDynamodbTable.Arn)
	ctx.Export(outputs.TableStreamArn, createdDynamodbTable.StreamArn)

	return createdDynamodbTable, nil
}
