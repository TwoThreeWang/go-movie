package spider

// Spider 采集接口
type Spider interface {
	ApiSearch(kw string) (data []Movie, msg string)
}
