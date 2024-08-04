# Tinder Clone

A full-featured Tinder-like application built with Go. This project includes user registration, profile management, geolocation-based matching, and more. The following sections outline the main features and TODOs for this project.

## Technologies
* Go
* Gin
* Gorm
* JWT
* Swagger/OpenApi
* PostgreSQL
* net/http
* API integration
* Redis
* Asynq
* Testing
* Sqlite3
* testify
* go-pagination
* go-playground
* smtp
* crypto/bcrypt

## Features

### 1. Registration and Authorization of Users
- **Registration**: Users can register by providing their basic information such as name, email address, phone number, and creating a password.
- **Authorization**: Users can log in using registered credentials.
- **Social Authentication (Planned)**: Possibility of registering and logging in via social networks (e.g., Facebook, Google).

### 2. Profile Creation and Management
- **Profile Information**: Users can add detailed information about themselves such as biography, interests, age, location, and gender.
- **Photos**: Ability to upload and manage profile photos.

### 3. Search and Selection of Pairs
- **Profile Display**: For now, the application just gives out profiles.
- **Swipe Functionality**: Ability to swipe profiles (left or right), expressing interest or refusal.
- **Display Users Who Liked Me**: Add logic for response interaction.
- **Swipe Limit for Unsubscribed Users**: Limit the number of swipes (e.g., 10 swipes) for unsubscribed users.
- **Profile Retrieval**: Get profiles without pagination and with pagination.

### 4. User Interaction (Planned)
- **Messages**: Built-in messaging system for communication between matched users (matches).
- **Notifications**: Push notifications and in-app notifications for new messages, matches, and likes.

### 5. Additional Features
- **Geolocation**: Using location to match nearby users.
- **Subscriptions and Premium Features (Planned)**: Paid features such as unlimited likes, the ability to see who has viewed a profile, and Rewind.

### 6. User Interface and Experience (UI/UX) (Planned)
- **Intuitive Design**: Easy to use and attractive interface for improved user experience.

### 7. Other Features
- **Return Profile for Current User**: Return the profile for the current user if the path does not have this parameter.
- **User Retrieval from Token**: Implement user retrieval from token, not path.
- **Bug Fixes**: Fix bugs in `EditProfileController`.
- **Test Photo Model**: Test the photo model.
- **Email Validation**: Add email validation.
- **Improve Database Queries**: Avoid "record not found" messages.
- **Location-Based Features**: Implement and test algorithms for suggesting users based on common interests, location, and other factors.
- **User Form Validation**: Ensure all user forms are validated.
- **Refresh Token (JWT)**: Implement valid JWT refresh tokens.
- **Redis and Asynchronous Tasks**: 
  - Configure Redis for caching and add basic operations.
  - Set up and test Asynq client for background tasks.
- **Periodic Tasks**: Implement a periodic task for deleting users with `is_active = false` if `CreatedAt` is older than 24 hours.
- **Refactor Controllers**: Refactor all controllers to use the repository pattern.
- **Simplify Main File**: Simplify the `main.go` file.
- **Admin Management**: Implement basic admin functionality without a dedicated admin role for simplicity in frontend development.
- **Admin User Creation**: Add a script for creating an admin user (optional).
- **Tests with TestDB**: Write tests using a test database (SQLite).
- **Log Analysis Service**: Add a service for log analysis (e.g., Prometheus).
- **Docker Support**: Create a `Dockerfile` and `docker-compose` configuration for easy deployment.

## Getting Started

### Prerequisites
- Go (latest version)
- PostgreSQL (or SQLite for testing)
- Redis
- Docker (for containerization)

### Installation
1. Clone the repository:
    ```sh
    git clone https://github.com/ilyaDyb/tinder-clone.git
    ```
2. Navigate to the project directory:
    ```sh
    cd tinder-clone
    ```
3. Install dependencies:
    ```sh
    go mod download
    ```

### Running the Application
1. Set up the database and Redis:
    ```sh
    docker-compose up -d
    ```
2. Run the application:
    ```sh
    go run main.go
    ```

### Running Tests
1. Set up the test database:
    ```sh
    export TEST_DB=sqlite3://test.db
    ```
2. Run the tests:
    ```sh
    go test ./...
    ```

## Contributing
1. Fork the repository.
2. Create a new branch (`git checkout -b feature-branch`).
3. Commit your changes (`git commit -am 'Add new feature'`).
4. Push to the branch (`git push origin feature-branch`).
5. Open a pull request.

## License
This project is licensed under the MIT License.

## Contact
For any questions or suggestions, feel free to open an issue or contact the repository owner.

`Telegram: http://t.me/wicki7`

## Frontend
[Frontend part for visualization of API](https://github.com/ilyaDyb/tinder_frontend)

This README file provides an overview of the project, outlines the main features, and includes instructions for setting up and running the application. It also contains guidelines for contributing to the project.
