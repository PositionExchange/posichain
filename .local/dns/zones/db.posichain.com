$TTL    604800
@       IN      SOA     ns1.posichain.com. root.posichain.com. (
                  3       ; Serial
             604800     ; Refresh
              86400     ; Retry
            2419200     ; Expire
             604800 )   ; Negative Cache TTL
;
; name servers - NS records
     IN      NS      ns1.posichain.com.

; name servers - A records
ns1.posichain.com.                      IN      A      172.189.0.3

s2.z.d.posichain.com.                   IN      A      172.189.0.9
s3.z.d.posichain.com.                   IN      A      172.189.0.10
s0.z.d.posichain.com.                   IN      A      172.189.0.11
s1.z.d.posichain.com.                   IN      A      172.189.0.12
s3.z.d.posichain.com.                   IN      A      172.189.0.13
_dnsaddr.bootstrap.d.posichain.com.     IN      TXT     "dnsaddr=/ip4/172.189.0.8/tcp/9876/p2p/Qmc1V6W7BwX8Ugb42Ti8RnXF1rY5PF7nnZ6bKBryCgi6cv"
