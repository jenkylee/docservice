package service

// 提供DocService 操作
type Service interface {
	Import(string) (string, error)
	Export(string) (int)
}