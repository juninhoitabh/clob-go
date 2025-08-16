package usecases

type (
	ISnapshotBookUseCase interface {
		Execute(input SnapshotBookInput) (*SnapshotBookOutput, error)
	}
)
