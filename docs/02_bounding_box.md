# Phase 2 — Bounding Box Optimization

### Algorithm
Instead of calculating the heavy Haversine distance for *every* user, we first draw a "bounding box" (a square) around the radius circle. We strictly filter out users who do not fall inside this bounding box using simple math, and only calculate the Haversine distance for users who are inside the box. 

### Pseudo Code

```pseudo
minLat, maxLat, minLon, maxLon = boundingBox(myLat, myLon, radius)

for each user:
    # Cheap check: Is user outside the square?
    if user.Lat < minLat or user.Lat > maxLat or user.Lon < minLon or user.Lon > maxLon:
        continue # Skip this user
        
    # Expensive check: Is user inside the circle?
    distance = haversine(...)
    if distance <= radius:
        add user
```

### Complexity
`O(N)`
Wait, it's still O(N)!

### Why ?
Because we are still iterating over **all** elements (`for each user`) in the array to check if their properties match the bounds.

### Improvement
Although the time complexity remains O(N), we have drastically reduced the constant overhead. `Haversine` is mathematically expensive (uses sine, cosine, arctangent), whereas a bounding box involves basic conditionals and additions. 

If a large chunk of your subset is outside the box, the bounding box safely removes the overhead of complex mathematical functions over thousands of items.

### Conclusion
A Bounding Box is an **optimization strategy**, not a scaling strategy. It speeds up the CPU math significantly but does not prevent the full table scan bottleneck.

### Visualization Workflow
We check if our math is faithful by viewing it on `http://localhost:8080/`.

```
User
  ↓
API
  ↓
Filter users inside Bounding Box (cheap conditional)
  ↓
Haversine calculation on remaining (expensive math)
```
