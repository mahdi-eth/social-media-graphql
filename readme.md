# Social Media GraphQL

## Overview

This project is a basic social media application built with Go, GraphQL, and MongoDB. It allows users to create accounts, post content, follow other users, and receive real-time updates on new posts from followed users.

## Features

- Create and manage user accounts.
- Post content and view posts from followed users.
- Real-time updates on new posts through GraphQL subscriptions.

## Getting Started

### Prerequisites

- Make sure you have **Docker** installed on your machine.

### Installation

1. **Clone the Repository**

   ```bash
   git clone https://github.com/mahdi-eth/social-media-graphql.git
   cd social-media-graphql
   ```

2. **Build and Run with Docker**

   ```bash
    docker-compose up --build
   ```
    - This command will set up and run MongoDB and the application. You can access the application at http://localhost:8080.

3. **Access GraphQL Playground**
    - Open your browser and go to http://localhost:8080 to interact with the GraphQL API.
