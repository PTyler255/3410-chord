# 3410-chord

Chord is a peer to peer distributed hash table that is self stabilizing and manages large swaths of key/value stores across multiple devices. Computers dropping in and out of the "ring" automatically updates where the data is stored in order to maintain where everything is located according to the key's hash. A ring can be created by one node and all others can drop in and communicate via RPC calls. Each device has its own command line that can manage the data distributed across the table.

![Example of Chord usage between three users](https://media.discordapp.net/attachments/352962574029029377/1103779832141185074/image.png?width=881&height=593)

**chord.go**

Main file that manages commands for chord program

**commands.go**

File maintaining all commands to be interacted with by the user.

- help

Display commands

- port <n>

Change Port to listen in on
  
- create

Create a new Distributed Has Table

- ping
  
Contact another Chord Node
  
- getip
  
Display local user address
  
- join <address>
  
Join a pre-existing Ring at the designated address
  
- put <key> <value>
  
Store a key/value pair on the ring
  
- putrandom <n>
  
*Not implemented*

- get <key>

Retrieve associated value from the key

- delete <key>
 
Delete key/value pair
  
- dump
  
Dump all data stored on local machine including neighboring nodes and all key/value pairs stored locally
  
- dumpkey <key>
- dumpaddr <address>
- dumpall
  
*Not implemented*
  
- quit
  
Transfer key/values to neighboring nodes and log off

**finger.go**

File managing all RPC based functions to be called by the commands.
