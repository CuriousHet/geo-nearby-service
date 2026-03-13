# Geo Nearby Service

Geo-spatial search service demonstrating how large-scale systems like Uber, Swiggy, and Pokemon Go implement proximity-based user discovery.

## 🚀 Quick Start

Run the entire service with a single command:

```bash
go run cmd/server/main.go
```

Once running, access the **Interactive Debug Console** at:
👉 **[http://localhost:8080](http://localhost:8080)**

## 🛠️ Features

- **100,000 Users**: Simulated high-scale user database.
- **Algorithm Sandbox**: Compare Naive, Bounding Box, Geohash, and Google S2 strategies live.
- **Visual S2 Polygons**: High-fidelity visualization of hierarchical spherical indexing cells.
- **Performance Metrics**: Microsecond-accurate latency benchmarks for each strategy.

## 📈 Project Evolution

The service evolved through 4 distinct architectural phases:

1. **Naive Haversine Scan**: Brute-force calculation for all users.
2. **Bounding Box Pre-filter**: Quick `$O(N)$` reduction using rectilinear bounds.
3. **Geohash Indexing**: `$O(1)$` dictionary lookup using 2D spatial strings.
4. **Google S2 hierarchical Indexing**: Advanced spherical geometry indexing with polygonal coverage.

Detailed architectural deep-dives for each phase can be found in the [docs/](docs/) directory.

## 🧪 API Benchmarking

You can also query the API directly to see low-level timing metrics:

```bash
curl "http://localhost:8080/nearby?lat=23.0225&lon=72.5714&radius=50&algorithm=s2"
```