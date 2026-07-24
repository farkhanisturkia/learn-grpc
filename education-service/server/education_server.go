package server

import (
	"context"
	"database/sql"
	"fmt"

	"education-service/pb"

	"connectrpc.com/connect"
)

type EducationServer struct {
	DB *sql.DB
}

func (s *EducationServer) CreateEducation(
	ctx context.Context,
	req *connect.Request[pb.CreateEducationRequest],
) (*connect.Response[pb.Education], error) {
	query := `INSERT INTO educations (user_id, level, program, university) VALUES (?, ?, ?, ?)`
	res, err := s.DB.ExecContext(ctx, query, req.Msg.UserId, req.Msg.Level, req.Msg.Program, req.Msg.University)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	id, _ := res.LastInsertId()
	return connect.NewResponse(&pb.Education{
		Id:      id,
		UserId:  req.Msg.UserId,
		Level:   req.Msg.Level,
		Program: req.Msg.Program,
		University: req.Msg.University,
	}), nil
}

func (s *EducationServer) GetEducation(
	ctx context.Context,
	req *connect.Request[pb.GetEducationRequest],
) (*connect.Response[pb.Education], error) {
	var education pb.Education
	query := `SELECT id, user_id, level, program, university FROM educations WHERE id = ?`
	err := s.DB.QueryRowContext(ctx, query, req.Msg.Id).Scan(&education.Id, &education.UserId, &education.Level, &education.Program, &education.University)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, connect.NewError(connect.CodeNotFound, fmt.Errorf("university tidak ditemukan"))
		}
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	return connect.NewResponse(&education), nil
}

func (s *EducationServer) GetEducationsByUser(
	ctx context.Context,
	req *connect.Request[pb.GetEducationsByUserRequest],
) (*connect.Response[pb.ListEducationsResponse], error) {
	query := `SELECT id, user_id, level, program, university FROM educations WHERE user_id = ?`
	rows, err := s.DB.QueryContext(ctx, query, req.Msg.UserId)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	defer rows.Close()

	var educations []*pb.Education
	for rows.Next() {
		var e pb.Education
		if err := rows.Scan(&e.Id, &e.UserId, &e.Level, &e.Program, &e.University); err != nil {
			return nil, connect.NewError(connect.CodeInternal, err)
		}
		educations = append(educations, &e)
	}

	return connect.NewResponse(&pb.ListEducationsResponse{Educations: educations}), nil
}

func (s *EducationServer) UpdateEducation(
	ctx context.Context,
	req *connect.Request[pb.UpdateEducationRequest],
) (*connect.Response[pb.Education], error) {
	query := `UPDATE educations SET level = ?, program = ?, university = ? WHERE id = ?`
	_, err := s.DB.ExecContext(ctx, query, req.Msg.Level, req.Msg.Program, req.Msg.University, req.Msg.Id)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	return s.GetEducation(ctx, connect.NewRequest(&pb.GetEducationRequest{Id: req.Msg.Id}))
}

func (s *EducationServer) DeleteEducation(
	ctx context.Context,
	req *connect.Request[pb.DeleteEducationRequest],
) (*connect.Response[pb.EmptyEducationResponse], error) {
	query := `DELETE FROM educations WHERE id = ?`
	_, err := s.DB.ExecContext(ctx, query, req.Msg.Id)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	return connect.NewResponse(&pb.EmptyEducationResponse{}), nil
}

func (s *EducationServer) DeleteEducationsByUser(
	ctx context.Context,
	req *connect.Request[pb.DeleteEducationsByUserRequest],
) (*connect.Response[pb.EmptyEducationResponse], error) {
	query := `DELETE FROM educations WHERE user_id = ?`
	_, err := s.DB.ExecContext(ctx, query, req.Msg.UserId)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	return connect.NewResponse(&pb.EmptyEducationResponse{}), nil
}