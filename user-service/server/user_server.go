package server

import (
	"context"
	"database/sql"
	"fmt"

	"user-service/pb"
	"user-service/pb/pbconnect"

	"connectrpc.com/connect"
	"golang.org/x/crypto/bcrypt"
)

type UserServer struct {
	DB        *sql.DB
	JobClient pbconnect.JobServiceClient
	EducationClient pbconnect.EducationServiceClient
}

func (s *UserServer) CreateUser(
	ctx context.Context,
	req *connect.Request[pb.CreateUserRequest],
) (*connect.Response[pb.User], error) {
	role := req.Msg.Role
    if role != "admin" && role != "user" {
        role = "user"
    }

	var hashedPassword string
    if req.Msg.Password != "" {
        hash, err := bcrypt.GenerateFromPassword([]byte(req.Msg.Password), bcrypt.DefaultCost)
        if err != nil {
            return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("gagal hash password: %w", err))
        }
        hashedPassword = string(hash)
    }

	query := `INSERT INTO users (name, email, password, role) VALUES (?, ?, ?, ?)`
    res, err := s.DB.ExecContext(ctx, query, req.Msg.Name, req.Msg.Email, hashedPassword, role)
    if err != nil {
        return nil, connect.NewError(connect.CodeInternal, err)
    }

	id, _ := res.LastInsertId()
	return connect.NewResponse(&pb.User{
		Id:    id,
		Name:  req.Msg.Name,
		Email: req.Msg.Email,
		Role:  role,
	}), nil
}

func (s *UserServer) GetUser(
	ctx context.Context,
	req *connect.Request[pb.GetUserRequest],
) (*connect.Response[pb.User], error) {
	var user pb.User
	query := `SELECT id, name, email, role FROM users WHERE id = ?`
	err := s.DB.QueryRowContext(ctx, query, req.Msg.Id).Scan(&user.Id, &user.Name, &user.Email, &user.Role)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, connect.NewError(connect.CodeNotFound, fmt.Errorf("user tidak ditemukan"))
		}
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	jobRes, err := s.JobClient.GetJobsByUser(ctx, connect.NewRequest(&pb.GetJobsByUserRequest{
		UserId: user.Id,
	}))
	if err == nil && jobRes.Msg != nil {
		user.Jobs = jobRes.Msg.Jobs
	}

	educationRes, err := s.EducationClient.GetEducationsByUser(ctx, connect.NewRequest(&pb.GetEducationsByUserRequest{
		UserId: user.Id,
	}))
	if err == nil && educationRes.Msg != nil {
		user.Educations = educationRes.Msg.Educations
	}

	return connect.NewResponse(&user), nil
}

func (s *UserServer) GetUserByEmailInternal(
	ctx context.Context,
	req *connect.Request[pb.GetUserByEmailRequest],
) (*connect.Response[pb.UserInternal], error) {
	var userInternal pb.UserInternal
	query := `SELECT id, name, email, role, password FROM users WHERE email = ?`
	err := s.DB.QueryRowContext(ctx, query, req.Msg.Email).Scan(&userInternal.Id, &userInternal.Name, &userInternal.Email, &userInternal.Role, &userInternal.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, connect.NewError(connect.CodeNotFound, fmt.Errorf("user tidak ditemukan"))
		}
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	return connect.NewResponse(&userInternal), nil
}

func (s *UserServer) UpdateUser(
	ctx context.Context,
	req *connect.Request[pb.UpdateUserRequest],
) (*connect.Response[pb.User], error) {
	query := `UPDATE users SET name = ?, email = ?, role = ? WHERE id = ?`
	_, err := s.DB.ExecContext(ctx, query, req.Msg.Name, req.Msg.Email, req.Msg.Role, req.Msg.Id)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	return s.GetUser(ctx, connect.NewRequest(&pb.GetUserRequest{Id: req.Msg.Id}))
}

func (s *UserServer) DeleteUser(
	ctx context.Context,
	req *connect.Request[pb.DeleteUserRequest],
) (*connect.Response[pb.EmptyUserResponse], error) {
	query := `DELETE FROM users WHERE id = ?`
	_, err := s.DB.ExecContext(ctx, query, req.Msg.Id)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	_, _ = s.JobClient.DeleteJobsByUser(ctx, connect.NewRequest(&pb.DeleteJobsByUserRequest{
		UserId: req.Msg.Id,
	}))

	_, _ = s.EducationClient.DeleteEducationsByUser(ctx, connect.NewRequest(&pb.DeleteEducationsByUserRequest{
		UserId: req.Msg.Id,
	}))

	return connect.NewResponse(&pb.EmptyUserResponse{}), nil
}