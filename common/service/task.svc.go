package service

import (
	"context"
	"log"

	"github.com/derkora/master/common/genproto/taskmaster"
	"github.com/derkora/master/common/model"
	"github.com/derkora/master/common/repository"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/golang/protobuf/ptypes/wrappers"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

type TaskService struct {
	repo *repository.TaskRepository
}

func NewTaskService(repo *repository.TaskRepository) *TaskService {
	return &TaskService{
		repo: repo,
	}
}

// buat task
func (s *TaskService) CreateTask(ctx context.Context, tm *taskmaster.Task) (*taskmaster.Task, error) {
	log.Printf("CreateTask(%v) \n", tm)

	newTask := &model.Task{
		Title:       tm.Title,
		Description: tm.Description,
	}

	task, err := s.repo.Save(newTask)
	if err != nil {
		log.Printf("%v", err)
		return nil, err
	}

	return s.toTask(&task), nil
}

// dapatin semua task
func (s *TaskService) ListTasks(ctx context.Context, e *empty.Empty) (*taskmaster.TaskList, error) {
	log.Printf("ListTasks() \n")

	var totas []*taskmaster.Task
	Tasks, err := s.repo.FindAll()
	if err != nil {
		log.Printf("%v", err)
		return nil, err
	}

	for _, u := range Tasks {
		totas = append(totas, s.toTask(&u))
	}

	TaskList := &taskmaster.TaskList{
		List: totas,
	}

	return TaskList, nil
}

// update task
func (s *TaskService) UpdateTask(ctx context.Context, tm *taskmaster.Task) (*taskmaster.Task, error) {
	log.Printf("UpdateTask(%v) \n", tm)

	if tm.Id == "" {
		return nil, status.Error(codes.FailedPrecondition, "UpdateTask must provide taskID")
	}

	taskID, err := primitive.ObjectIDFromHex(tm.Id)
	if err != nil {
		log.Printf("Invalid TaskID(%s) \n", tm.Id)
		return nil, status.Errorf(codes.Internal, "%v", err)
	}

	updateTask := &model.Task{
		ID:          taskID,
		Title:       tm.Title,
		Description: tm.Description,
	}

	task, err := s.repo.Update(updateTask)
	if err != nil {
		log.Printf("Fail UpdateTask %v \n", err)
		return nil, status.Errorf(codes.Internal, "%v", err)
	}

	return s.toTask(&task), nil
}

// delete task
func (s *TaskService) DeleteTask(ctx context.Context, id *wrappers.StringValue) (*wrappers.BoolValue, error) {
	log.Printf("DeleteTask(%s) \n", id.GetValue())

	deleted, err := s.repo.Delete(id.GetValue())
	if err != nil {
		log.Printf("%v", err)
		return nil, err
	}

	return &wrapperspb.BoolValue{Value: deleted}, nil
}

// map task ke toTask
func (s *TaskService) toTask(u *model.Task) *taskmaster.Task {
	tota := &taskmaster.Task{
		Id:          u.ID.Hex(),
		Title:       u.Title,
		Description: u.Description,
	}
	return tota
}
