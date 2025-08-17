package usecases

type (
	SnapshotBookInput struct {
		Instrument string
	}
	SnapshotBookOutput struct {
		Instrument string
		Bids       []Level
		Asks       []Level
	}
	ISnapshotBookUseCase interface {
		Execute(input SnapshotBookInput) (*SnapshotBookOutput, error)
	}
)
