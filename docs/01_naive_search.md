# Phase 1 — Naive Search

### Algorithm
Scan all users and calculate the distance.

### Pseudo Code

```pseudo
for each user:
    distance = haversine(...)
    if distance < radius:
        add user
```

### Complexity
`O(N)`

### Scaling Problem

Assume:
- 100 million users
- 10k requests/sec

### Total Calculations

```
100M × 10k = 1 trillion
```

### Conclusion
Full table scan is **not scalable**.

We must **reduce the search space**.


## Workflow

```
User
  ↓
API
  ↓
Scan all users
  ↓
Haversine calculation
```