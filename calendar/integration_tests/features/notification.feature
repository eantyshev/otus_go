# file: features/notification.feature

Feature: Notification sending
	Appointment owner are notified
	when starting time is approaching

	Scenario: Notification is sent to RabbitMQ
		When appointment has start time at now
		Then notification is received within 60 seconds
