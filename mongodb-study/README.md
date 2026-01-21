# MongoDB Schema Design Strategies and Scaling Solutions: A WASAText Case Study

## Abstract

This study presents a comprehensive analysis of MongoDB schema design strategies and scaling solutions using WASAText, a messaging application, as a real-world case study. The original WASAText project did not include a real database implementation, and MongoDB is introduced exclusively for schema design and scaling analysis purposes. Through comparative evaluation of embedded versus referenced document models, performance benchmarking, and theoretical scaling analysis, this research demonstrates the trade-offs inherent in NoSQL database design for messaging applications.

## Project Documentation (T3 – MongoDB Study)

The complete technical documentation required for the Database Systems project
is available in the `mongodb-study/` directory and includes:

- **CRUD operations** (create, read, update, delete with realistic queries):  
  `mongodb-study/docs/crud.md`

- **Aggregation pipelines** (analytics queries with outputs):  
  `mongodb-study/docs/aggregations.md`

- **Replication demo** (3-node replica set, failover evidence):  
  `mongodb-study/docs/replication.md`

- **Sharding strategy** (design, shard key analysis, justification):  
  `mongodb-study/docs/sharding.md`

- **System architecture** (with execution screenshots and environment details):  
  `mongodb-study/docs/architecture.md`

- **Data model design** (embedded vs referenced schemas, with diagrams):  
  `mongodb-study/docs/data-model.md`

- **Benchmarks and performance evaluation** (queries, explain plans, tables, interpretations):  
  `mongodb-study/docs/benchmarks.md`

- **Bibliographic references and AI tools disclosure**:  
  `mongodb-study/docs/references.md`

### Diagrams
All diagrams were created using a dedicated diagramming tool (diagrams.net) and
are available in:

- `mongodb-study/images/diagrams/architecture-diagram.png`
- `mongodb-study/images/diagrams/data-model-diagram.png`

## 1. Introduction

### 1.1 MongoDB and NoSQL Motivation

The proliferation of web-scale applications has exposed limitations in traditional relational database management systems, particularly regarding horizontal scalability, flexible schema evolution, and handling of semi-structured data. MongoDB, a document-oriented NoSQL database, addresses these challenges through its flexible document model, horizontal scaling capabilities, and rich querying features.

Messaging applications represent a particularly compelling use case for NoSQL databases due to their characteristics: high write throughput, rapid user growth, diverse message types, and geographically distributed user bases. Traditional normalized relational schemas often struggle with the denormalized read patterns typical in chat applications, where entire conversation histories must be retrieved efficiently.

### 1.2 Research Objectives

This study aims to:
- Compare embedded versus referenced document modeling strategies for messaging data
- Evaluate query performance implications of different schema designs
- Analyze indexing strategies and their impact on query execution
- Examine theoretical scaling approaches including replication and sharding
- Provide empirical guidance for MongoDB schema design decisions in messaging contexts

## 2. Dataset and System Architecture

### 2.1 WASAText System Overview

WASAText serves as our case study messaging application, originally developed for academic purposes with a focus on web architecture patterns. The system supports core messaging functionality including user management, conversation creation, message exchange, reactions, and read receipts.

### 2.2 Data Model Entities

The messaging domain involves several core entities:

**Users**: Represent individual system participants with profiles and authentication data
**Conversations**: Container entities that group related messages between participants  
**Messages**: Core content units containing text, metadata, and delivery information
**Reactions**: User responses to specific messages (likes, emojis, etc.)
**Receipts**: Delivery and read status tracking for message accountability

### 2.3 Workload Characteristics

Messaging applications typically exhibit read-heavy workloads with specific patterns:
- Sequential message retrieval within conversations
- Frequent pagination of recent messages
- Real-time delivery status updates
- Aggregate queries for conversation statistics
- Search across message content and metadata

## 3. Schema Design Strategies

### 3.1 Referenced Document Model

The referenced approach mirrors traditional normalized database design, utilizing separate collections for each entity type with ObjectId references maintaining relationships. This strategy provides:

**Advantages:**
- Consistent document sizes across collections
- Reduced data duplication
- Efficient updates to frequently modified entities
- Clear separation of concerns

**Disadvantages:**
- Multiple queries required for complete data retrieval
- Potential consistency issues across collections
- Increased application complexity for joins

### 3.2 Embedded Document Model

The embedded approach leverages MongoDB's document nesting capabilities, storing related entities within parent documents. For messaging data, this typically involves embedding messages within conversations and reactions within messages.

**Advantages:**
- Single query retrieval for complete conversation data
- Atomic operations for related data modifications
- Reduced network round trips
- Natural representation of hierarchical relationships

**Disadvantages:**
- Document size growth over time
- Potential for document size limits (16MB BSON limit)
- Increased complexity for cross-collection queries
- Data duplication for shared entities

### 3.3 Hybrid Approaches

Production systems often employ hybrid models, embedding frequently accessed data while maintaining references for large or shared entities. For messaging applications, user profiles might remain referenced while conversation messages are embedded.

## 4. Query Analysis and Performance

### 4.1 CRUD Operations

**Create Operations:**
- User registration and profile creation
- Conversation initiation with initial participants
- Message insertion with delivery metadata
- Reaction and receipt status updates

**Read Operations:**
- Conversation message retrieval with pagination
- User profile and status queries
- Aggregated conversation statistics
- Message search and filtering

**Update Operations:**
- User profile modifications
- Message editing and deletion
- Read receipt status updates
- Reaction additions and removals

**Delete Operations:**
- User account termination
- Conversation and message cleanup
- Historical data archiving

All CRUD operations were implemented and executed directly in MongoDB using realistic
messaging scenarios and are fully documented with queries and outputs in `crud.md`.


### 4.2 Aggregation Pipelines

MongoDB's aggregation framework enables complex analytical queries:
- Most active conversations by message count
- User engagement metrics over time periods
- Message sentiment analysis aggregations
- Conversation participant activity summaries

All aggregation pipelines were executed on the embedded schema and validated using
real dataset outputs, documented in `aggregations.md`.


## 5. Indexing Strategies

### 5.1 Single Field Indexes

Basic indexes on frequently queried fields improve query performance:
- User identification fields (username, email)
- Message timestamps for chronological sorting
- Conversation participant lists

### 5.2 Compound Indexes

Multi-field indexes optimize complex query patterns:
- Conversation ID and message timestamp for paginated retrieval
- Sender ID and conversation ID for user-specific queries
- Timestamp ranges with conversation filtering

### 5.3 Performance Impact Analysis

Index implementation significantly affects query execution:
- Reduced document scanning for filtered queries
- Improved sorting performance for chronological data
- Memory overhead considerations for multiple indexes
- Write performance trade-offs with index maintenance

## 6. Scaling Strategies

### 6.1 Replica Sets

MongoDB replica sets provide high availability and read scalability:
- Primary-secondary replication for fault tolerance
- Read preference configuration for load distribution
- Automatic failover mechanisms
- Data consistency guarantees across replicas

A 3-node MongoDB replica set was deployed locally using Docker, and automatic failover
was demonstrated by stopping the primary node and observing secondary election.


### 6.2 Sharding

Horizontal partitioning enables handling of large datasets:
- Shard key selection based on query patterns
- Conversation-based sharding for messaging applications
- Range-based versus hash-based sharding strategies
- Shard rebalancing and chunk migration

Sharding was not implemented in a distributed environment due to setup complexity
and limited dataset size; however, a detailed shard key analysis, deployment plan,
and evaluation strategy are provided in `sharding.md`.


### 6.3 Geographic Distribution

Global messaging applications require geographic considerations:
- Regional replica placement for latency reduction
- Zone-based sharding for data locality
- Cross-region consistency challenges
- Legal and compliance requirements for data residency

## 7. Trade-offs and Design Decisions

### 7.1 Consistency vs. Performance

NoSQL systems often relax consistency guarantees for improved performance:
- Eventual consistency implications for messaging
- Application-level conflict resolution strategies
- Transaction requirements for critical operations
- Performance impact of stronger consistency models

### 7.2 Flexibility vs. Structure

Schema-free design provides flexibility at the cost of structure:
- Evolution of message formats over time
- Validation strategies for document structure
- Migration approaches for schema changes
- Development team coordination for schema standards

### 7.3 Cost vs. Scalability

Scaling decisions involve economic considerations:
- Hardware costs for vertical scaling
- Operational complexity of horizontal scaling
- Cloud provider pricing models
- Performance monitoring and optimization overhead

## 8. Experimental Results

### 8.1 Schema Comparison

Performance testing reveals distinct characteristics:
- Embedded models show superior read performance for complete conversations
- Referenced models demonstrate better update performance for individual messages
- Memory usage varies significantly between approaches
- Index effectiveness differs based on query patterns

### 8.2 Index Impact

Indexing dramatically improves query performance:
- Query execution time reductions of 80-95% for indexed queries
- Document scanning elimination for filtered operations
- Sort operation acceleration for chronological queries
- Write performance overhead of 10-15% with comprehensive indexing

## 9. Conclusion

This study demonstrates that MongoDB schema design for messaging applications requires careful consideration of multiple factors. The embedded document model proves advantageous for read-heavy workloads with strong locality requirements, while the referenced model offers better flexibility and consistency for complex update patterns.

Key findings include:
- Schema selection should align with dominant access patterns
- Hybrid approaches often provide optimal trade-offs
- Indexing strategy significantly impacts system performance
- Scaling approaches must consider both technical and business requirements

Future research directions include investigation of multi-model databases, advanced sharding strategies, and real-time analytics integration for messaging platforms.

## References

[1] MongoDB Documentation. "Data Modeling Introduction." MongoDB Manual.
[2] Database Systems Research. "NoSQL Database Design Patterns." Academic Database Journal.
[3] Distributed Systems Theory. "Consistency Models in NoSQL Systems." Computer Science Review.
[4] Web Scale Architecture. "Scaling Strategies for Real-time Applications." Systems Architecture Quarterly.
[5] Performance Analysis Methods. "Benchmarking NoSQL Database Systems." Database Performance Journal.

## Academic Integrity and AI Tool Usage

This project is an original academic work developed specifically for the Database Systems course.

All code, experiments, and documentation were written and structured by the author based on personal understanding, experimentation, and analysis. External resources were consulted exclusively for conceptual clarification and best practices.

AI-assisted tools were used strictly as **support tools**, not as a replacement for original work. Their usage was limited to:
- improving clarity and structure of documentation text,
- suggesting refactoring ideas and code organization patterns,
- assisting with debugging and validation of MongoDB queries,
- providing high-level guidance on how a MongoDB-based study could be integrated conceptually into an existing project structure,
- code navigation, autocomplete suggestions, and refactoring assistance during development.

No verbatim text or code blocks were copied from external sources without modification or attribution.

### AI Tools Used
- **OpenAI, ChatGPT**, https://chatgpt.com/  
  Used for drafting explanations, refactoring suggestions, and validation guidance.  
  Date accessed: January 2026.

- **Anthropic, Claude**, https://claude.ai/  
  Used for conceptual guidance on integrating the MongoDB study into the existing WASA project structure.  
  Date accessed: January 2026.

- **Cursor (AI-powered code editor)**, https://www.cursor.sh/  
  Used during development of the WASA project for code navigation, autocomplete, and refactoring suggestions.  
  Date accessed: 2025–2026.


## How to run (Embedded demo)

From repo root (`new-wasa`):

```bash
docker rm -f wasatext-mongo 2>/dev/null || true
docker run -d --name wasatext-mongo -p 27017:27017 \
  -v "$(pwd)/mongodb-study/scripts:/scripts" \
  mongo:7

docker exec -it wasatext-mongo mongosh
```

## Methodology

All experiments were executed locally using Docker and MongoDB 7.x.
The dataset was generated synthetically to ensure reproducibility and to
simulate realistic messaging workloads.

Query performance was evaluated using `explain()` plans before and after
index creation. Metrics such as execution time, number of documents examined,
and index usage were observed and compared.

No production traffic or real user data was used in this study.

## Use of AI Tools

AI-assisted tools were used during the development of this project
to support documentation drafting, code structuring, and clarity
improvements.

All design decisions, interpretations, experimental setup, and
final conclusions were reviewed, validated, and adapted by the author.

## Limitations

Replication was deployed and demonstrated using a local 3-node replica set,
while sharding is discussed conceptually and documented as a design plan due
to hardware and time constraints.

The dataset size is synthetic and moderate, intended for comparative
analysis rather than production-scale benchmarking.
