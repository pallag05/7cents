# 7cents - Group Study Management System

A Go-based backend service for managing study groups, user matching, and collaborative learning.

## Features

- Create and manage study groups
- User matching based on academic performance
- Group messaging system
- Question bank integration
- Tag-based group search
- Activity scoring system

## Prerequisites

- Go 1.16 or higher
- Git

## Setup

1. Clone the repository:
```bash
git clone https://github.com/Allen-Career-Institute/7cents.git
cd 7cents
```

2. Install dependencies:
```bash
go mod download
```

3. Run the server:
```bash
go run main.go
```

The server will start on port 96 by default.

## API Endpoints

### Groups

#### Create Group
- **POST** `/api/groups`
- Creates a new study group
```json
{
    "title": "Physics Study Group",
    "description": "Advanced physics study group",
    "tag": "physics",
    "type": "study",
    "private": false,
    "capacity": 10
}
```

#### Get User's Groups
- **GET** `/api/groups/user/:user_id`
- Returns all groups associated with a user

#### Get Group by ID
- **GET** `/api/groups/:id`
- Returns details of a specific group

#### Join Group
- **POST** `/api/groups/:id/join/:user_id`
- Adds a user to a group

#### Leave Group
- **POST** `/api/groups/:id/leave/:user_id`
- Removes a user from a group

#### Update Group
- **PUT** `/api/groups/:id`
- Updates group information
```json
{
    "message": {
        "content": "New message content",
        "sender_id": "user_id"
    }
}
```
or
```json
{
    "action": {
        "type": "action_type",
        "content": "action content"
    }
}
```

#### Search Groups by Tag
- **POST** `/api/groups/search`
- Searches for groups based on tag
```json
{
    "tag": "physics"
}
```

## Data Models

### User
```go
type User struct {
    ID    string
    Email string
    Score []Score
}
```

### Group
```go
type Group struct {
    ID            string
    Title         string
    Description   string
    Members       []string
    Tag           string
    Type          string
    Private       bool
    Messages      []Message
    Actions       []Action
    CreateBy      string
    Capacity      int
    ActivityScore int
    Questions     []Question
}
```

### Question
```go
type Question struct {
    ID        string
    Content   string
    Options   []string
    Timestamp time.Time
}
```

## Features in Detail

### User Matching System
- Automatically matches users based on academic performance
- Calculates similarity scores between users
- Creates paired study groups for highly compatible users

### Activity Scoring
- Groups are ranked by activity score
- Scores are influenced by member participation and group interactions

### Question Bank
- Integrated question system for study groups
- Multiple choice questions across different subjects
- Real-time question updates in groups

## Error Handling

The API returns appropriate HTTP status codes:
- 200: Success
- 400: Bad Request
- 404: Not Found
- 500: Internal Server Error

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the LICENSE file for details. 