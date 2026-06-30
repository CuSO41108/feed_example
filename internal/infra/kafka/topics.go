package kafka

const (
	TopicPostPublished    = "post.published"
	TopicPostDeleted      = "post.deleted"
	TopicFeedFanoutChunk  = "feed.fanout.chunk"
	TopicFeedFanoutRetry  = "feed.fanout.chunk.retry"
	TopicFeedFanoutDLQ    = "feed.fanout.chunk.dlq"
	DefaultKafkaPartition = 0
)
