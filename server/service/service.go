package service

// 提供DocService 操作
type Service interface {
	Login(name, pwd string) (string, error)
	Import(string) (string, error)
	Export(string) (int)
}