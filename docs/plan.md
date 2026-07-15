

Phase 1 (MVP) : image metadata analysis

These are all based on the Docker Engine API (docker image inspect) and are relatively straightforward to implement.

```
docker-image-analyzer analyze nginx:latest
```

```
Image Name
Image ID
Digest
Created Date
OS
Architecture
Entrypoint
CMD
Working Directory
Exposed Ports
Environment Variables
Labels
Total Size
Virtual Size
Number of Layers
```


Phase 2 : layer analysis
```
docker-image-analyzer layers nginx
```

```
Layer 1  25 MB
Layer 2  10 MB
Layer 3 120 MB  ⚠️ Largest layer
Layer 4   2 MB

Total Layers : 9
Largest Layer : Layer 3
Average Layer Size : 18 MB
```






Phase 4 : Image Score - Identify critical vulnerability

```
docker-image-analyzer score nginx
```

```
Image Score

Overall : 86/100

✔ Small image

✔ Few layers

✔ Uses official image

⚠ Runs as root

⚠ Latest tag used

Suggestions

Use non-root user
Pin image version
Reduce layer count
```