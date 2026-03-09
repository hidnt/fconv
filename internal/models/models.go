package models

type Config struct {
	DstExt        string
	DstFolder     string
	NeedRecursion bool
	LevelOfRec    int
	Delete        bool
	Force         bool
}
