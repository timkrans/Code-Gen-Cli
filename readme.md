# Code-Gen-CLI

A command-line interface for generating code using LLM (Large Language Model) agents. This tool allows you to describe what code you want, and it will generate the corresponding Go code files for you.

## Features

- **Interactive CLI**: Simple command-line interface for code generation
- **LLM Integration**: Uses Ollama with Llama 3.2 model for code generation
- **Multi-file Support**: Generates multiple files in a single request
- **File System Management**: Automatically creates directories and writes files
- **Streaming Response**: Real-time streaming from the LLM for better user experience

## Prerequisites

- Go 1.24.5 or later
- [Ollama](https://ollama.ai/) installed and running locally
- Llama 3.2 model pulled in Ollama

## Installation

1. Clone the repository:
```bash
git clone <repository-url>
cd Code-Gen-CLI
```

2. Install dependencies:
```bash
go mod tidy
```

3. Build the application:
```bash
go build -o code-gen-cli cmd/main.go
```

## Setup Ollama

1. Install Ollama from [https://ollama.ai/](https://ollama.ai/)
2. Start Ollama service:
```bash
ollama serve
```

3. Pull the Llama 3.2 model:
```bash
ollama pull llama3.2
```

## Usage

Run the application:
```bash
./code-gen-cli
```

The application will start an interactive session where you can:

1. Type your code generation prompt
2. The tool will connect to Ollama and generate code
3. Generated files will be written to the `./test` directory
4. Type `quit` or `exit` to leave the application

### Example Usage

```
Welcome to LLM Codegen Agent!
Type a prompt and I'll generate code for you.
Type 'quit' or 'exit' to leave.

You: Create a simple HTTP server with a health check endpoint
Generating code...
Connecting to Ollama...
Written: ./example/main.go
Written: ./example/handlers.go
Code written to './example'
--------------------------------------------------
```

## Project Structure

```
Code-Gen-CLI/
├── agents/
│   ├── code-gen/
│   │   └── llm.go          # LLM integration with Ollama
│   └── fs/
│       └── fs.go           # File system operations
├── cmd/
│   └── main.go             # Main application entry point
├── go.mod                  # Go module definition
└── README.md              # This file
```

## Architecture

The application is structured with a modular architecture:

- **`cmd/main.go`**: Main application logic and user interaction
- **`agents/code-gen/llm.go`**: Handles communication with Ollama API and code generation
- **`agents/fs/fs.go`**: Manages file system operations for writing generated code

## API Integration

The tool integrates with Ollama's API using the following endpoints:
- **Endpoint**: `http://localhost:11434/api/generate`
- **Model**: `llama3.2`
- **Streaming**: Enabled for real-time response

## Generated Code Format

The LLM generates code using a specific format:
```
/// FILE: <relative_path>
<code content>
```

This format allows the tool to parse multiple files from a single LLM response and write them to the appropriate locations.

## Configuration

- **Output Directory**: Currently hardcoded to `./test` (can be modified in `cmd/main.go`)
- **Ollama URL**: `http://localhost:11434` (default Ollama installation)
- **Model**: `llama3.2` (can be changed in `agents/code-gen/llm.go`)

## Troubleshooting

### Common Issues

1. **Connection Error**: Make sure Ollama is running (`ollama serve`)
2. **Model Not Found**: Ensure Llama 3.2 is pulled (`ollama pull llama3.2`)
3. **Permission Errors**: Check write permissions for the output directory

## Future Goal

1. Build out to allow for real time context of projects
2. Build an optimization model to optimize the per token context of above
3. Incorperate with other models
4. Add commands such as npm and pythons pip and venv
5. Add realtime error handing giving suggested ideas