package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

var (
	all    = flag.Bool("a", false, "get all projects")
	output = flag.String("o", ".", "output directory")

	status = 0
)

func main() {
	flag.Usage = usage
	flag.Parse()
	if flag.NArg() < 1 {
		usage()
	}

	user, err := getUser(flag.Arg(0))
	ck(err)

	repos, err := getRepos(user)
	ck(err)

	if flag.NArg() == 1 && !*all {
		fmt.Println(user)
		fmt.Println("Repos:\n")
		for i, r := range repos {
			fmt.Printf("%d %s %q %s", i+1, r.Name, r.Description, r.Url)
			if r.Stargazers_Count != 0 {
				fmt.Printf(" (%d stars) ", r.Stargazers_Count)
			}
			if r.Forks_Count != 0 {
				fmt.Printf("(%d forks) ", r.Forks_Count)
			}
			fmt.Printf("\n")
		}
		return
	}

	os.MkdirAll(*output, 0755)
	ck(os.Chdir(*output))

	if *all {
		for _, r := range repos {
			clup(r.Name, r.Clone_Url)
		}
	} else {
		for _, name := range flag.Args()[1:] {
			for _, r := range repos {
				if strings.ToLower(r.Name) == strings.ToLower(name) {
					clup(r.Name, r.Clone_Url)
					break
				}
			}
		}
	}
	os.Exit(status)
}

func clup(name string, url string) {
	_, err := os.Stat(name)
	if os.IsNotExist(err) {
		execl("git", "clone", url)
	} else {
		execl("git", "-C", name, "pull")
	}
}

func usage() {
	fmt.Fprintln(os.Stderr, "usage: get-github [options] user [projects ...]")
	flag.PrintDefaults()
	os.Exit(2)
}

func ck(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, "get-github:", err)
		os.Exit(1)
	}
}

func ek(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, "get-github:", err)
		status = 1
	}
}

func get(url string, v interface{}) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	buf, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		e := &Error{}
		err := json.Unmarshal(buf, &e)
		if err != nil {
			return err
		}
		return e
	}

	return json.Unmarshal(buf, v)
}

type Error struct {
	Message           string
	Documentation_Url string
}

func (e Error) Error() string {
	return e.Message
}

type User struct {
	Login        string
	Id           uint64
	Url          string
	Repos_Url    string
	Name         string
	Public_Repos uint64
	Public_Gists uint64
	Followers    uint64
	Following    uint64
	Type         string
	Site_Admin   bool
	Hireable     bool
	Bio          string
	Created_At   string
	Updated_At   string
}

func (u *User) String() string {
	w := new(bytes.Buffer)
	fmt.Fprintf(w, "Login:         %v\n", u.Login)
	fmt.Fprintf(w, "ID:            %v\n", u.Id)
	fmt.Fprintf(w, "URL:           %v\n", u.Url)
	fmt.Fprintf(w, "Repos URL:     %v\n", u.Repos_Url)
	fmt.Fprintf(w, "Name:          %v\n", u.Name)
	fmt.Fprintf(w, "Public Repos:  %v\n", u.Public_Repos)
	fmt.Fprintf(w, "Public Gists:  %v\n", u.Public_Gists)
	fmt.Fprintf(w, "Followers:     %v\n", u.Followers)
	fmt.Fprintf(w, "Following:     %v\n", u.Following)
	fmt.Fprintf(w, "Type:          %v\n", u.Type)
	fmt.Fprintf(w, "Site Admin:    %v\n", u.Site_Admin)
	fmt.Fprintf(w, "Hireable:      %v\n", u.Hireable)
	fmt.Fprintf(w, "Bio:           %v\n", u.Bio)
	fmt.Fprintf(w, "Created At:    %v\n", u.Created_At)
	fmt.Fprintf(w, "Updated At:    %v\n", u.Updated_At)
	return w.String()
}

func getUser(name string) (*User, error) {
	url := fmt.Sprintf("https://api.github.com/users/%v", name)
	user := &User{}
	err := get(url, user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

type Repo struct {
	Id               uint64
	Name             string
	Full_Name        string
	Description      string
	Url              string
	Clone_Url        string
	Stargazers_Count uint64
	Forks_Count      uint64
}

func getRepos(u *User) ([]Repo, error) {
	var repos []Repo

	tag := regexp.MustCompile("<(.*?)>")
	rel := regexp.MustCompile("(rel=\".*?\")")
	url := u.Repos_Url
loop:
	for {
		resp, err := http.Get(url)
		if err != nil {
			return nil, err
		}

		buf, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}

		if resp.StatusCode != 200 {
			e := &Error{}
			err := json.Unmarshal(buf, &e)
			if err != nil {
				return nil, err
			}
			return nil, e
		}

		repo := make([]Repo, 0)
		err = json.Unmarshal(buf, &repo)
		if err != nil {
			return nil, err
		}
		repos = append(repos, repo...)

		link := resp.Header.Get("Link")
		m := tag.FindAllString(link, -1)
		n := rel.FindAllString(link, -1)
		if len(m) != len(n) {
			break
		}

		for i, s := range m {
			if len(s) >= 2 {
				m[i] = s[1 : len(s)-1]
			}
		}
		for i, s := range n {
			if len(s) >= 6 {
				n[i] = s[5 : len(s)-1]
			}
		}

		for i := range n {
			if n[i] == "next" {
				url = m[i]
				continue loop
			}
		}
		break
	}

	return repos, nil
}

func execl(name string, args ...string) {
	cmd := exec.Command(name, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	fmt.Println(name, args)
	ek(cmd.Run())
}
