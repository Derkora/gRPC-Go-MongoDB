package main

import (
	"context"
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/derkora/master/common/genproto/taskmaster"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

type httpServer struct {
	addr string
}

func NewHttpServer(addr string) *httpServer {
	return &httpServer{addr: addr}
}

func init() {
	args := os.Args[1:]
	var configname string = "default-config"
	if len(args) > 0 {
		configname = args[0] + "-config"
	}
	log.Printf("loading config file %s.yml \n", configname)

	viper.SetConfigName(configname)
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./configs")
	err := viper.ReadInConfig()
	if err != nil {
		panic("Fatal error config file: " + err.Error())
	}
}

func (s *httpServer) Run() error {
	http.HandleFunc("/", s.handleIndex)
	http.HandleFunc("/create", s.handleCreate)
	http.HandleFunc("/update", s.handleUpdate)
	http.HandleFunc("/delete", s.handleDelete)
	http.HandleFunc("/list", s.handleList)
	log.Printf("Starting HTTP server on localhost%s/list\n", s.addr)
	return http.ListenAndServe(s.addr, nil)
}

func (s *httpServer) handleIndex(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.New("index").Parse(tasksTemplate))
	if err := tmpl.Execute(w, nil); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// Handlers for CRUD operations
func (s *httpServer) handleCreate(w http.ResponseWriter, r *http.Request) {
	// Mendapatkan data dari form HTML
	title := r.FormValue("title")
	description := r.FormValue("description")

	// Menginisialisasi koneksi gRPC
	port := ":" + viper.GetString("app.grpc.port")
	client, err := grpc.Dial(port, grpc.WithInsecure())
	if err != nil {
		http.Error(w, "Could not connect to gRPC server", http.StatusInternalServerError)
		return
	}
	defer client.Close()

	// Membuat klien API tugas
	taskClient := taskmaster.NewTaskApiClient(client)

	// Membuat task baru
	_, err = taskClient.CreateTask(context.Background(), &taskmaster.Task{
		Title:       title,
		Description: description,
	})
	if err != nil {
		http.Error(w, "Failed to create task", http.StatusInternalServerError)
		return
	}

	// Tetapkan pengalihan arahan ke halaman list
	http.Redirect(w, r, "/list", http.StatusSeeOther)
}

func (s *httpServer) handleUpdate(w http.ResponseWriter, r *http.Request) {
	// Mendapatkan data dari form HTML
	id := r.FormValue("id")
	title := r.FormValue("title")
	description := r.FormValue("description")

	// Menginisialisasi koneksi gRPC
	port := ":" + viper.GetString("app.grpc.port")
	client, err := grpc.Dial(port, grpc.WithInsecure())
	if err != nil {
		http.Error(w, "Could not connect to gRPC server", http.StatusInternalServerError)
		return
	}
	defer client.Close()

	// Membuat klien API tugas
	taskClient := taskmaster.NewTaskApiClient(client)

	// Update task
	_, err = taskClient.UpdateTask(context.Background(), &taskmaster.Task{
		Id:          id,
		Title:       title,
		Description: description,
	})
	if err != nil {
		http.Error(w, "Failed to update task", http.StatusInternalServerError)
		return
	}

	// Tetapkan pengalihan arahan ke halaman list
	http.Redirect(w, r, "/list", http.StatusSeeOther)
}

func (s *httpServer) handleDelete(w http.ResponseWriter, r *http.Request) {
	// Mendapatkan ID task dari parameter URL
	id := r.URL.Query().Get("id")

	// Menginisialisasi koneksi gRPC
	port := ":" + viper.GetString("app.grpc.port")
	client, err := grpc.Dial(port, grpc.WithInsecure())
	if err != nil {
		http.Error(w, "Could not connect to gRPC server", http.StatusInternalServerError)
		return
	}
	defer client.Close()

	// Membuat klien API tugas
	taskClient := taskmaster.NewTaskApiClient(client)

	// Menghapus task
	_, err = taskClient.DeleteTask(context.Background(), &wrapperspb.StringValue{Value: id})
	if err != nil {
		http.Error(w, "Failed to delete task", http.StatusInternalServerError)
		return
	}

	// Tetapkan pengalihan arahan ke halaman list
	http.Redirect(w, r, "/list", http.StatusSeeOther)
}

func (s *httpServer) handleList(w http.ResponseWriter, r *http.Request) {
	// Menginisialisasi koneksi gRPC
	port := ":" + viper.GetString("app.grpc.port")
	client, err := grpc.Dial(port, grpc.WithInsecure())
	if err != nil {
		http.Error(w, "Could not connect to gRPC server", http.StatusInternalServerError)
		return
	}
	defer client.Close()

	// Membuat klien API tugas
	taskClient := taskmaster.NewTaskApiClient(client)

	// Mengambil daftar task dari server
	tasks, err := taskClient.ListTasks(context.Background(), &emptypb.Empty{})
	if err != nil {
		http.Error(w, "Failed to fetch tasks: "+err.Error(), http.StatusInternalServerError)
		log.Printf("Failed to fetch tasks: %v\n", err)
		return
	}

	// Menyiapkan data task untuk ditampilkan di halaman HTML
	type ViewData struct {
		Tasks []*taskmaster.Task
	}
	data := ViewData{
		Tasks: tasks.List,
	}

	// Membuat template HTML
	tmpl := template.Must(template.New("index").Parse(tasksTemplate))

	// Menampilkan template HTML dengan data task yang telah disiapkan
	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func main() {
	httpServer := NewHttpServer(":1000")
	httpServer.Run()
}

var tasksTemplate = `
<!DOCTYPE html>
<html>
<head>
    <title>Task Master</title>
    <style>
        body {
            font-family: Arial, sans-serif;
            margin: 0;
            padding: 0;
            background-color: #f2f2f2;
        }
        .container {
            max-width: 800px;
            margin: 20px auto;
            padding: 20px;
            background-color: #fff;
            border-radius: 5px;
            box-shadow: 0 0 10px rgba(0, 0, 0, 0.1);
        }
        h1 {
            color: #333;
        }
        form {
            margin-bottom: 20px;
        }
        label {
            display: block;
            margin-bottom: 5px;
        }
        input[type="text"],
        textarea {
            width: 100%;
            padding: 10px;
            margin-bottom: 10px;
            border: 1px solid #ccc;
            border-radius: 4px;
            box-sizing: border-box;
        }
        input[type="submit"] {
            background-color: #4caf50;
            color: white;
            border: none;
            border-radius: 4px;
            padding: 10px 20px;
            cursor: pointer;
        }
        input[type="submit"]:hover {
            background-color: #45a049;
        }
        ul {
            list-style-type: none;
            padding: 0;
        }
        li {
            padding: 10px;
            margin-bottom: 5px;
            background-color: #f9f9f9;
            border-radius: 4px;
            box-shadow: 0 0 5px rgba(0, 0, 0, 0.1);
        }
        a {
            text-decoration: none;
            color: #4caf50;
        }
        a:hover {
            text-decoration: underline;
        }
        .refresh-btn {
            background-color: #008CBA;
            color: white;
            border: none;
            border-radius: 4px;
            padding: 10px 20px;
            cursor: pointer;
            text-decoration: none;
        }
        .refresh-btn:hover {
            background-color: #005f6b;
        }
        .update-form {
            display: none;
            margin-bottom: 10px;
        }
        .update-form input[type="submit"] {
            background-color: #4caf50;
            color: white;
            border: none;
            border-radius: 4px;
            padding: 10px 20px;
            cursor: pointer;
        }
		.back-btn {
			background-color: #f44336;
			color: white;
			border: none;
			border-radius: 4px;
			padding: 10px 20px;
			cursor: pointer;
			text-decoration: none;
			margin-right: 10px; /* Margin kanan untuk memberi jarak dari tombol "Update" */
		}
		
		.back-btn:hover {
			background-color: #d32f2f;
		}		
    </style>
</head>
<body>
    <div class="container">
        <h1>Task Master</h1>
        <form action="/create" method="post">
            <label for="title">Title:</label><br>
            <input type="text" id="title" name="title"><br>
            <label for="description">Description:</label><br>
            <textarea id="description" name="description"></textarea><br><br>
            <input type="submit" value="Create Task">
        </form>
        <hr>
        <h2>Task List</h2>
        <a href="/list" class="refresh-btn">Refresh Task List</a>
        <ul>
            {{if not (eq (len .Tasks) 0)}}
                {{range .Tasks}}
                <li>
                    <div>
                        <span>{{.Title}} - {{.Description}}</span>
                        <a href="#" onclick="showUpdateForm('{{.Id}}')">Update</a> 
                        <a href="/delete?id={{.Id}}">Delete</a>
                    </div>
                    <form id="updateForm{{.Id}}" class="update-form" action="/update" method="post">
						<input type="hidden" name="id" value="{{.Id}}">
						<label for="title{{.Id}}">New Title:</label><br>
						<input type="text" id="title{{.Id}}" name="title" value="{{.Title}}"><br>
						<label for="description{{.Id}}">New Description:</label><br>
						<textarea id="description{{.Id}}" name="description">{{.Description}}</textarea><br><br>
						<input type="submit" value="Update Task">
						<a href="/list" class="back-btn">Back</a>
					</form>
                </li>
                {{end}}
            {{else}}
                <li>No tasks available</li>
            {{end}}
        </ul>
    </div>

    <script>
        // Fungsi untuk menampilkan formulir update task
        function showUpdateForm(taskId) {
            var formId = 'updateForm' + taskId;
            var form = document.getElementById(formId);
            if (form.style.display === 'none') {
                form.style.display = 'block';
            } else {
                form.style.display = 'none';
            }
        }
    </script>
</body>
</html>
`
