package ws

type Trainer struct {
	Content   string `bson:"content"`
	StartTime int64  `bson:"starttime"`
	EndTime   int64  `bson:"endtime"`
	Read      int64  `bson:"read"`
}

type Result struct {
	StartTime int64
	Msg       string
	Content   interface{}
	From      string
}
