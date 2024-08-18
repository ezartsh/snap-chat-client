## Snap Chat CLI
> This is a simple chat application using golang in command line interfaces.

Application are tested using Golang 1.22 and PostgreSQL 16.

## Clone the application 
```sh
git clone https://github.com/ezartsh/snap-chat-client.git
cd snap-chat-client
```

### Installation
```sh
# No configuration needed. Just run this.
docker build --tag snap-chat .
```

### Manual
```sh
# After clone the repo and move inside the project folder.
go mod tidy
go mod download
# No configuration needed. Just run this.
go run main.go localhost:3000
```

### Start The CLI
```sh
# The only things that need to adjust is the server host address.
# For this example : localhost:3000
# Change that based on the server exposed address
docker run --rm -it --network="host" snap-chat localhost:3000
```

After start the cli application, you will be prompted to choose between login, register or exit the application.
You need to loggedin an account to use the chat app.
```sh
? Login if you already have an account, or register for new one.  [Use arrows to move, type to filter]
> Login # Will prompt an input username and password
  Register # Will prompt an input to fill name, username and password
  Exit # Exit the application

# Navigation [Up] - Arrow Up / k
# Navigation [Down] - Arrow Down / j
# Navigation [Submit] - Enter
```

#### Chat Editor
```sh
# After logged in, you will be shown this prompt layout.

You are Loggedin !
                 
Chat successfully connected.
                     
Welcome to Snap chat.

[15:30] » [Lobby] [Your Name] # type here
```

This cli comes with pre-built command to help you navigate through the feature.
> ***Lobby*** are an indicator that showed you're not inside the room. You can type anything other than the pre-built command, but it wont send anything to the server. You should use /dm (direct message) or /gm (group message) to send chat from lobby.
```sh
# This is example send message from the lobby
[15:30] » [Lobby] [Your Name] /dm <username> <message> # Direct message 
[15:30] » [Lobby] [Your Name] /gm <group_key> <message> # Group message 
```
Or you can enter the room :
```sh
# This is example send message from the lobby
[15:30] » [Lobby] [Your Name] /dm timoty # Enter the private room
[15:30] » [Private] [Your Name] Hi, timoty..How are you ?
# Exit will move you back to the lobby
[15:30] » [Private] [Your Name] /exit

[15:30] » [Lobby] [Your Name] /gm friends # Enter Group room
[15:30] » [Friends] [Your Name] Hi everyone.
```

You have to know the username of the account of the one you want to private chat with, or the group key to chat to all members in the group.
Here all the rest of pre-built command availables :
```sh
... [Your Name] /contact add <username> # add user to your contact
... [Your Name] /contact list # list all user in your contact
... [Your Name] /contact remove <username> # remove user from your contact
... [Your Name] /group list # list all available groups
... [Your Name] /group create <group_name> # create group
... [Your Name] /group join <group_key> # join group
... [Your Name] /group leave <group_key> # leave group
```
