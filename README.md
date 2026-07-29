# 🚀 Go Microservices with Connect RPC

Project ini adalah implementasi *microservices* sederhana menggunakan **Golang** dan **Connect RPC** (tanpa Envoy Proxy / API Gateway). Terdiri dari 2 service independen yang berkomunikasi secara *inter-service*:

1. **`user-service`** (`:50051`): Mengelola data user & berkomunikasi dengan `job-service` dan `education-service` untuk agregasi data.
2. **`job-service`** (`:50052`): Mengelola data pekerjaan berbasis `userId`.
2. **`education-service`** (`:50053`): Mengelola data pendidikan berbasis `userId`.

---

## 🛠️ Prasyarat (Prerequisites)

Pastikan kamu sudah menginstal tools berikut di sistem kamu:

* **Go** (v1.20+)
* **Protocol Buffers Compiler (`protoc`)**
* **`protoc-gen-go`** & **`protoc-gen-connect-go`**
  ```bash
  go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
  go install [connectrpc.com/connect/cmd/protoc-gen-connect-go@latest](https://connectrpc.com/connect/cmd/protoc-gen-connect-go@latest)
  ```
---

## 🏗️ Struktur Project

  ```bash
  LEARN-GRPC/
  ├── proto/                  # Definisi gRPC/Connect Schema
  │   ├── user.proto
  │   ├── job.proto
  │   └── education.proto
  ├── user-service/           # User Service App
  │   ├── client/             # Internal Client untuk panggil Job Service
  │   ├── db/                 # Koneksi DB SQLite
  │   ├── pb/                 # Generated Protobuf Code (User + Job)
  │   ├── server/             # Implementation Server UserService
  │   ├── go.mod
  │   └── main.go
  ├── job-service/            # Job Service App
  │   ├── db/                 # Koneksi DB SQLite
  │   ├── pb/                 # Generated Protobuf Code (Job)
  │   ├── server/             # Implementation Server JobService
  │   ├── go.mod
  │   └── main.go
  └── education-service/      # Education Service App
      ├── db/                 # Koneksi DB SQLite
      ├── pb/                 # Generated Protobuf Code (Education)
      ├── server/             # Implementation Server EducationService
      ├── go.mod
      └── main.go
  ```
---

## ⚙️ 1. Generate Protobuf Code

Jalankan perintah ini dari root folder project (LEARN-GRPC) untuk mengompilasi file .proto ke masing-masing service:

* Compile Protobuf untuk user-service
  ```bash
  protoc -I=proto \
  --go_out=user-service/pb --go_opt=paths=source_relative \
  --go_opt=Mjob.proto=user-service/pb \
  --go_opt=Meducation.proto=user-service/pb \
  --connect-go_out=user-service/pb --connect-go_opt=paths=source_relative \
  --connect-go_opt=Mjob.proto=user-service/pb \
  --connect-go_opt=Meducation.proto=user-service/pb \
  user.proto job.proto education.proto
  ```

* Compile Protobuf untuk job-service
  ```bash
  protoc -I=proto \
  --go_out=job-service/pb --go_opt=paths=source_relative \
  --plugin=protoc-gen-connect-go=$(go env GOPATH)/bin/protoc-gen-connect-go \
  --connect-go_out=job-service/pb --connect-go_opt=paths=source_relative \
  job.proto
  ```

* Compile Protobuf untuk education-service
  ```bash
  protoc -I=proto \
  --go_out=education-service/pb --go_opt=paths=source_relative \
  --plugin=protoc-gen-connect-go=$(go env GOPATH)/bin/protoc-gen-connect-go \
  --connect-go_out=education-service/pb --connect-go_opt=paths=source_relative \
  education.proto
  ```
---

## 🏃 2. Menjalankan Service

Buka 2 tab terminal terpisah dan jalankan kedua service:

* Tab 1: Running User Service
  ```bash
  cd user-service
  go run .
  ```
    (Berjalan di http://localhost:50051)

* Tab 2: Running Job Service  
  ```bash
  cd job-service
  go run .
  ```
    (Berjalan di http://localhost:50052)

* Tab 2: Running Education Service  
  ```bash
  cd education-service
  go run .
  ```
    (Berjalan di http://localhost:50053)

---

## 🧪 3. Pengujian API (cURL + jq)
#### 👤 User Service
* Create User
  ```bash
  curl -s -X POST \
  -H "Content-Type: application/json" \
  -d '{
    "name": "User",
    "email": "user@example.id",
    "password": "userpassword"
  }' \
  http://localhost:50051/pb.UserService/CreateUser | jq
  ```

  ```bash
  curl -s -X POST \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Admin",
    "email": "admin@example.id",
    "password": "adminpassword",
    "role": "admin"
  }' \
  http://localhost:50051/pb.UserService/CreateUser | jq
  ```

* Get User
  ```bash
  curl -s -X POST \
  -H "Content-Type: application/json" \
  -d '{
    "id": 1
  }' \
  http://localhost:50051/pb.UserService/GetUser | jq
  ```

* Get User By Email
  ```bash
  curl -s -X POST \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.id"
  }' \
  http://localhost:50051/pb.UserService/GetUserByEmailInternal | jq
  ```

* Update User
  ```bash
  curl -s -X POST \
  -H "Content-Type: application/json" \
  -d '{
    "id": 1,
    "name": "User New",
    "email": "user.new@example.id"
  }' \
  http://localhost:50051/pb.UserService/UpdateUser | jq
  ```

* Delete User (Cascade Delete)
  ```bash
  curl -s -X POST \
  -H "Content-Type: application/json" \
  -d '{
    "id": 1
  }' \
  http://localhost:50051/pb.UserService/DeleteUser | jq
  ```

#### 💼 Job Service
* Create Job
  ```bash
  curl -s -X POST \
  -H "Content-Type: application/json" \
  -d '{
    "userId": 1,
    "title": "Backend Go Developer",
    "company": "MSRoot Indonesia"
  }' \
  http://localhost:50052/pb.JobService/CreateJob | jq
  ```

* Get Job
  ```bash
  curl -s -X POST \
  -H "Content-Type: application/json" \
  -d '{
    "id": 1
  }' \
  http://localhost:50052/pb.JobService/GetJob | jq
  ```

* Get Jobs By User ID
  ```bash
  curl -s -X POST \
  -H "Content-Type: application/json" \
  -d '{
    "userId": 1
  }' \
  http://localhost:50052/pb.JobService/GetJobsByUser | jq
  ```

* Update Job
  ```bash
  curl -s -X POST \
  -H "Content-Type: application/json" \
  -d '{
    "id": 1,
    "userId": 1,
    "title": "Senior Go Engineer",
    "company": "MSRoot.id Tech"
  }' \
  http://localhost:50052/pb.JobService/UpdateJob | jq
  ```

* Delete Job
  ```bash
  curl -s -X POST \
  -H "Content-Type: application/json" \
  -d '{
    "id": 1
  }' \
  http://localhost:50052/pb.JobService/DeleteJob | jq
  ```

* Delete Job By User ID
  ```bash
  curl -s -X POST \
  -H "Content-Type: application/json" \
  -d '{
    "userId": 1
  }' \
  http://localhost:50052/pb.JobService/DeleteJobsByUser | jq
  ```

#### 🎓 Education Service
* Create Education
  ```bash
  curl -s -X POST \
  -H "Content-Type: application/json" \
  -d '{
    "userId": 1,
    "level": "S1",
    "program": "Teknik Informatika",
    "university": "Universitas Negeri Surabaya"
  }' \
  http://localhost:50053/pb.EducationService/CreateEducation | jq
  ```

* Get Education
  ```bash
  curl -s -X POST \
  -H "Content-Type: application/json" \
  -d '{
    "id": 1
  }' \
  http://localhost:50053/pb.EducationService/GetEducation | jq
  ```

* Get Educations By User ID
  ```bash
  curl -s -X POST \
  -H "Content-Type: application/json" \
  -d '{
    "userId": 1
  }' \
  http://localhost:50053/pb.EducationService/GetEducationsByUser | jq
  ```

* Update Education
  ```bash
  curl -s -X POST \
  -H "Content-Type: application/json" \
  -d '{
    "id": 1,
    "userId": 1,
    "level": "S1",
    "program": "Ilmu Administrasi Negara",
    "university": "Universitas Negeri Surabaya"
  }' \
  http://localhost:50053/pb.EducationService/UpdateEducation | jq
  ```

* Delete Education
  ```bash
  curl -s -X POST \
  -H "Content-Type: application/json" \
  -d '{
    "id": 1
  }' \
  http://localhost:50053/pb.EducationService/DeleteEducation | jq
  ```

* Delete Education By User ID
  ```bash
  curl -s -X POST \
  -H "Content-Type: application/json" \
  -d '{
    "userId": 1
  }' \
  http://localhost:50053/pb.EducationService/DeleteEducationsByUser | jq
  ```
