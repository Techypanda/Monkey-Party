# /bin/sh
go test ./... -cover 2>&1 | 
  perl -p -e 'if(/coverage: (\d+.\d)%/) { 
    die "coverage too low: $1" if ($1 < 100) 
  }'