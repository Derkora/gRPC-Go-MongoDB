package repository

import (
	"context"
	"log"
	"time"

	"github.com/derkora/master/common/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type TaskRepository struct {
	db  *mongo.Database
	col *mongo.Collection
}

func NewTaskRepo(db *mongo.Database) *TaskRepository {
	return &TaskRepository {
		db: db,
		col: db.Collection(model.TaskCollection),
	}
}

// Save task
func (r *TaskRepository) Save(u *model.Task) (model.Task, error) {
	log.Printf("Save(%v) \n", u)
	ctx, cancel := timeoutContext()
	defer cancel()

	var task model.Task
	res, err := r.col.InsertOne(ctx, u)
	if err != nil {
		log.Println(err)
		return task, err
	}

	err = r.col.FindOne(ctx, bson.M{"_id": res.InsertedID}).Decode(&task)
	if err != nil {
		log.Println(err)
		return task, err
	}

	return task, nil
}

// find all tasks
func (r *TaskRepository) FindAll() ([]model.Task, error) {
	log.Println("FindAll()")
	ctx, cancel := timeoutContext()
	defer cancel()

	var tasks []model.Task
	cur, err := r.col.Find(ctx, bson.M{})
	if err != nil {
		log.Println(err)
		return tasks, err
	}

	defer cur.Close(ctx)
	for cur.Next(ctx) {
		var task model.Task
		err := cur.Decode(&task)
		if err != nil {
			log.Println(err)
		}
		tasks = append(tasks, task)
	}

	if err := cur.Err(); err != nil {
		log.Println(err)
		return nil, err
	}

	return tasks, nil
}

// update task 
func (r *TaskRepository) Update(u *model.Task) (model.Task, error) {
	log.Printf("Update(%v) \n", u)
	ctx, cancel := timeoutContext()
	defer cancel()

	filter := bson.M{"_id": u.ID}
	update := bson.M{
		"$set": bson.M{
			"title":  u.Title,
			"description": u.Description,
		},
	}

	var task model.Task
	err := r.col.FindOneAndUpdate(ctx, filter, update).Decode(&task)
	if err != nil {
		log.Printf("ERR 115 %v", err)
		return task, err
	}

	return task, nil
}

// delete task by id
func (r  *TaskRepository) Delete(id string) (bool, error) {
	log.Printf("Delete(%s) \n", id)
	ctx, cancel := timeoutContext()
	defer cancel()

	var task model.Task
	oid, _ := primitive.ObjectIDFromHex(id)
	err := r.col.FindOneAndDelete(ctx, bson.M{"_id": oid}).Decode(&task)
	if err != nil {
		log.Printf("Fail to delete task: %v \n", err)
		return false, err
	}
	log.Printf("Deleted_task(%v) \n", task)
	return true, nil
}

// buat timeout
func timeoutContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), time.Duration(60)*time.Second)
}