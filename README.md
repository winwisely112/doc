

The project is designed in 3 high level layers.

1. Business GUI

- This is the applied Vertical domain for which the architecture is applied.

- It uses IOC in order to use the Flutter components from the Embed layer

2. Embed

- This is reusable Flutter and Golang code that is client side.

- The embed repo is the golang that sits embedded behind the Flutter code. 

	- It does:

		- Security

		- Storage

		- Networking

- The packages repo is the Flutter

	- It does:

		- Gsuite style functionality.

		- i18n embedded via Flutter 

3. Backend

This is the reusable Golang backend that supports the Embed layer.



This is our Best Practice Architecture that meets our goals.

Packages ( See repo ) use embed and the Architecture is provides.

Harness ( this repo - see Folder) is our Best practice app to show how it works and refine it.

Code written in golang with Flutter.

Why ?

- Flutter is great at GUI. (MANY Foreground)
- Golang is great at networking. (Singleton Background)


## Run as a Background Service

We are moving to a Foreground and Background process model and so embedding is not needed.

Why ? 

- UX Ergonomics

	- Apps can be segregated to all the UX multi tasking experience to be high ergonomics.

	- Allows the GUI to not stutter because of the multi threading.

- Security

	- Having the data and networking in a separate process provides more security at development time and runtime. 

	- SSO (single Sign on) is centralised to one place.

	- Multi-tenant can be done in one place independent of the OS.

- Notifications

	- Background is always running ( OS ) and so notification wake ups can be do not need the Google and Apple notification gateways and leak data.

- Data synchronisation

	- Runs in the background and so UX on app started up is instant.
  
	- Allows offline editing to work much better because the data is less old.

- Alignment

	- Mobile and Desktop have same Topology.
	
	- Web just bypasses the local golang background service and uses the Cloud Server that acts on the users behalf.

	- The same Services can run on the devices as in the cloud such as:
		- Search indexer
		- KV Store
		- VPN 
		- and other goodies that Flutter cant do.

- Interoperability

	- Others can build Apps in any languages they want.


## Architecture:

- Foreground runs the Flutter Apps
- Background singleton service runs the golang code for all networking.
- IPC options

	- GRPC based IO between them OS IPC mechanism.
		- Desktop will be easy
  		- Mobile not so easy
		- Protocol buffers not GRPC

  	- Or GRPC uses a GRPC socket and so we do not have to get into the weeds of IPC.
    	- Often used approach
    	- Will be much secure if we use a Cert between Foreground and Background. Can be bootstrapped from the Server on startup.


## Dev aspects


- The background process is best thought of as a PUB SUB topic based Network.

- You are essentially working at the Domain Model layer and so your models do not have to match one to one to the PUB SUB topics.

- You can remap them how you want. Its gives the Flutter app the ability to compose their Domain Models to not be a one for one match as the PUB SUB topics in the network.

- When a record changes in the Background process it will just tell the Flutter foreground app over the GRPC stream.

- We can set it up so the common Flutter code just puts it directly into HIVE, and so you just forget about it.
- Pull Request
