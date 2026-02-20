package repository

import (
	"context"

	/*log "github.com/Elmar006/todo_grpc/internal/logger"*/
	model "github.com/Elmar006/todo_grpc/internal/model"
)

func CreateTaskDB(ctx context.Context, title, description string) error {

	return nil
}

func GetIDByDB(ctx context.Context, id string) (*model.Model, error) {

	return nil, nil
}

func ListDB(ctx context.Context) (*model.Model, error) {

	return nil, nil
}

func UpdateDB(ctx context.Context, task *model.Model) (*model.Model, error) {

	return nil, nil
}

func DeleteByID(ctx context.Context, id string) error {

	return nil
}
