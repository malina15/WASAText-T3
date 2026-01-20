# WASAText Dataset Description

## Overview

The WASAText dataset represents a comprehensive messaging application data model designed to capture the essential entities and relationships typical of modern chat platforms. This synthetic dataset serves as the foundation for our MongoDB schema design and performance analysis study.

## Core Entities

### Users
- **Purpose**: Represent individual system participants
- **Key Attributes**: 
  - Unique identifier (ObjectId)
  - Username (unique, indexed)
  - Email address (unique, indexed) 
  - Display name
  - Profile avatar URL
  - Account creation timestamp
  - Last activity timestamp
  - Online status indicator
- **Estimated Volume**: 10,000 users
- **Growth Pattern**: Steady user registration with periodic peaks

### Conversations
- **Purpose**: Container entities grouping related messages
- **Key Attributes**:
  - Unique identifier (ObjectId)
  - Conversation type (direct, group)
  - Participant list (user references)
  - Conversation title (for group chats)
  - Creation timestamp
  - Last message timestamp
  - Message count
- **Estimated Volume**: 2,500 conversations
- **Distribution**: 70% direct messages, 30% group conversations

### Messages
- **Purpose**: Core content units representing individual communications
- **Key Attributes**:
  - Unique identifier (ObjectId)
  - Conversation reference
  - Sender reference
  - Message content (text)
  - Message type (text, image, file, system)
  - Timestamp
  - Edit history
  - Reply-to reference (for threaded messages)
- **Estimated Volume**: 100,000 messages
- **Distribution**: Average 40 messages per conversation, heavy tail distribution

### Reactions
- **Purpose**: User emotional responses to specific messages
- **Key Attributes**:
  - Unique identifier (ObjectId)
  - Message reference
  - User reference
  - Reaction type (like, love, laugh, angry, sad)
  - Timestamp
- **Estimated Volume**: 25,000 reactions
- **Pattern**: 25% message reaction rate, clustered around popular messages

### Receipts
- **Purpose**: Track message delivery and read status
- **Key Attributes**:
  - Unique identifier (ObjectId)
  - Message reference
  - User reference
  - Status type (delivered, read)
  - Timestamp
- **Estimated Volume**: 200,000 receipts
- **Pattern**: Near 100% delivery rate, 85% read rate

## Entity Relationships

### Primary Relationships
- **User ↔ Conversation**: Many-to-many through participant lists
- **Conversation ↔ Message**: One-to-many, conversation contains multiple messages
- **User ↔ Message**: Many-to-many, users send messages and can be recipients
- **Message ↔ Reaction**: One-to-many, messages can have multiple reactions
- **Message ↔ Receipt**: One-to-many, messages have delivery receipts per recipient

### Referential Integrity Considerations
- Orphaned messages when conversations are deleted
- User reference cleanup when accounts are terminated  
- Cascade deletion policies for conversation removal
- Historical data preservation requirements

## Dataset Size Estimations

### Storage Requirements
- **Referenced Schema**: Approximately 250MB for complete dataset
- **Embedded Schema**: Approximately 180MB due to reduced overhead
- **Index Overhead**: Additional 50-75MB depending on indexing strategy
- **Growth Rate**: 10-15% monthly increase in message volume

### Memory Requirements
- **Working Set**: 100-150MB for active queries
- **Index Memory**: 25-40MB for comprehensive indexing
- **Connection Overhead**: 10-20MB for typical connection pools

## Workload Characteristics

### Read Patterns
- **Conversation Browsing**: Sequential message retrieval (80% of queries)
- **Search Operations**: Full-text search across message content (10% of queries)
- **User Lookups**: Profile and status queries (5% of queries)
- **Analytics**: Aggregation queries for statistics (5% of queries)

### Write Patterns
- **Message Creation**: Primary write operation (70% of writes)
- **Status Updates**: Receipt and online status modifications (20% of writes)
- **User Management**: Profile updates and authentication (10% of writes)

### Temporal Patterns
- **Peak Hours**: 18:00-23:00 local time zones
- **Daily Cycle**: 5x variance between peak and off-peak periods
- **Weekly Cycle**: 40% higher activity on weekends
- **Seasonal Variations**: Holiday periods show 2-3x increased activity

### Geographic Distribution
- **Primary Regions**: North America (40%), Europe (35%), Asia (25%)
- **Latency Requirements**: Sub-100ms response times for interactive operations
- **Consistency Requirements**: Strong consistency for financial or critical messages

## Data Quality Characteristics

### Completeness
- **Required Fields**: 100% completion for core entity attributes
- **Optional Fields**: 60-80% completion for profile and metadata fields
- **Historical Data**: Complete retention for 2+ years of message history

### Consistency
- **Referential Integrity**: Maintained through application logic
- **Data Validation**: Schema validation for critical fields
- **Duplicate Prevention**: Unique constraints on usernames and emails

### Performance Implications
- **Cache Locality**: Recent conversations accessed frequently
- **Data Skew**: Power-law distribution of user activity
- **Query Optimization**: Index selection based on access patterns

## Testing and Validation

### Synthetic Data Generation
- **Realistic Patterns**: Based on observed messaging application behaviors
- **Edge Cases**: Large conversations, inactive users, bulk operations
- **Performance Testing**: Load testing with representative query patterns

### Data Validation
- **Schema Compliance**: Automated validation of document structure
- **Relationship Integrity**: Verification of cross-collection references
- **Performance Benchmarks**: Baseline measurements for optimization comparison

This dataset description provides the foundation for comprehensive MongoDB schema design analysis, enabling meaningful comparison between embedded and referenced document modeling approaches while maintaining realistic characteristics of production messaging systems.
