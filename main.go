package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gorilla/csrf"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	mailgun "gopkg.in/mailgun/mailgun-go.v1"
)

var csrfToken = "DdgSvx83mA2C3aSW7Xudh3yPXusUtngd"

var mgDomain string
var mgPublicAPIKey string
var mgPrivateAPIKey string

var productionFlag bool

var domainName string

var serverPort string

func homeHandler(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("templates/index.html")
	t.Execute(w, map[string]interface{}{
		csrf.TemplateTag: csrf.TemplateField(r),
	})
}

func cvHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Disposition", "attachment; filename=James_Bannister_CV.pdf")
	w.Header().Set("Content-Type", "application/pdf")
	http.ServeFile(w, r, "./download/James_Bannister_CV.pdf")
}

func contactHandler(w http.ResponseWriter, r *http.Request) {

	mg := mailgun.NewMailgun(mgDomain, mgPrivateAPIKey, mgPublicAPIKey)
	message := mailgun.NewMessage(
		"mailrobot@bannister.me",
		"Someone wants to say hello!",
		fmt.Sprintf(`Hey James,

Someone touched base with you through your website.

If you want to know more, their details are below:

Name - %v
Email - %v
Phone - %v
Message -
%v

Thanks,
Your Friendly Mail Robot`,
			r.FormValue("name"),
			r.FormValue("email"),
			r.FormValue("phone"),
			r.FormValue("message"),
		),
		"james@bannister.me")
	_, _, err := mg.Send(message)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Someone contacted you! An email is on the way.")
}

func sitemapHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/xml")
	http.ServeFile(w, r, "./asset/sitemap.xml")
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	mgDomain = os.Getenv("MG_DOMAIN")
	mgPublicAPIKey = os.Getenv("MG_PUBLIC_API_KEY")
	mgPrivateAPIKey = os.Getenv("MG_API_KEY")

	productionFlag, err = strconv.ParseBool(os.Getenv("PRODUCTION"))

	domainName = os.Getenv("DOMAIN_NAME")

	serverPort = fmt.Sprintf(":%v", os.Getenv("SERVER_PORT"))

	if err != nil {
		log.Fatal(err)
	}

	r := mux.NewRouter()
	s := http.StripPrefix("/assets/", http.FileServer(http.Dir("./asset/")))
	r.PathPrefix("/assets/").Handler(s).Methods("GET")
	r.HandleFunc("/", homeHandler).Methods("GET")
	r.HandleFunc("/cv", cvHandler).Methods("GET")
	r.HandleFunc("/send", contactHandler).Methods("POST")
	r.HandleFunc("/sitemap.xml", sitemapHandler).Methods("GET")
	http.Handle("/", r)

	srv := &http.Server{
		Handler: csrf.Protect([]byte(csrfToken), csrf.FieldName("_token"), csrf.Secure(productionFlag))(r),
		Addr:    serverPort,
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())

}
