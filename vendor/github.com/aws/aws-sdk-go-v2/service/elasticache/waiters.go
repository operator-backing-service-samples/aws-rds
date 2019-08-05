// Code generated by private/model/cli/gen-api/main.go. DO NOT EDIT.

package elasticache

import (
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
)

// WaitUntilCacheClusterAvailable uses the Amazon ElastiCache API operation
// DescribeCacheClusters to wait for a condition to be met before returning.
// If the condition is not met within the max attempt window, an error will
// be returned.
func (c *ElastiCache) WaitUntilCacheClusterAvailable(input *DescribeCacheClustersInput) error {
	return c.WaitUntilCacheClusterAvailableWithContext(aws.BackgroundContext(), input)
}

// WaitUntilCacheClusterAvailableWithContext is an extended version of WaitUntilCacheClusterAvailable.
// With the support for passing in a context and options to configure the
// Waiter and the underlying request options.
//
// The context must be non-nil and will be used for request cancellation. If
// the context is nil a panic will occur. In the future the SDK may create
// sub-contexts for http.Requests. See https://golang.org/pkg/context/
// for more information on using Contexts.
func (c *ElastiCache) WaitUntilCacheClusterAvailableWithContext(ctx aws.Context, input *DescribeCacheClustersInput, opts ...aws.WaiterOption) error {
	w := aws.Waiter{
		Name:        "WaitUntilCacheClusterAvailable",
		MaxAttempts: 40,
		Delay:       aws.ConstantWaiterDelay(15 * time.Second),
		Acceptors: []aws.WaiterAcceptor{
			{
				State:   aws.SuccessWaiterState,
				Matcher: aws.PathAllWaiterMatch, Argument: "CacheClusters[].CacheClusterStatus",
				Expected: "available",
			},
			{
				State:   aws.FailureWaiterState,
				Matcher: aws.PathAnyWaiterMatch, Argument: "CacheClusters[].CacheClusterStatus",
				Expected: "deleted",
			},
			{
				State:   aws.FailureWaiterState,
				Matcher: aws.PathAnyWaiterMatch, Argument: "CacheClusters[].CacheClusterStatus",
				Expected: "deleting",
			},
			{
				State:   aws.FailureWaiterState,
				Matcher: aws.PathAnyWaiterMatch, Argument: "CacheClusters[].CacheClusterStatus",
				Expected: "incompatible-network",
			},
			{
				State:   aws.FailureWaiterState,
				Matcher: aws.PathAnyWaiterMatch, Argument: "CacheClusters[].CacheClusterStatus",
				Expected: "restore-failed",
			},
		},
		Logger: c.Config.Logger,
		NewRequest: func(opts []aws.Option) (*aws.Request, error) {
			var inCpy *DescribeCacheClustersInput
			if input != nil {
				tmp := *input
				inCpy = &tmp
			}
			req := c.DescribeCacheClustersRequest(inCpy)
			req.SetContext(ctx)
			req.ApplyOptions(opts...)
			return req.Request, nil
		},
	}
	w.ApplyOptions(opts...)

	return w.WaitWithContext(ctx)
}

// WaitUntilCacheClusterDeleted uses the Amazon ElastiCache API operation
// DescribeCacheClusters to wait for a condition to be met before returning.
// If the condition is not met within the max attempt window, an error will
// be returned.
func (c *ElastiCache) WaitUntilCacheClusterDeleted(input *DescribeCacheClustersInput) error {
	return c.WaitUntilCacheClusterDeletedWithContext(aws.BackgroundContext(), input)
}

// WaitUntilCacheClusterDeletedWithContext is an extended version of WaitUntilCacheClusterDeleted.
// With the support for passing in a context and options to configure the
// Waiter and the underlying request options.
//
// The context must be non-nil and will be used for request cancellation. If
// the context is nil a panic will occur. In the future the SDK may create
// sub-contexts for http.Requests. See https://golang.org/pkg/context/
// for more information on using Contexts.
func (c *ElastiCache) WaitUntilCacheClusterDeletedWithContext(ctx aws.Context, input *DescribeCacheClustersInput, opts ...aws.WaiterOption) error {
	w := aws.Waiter{
		Name:        "WaitUntilCacheClusterDeleted",
		MaxAttempts: 40,
		Delay:       aws.ConstantWaiterDelay(15 * time.Second),
		Acceptors: []aws.WaiterAcceptor{
			{
				State:   aws.SuccessWaiterState,
				Matcher: aws.PathAllWaiterMatch, Argument: "CacheClusters[].CacheClusterStatus",
				Expected: "deleted",
			},
			{
				State:    aws.SuccessWaiterState,
				Matcher:  aws.ErrorWaiterMatch,
				Expected: "CacheClusterNotFound",
			},
			{
				State:   aws.FailureWaiterState,
				Matcher: aws.PathAnyWaiterMatch, Argument: "CacheClusters[].CacheClusterStatus",
				Expected: "available",
			},
			{
				State:   aws.FailureWaiterState,
				Matcher: aws.PathAnyWaiterMatch, Argument: "CacheClusters[].CacheClusterStatus",
				Expected: "creating",
			},
			{
				State:   aws.FailureWaiterState,
				Matcher: aws.PathAnyWaiterMatch, Argument: "CacheClusters[].CacheClusterStatus",
				Expected: "incompatible-network",
			},
			{
				State:   aws.FailureWaiterState,
				Matcher: aws.PathAnyWaiterMatch, Argument: "CacheClusters[].CacheClusterStatus",
				Expected: "modifying",
			},
			{
				State:   aws.FailureWaiterState,
				Matcher: aws.PathAnyWaiterMatch, Argument: "CacheClusters[].CacheClusterStatus",
				Expected: "restore-failed",
			},
			{
				State:   aws.FailureWaiterState,
				Matcher: aws.PathAnyWaiterMatch, Argument: "CacheClusters[].CacheClusterStatus",
				Expected: "snapshotting",
			},
		},
		Logger: c.Config.Logger,
		NewRequest: func(opts []aws.Option) (*aws.Request, error) {
			var inCpy *DescribeCacheClustersInput
			if input != nil {
				tmp := *input
				inCpy = &tmp
			}
			req := c.DescribeCacheClustersRequest(inCpy)
			req.SetContext(ctx)
			req.ApplyOptions(opts...)
			return req.Request, nil
		},
	}
	w.ApplyOptions(opts...)

	return w.WaitWithContext(ctx)
}

// WaitUntilReplicationGroupAvailable uses the Amazon ElastiCache API operation
// DescribeReplicationGroups to wait for a condition to be met before returning.
// If the condition is not met within the max attempt window, an error will
// be returned.
func (c *ElastiCache) WaitUntilReplicationGroupAvailable(input *DescribeReplicationGroupsInput) error {
	return c.WaitUntilReplicationGroupAvailableWithContext(aws.BackgroundContext(), input)
}

// WaitUntilReplicationGroupAvailableWithContext is an extended version of WaitUntilReplicationGroupAvailable.
// With the support for passing in a context and options to configure the
// Waiter and the underlying request options.
//
// The context must be non-nil and will be used for request cancellation. If
// the context is nil a panic will occur. In the future the SDK may create
// sub-contexts for http.Requests. See https://golang.org/pkg/context/
// for more information on using Contexts.
func (c *ElastiCache) WaitUntilReplicationGroupAvailableWithContext(ctx aws.Context, input *DescribeReplicationGroupsInput, opts ...aws.WaiterOption) error {
	w := aws.Waiter{
		Name:        "WaitUntilReplicationGroupAvailable",
		MaxAttempts: 40,
		Delay:       aws.ConstantWaiterDelay(15 * time.Second),
		Acceptors: []aws.WaiterAcceptor{
			{
				State:   aws.SuccessWaiterState,
				Matcher: aws.PathAllWaiterMatch, Argument: "ReplicationGroups[].Status",
				Expected: "available",
			},
			{
				State:   aws.FailureWaiterState,
				Matcher: aws.PathAnyWaiterMatch, Argument: "ReplicationGroups[].Status",
				Expected: "deleted",
			},
		},
		Logger: c.Config.Logger,
		NewRequest: func(opts []aws.Option) (*aws.Request, error) {
			var inCpy *DescribeReplicationGroupsInput
			if input != nil {
				tmp := *input
				inCpy = &tmp
			}
			req := c.DescribeReplicationGroupsRequest(inCpy)
			req.SetContext(ctx)
			req.ApplyOptions(opts...)
			return req.Request, nil
		},
	}
	w.ApplyOptions(opts...)

	return w.WaitWithContext(ctx)
}

// WaitUntilReplicationGroupDeleted uses the Amazon ElastiCache API operation
// DescribeReplicationGroups to wait for a condition to be met before returning.
// If the condition is not met within the max attempt window, an error will
// be returned.
func (c *ElastiCache) WaitUntilReplicationGroupDeleted(input *DescribeReplicationGroupsInput) error {
	return c.WaitUntilReplicationGroupDeletedWithContext(aws.BackgroundContext(), input)
}

// WaitUntilReplicationGroupDeletedWithContext is an extended version of WaitUntilReplicationGroupDeleted.
// With the support for passing in a context and options to configure the
// Waiter and the underlying request options.
//
// The context must be non-nil and will be used for request cancellation. If
// the context is nil a panic will occur. In the future the SDK may create
// sub-contexts for http.Requests. See https://golang.org/pkg/context/
// for more information on using Contexts.
func (c *ElastiCache) WaitUntilReplicationGroupDeletedWithContext(ctx aws.Context, input *DescribeReplicationGroupsInput, opts ...aws.WaiterOption) error {
	w := aws.Waiter{
		Name:        "WaitUntilReplicationGroupDeleted",
		MaxAttempts: 40,
		Delay:       aws.ConstantWaiterDelay(15 * time.Second),
		Acceptors: []aws.WaiterAcceptor{
			{
				State:   aws.SuccessWaiterState,
				Matcher: aws.PathAllWaiterMatch, Argument: "ReplicationGroups[].Status",
				Expected: "deleted",
			},
			{
				State:   aws.FailureWaiterState,
				Matcher: aws.PathAnyWaiterMatch, Argument: "ReplicationGroups[].Status",
				Expected: "available",
			},
			{
				State:    aws.SuccessWaiterState,
				Matcher:  aws.ErrorWaiterMatch,
				Expected: "ReplicationGroupNotFoundFault",
			},
		},
		Logger: c.Config.Logger,
		NewRequest: func(opts []aws.Option) (*aws.Request, error) {
			var inCpy *DescribeReplicationGroupsInput
			if input != nil {
				tmp := *input
				inCpy = &tmp
			}
			req := c.DescribeReplicationGroupsRequest(inCpy)
			req.SetContext(ctx)
			req.ApplyOptions(opts...)
			return req.Request, nil
		},
	}
	w.ApplyOptions(opts...)

	return w.WaitWithContext(ctx)
}
