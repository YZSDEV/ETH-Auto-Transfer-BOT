Clone this Repository

```
git clone https://github.com/YZSDEV/ETH-Auto-Transfer-BOT.git

```
Download Go: Run the following command to download the latest Go version:

```
wget https://go.dev/dl/go1.20.7.linux-amd64.tar.gz

```
Extract the tarball:
```
sudo tar -xvf go1.20.7.linux-amd64.tar.gz -C /usr/local
```
Set up the Go environment: Add Go to your $PATH by editing ~/.profile or ~/.bashrc

```
export PATH=$PATH:/usr/local/go/bin
```
Reload the profile

```
source ~/.profile
```

Verify installation
```
go version
```

Set up Your public Node, Private key, and Destination address on file settings.env
```
nano settings.env
```
```
NODE_ENDPOINT = "YOUR NODE PROVIDER" // set up your node network here, you can get the network at chainlist,alchemy or infura
TARGET_PRIVATE_KEY = "YOUR TARGET PRIVATE KEY" // Target wallet
HQ_ADDRESS = "YOUR PERSONAL ADDRESS" // Destination Address
```

Run this bot with command
```
go run go.main
```









