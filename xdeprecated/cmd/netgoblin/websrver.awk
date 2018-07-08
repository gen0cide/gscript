#!/usr/local/bin/gawk -f
BEGIN {
if (ARGC < 2) { print "Usage: wwwawk  file.html"; exit 0 }
  Concnt = 1;
        while (1) {
        RS = ORS = "\r\n";
        HttpService = "/inet/tcp/8080/0/0";
        getline Dat < ARGV[1];
        Datlen = length(Dat) + length(ORS);
        while (HttpService |& getline ){
    if (ERRNO) { print "Connection error: " ERRNO; exit 1}
                print "client: " $0;
                if ( length($0) < 1 ) break;
        }
        print "HTTP/1.1 200 OK"             |& HttpService;
        print "Content-Type: text/html"     |& HttpService;
        print "Server: wwwawk/1.0"          |& HttpService;
        print "Connection: close"           |& HttpService;
        print "Content-Length: " Datlen ORS |& HttpService;
        print Dat                           |& HttpService;
        close(HttpService);
        print "OK: served file " ARGV[1] ", count " Concnt;
        Concnt++;
      }
} 
