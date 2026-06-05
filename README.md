# Memory Infrastructure System

## Overview

This project is a backend memory infrastructure system built in Go. It is designed to store, retrieve, and intelligently search user memories using a combination of relational storage, vector search, caching, and ranking systems.

The system evolves in multiple versions, starting from a simple CRUD backend and progressing toward a production-grade hybrid retrieval engine.

---

## Goals

- Provide persistent memory storage for users
- Support session-based context grouping
- Enable hybrid search (keyword + semantic)
- Implement ranking over retrieved memories
- Support scalable and modular backend design
- Demonstrate production-ready system design principles

---

## Tech Stack

- Go (Backend API)
- PostgreSQL (Relational storage)
- Qdrant (Vector database for embeddings)
- Redis (Caching layer)
- Docker (Containerization)
- Next.js (frontend dashboard)

---

## System Architecture

The system is designed in layered components:

1. API Layer (Go HTTP server)
2. Service Layer (Business logic)
3. Repository Layer (Database access)
4. Storage Layer (PostgreSQL, Qdrant, Redis)
5. Async Worker (Embedding pipeline)
6. Search Engine (Hybrid retrieval + ranking)

---

## Core Features

### Memory Management
- Create memory
- Read memory
- Delete memory
- List memories by user or session

### Session Management
- Create sessions
- Attach memories to sessions
- Retrieve session history
- Context-based grouping of memories

### Hybrid Search System
- Semantic search using vector embeddings (Qdrant)
- Keyword search using PostgreSQL full-text search
- Merging and deduplication of results

### Ranking Engine
- Similarity score
- Recency score
- Importance score
- Weighted ranking formula

### Embedding Pipeline
- Async worker-based processing
- Generates embeddings on memory creation
- Stores vectors in Qdrant

### Caching Layer
- Redis-based caching for:
  - search results
  - session data
  - frequent queries

### Authentication
- API key-based authentication
- Middleware-based request validation
- User-scoped data access

---

## Retrieval Flow

When a search request is made:

1. Check Redis cache
2. Query vector database (Qdrant)
3. Query PostgreSQL full-text search
4. Merge results
5. Apply ranking engine
6. Return top-K results

---

## Project Versions

### Version 1 — Core Backend System
Focus: Basic functionality

- Memory CRUD
- Session system
- PostgreSQL integration
- REST API
- API key authentication

Outcome:
A fully functional backend capable of storing and retrieving structured memories.

---

### Version 2 — Search Engine Layer
Focus: Intelligence layer

- Embedding pipeline
- Qdrant integration
- Semantic search
- Keyword search (PostgreSQL FTS)
- Hybrid search system

Outcome:
The system becomes a retrieval engine rather than just a database.

---

### Version 3 — Ranking + System Design
Focus: Intelligence and optimization

- Ranking engine (similarity, recency, importance)
- Redis caching layer
- Session-aware boosting
- Performance improvements

Outcome:
A context-aware retrieval system with optimized response quality.

---

### Version 4 — Production Readiness
Focus: Scalability and observability

- Async job improvements
- Rate limiting
- Logging and latency tracking
- Search tracing
- Optional Next.js dashboard

Outcome:
A production-like distributed system with observability and monitoring.

---

## API Endpoints

### Memories
- POST /memories
- GET /memories
- DELETE /memories/:id

### Sessions
- POST /sessions
- GET /sessions
- GET /sessions/:id
- GET /sessions/:id/memories

### Search
- POST /search

---

## Future Improvements

- Distributed worker system (Kafka or NATS)
- Advanced ranking models (ML-based scoring)
- Multi-tenant architecture
- Real-time streaming memory ingestion
- Analytics dashboard

---

## Summary

This project demonstrates how to evolve a simple CRUD backend into a full hybrid retrieval system using modern backend engineering practices including vector search, caching, async processing, and ranking systems.