package main

import (
	"io"
	"os"
	"net/http"
	"html/template"
	"io/ioutil"
	"fmt"
	"strings"
	"bytes"
	"encoding/json"
)



/************ MAPPING STRUCTURES  ******************/

type DICT map[string]interface{}
type TemplateData struct {
   Template      string
   Data          DICT
}

var webPages map[string]TemplateData
var menu []MenuItem
var projects []ProjectPage
var userData PersonalData

var projectIndex map[string]int



type MenuItem struct {
    Id        string    `json: "ID"`
    Caption   string    `json: "Caption"`
    Icon      string   	`json: "Icon"`
    Active    bool      `json: "Active"`
}

type ProjectPage struct {
    Id        string    `json: "ID"`
    Title     string  	`json: "Title"`
    Timestamp string   	`json: "Timestamp"`
    Image     string   	`json: "Image"`
    Summary   string	`json: "Summary"`
    Tags      string    `json: "Tags"`
	Content   string    `json: "Content"`
	Page      template.HTML
	
	OddIndex  bool
}


type RequestData struct {
    Path      string  	`json: "path"`
    Method    string   	`json: "method"`
    User      string   	`json: "user"`
    Password  string	`json: "password"`
    Data      string	`json: "data"`
}


type PersonalData struct {

    Name        string
    Profession  string
    Linkedin    string
    Github      string
    Email       string
    Phone       string
    
    Header      string


    Home DICT `json:"Home"`
    
    AboutMe struct {
        Profile        []string  `json:"Profile"`
        Experience    []DICT    `json:"Experience"`
        Education     []DICT    `json:"Education"`
        Publications  []DICT    `json:"Publications"`
    } `json:"AboutMe"`
    
    ProjectsInfo    DICT `json:"ProjectsInfo"`
    ContactInfo     DICT `json:"ContactInfo"`
}


/*************************************************/





/********** FUNCTIONS TO BE CALLED IN TEMPLATES ***/

var templateFunctions = template.FuncMap{
	
	// functions to call inside html templates
	"contains"         : strings.Contains,
	"isIndexOddNumber" : func(index int) bool {
	                      if (index%2 == 0) {
	                      return false
	                      } else{ return true}},
	
	"insertHref"       : func(text string, ht map[string]interface{}, substr string) template.HTML {

	        for key, value := range ht {
	            if strings.Contains(text, key) {	            	    
        	        r := strings.NewReplacer("$key", key, "$value", value.(string))
	                text = strings.Replace(text, key, r.Replace(substr), 1)
	            } 
	        }
	        return template.HTML(text)}}
                      
/***************************************************/







/***********  DISABLE DIRECTORY LISTING  **********/
/** https://stackoverflow.com/questions/49589685/good-way-to-disable-directory-listing-with-http-fileserver-in-go  **/
type neuteredStatFile struct {
	http.File
	readDirBatchSize int
}

type justFilesFilesystem struct {
	fs               http.FileSystem
	readDirBatchSize int
}

func (fs justFilesFilesystem) Open(name string) (http.File, error) {
	f, err := fs.fs.Open(name)
	if err != nil {
	    return nil, err
	}
	return neuteredStatFile{File: f, readDirBatchSize: fs.readDirBatchSize}, nil
}


func (e neuteredStatFile) Stat() (os.FileInfo, error) {
	s, err := e.File.Stat()
	if err != nil {
	    return nil, err
	}
	if s.IsDir() {
	LOOP:
	    for {
		fl, err := e.File.Readdir(e.readDirBatchSize)
		switch err {
		case io.EOF:
		    break LOOP
		case nil:
		    for _, f := range fl {
			if f.Name() == "index.html" {
			    return s, err
			}
		    }
		default:
		    return nil, err
		}
	    }
	    return nil, os.ErrNotExist
	}
	return s, err
}
/*********************************************************************/







/*****************    HANDLERS    *******************/

func handler(w http.ResponseWriter, r *http.Request) {
	var request RequestData
	if r.Method == "POST" {
		decoder := json.NewDecoder(r.Body)
		decoder.Decode(&request)
	} else {
		 request.Path = r.URL.Path
	}

	response := webPages["/"].Template
	data     := webPages["/"].Data
	data["IsLoading"] = false
	
	if request.Path == "/" {
	    request.Path = "/home"
	    data["ActiveMenuItem"] = "home"
		data["IsLoading"] = true
	}
	
	
	path := strings.Trim(request.Path, "/")
	splits := strings.Split(path, "/")
	path = splits[0]
	if path == "projects" && len(splits) > 1 {
		path = "project"
	}
	if content, ok := webPages[path]; ok {

		if path == "project" {
			// for project pages retrieve ids
			data["ActiveMenuItem"] = ""
			_id := splits[1]
			if _, ok := projectIndex[_id]; ok {
				data["Content"] = HtmlTemplate(webPages[path].Template, DICT{"Project": projects[projectIndex[_id]]})
			} else {
				data["Content"] = HtmlTemplate("<center>404 Page not found</center>", nil)
			}
		} else {
		data["ActiveMenuItem"] = splits[0]
		data["Content"] = HtmlTemplate(content.Template, content.Data)
	}
	} else {
		data["Content"] = HtmlTemplate("<center>404 Page not found</center>", nil)
	}


	t := template.New("response").Funcs(templateFunctions)
                                      
	t.Parse(response)
	
	t.Execute(w, data)
}


func formatRequest(r *http.Request) string {
	 // Create return string
	 var request []string

	 // Add the request string
	 url := fmt.Sprintf("%v %v %v", r.Method, r.URL, r.Proto)
	 request = append(request, url)

	 // Add the host
	 request = append(request, fmt.Sprintf("Host: %v", r.Host))

	 // Loop through headers
	 for name, headers := range r.Header {
	   name = strings.ToLower(name)
	   for _, h := range headers {
	     request = append(request, fmt.Sprintf("%v: %v", name, h))
	   }
	 }

	 // If this is a POST, add post data
	 if r.Method == "POST" {
	    r.ParseForm()
	    request = append(request, "\n")
	    request = append(request, r.Form.Encode())
	 }

	  // Return the request as a string
	  return strings.Join(request, "\n")
}

/***************************************************************/




/******************    DATA LOAD   ****************************/
func loadMenuItems() {

	jsonFile, err := os.Open("menu.json")
	if err != nil {
		fmt.Println(err)
	}
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)
	json.Unmarshal(byteValue, &menu)
}



func loadJsonFile(fname string, s interface{}) {

	jsonFile, err := os.Open(fname)
	if err != nil {
		fmt.Println(err)
	}
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)
	json.Unmarshal(byteValue, s)
}

func loadUserData() {

    loadJsonFile("personal.json", &userData)
	
}

//loads the list with projects and the content of the page for each project
func loadProjectData() {

    loadJsonFile("project_list.json", &projects)
	for index, _ := range projects {
		projects[index].Page = HtmlTemplate(loadHtmlTemplate(projects[index].Content), DICT{"Project" : projects[index]})
		
		if (index %2 == 0) {
          projects[index].OddIndex = false
        } else {
          projects[index].OddIndex = true;
        }
	}

	projectIndex = setprojectIndexMap()
}


func loadHtmlTemplate(path string) string {
	html, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Printf("\nFailed to read %s. Err: %s \n", path, err)
		return "<center>404 Page not found</center>"
	}
	html_str := string(html[:])
	//remove the damned invisible &#65279
	bom := []byte{0xef, 0xbb, 0xbf} // UTF-8
	if bytes.Equal(html[:3], bom) {
		html_str = string(html[3:])
	}
	return html_str
}


func HtmlTemplate(tmpl string, data map[string]interface{}) template.HTML {
  if data != nil {
    t := template.New("partial").Funcs(templateFunctions)
    t, err := t.Parse(tmpl)
    if err != nil {
      fmt.Printf("\nFailed to convert string to html:\n%s\n\nErr: %s \n", tmpl, err.Error())
      return template.HTML("404 Page not found")
    }
    var tpl bytes.Buffer
    if err := t.Execute(&tpl, data); err != nil {
      fmt.Printf("\nFailed to convert string to html:\n%s\n\nErr: %s \n", tmpl, err.Error())
      return template.HTML("<center>404 Page not found</center>")
    }
    return template.HTML(tpl.String())
  }
  return template.HTML(tmpl)
}


func setprojectIndexMap() map[string]int {
  m := make(map[string]int, len(projects))
  for index, elem := range projects {
    m[elem.Id] = index
  }
	return m
}

func loadTemplates() {

	webPages = map[string]TemplateData{
		"/"       : TemplateData{
				loadHtmlTemplate("templates/main.tmpl"),
				DICT{
					"Name"           : userData.Name,
					"Profession"     : userData.Profession,
					"Header"         : userData.Header,
					"Linkedin"       : userData.Linkedin,
					"Github"         : userData.Github,
					"Email"          : userData.Email,
					"Phone"          : userData.Phone,
					
					"Menu"           : menu,
					"Content"        : ""}},
		"home"    : TemplateData{ 
		        loadHtmlTemplate("templates/home.tmpl"), 
		        DICT{
		            "Home" : userData.Home, 
		            "Projects" : projects[0:3]}},
		"projects": TemplateData{
				loadHtmlTemplate("templates/projects.tmpl"),
				DICT{
				    "ProjectsInfo" : userData.ProjectsInfo,
					"Projects" : projects}},
		"about"   : TemplateData{
				loadHtmlTemplate("templates/about.tmpl"),
				DICT{
				    "Name"         : userData.Name,
				    "Profile"       : userData.AboutMe.Profile,
				    "Experience"   : userData.AboutMe.Experience,
				    "Education"    : userData.AboutMe.Education,
				    "Publications" : userData.AboutMe.Publications}},
		"setup"   : TemplateData{
				loadHtmlTemplate("templates/projects/test.tmpl"),
				nil},
		"project": TemplateData{
				loadHtmlTemplate("templates/page.tmpl"),
				nil},
		"contact" : TemplateData{
				loadHtmlTemplate("templates/contact.tmpl"),
				DICT{
				    "ContactInfo"    : userData.ContactInfo,
				    "Linkedin"       : userData.Linkedin,
					"Github"         : userData.Github,
					"Email"          : userData.Email,
					"Phone"          : userData.Phone}}}
	
}



func loadWebsite() {
	loadMenuItems()
	loadUserData()
	loadProjectData()
	loadTemplates()
}

/************************************************************************/




func main() {

	loadWebsite()

	server := http.Server{
		Addr: ":8090",
	}

	fss := justFilesFilesystem{fs: http.Dir("public"), readDirBatchSize: 2}
	fs := http.FileServer(fss)
	http.Handle("/public/", http.StripPrefix("/public/", fs))
	http.HandleFunc("/", handler)
	fmt.Printf("Listening on 8090 ...\n")
	server.ListenAndServe()
}





