package server

import (
	"context"
	"database/sql"
	"fmt"

	"user-service/pb"
	"user-service/pb/pbconnect"

	"connectrpc.com/connect"
)

type UserServer struct {
	DB        *sql.DB
	JobClient pbconnect.JobServiceClient
	EducationClient pbconnect.EducationServiceClient
}

// 1. CreateUser
func (s *UserServer) CreateUser(
	ctx context.Context,
	req *connect.Request[pb.CreateUserRequest],
) (*connect.Response[pb.User], error) {
	query := `INSERT INTO users (name, email) VALUES (?, ?)`
	res, err := s.DB.ExecContext(ctx, query, req.Msg.Name, req.Msg.Email)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	id, _ := res.LastInsertId()
	return connect.NewResponse(&pb.User{
		Id:    id,
		Name:  req.Msg.Name,
		Email: req.Msg.Email,
	}), nil
}

// 2. GetUser (Ambil data user + List job-nya dari Job Service + List Education-nya dari Education Service)
func (s *UserServer) GetUser(
	ctx context.Context,
	req *connect.Request[pb.GetUserRequest],
) (*connect.Response[pb.User], error) {
	var user pb.User
	query := `SELECT id, name, email FROM users WHERE id = ?`
	err := s.DB.QueryRowContext(ctx, query, req.Msg.Id).Scan(&user.Id, &user.Name, &user.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, connect.NewError(connect.CodeNotFound, fmt.Errorf("user tidak ditemukan"))
		}
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	// Panggil Job Service via Connect RPC Client
	jobRes, err := s.JobClient.GetJobsByUser(ctx, connect.NewRequest(&pb.GetJobsByUserRequest{
		UserId: user.Id,
	}))
	if err == nil && jobRes.Msg != nil {
		user.Jobs = jobRes.Msg.Jobs
	}

	// Panggil Education Service via Connect RPC Client
	educationRes, err := s.EducationClient.GetEducationsByUser(ctx, connect.NewRequest(&pb.GetEducationsByUserRequest{
		UserId: user.Id,
	}))
	if err == nil && educationRes.Msg != nil {
		user.Educations = educationRes.Msg.Educations
	}

	return connect.NewResponse(&user), nil
}

// 3. UpdateUser
func (s *UserServer) UpdateUser(
	ctx context.Context,
	req *connect.Request[pb.UpdateUserRequest],
) (*connect.Response[pb.User], error) {
	query := `UPDATE users SET name = ?, email = ? WHERE id = ?`
	_, err := s.DB.ExecContext(ctx, query, req.Msg.Name, req.Msg.Email, req.Msg.Id)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	return s.GetUser(ctx, connect.NewRequest(&pb.GetUserRequest{Id: req.Msg.Id}))
}

// 4. DeleteUser (Hapus User + Cascade Hapus Job terkait via Job Service + Cascade Hapus Education terkait via Education Service)
func (s *UserServer) DeleteUser(
	ctx context.Context,
	req *connect.Request[pb.DeleteUserRequest],
) (*connect.Response[pb.EmptyUserResponse], error) {
	// Hapus user dari DB local
	query := `DELETE FROM users WHERE id = ?`
	_, err := s.DB.ExecContext(ctx, query, req.Msg.Id)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	// Hapus seluruh job terkait di Job Service
	_, _ = s.JobClient.DeleteJobsByUser(ctx, connect.NewRequest(&pb.DeleteJobsByUserRequest{
		UserId: req.Msg.Id,
	}))

	// Hapus seluruh education terkait di Education Service
	_, _ = s.EducationClient.DeleteEducationsByUser(ctx, connect.NewRequest(&pb.DeleteEducationsByUserRequest{
		UserId: req.Msg.Id,
	}))

	return connect.NewResponse(&pb.EmptyUserResponse{}), nil
}