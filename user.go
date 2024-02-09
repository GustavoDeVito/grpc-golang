package main

import (
	"context"
	"database/sql"
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/GustavoDeVito/grpc-golang/proto/gen"
	_ "github.com/mattn/go-sqlite3"
)

type UserServer struct {
	pb.UnimplementedUserServiceServer
	db *sql.DB
}

type User struct {
	ID     int32
	Name   string
	Status bool
}

func (s *UserServer) InitializeDB() error {
	db, err := sql.Open("sqlite3", "./users.db")
	if err != nil {
		return err
	}

	s.db = db
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT,
		status BOOLEAN
	)`)

	return err
}

func (s *UserServer) FindAll(ctx context.Context, req *pb.FindAllRequest) (*pb.FindAllResponse, error) {
	rows, err := s.db.Query("SELECT id, name, status FROM users")
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var users []*pb.User
	for rows.Next() {
		var user User

		err := rows.Scan(&user.ID, &user.Name, &user.Status)
		if err != nil {
			return nil, err
		}

		users = append(users, &pb.User{Id: user.ID, Name: user.Name, Status: user.Status})
	}

	return &pb.FindAllResponse{Users: users}, nil
}

func (s *UserServer) FindOne(ctx context.Context, req *pb.FindOneRequest) (*pb.FindOneResponse, error) {
	var user User

	err := s.db.QueryRow("SELECT id, name, status FROM users WHERE id = ?", req.Id).Scan(&user.ID, &user.Name, &user.Status)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, status.Errorf(codes.NotFound, "User not found")
		}
		return nil, err
	}

	return &pb.FindOneResponse{User: &pb.User{Id: user.ID, Name: user.Name, Status: user.Status}}, nil

}

func (s *UserServer) Create(ctx context.Context, req *pb.CreateRequest) (*pb.CreateResponse, error) {
	result, err := s.db.Exec("INSERT INTO users (name, status) VALUES (?, ?)", req.Name, req.Status)
	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	return &pb.CreateResponse{Id: int32(id)}, nil
}

func (s *UserServer) Update(ctx context.Context, req *pb.UpdateRequest) (*pb.UpdateResponse, error) {
	_, err := s.db.Exec("UPDATE users SET name = ?, status = ? WHERE id = ?", req.Name, req.Status, req.Id)
	if err != nil {
		return nil, err
	}

	return &pb.UpdateResponse{}, nil
}

func (s *UserServer) Delete(ctx context.Context, req *pb.DeleteRequest) (*pb.DeleteResponse, error) {
	_, err := s.db.Exec("DELETE FROM users WHERE id = ?", req.Id)
	if err != nil {
		return nil, err
	}

	return &pb.DeleteResponse{}, nil
}
