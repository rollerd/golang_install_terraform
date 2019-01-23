## golang_install_terraform

##### What is this:
A simple Golang practice script that does the following:
   - Create .hashicorp directory under user's home dir if it doesn't already exist
   - Downloads Terraform zip file
   - Unzips Terraform zip file
   - Creates .bashrc if it doesn't exist and appends .hashicorp path to PATH
   - Creates .bash_profile if doesn't exist and appends `source ~/.bashrc`

##### Requirements:
   - [Golang](https://golang.org/)

##### How to use:
   - Simply run: `go build install_tf.go` to build the binary
   - and then: `./install_tf`
