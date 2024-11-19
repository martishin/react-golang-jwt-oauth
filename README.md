# React-Golang OAuth
An example project demonstrating user authentication and authorization with a React.js frontend and a Go backend. It integrates Google OAuth 2.0 for secure login and uses JWT with refresh tokens for session management.

This project uses JWT tokens, for a simpler setup that uses session cookies, please check [this project](https://github.com/martishin/react-golang-goth-oauth).

<img src="https://i.giphy.com/media/v1.Y2lkPTc5MGI3NjExbnRwbDB6cmN2emtiaXhpY3hydWI3ZGJtbGM0cHZ2dzEzZXAxaHA5dCZlcD12MV9pbnRlcm5hbF9naWZfYnlfaWQmY3Q9Zw/KtKvOlylZtd9oOJQNF/giphy.gif" width="400"/>

## Setup
### Prerequisites
1. Install the dependencies:
   - [Docker](https://www.docker.com/products/docker-desktop)
   - [Node.js and npm](https://nodejs.org/) (for the client)
   - [Go](https://golang.org/) (for the server)
2. Obtain a valid **Google OAuth Client ID** and **Client Secret**:
   - [Follow this guide](https://developers.google.com/identity/protocols/oauth2) to set up credentials.

### Steps
1. Clone the repository:
   ```bash
   git clone git@github.com:martishin/react-golang-jwt-oauth.git
   cd react-golang-jwt-oauth
   ```
2. Set up environment variables YOUR_GOOGLE_APP_CLIENT_ID and YOUR_GOOGLE_APP_CLIENT_SECRET in the code
3. Start the database:
   ```bash
   docker-compose up db
   ```
4. Start the server:  
   Open a new terminal window and navigate to the `server` directory:
   ```bash
   cd server
   ```
   Build and start the server:
   ```bash
   go run main.go
   ```
   The server should now be running on [http://localhost:8000](http://localhost:8000)

5. Start the client:  
   Navigate to the `client` directory:
   ```bash
   cd client
   ```
   Install the dependencies:
   ```bash
   npm install
   ```
   Start the development server:

   ```bash
   npm run dev
   ```
   The client should now be accessible at [http://localhost:5173](http://localhost:5173)

### API Endpoints
| Method | Endpoint          | Description                    |
|--------|-------------------|--------------------------------|
| POST   | `/api/auth/google` | Google OAuth login             |
| POST   | `/api/auth/logout` | Log out and invalidate session |
| POST   | `/api/auth/refresh`| Refresh access token           |
| GET    | `/api/user`        | Get authenticated user details |

## Tech Stack
- **Frontend**: React.js, Tailwind CSS
- **Backend**: Go, Chi
- **Database**: PostgreSQL
- **Auth Provider**: Google OAuth 2.0
- **Others**: Docker, JWT

## License
This project is licensed under the MIT License.


