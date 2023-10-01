package pkg

type Roomiface interface {
	Join(password *string) error
	IsStillValid() bool
	Name() string
}
