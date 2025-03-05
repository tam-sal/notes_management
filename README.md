# Notes@v1  
## Project Overview  
Full-stack application with Go backend and React.js frontend. Includes authentication, note management, and filtering.

## Requirements  
### Backend  
- Go v1.23.0  
- PostgreSQL v16.3  
- Prometheus v3.2.1

### Frontend  
- Node.js v22.13.1  
- npm v10.9.2  

## Running the Project Locally  
1. Clone repository  
2. Navigate to `root-directory` where `run-project.sh` appears.     
3. Run the project `chmod +x run-project.sh && ./run-project.sh`  
4. Fornt End expected local url `http://localhost:5173`  
5. Back End expected local url `http://localhost:8025` 


## Deployed Version  
**Live URL**: [https://note-it-quick.vercel.app/](https://note-it-quick.vercel.app/)  
- Backend: AWS EC2  
- Frontend: Vercel  

## Preset Credentials  
`Username`: superuser  
`Password`: superuser  

## API Documentation  
Swagger: [SWAGGER DOCS](https://backend-tamers.lat/swagger/index.html)  

## Metrics  
**Prometheus**: [PROMETHEUS LIVE MERICS SCRAPPING](http://18.219.246.232:9090/targets)  
**Grafana**: [GRAFANA LIVE METRICS GRAPHING](http://18.219.246.232:3000/)  
**Grafana CREDENTIALS AND DASHBOARD DIRECT ACCESS**:
- Credentials: user: `admin` | password: `notes`  
- Dashboard: [NOTES_APP](http://18.219.246.232:3000/d/aec64bmsh7w8wd/notes-app)  

## Postman Collection  
[NOTES_POSTMAN_COLELCTION](https://tamer0.postman.co/workspace/TAMER-Workspace~fc865e6c-d397-405b-871b-19f94759fb75/collection/38531411-08766622-edfc-4876-a863-789904bf42f7)  

# App Technologies and Features

## Backend
- **Language & Framework:** Golang | Gorilla Mux
- **Database:** PostgreSQL with GORM
- **Architecture:** Follows a layered architecture similar to a Spring Boot project, with Repository, Service, and Handlers (Controllers) layers
- **Authentication:** JWT-based authentication, with tokens sent via cookies
- **Middleware:** Custom middleware for request validation, authentication, and route protection
- **Rate Limiting:** Implements a rate limiter middleware to prevent abuse
- **Logging:** Access logging with Go-coded feature for better tracking
- **Metrics:** Built-in metrics using Go along with Prometheus, both are consumed by Grafana for live graphing
- **Response Structure:** Consistent structured response across the application
- **API Documentation:** Swagger for API documentation

## Frontend
- **Library & Framework:** React with Vite, deployed on Vercel and synced via git push
- **Styling:** Tailwind CSS with DaisyUI preset
- **State Management:** React Context API
- **Notifications:** Toast notifications using `react-hot-toast`
- **Theme & Responsiveness:** Responsive design with dark and light mode, utilizing standard dark/light themes

## Deployment   
- **EC2:** Backend is fully deployed on AWS EC2 t2.micro free tier Linux AMI instance with setup of security group and role.    
- **Elastic IP:** Secured static public IP avioding rotation
- **Containerization:** Dockerized backend and services
- **Orchestration:** Docker Compose for managing multi-container setup and resource consuption of each service
- **Persistence:**
  - PostgreSQL data stored in a named Docker volume (`postgres-data`)
  - Prometheus data persisted in a bind mount (`./backend/deploy/prometheus-data` directory)
  - Grafana data persisted in a bind mount (`./backend/deploy/grafana-data` directory)
- **HTTP to HTTPS Redirection / SSL-TLS Termination:**
  - HTTP to HTTPS redirection by configuring Nginx to terminate SSL/TLS encryption, SSL and DNS are provided by Namecheap [check deployed backend api status here](https://backend-tamers.lat/status)
- **Automation:**
  - EC2 instance has an automated script that runs `docker-compose up -d` on machine restart, ensuring all services restart automatically 


## Documentation and Monitoring  
**Swagger**  

![Image](https://github.com/user-attachments/assets/0e7091ce-d5ec-402a-a3d4-c12c45c39631)

**Prometheus**  

![Image](https://github.com/user-attachments/assets/53bfa5ca-cd4e-42e1-9eb8-ba11c481b640)

**Grafana**  

![Image](https://github.com/user-attachments/assets/b324d59d-ff09-4b20-84f6-fd88693003d1)
![Image](https://github.com/user-attachments/assets/58c260aa-1157-4e0d-a5a6-a5091a960bf5)
![Image](https://github.com/user-attachments/assets/3431c869-4d95-4cf1-acab-ad52427c1899)



## User Story  

1. **Sign Up or Log In**  
   🔐 Create an account or sign in to access your notes.  

2. **Create & Manage Notes**  
   - ✏️ **Add**: Write notes with titles and content.  
   - 📝 **Edit/Delete**: Update or remove notes anytime.  

3. **Organize with Categories**  
   🏷️ Assign labels (e.g., Work, Personal) for quick filtering.  

4. **Archive/Unarchive**  
   📂 Toggle notes between active and archived states.  

5. **Filter Notes**  
   🔍 Search by:  
   - **Status** (Active/Archived)  
   - **Categories**
   - **Status and Categories at the same time**

## Application Flow   
**Log in or Register**  
![Image](https://github.com/user-attachments/assets/e6bd1de6-6f5e-42fe-b346-5373602782eb)

**Create**    
![Image](https://github.com/user-attachments/assets/640031c5-eb60-4336-9541-e86891d133cd)

**List and access notes' functionalities**   
![Image](https://github.com/user-attachments/assets/3b609b7f-10b4-4a42-944e-46ab904e8c2b)   

**Advanced Filters**  
![Image](https://github.com/user-attachments/assets/9215608a-a9af-4828-b005-ded2873cb7f8)
![Image](https://github.com/user-attachments/assets/c823cbf9-5b72-46b7-a631-ab7ab52593b8)


## BACKEND STRUCTURE 
```
backend
├── cmd/
│   └── api/
│       ├── main.go
│       ├── routes.go
│       ├── server.go
│       └── docs/
│           ├── docs.go
│           ├── swagger.json
│           └── swagger.yaml
│
├── deploy/
│   ├── prometheus.yml
│   ├── grafana-config/
│   │   └── grafana.ini
│   ├── grafana-data/          
│   └── prometheus-data/       
│
├── internal/
│   ├── api/
│   │   └── handlers/
│   │       ├── categories.go
│   │       ├── errors.go
│   │       ├── handlers.go
│   │       ├── middlewares.go
│   │       ├── notes.go
│   │       ├── ratelimiter.go
│   │       ├── types.go
│   │       ├── user.go
│   │       └── metrics/
│   │           └── prometheus.go
│   ├── configs/
│   │   └── configs.go
│   ├── models/
│   │   ├── category.go
│   │   ├── note.go
│   │   └── user.go
│   ├── repositories/
│   │   ├── category.go
│   │   ├── interface.go
│   │   ├── note.go
│   │   └── user.go
│   └── services/
│       ├── category.go
│       ├── note.go
│       └── user.go
│
├── pkg/
│   ├── date/
│   │   └── date.go
│   ├── request/
│   │   └── json.go
│   ├── response/
│   │   ├── json.go
│   │   └── metrics.go
│   ├── utils/
│   │   ├── auth.go
│   │   └── helpers.go
│   └── validations/
│       └── errors.go
│
├── scripts/
│   └── aws-scripts/
│       ├── docker-compose-app.service
│       └── info.md
│
├── .env
├── .env.local
├── docker-compose.yml
├── Dockerfile
├── go.mod
├── go.sum
└── package-lock.json
```
## FRONTEND STRUCTURE   
```
frontend
├── public/
│   └── logo.png
│
├── src/
│   ├── App.css
│   ├── App.jsx
│   ├── index.css
│   ├── main.jsx
│   │
│   ├── assets/
│   │   ├── bg-brown-1.png
│   │   ├── bg-brown.png
│   │   ├── bg-transparent.png
│   │   ├── bg.png
│   │   ├── logo.png
│   │   ├── logo_aux.png
│   │   ├── notes-2bg.jpg
│   │   ├── notes-bg.jpg
│   │   └── react.svg
│   │
│   ├── components/
│   │   ├── toggle-theme.jsx
│   │   │
│   │   ├── Fallback/
│   │   │   └── Fallback.jsx
│   │   │
│   │   ├── Footer/
│   │   │   └── Footer.jsx
│   │   │
│   │   ├── GetStarted/
│   │   │   └── Register.jsx
│   │   │
│   │   ├── Header/
│   │   │   ├── Header.jsx
│   │   │   ├── LogoTitle.jsx
│   │   │   └── NavBar.jsx
│   │   │
│   │   ├── Landing.jsx/
│   │   │   └── Landing.jsx
│   │   │
│   │   ├── Loader/
│   │   │   └── Loader.jsx
│   │   │
│   │   ├── Login/
│   │   │   └── Login.jsx
│   │   │
│   │   ├── Note/
│   │   │   ├── CreateNote.jsx
│   │   │   ├── Note.jsx
│   │   │   ├── Notes.jsx
│   │   │   └── UpdateNote.jsx
│   │   │
│   │   └── Protect/
│   │       └── ProtectedRoute.jsx
│   │
│   ├── config/
│   │   └── axios.jsx
│   │
│   └── context/
│       ├── AuthContext.jsx
│       ├── ThemeContext.jsx
│       └── toast-utils.jsx
│
├── .env
├── eslint.config.js
├── index.html
├── package-lock.json
├── package.json
├── postcss.config.js
├── tailwind.config.js
└── vite.config.js
```