# Groupie Tracker

## Overview
Groupie Tracker is a web application built with Go (Golang) that consumes a public REST API and displays information about music bands and artists. The application retrieves data such as artist details, concert locations, and tour dates, then presents it in a user-friendly website with interactive navigation and search functionality.

The goal of the project is to practice:  
- Client–server architecture  
- JSON data manipulation  
- REST API consumption  
- HTML templating in Go  
- Event handling  
- Error handling  
- Data visualization in a web interface  

---

## Features

### Artist Directory
Displays a list of artists with:  
- Artist image  
- Band name  
- Creation year  
- First album release  
- Band members  

Each artist is displayed using interactive cards.

### Artist Details Page
Clicking on an artist card opens a detailed page that shows:  
- Artist image  
- Formation year  
- First album  
- Band members  
- Tour locations  
- Tour dates  

### Search Functionality
Users can search artists by:  
- Band name  
- Band member name  

The search filters results dynamically.

### Error Handling
The server gracefully handles errors and displays custom error pages for:  
- 404 – Page not found  
- 405 – Method not allowed  
- 400 – Bad request  
- 500 – Internal server error  

The application is designed so the server never crashes.

### API Data Integration
The application retrieves and combines data from the following API endpoints:  

| Endpoint    | Description                  |
|------------|-------------------------------|
| artists    | Artist information           |
| locations  | Concert locations            |
| dates      | Concert dates                |
| relation   | Links artists with dates and locations |

### Data Caching
To improve performance, the application caches API data after the first request, avoiding repeated API calls.

---

## Technologies Used

**Backend:**  
- Go (Golang)  
- `net/http`  
- `html/template`  
- `encoding/json`  

**Frontend:**  
- HTML5  
- CSS3  

**Other:**  
- REST API  
- JSON  

---

## Project Structure

groupie-tracker
│
├── backend
│ ├── services.go # Fetches data from API
│ ├── handlers.go # HTTP handlers
│ └── models.go # Data structures
│
├── frontend
│ ├── index.html # Home page
│ ├── details.html # Artist details page
│ ├── error.html # Error page
│ │
│ └── static
│ └── styles.css # Website styling
│
├── main.go # Server entry point
└── README.md


---

## How It Works

### 1. Server Startup
Start the application:

```bash
go run main.go

http://localhost:8080

### 2. Request Flow
Browser Request
      ↓
Go HTTP Server
      ↓
Handlers (Home / Details)
      ↓
API Requests
      ↓
Data Processing
      ↓
HTML Templates
      ↓
Response to Browser
### 3. Template Rendering

Go uses HTML templates to generate dynamic pages. Example:

{{ .Name }}

This placeholder is replaced with artist data received from the server.

## Installation & Usage

Clone the Repository:

git clone https://learn.zone01kisumu.ke/git/evajuma/groupie-tracker

Navigate to the Project:

cd groupie-tracker

Run the Server:

go run main.go

Open the Website

Visit:

http://localhost:8080
Example Features
Home Page

Displays all artists in a card layout. Users can:

Browse artists

Search by name or member

Click cards to view more details

Artist Page

Displays:

Artist image

Members

Album information

Concert locations and dates

Learning Objectives

This project helps develop understanding of:

Go web servers

REST API integration

JSON decoding

Template rendering

Client-server communication

Error handling

Search filtering

Data caching

## Authors

Clare Gisore

Ibrahim Samwel

Evans Juma
