# watt-flow

## Authors

- Aleksa Perovic SV24/2021
- Danilo Cvijetic SV25/2021
- Vladimir Cornenki SV53/2021

## Init App:

- When initializing app for the first time, in config.env file set RESTART field to True
- After that SuperAdmin account is created
- To be able to activate SuperAdmin account you will have to change account's default password
- That default password is stored in txt file : ./data/admin_password.txt
- To change that password write url: http://localhost:5173/superadmin, there you change default password
- Now you can log in with SuperAdmin account with: username: admin, password: your_new_password
- That is it. Use this steps when running the app for the first time or when you want to restart it for some reason.
- Warning: Set RESTART field in config.env back to false for future.


## Project Description
Watt-Flow is a real-time electricity consumption monitoring and management system designed to handle thousands of household energy meters simultaneously. The platform collects, processes, and visualizes electricity usage data while providing comprehensive billing and administrative capabilities. 

## Core Functionality
The system captures real-time electricity consumption data from thousands of households using IoT simulators that generate continuous power consumption readings. It processes these high-volume data streams through message queues and stores them efficiently for both real-time monitoring and historical analysis.  Users can track their consumption patterns, manage properties, schedule appointments, and receive automated electricity bills.

## Technology Stack

### Backend (Go)
- **Gin** - High-performance HTTP web framework for RESTful API endpoints
- **PostgreSQL** - Relational database for transactional data (users, properties, households)
- **InfluxDB 2** - Time-series database optimized for fast writes and queries of electricity consumption metrics
- **Redis** - In-memory caching layer for rapid data retrieval and reducing database load
- **RabbitMQ** - Message broker for asynchronous processing of device readings
- **GORM** - Object-relational mapping for database operations

### Frontend (Vue. js 3)
- **Vite** - Fast build tool and development server
- **TypeScript** - Type-safe JavaScript for robust frontend development
- **Pinia** - State management
- **Vue Router** - Client-side routing
- **Axios** - HTTP client for API communication
- **Tailwind CSS** - Utility-first CSS framework

### Infrastructure
- **Docker Compose** - Containerized microservices architecture
- **Nginx** - Reverse proxy and static file serving
- **JWT Authentication** - Secure token-based authentication

## Key Advantages
1. **High Performance**: Redis caching and InfluxDB enable fast searches across millions of time-series data points
2. **Scalability**: Microservices architecture with RabbitMQ message queuing handles thousands of concurrent device connections
3. **Real-time Processing**: Asynchronous consumers process electricity readings instantly without blocking
4. **Data Integrity**: PostgreSQL transactions ensure consistent property ownership and billing records
5. **Optimized Queries**: Time-series database (InfluxDB) specifically designed for rapid aggregation of consumption data
6. **Reliable Messaging**: RabbitMQ ensures no data loss during high-volume periods
7. **Efficient Caching**: Redis reduces database load for frequently accessed data like device status

The system demonstrates enterprise-grade architecture with separation of concerns between data ingestion (simulators), message processing (consumers), business logic (Gin server), and presentation (Vue.js frontend).
