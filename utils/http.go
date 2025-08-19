package utils

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"time"
)

// OptimizedHTTPClient provides a high-performance HTTP client for schema downloads
type OptimizedHTTPClient struct {
	client       *http.Client
	maxRetries   int
	retryBackoff time.Duration
}

// HTTPClientOptions contains configuration for the optimized HTTP client
type HTTPClientOptions struct {
	Timeout               time.Duration
	MaxIdleConns          int
	MaxIdleConnsPerHost   int
	IdleConnTimeout       time.Duration
	TLSHandshakeTimeout   time.Duration
	ResponseHeaderTimeout time.Duration
	KeepAlive             time.Duration
	MaxRetries            int
	RetryBackoff          time.Duration
}

// DefaultHTTPClientOptions returns optimized default options for schema downloads
func DefaultHTTPClientOptions() *HTTPClientOptions {
	return &HTTPClientOptions{
		Timeout:               30 * time.Second,
		MaxIdleConns:          100,
		MaxIdleConnsPerHost:   10,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ResponseHeaderTimeout: 10 * time.Second,
		KeepAlive:             30 * time.Second,
		MaxRetries:            3,
		RetryBackoff:          1 * time.Second,
	}
}

// NewOptimizedHTTPClient creates a new optimized HTTP client
func NewOptimizedHTTPClient(opts *HTTPClientOptions) *OptimizedHTTPClient {
	if opts == nil {
		opts = DefaultHTTPClientOptions()
	}

	// Create custom transport with optimized settings
	transport := &http.Transport{
		// Connection pooling settings
		MaxIdleConns:        opts.MaxIdleConns,
		MaxIdleConnsPerHost: opts.MaxIdleConnsPerHost,
		IdleConnTimeout:     opts.IdleConnTimeout,

		// Timeout settings
		TLSHandshakeTimeout:   opts.TLSHandshakeTimeout,
		ResponseHeaderTimeout: opts.ResponseHeaderTimeout,

		// Connection settings with keep-alive
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: opts.KeepAlive,
		}).DialContext,

		// Performance optimizations
		DisableCompression: false, // Keep compression enabled for schemas
		ForceAttemptHTTP2:  true,  // Use HTTP/2 when available
	}

	client := &http.Client{
		Transport: transport,
		Timeout:   opts.Timeout,
		// Don't follow redirects automatically - we want to control this
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			if len(via) >= 3 {
				return fmt.Errorf("too many redirects")
			}
			return nil
		},
	}

	return &OptimizedHTTPClient{
		client:       client,
		maxRetries:   opts.MaxRetries,
		retryBackoff: opts.RetryBackoff,
	}
}

// Get performs an optimized GET request with retry logic
func (c *OptimizedHTTPClient) Get(ctx context.Context, url string) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set optimal headers for schema downloads
	req.Header.Set("User-Agent", "NetEX-Validator-Go/1.0 (optimized)")
	req.Header.Set("Accept", "application/xml, text/xml, */*")
	req.Header.Set("Accept-Encoding", "gzip, deflate")
	req.Header.Set("Connection", "keep-alive")

	return c.doWithRetry(req)
}

// doWithRetry performs an HTTP request with exponential backoff retry logic
func (c *OptimizedHTTPClient) doWithRetry(req *http.Request) (*http.Response, error) {
	var lastErr error

	for attempt := 0; attempt <= c.maxRetries; attempt++ {
		// Clone request for retry attempts (body might be consumed)
		reqClone := req.Clone(req.Context())

		resp, err := c.client.Do(reqClone)
		if err != nil {
			lastErr = err

			// Check if this is a retryable error
			if !c.isRetryableError(err) {
				return nil, err
			}

			// Don't retry on the last attempt
			if attempt == c.maxRetries {
				break
			}

			// Wait before retry with exponential backoff
			// #nosec G115: attempt is small and non-negative; shift is bounded by maxRetries
			backoff := c.retryBackoff * time.Duration(1<<uint64(attempt))
			if backoff > 30*time.Second {
				backoff = 30 * time.Second // Cap at 30 seconds
			}

			select {
			case <-req.Context().Done():
				return nil, req.Context().Err()
			case <-time.After(backoff):
				continue
			}
		}

		// Check response status
		if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			return resp, nil
		}

		// Close response body for non-success responses
		_ = resp.Body.Close()

		// Check if status code is retryable
		if c.isRetryableStatusCode(resp.StatusCode) {
			lastErr = fmt.Errorf("HTTP %d: %s", resp.StatusCode, resp.Status)

			// Don't retry on the last attempt
			if attempt == c.maxRetries {
				break
			}

			// Wait before retry
			// #nosec G115: attempt is small and non-negative; shift is bounded by maxRetries
			backoff := c.retryBackoff * time.Duration(1<<uint64(attempt))
			if backoff > 30*time.Second {
				backoff = 30 * time.Second
			}

			select {
			case <-req.Context().Done():
				return nil, req.Context().Err()
			case <-time.After(backoff):
				continue
			}
		} else {
			// Non-retryable status code
			return nil, fmt.Errorf("HTTP %d: %s", resp.StatusCode, resp.Status)
		}
	}

	return nil, fmt.Errorf("request failed after %d attempts: %w", c.maxRetries+1, lastErr)
}

// isRetryableError determines if an error is worth retrying
func (c *OptimizedHTTPClient) isRetryableError(err error) bool {
	// Retry on network timeouts
	if netErr, ok := err.(net.Error); ok {
		return netErr.Timeout()
	}

	// Retry on context deadline exceeded (but not context canceled)
	if err == context.DeadlineExceeded {
		return true
	}

	return false
}

// isRetryableStatusCode determines if an HTTP status code is worth retrying
func (c *OptimizedHTTPClient) isRetryableStatusCode(statusCode int) bool {
	switch statusCode {
	case http.StatusRequestTimeout, // 408
		http.StatusTooManyRequests,     // 429
		http.StatusInternalServerError, // 500
		http.StatusBadGateway,          // 502
		http.StatusServiceUnavailable,  // 503
		http.StatusGatewayTimeout:      // 504
		return true
	default:
		return false
	}
}

// Close closes the HTTP client and cleans up resources
func (c *OptimizedHTTPClient) Close() {
	if transport, ok := c.client.Transport.(*http.Transport); ok {
		transport.CloseIdleConnections()
	}
}

// GetStats returns statistics about the HTTP client connections
func (c *OptimizedHTTPClient) GetStats() map[string]interface{} {
	stats := map[string]interface{}{
		"maxRetries":   c.maxRetries,
		"retryBackoff": c.retryBackoff.String(),
	}

	if transport, ok := c.client.Transport.(*http.Transport); ok {
		stats["maxIdleConns"] = transport.MaxIdleConns
		stats["maxIdleConnsPerHost"] = transport.MaxIdleConnsPerHost
		stats["idleConnTimeout"] = transport.IdleConnTimeout.String()
		stats["tlsHandshakeTimeout"] = transport.TLSHandshakeTimeout.String()
		stats["responseHeaderTimeout"] = transport.ResponseHeaderTimeout.String()
	}

	return stats
}
