package usecases

type (
	ISnapshotBookUseCase interface {
		Execute(instrument string) (*SnapshotBookOutput, error)
	}
)
