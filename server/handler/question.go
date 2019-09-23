package handler

import (
	"yokitalk.com/docservice/server/repository"
)



type Question struct {
	Repo *repository.QuestionRepository
}

func (q *Question) Import(fileName string, dir string) error {



	return nil
}

