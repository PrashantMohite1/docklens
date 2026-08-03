

from blobs/sha256/ we get this 
custom-alpine/blobs/sha256/30f0a0905a2f95c55b83fff834b142f1b86a2d687618ffa44650b15fafae1c04: JSON text data
custom-alpine/blobs/sha256/55afa1ecc21d2bb5e5045f32dafee56272ffd89860bac26f6c32123439af26a4: gzip compressed data, original size modulo 2^32 8697856
custom-alpine/blobs/sha256/6f667f7929b17e7e002e38ccf49d68266add00f9a41723518b9e7a06f6487166: JSON text data
custom-alpine/blobs/sha256/94bd1145481b11b4fe4a3c627e30bc2b39fe7f6ae89aa4a42a1f7277a3af5265: JSON text data
custom-alpine/blobs/sha256/a9fe4aefa30767a7123f6887174198afd93836f71a8072cb0b1a156f5600143b: gzip compressed data, original size modulo 2^32 2560
custom-alpine/blobs/sha256/b2b0c34e332897f3b5455566da7f2af89280ea18db2078ffabb6744fcd1f4d1e: JSON text data
custom-alpine/blobs/sha256/c881c5b85894c69ec969b1d4adc47f9f20b402eb96a0d453a2f2223d8d7cdd68: JSON text data
custom-alpine/blobs/sha256/d10a1afecd3712c1f0b2027cf0e76d424ec43c6d7a6a802df70b96517a5cc8eb: gzip compressed data, original size modulo 2^32 1536
custom-alpine/blobs/sha256/ef48ec11c41b87684eb9d434f6107da03d671ef9fb6d588253d08a832705f1f8: JSON text data


gzip -dc custom-alpine/blobs/sha256/55afa1ecc21d2bb5e5045f32dafee56272ffd89860bac26f6c32123439af26a4 | sha256sum
34884abbe92863fce933ed7c39c0e045631af0ed86d5cc0dfbdf9fdca426ce3c *-

docker image inspect custom-alpine --format '{{.RootFS.Layers}}'
[sha256:34884abbe92863fce933ed7c39c0e045631af0ed86d5cc0dfbdf9fdca426ce3c sha256:9efc4a8a5fb13d7513d949e0a342f2a8c7864504ac4d9f1e097bb02f588316f3 sha256:41d4f1bce6c69e0724dea1f4cbf2ba48f494f8a8211c5efac3ae0a5e82dd0377]

so in a nut shell docker layers gives us  uncompressed (decompressed) layer contents. 
which means if unzip 



#### does file exist in image or not  ( THIS OPTION WILL NOT WORK FIND DIFFERENT )


docklens verify \
    --image myimage:latest \
    --local ./myapp \
    --path /usr/local/bin/myapp

A) docker image verify image-name command-to-run
B) docker image verify --local file.txt image-name:/app/file.txt


## file hash Check 

run sha256sum on host 
sha256sum /app/temp.txt

run sha256sum inside the container of image
docker run --rm img-name sha256sum /app/temp.txt 

## folder hash check 

on host 
```
find releaseDir -type f -exec sha256sum {} + \
| awk '{print $1}' \
| sort \
| sha256sum
```

on docker 
```
docker run --rm alpine-dir \
sh -c 'find releaseDir -type f -exec sha256sum {} + | awk "{print \$1}" | sort | sha256sum'
```


