package tui

type CreateConsumerMsg struct {
	Topic string
}

type FocusMsg struct {
	Focus bool
}
