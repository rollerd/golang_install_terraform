/* Practice Golang script that does the following:
   - Create .hashicorp directory under user's home dir if it doesn't already exist
   - Downloads Terraform zip file
   - Unzips Terraform zip file
   - Creates .bashrc if it doesn't exist and appends .hashicorp path to PATH
   - Creates .bash_profile if doesn't exist and appends `source ~/.bashrc`
*/
package main

import "fmt"
import "net/http"
import "os"
import "os/user"
import "io"
import "archive/zip"
import "bufio"
import "log"

var terraform_version string = "0.11.11"

func main() {
	usr, _ := user.Current()
	home_dir := usr.HomeDir
	hashi_path := home_dir + "/.hashicorp"

	fmt.Printf("Creating directory: %s\n", hashi_path)
	create_dir(hashi_path)

	fmt.Println("Downloading Terraform")
	terraform_download_url := fmt.Sprintf("https://releases.hashicorp.com/terraform/%s/terraform_%s_darwin_amd64.zip", terraform_version, terraform_version)
	download_file(terraform_download_url, hashi_path+"/terraform.zip")

	fmt.Println("Unzipping Terraform")
	unzip(hashi_path, "terraform.zip")

	fmt.Printf("Adding %s to PATH via .bashrc\n", hashi_path)
	target_string := fmt.Sprintf("export PATH=$PATH:%s", hashi_path)
	append_or_create(home_dir+"/.bashrc", target_string)

	fmt.Println("Adding sourcing for .bashrc to .bash_profile")
	target_string = "source ~/.bashrc"
	append_or_create(home_dir+"/.bash_profile", target_string)
}

func search_string(target_path, target_string string) bool {
	f, err := os.Open(target_path)
	if err != nil {
		fmt.Printf("WARN: %s\n", err)
		return false
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		text := scanner.Text()
		if text == target_string {
			fmt.Printf("Found target string: %s - skipping\n", text)
			return true
		}
	}
	return false
}

func append_or_create(target_path, target_string string) {
	search_result := search_string(target_path, target_string)

	if search_result != true {
		fmt.Printf("Creating file: '%s'\n", target_path)
		f, err := os.OpenFile(target_path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()

		_, err = f.WriteString(target_string)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func unzip(path, zipfile string) {
	zip_file_path := path + "/" + zipfile
	readcloser, err := zip.OpenReader(zip_file_path)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer readcloser.Close()

	for _, zipped_file := range readcloser.File {
		zipped_filename := zipped_file.Name
		write_zip_contents(zipped_file, path+"/"+zipped_filename)
	}
}

func write_zip_contents(src *zip.File, dest string) {
	openfile, err := os.Create(dest)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer openfile.Close()

	err = openfile.Chmod(0755)
	if err != nil {
		fmt.Println(err)
	}

	dest_file, err := src.Open()
	if err != nil {
		fmt.Println(err)
	}
	defer dest_file.Close()

	written, err := io.Copy(openfile, dest_file)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("zip file bytes written: %v\n", written)
}

func create_dir(path string) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.Mkdir(path, 0755)
	}
}

func download_file(url, filename string) {

	filepath, err := os.Create(filename)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer filepath.Close()

	resp, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	written, err := io.Copy(filepath, resp.Body)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Printf("downloaded file bytes written: %v\n", written)
}
