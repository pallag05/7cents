# Learning Activities Streak System

A robust Go-based backend system for tracking learning activities, managing user streaks, and implementing a reward system based on user performance.

## Features

- **User Management**
  - User registration and profile management
  - Batch-based user organization
  - User activity tracking

- **Streak System**
  - Daily activity tracking
  - Streak counting and maintenance
  - Streak break detection
  - Activity history

- **Rating System**
  - Dynamic rating calculation based on:
    - Streak consistency
    - Activity duration
    - Activity frequency
  - Rating tiers: Beginner → Novice → Intermediate → Expert
  - Rating breakdown and progress tracking

- **Leaderboard System**
  - Batch-specific leaderboards
  - Global top performers
  - Rating distribution visualization
  - Streak distribution analysis

- **Reward System**
  - Rating-based reward tiers
  - Progress tracking towards rewards
  - Reward history
  - Available rewards display

## System Architecture

The system follows a clean architecture pattern with the following components:

- **Handlers**: HTTP request handlers
- **Services**: Business logic layer
- **Models**: Data structures and entities
- **Storage**: Data persistence layer (in-memory for this implementation)

## API Endpoints

### User Management
- `POST /api/users` - Create new user
- `GET /api/users/:user_id` - Get user details

### Streak Management
- `POST /api/streaks/activity` - Record learning activity
- `GET /api/streaks/user/:user_id` - Get user streak information

### Rating System
- `GET /api/ratings/user/:user_id` - Get user rating
- `GET /api/ratings/user/:user_id/breakdown` - Get rating calculation breakdown

### Leaderboard
- `GET /api/leaderboards/batch/:batch_id` - Get batch leaderboard
- `GET /api/leaderboards/top` - Get top performers
- `GET /api/leaderboards/batch/:batch_id/stats` - Get leaderboard statistics
- `GET /api/leaderboards/batch/:batch_id/rating-distribution` - Get rating distribution
- `GET /api/leaderboards/batch/:batch_id/streak-distribution` - Get streak distribution

### Rewards
- `GET /api/rewards/user/:user_id` - Get user rewards
- `GET /api/rewards/reward/:reward_id` - Get reward details
- `GET /api/rewards/available/:rating` - Get available rewards for rating
- `GET /api/rewards/progress/:user_id` - Get reward progress

## System Flow

```mermaid
sequenceDiagram
    participant User
    participant Frontend
    participant API
    participant UserService
    participant StreakService
    participant RewardService
    participant MemoryStore

    %% User Registration Flow
    User->>Frontend: Register with name, phone, batch
    Frontend->>API: POST /api/users
    API->>UserService: CreateUser(user)
    UserService->>MemoryStore: SaveUser(user)
    MemoryStore-->>UserService: User created
    UserService-->>API: Return user
    API-->>Frontend: User response
    Frontend-->>User: Show confirmation

    %% Activity Recording Flow
    User->>Frontend: Record learning activity
    Frontend->>API: POST /api/streaks/activity
    API->>StreakService: RecordActivity(activity)
    StreakService->>MemoryStore: Update user streak
    MemoryStore-->>StreakService: Streak updated
    StreakService->>StreakService: Calculate rating
    StreakService-->>API: Activity recorded
    API-->>Frontend: Success response
    Frontend-->>User: Show updated streak

    %% Leaderboard Flow
    User->>Frontend: View batch leaderboard
    Frontend->>API: GET /api/leaderboards/batch/:batch_id
    API->>StreakService: GetBatchLeaderboard(batchID)
    StreakService->>MemoryStore: Get users by batch
    MemoryStore-->>StreakService: Users list
    StreakService->>StreakService: Calculate scores
    StreakService->>StreakService: Sort by score
    StreakService-->>API: Leaderboard data
    API-->>Frontend: Leaderboard response
    Frontend-->>User: Display leaderboard

    %% Reward System Flow
    User->>Frontend: Check available rewards
    Frontend->>API: GET /api/rewards/available/:rating
    API->>RewardService: GetAvailableRewards(rating)
    RewardService->>MemoryStore: Get rewards for rating
    MemoryStore-->>RewardService: Rewards list
    RewardService-->>API: Available rewards
    API-->>Frontend: Rewards response
    Frontend-->>User: Display rewards

    %% Rating Calculation Flow
    User->>Frontend: View rating details
    Frontend->>API: GET /api/ratings/user/:user_id/breakdown
    API->>StreakService: GetRatingBreakdown(userID)
    StreakService->>MemoryStore: Get user activities
    MemoryStore-->>StreakService: Activities list
    StreakService->>StreakService: Calculate rating components
    StreakService-->>API: Rating breakdown
    API-->>Frontend: Rating response
    Frontend-->>User: Display rating details
```

## Getting Started

### Prerequisites

- Go 1.16 or higher
- Git

### Installation

1. Clone the repository:
```bash
git clone https://github.com/yourusername/learning-activities-streak.git
cd learning-activities-streak
```

2. Install dependencies:
```bash
go mod download
```

3. Run the application:
```bash
go run main.go
```

The server will start on `http://localhost:96`

## Testing

Run the test suite:
```bash
go test ./...
```

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/AmazingFeature`)
3. Commit your changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Acknowledgments

- Built with [Gin](https://github.com/gin-gonic/gin) web framework
- Uses [UUID](https://github.com/google/uuid) for unique identifiers
- Inspired by gamification principles in learning platforms