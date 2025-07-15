package dagcuter

import (
	"context"
	"fmt"
	"math"
	"time"
)

type RetryPolicy struct {
	Interval    time.Duration `json:"interval" yaml:"interval"`
	MaxInterval time.Duration `json:"maxInterval" yaml:"maxInterval"`
	MaxAttempts int           `json:"maxAttempts" yaml:"maxAttempts"`
	Multiplier  float64       `json:"multiplier" yaml:"multiplier"`
}

type RetryExecutor struct {
	policy *RetryPolicy
}

func (d *Dagcuter) newRetryExecutor(policy *RetryPolicy) *RetryExecutor {
	if policy == nil {
		// 默认策略：不重试, 只执行一次
		policy = &RetryPolicy{
			MaxAttempts: -1,
		}
	}
	if policy.Interval <= 0 {
		policy.Interval = 1 * time.Second // 默认间隔1秒
	}
	if policy.MaxInterval <= 0 {
		policy.MaxInterval = 30 * time.Second // 默认最大间隔30秒
	}
	if policy.Multiplier <= 0 {
		policy.Multiplier = 2.0 // 默认乘数为2
	}
	return &RetryExecutor{policy: policy}
}

// ExecuteWithRetry 带重试的执行函数
func (r *RetryExecutor) ExecuteWithRetry(ctx context.Context, taskName string, fn func(n int) error) error {
	if r.policy.MaxAttempts <= 0 {
		return fn(0) // 无重试策略，直接执行
	}

	if r.policy.MaxInterval > 150*time.Second {
		r.policy.MaxInterval = 150 * time.Second // 默认最大间隔150秒
	}

	var lastErr error
	for attempt := 1; attempt <= r.policy.MaxAttempts; attempt++ {
		select {
		case <-ctx.Done():
			return fmt.Errorf("context cancelled during retry attempt %d: %w", attempt, ctx.Err())
		default:
		}

		// 执行任务
		if err := fn(attempt); err == nil {
			return nil // 成功执行
		} else {
			lastErr = err
		}

		// 如果是最后一次尝试，不需要等待
		if attempt == r.policy.MaxAttempts {
			break
		}

		// 计算等待时间（指数退避）
		waitTime := r.calculateBackoff(attempt, r.policy.MaxInterval)

		// 等待重试
		select {
		case <-ctx.Done():
			return fmt.Errorf("context cancelled during retry wait: %w", ctx.Err())
		case <-time.After(waitTime):
			// 继续重试
		}
	}

	return fmt.Errorf("task %s failed after %d attempts, last error: %w",
		taskName, r.policy.MaxAttempts, lastErr)
}

// calculateBackoff 计算指数退避时间，使用 math.Pow 处理浮点数
func (r *RetryExecutor) calculateBackoff(attempt int, maxInterval time.Duration) time.Duration {
	// 使用 math.Pow 进行精确的浮点数幂运算
	// 公式: baseInterval * (multiplier ^ (attempt - 1))
	backoff := float64(r.policy.Interval) * math.Pow(r.policy.Multiplier, float64(attempt-1))

	result := time.Duration(backoff)

	// 确保不超过最大间隔
	if result > maxInterval {
		return maxInterval
	}

	return result
}
