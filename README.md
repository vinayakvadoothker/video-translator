# 🚀 **Video Translator** 🌐

**Video Translator** is a simulation of an AI-powered video translation service. It demonstrates a client-server architecture with logging, configuration, and seamless testing workflows. This project can help you understand how to build and simulate asynchronous job processing systems.

## **📌 Key Features**

✅ **Simulated Video Translation Server**  
   - Asynchronous job processing with `pending`, `completed`, or `error` statuses.  
   - Simulated delay and error probabilities for realism.  
   - RESTful API design for extensibility.

✅ **Client Library**  
   - Interact with the server to create and track translation jobs.  
   - Exponential backoff for efficient polling and reduced server load.

✅ **Modern Logging**  
   - Structured, JSON-formatted logs powered by [Logrus](https://github.com/sirupsen/logrus).  

✅ **Easy Configuration**  
   - `.env`-based settings for server port, timeout durations, etc.

✅ **Comprehensive Testing**  
   - Unit tests for individual components.  
   - End-to-end integration tests with `httptest`.

---

## **🎯 Use Cases**

📽️ **AI Simulation**: Simulate time-intensive AI tasks such as video translation, processing, or rendering.  
📡 **API Prototyping**: Learn how to design and consume APIs effectively.  
🎓 **Education**: Understand client-server architecture and async processing.  

---

## **🛠️ Getting Started**

### **1. Prerequisites**
- **Go** (1.20+)

### **2. Installation**

Clone the repository

```bash
git clone https://github.com/vinayakvadoothker/video-translator.git
cd video-translator
```

Install dependencies

```bash
go mod tidy
```
---

## **🚀 Usage**

### Running the Server

To start the server:

```bash
make run-server
```
The server will run on http://localhost:8080 by default. You can configure the port and timeout in the `.env` file.

---

### Using the Client

1. **Create a Client**  

```bash
client := client.NewClient("http://localhost:8080")
```
2. **Start a New Job**  

   ```bash
   jobID, err := client.StartJob()
   if err != nil {
       log.Fatalf("Failed to start job: %v", err)
   }
   fmt.Printf("Job started with ID: %s\n", jobID)
   ```

3. **Poll for Job Status**  

   ```bash
   status, err := client.PollStatus(jobID)
   if err != nil {
       log.Fatalf("Failed to fetch job status: %v", err)
   }
   fmt.Printf("Job %s completed with status: %s\n", jobID, status)
   ```

---

## **🔍 Testing**

Run all tests:

```bash
make test
```
- **Unit Tests**: Validate individual client and server components.
- **Integration Tests**: End-to-end testing of client-server interaction.

---

## **⚙️ Configuration**

The application can be configured using a `.env` file placed in the root directory.

### Example `.env`

```bash
PORT=8080
JOB_TIMEOUT=10
```

- **`PORT`**: The port on which the server runs.  
- **`JOB_TIMEOUT`**: Simulated job completion delay in seconds.

---

## **🗂️ Project Structure**

```bash
video-translator/
├── client/                # Client library
│   ├── client.go          # Client implementation
│   ├── client_test.go     # Unit tests for client
├── cmd/
│   ├── server/            # Main entry point for server
│   │   ├── main.go
├── server/                # Server library
│   ├── main.go            # Server implementation
├── tests/                 # Integration tests
│   ├── integration_test.go
├── .env                   # Environment variables 
├── .gitignore             # Ignored files
├── go.mod                 # Go module dependencies
├── go.sum                 # Dependency checksums
├── LICENSE                # License for the project
├── Makefile               # Make commands for common tasks
└── README.md              # Project documentation
```
---

## **🛠️ Makefile Commands**

- **Run the Server**:

    ```make run-server```

- **Run Tests**:

    ```make test```

- **Build the Binary**:

  ```make build```

- **Run Built Binary**:

  ```make run-binary```

- **Clean Up**:

    ```make clean```

---

## **📜 License**

This project is licensed under the MIT License.

---

## **🙌 Acknowledgements**

- **Go** for the programming language.  
- **Logrus** for structured logging.  
- **Godotenv** for environment variable management.  

---

## **👨‍💻 Author**

[Vinayak Vadoothker](https://github.com/vinayakvadoothker)
