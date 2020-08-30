Simple cli tool for communicating with a Modbus server

E.g.
```
go build
./modbus-cli --address=0 --host=localhost --operation=writeSingleRegister --size=1 0
```

e.g. 

set output high on wago 750-5xx 2do

modbus-cli --address=0 --host=192.168.226.42 --operation=writeSingleCoil --size=1 1
