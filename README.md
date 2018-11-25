# go-froxlor-dyndns
Application written in Go to target the [HTTP Froxlor Dynamic DNS API](https://github.com/ServiusHack/Froxlor/blob/dynamic-dns/updateip.php).

## Usage
```bash
# simply run it
froxlor-dyndns/ >> cd main && go run main.go -cfg=<path-to-config-file.json>
# build it then run it
froxlor-dyndns/main >> cd main && go build -o fdyndns main.go
froxlor-dyndns/main >> ./fdyndns -cfg=<path-to-config-file.json>
```

Configuration file example
```json
{
  "target": "https://your-froxlor.com", // root url of your froxlor instance
  "interval": 60, // in seconds
  "username": "username", 
  "password": "password", // will be used as fallback when an update does not provide credentials
  "updates": [
    {
      "domains": ["foo.example.com", "bar.example.com"],
      "ipv4": "127.0.0.1",
      "ipv6": "2001:0db8:85a3:0000:0000:8a2e:0370:7334."
    }, // will use the provided IPs and the fallback credentials
    {
        "domains": ["foo.example.com", "bar.example.com"],
        "username": "Bob",
        "password": "pw"
    } // will use the auto detection flag when calling froxlor and the provided credentials
  ]
}
```

## Contribution
Feel free to enhance or improve the modules and create pull-requests.

## License
MBDev is Open Source software released under the [MIT license](https://opensource.org/licenses/MIT).
