Example Elasticsearch configuration
===================================

Before actually using it in production you will need to edit the content of the mapping.json and change TLD.

### Configure your whoisd.conf

Change the host to whatever hostname elasticsearch is running on, as well as the index and the type names (if different)
```json
{
  "host": "0.0.0.0",
  "port": 43,
  "workers": 1000,
  "connections": 1000,
  "storage": {
    "storageType": "Elasticsearch",
    "host": "localhost",
    "port": 9200,
    "indexBase": "whois",
    "typeTable": "domain"
  },
  "secondaryWhois": "",
  "secondaryRegexp": ""
}
```

### Start the whoisd with -config and -mapping parameters before installing it as service and creating /etc/whoisd, in order to test if everything is working fine.
```sh
cd /path/to/your/whoisd/binary
./whoisd -config=/path/to/your/whoisd.conf -mapping=/path/to/mapping.json
```
_NOTE_: You can test now with:
```sh
whois -h localhost example.com
```
