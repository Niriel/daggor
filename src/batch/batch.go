// batch project batch.go
package batch

type Batcher interface {
	Enter()
	Exit()
	Run()
}

// BaseBatch is a base object for Batches.
// It defines a Run method which is likely to be used by every
// other batch ever.  This justifies the creation of BaseBatch.
// Other batches can inherit the Run method by using BaseBatch
// as anonymous field.
type BaseBatch struct {
	Batches []Batcher
}

func (batch BaseBatch) Enter() {
}

func (batch BaseBatch) Exit() {
}

func (batch BaseBatch) Run() {
	for _, subBatch := range batch.Batches {
		subBatch.Enter()
		subBatch.Run()
		subBatch.Exit()
	}
}
