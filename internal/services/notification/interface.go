package notification

type NotificationService interface { //nolint
	Publish(topic string, message any) error
}
