# Milkyteadrop Fileserver

Milkyteadrop Fileserver is a robust and efficient file server built in Go. It is designed to handle and store a variety of file types, including images and text, and offers a RESTful API for easy file creation and retrieval. This makes it an ideal solution for integrating with other services and applications that require a reliable file management system.

## Features

- **File Handling**: Supports different file types including images and text.
- **RESTful API**: Provides endpoints for file creation and retrieval.
- **Base64 Encoding**: Handles base64 encoded data for image files.
- **File Existence Check**: Verifies if a file already exists before creation to prevent duplicates.

## Getting Started

### Prerequisites

- Go **1.21.5**

### Installation

1. Clone the repository to your local machine:
   git clone https://github.com/CoffeeeAtNight/MilkyTeadrop_FileServer

2. Navigate to the project directory:
   cd milkyteadrop-fileserver

3. Compile and run the server:
   go build
   sudo ./milkyteadrop-fileserver

## API Endpoints

### Create File

- **POST** `/api/v1/create/file`
- **Request Body**:
  {
    "filename": "example.png",
    "filetype": "image",
    "fileContent": "base64_encoded_data"
  }
- **Response**:
  - `200 OK` on success:
    {
      "status": 200,
      "message": "Successfully created file",
      "body": "example.png"
    }
  - `400 Bad Request` on invalid input.
  - `500 Internal Server Error` on server error.

### Retrieve File

- **GET** `/api/v1/file/[filename]`
- Retrieves the specified file if it exists.
- **Response**:
  - File content in the specified format on success.
  - `404 Not Found` if the file does not exist.

## Configuration

The file server uses a default root path defined as `/usr/local/bin/milkyteadrop-fs/`. You can modify this path in the source code to suit your deployment environment.

## Usage

To interact with the file server, use HTTP requests to the provided endpoints. You can use tools like `curl` or Postman, or integrate the API calls into your application.


## Note

THIS FILESERVER IS NOT INTENDED FOR BEING USED IN PRODUCTION ENVIRONMENTS SINCE IT'S EXPOSES ENDPOINTS WHICH ARE NOT SECURE, USE AT OWN RISK. IT'S ONLY A FUN PROJECT OF MINE!
