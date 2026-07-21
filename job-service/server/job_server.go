package server

import (
	"context"
	"database/sql"
	"fmt"

	"job-service/pb"

	"connectrpc.com/connect"
)

type JobServer struct {
	DB *sql.DB
}

// 1. CreateJob
func (s *JobServer) CreateJob(
	ctx context.Context,
	req *connect.Request[pb.CreateJobRequest],
) (*connect.Response[pb.Job], error) {
	query := `INSERT INTO jobs (user_id, title, company) VALUES (?, ?, ?)`
	res, err := s.DB.ExecContext(ctx, query, req.Msg.UserId, req.Msg.Title, req.Msg.Company)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	id, _ := res.LastInsertId()
	return connect.NewResponse(&pb.Job{
		Id:      id,
		UserId:  req.Msg.UserId,
		Title:   req.Msg.Title,
		Company: req.Msg.Company,
	}), nil
}

// 2. GetJob
func (s *JobServer) GetJob(
	ctx context.Context,
	req *connect.Request[pb.GetJobRequest],
) (*connect.Response[pb.Job], error) {
	var job pb.Job
	query := `SELECT id, user_id, title, company FROM jobs WHERE id = ?`
	err := s.DB.QueryRowContext(ctx, query, req.Msg.Id).Scan(&job.Id, &job.UserId, &job.Title, &job.Company)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, connect.NewError(connect.CodeNotFound, fmt.Errorf("job tidak ditemukan"))
		}
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	return connect.NewResponse(&job), nil
}

// 3. GetJobsByUser (Daftar semua job milik user_id)
func (s *JobServer) GetJobsByUser(
	ctx context.Context,
	req *connect.Request[pb.GetJobsByUserRequest],
) (*connect.Response[pb.ListJobsResponse], error) {
	query := `SELECT id, user_id, title, company FROM jobs WHERE user_id = ?`
	rows, err := s.DB.QueryContext(ctx, query, req.Msg.UserId)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	defer rows.Close()

	var jobs []*pb.Job
	for rows.Next() {
		var j pb.Job
		if err := rows.Scan(&j.Id, &j.UserId, &j.Title, &j.Company); err != nil {
			return nil, connect.NewError(connect.CodeInternal, err)
		}
		jobs = append(jobs, &j)
	}

	return connect.NewResponse(&pb.ListJobsResponse{Jobs: jobs}), nil
}

// 4. UpdateJob
func (s *JobServer) UpdateJob(
	ctx context.Context,
	req *connect.Request[pb.UpdateJobRequest],
) (*connect.Response[pb.Job], error) {
	query := `UPDATE jobs SET title = ?, company = ? WHERE id = ?`
	_, err := s.DB.ExecContext(ctx, query, req.Msg.Title, req.Msg.Company, req.Msg.Id)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	return s.GetJob(ctx, connect.NewRequest(&pb.GetJobRequest{Id: req.Msg.Id}))
}

// 5. DeleteJob (Hapus 1 Job)
func (s *JobServer) DeleteJob(
	ctx context.Context,
	req *connect.Request[pb.DeleteJobRequest],
) (*connect.Response[pb.EmptyJobResponse], error) {
	query := `DELETE FROM jobs WHERE id = ?`
	_, err := s.DB.ExecContext(ctx, query, req.Msg.Id)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	return connect.NewResponse(&pb.EmptyJobResponse{}), nil
}

// 6. DeleteJobsByUser (Cascade delete seluruh job milik user_id)
func (s *JobServer) DeleteJobsByUser(
	ctx context.Context,
	req *connect.Request[pb.DeleteJobsByUserRequest],
) (*connect.Response[pb.EmptyJobResponse], error) {
	query := `DELETE FROM jobs WHERE user_id = ?`
	_, err := s.DB.ExecContext(ctx, query, req.Msg.UserId)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	return connect.NewResponse(&pb.EmptyJobResponse{}), nil
}