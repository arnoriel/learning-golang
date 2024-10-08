package main

import (
    "encoding/json"
    "fmt"
    "html/template"
    "io/ioutil"
    "log"
    "net/http"
    "os"
    "strconv"

    "github.com/joho/godotenv"
)

// Struct for FAQ
type FAQ struct {
    ID       int    `json:"id"`
    Question string `json:"question"`
    Answer   string `json:"answer"`
}

var faqs []FAQ
const filePath = "faq.json"

func loadFAQs() error {
    file, err := ioutil.ReadFile(filePath)
    if err != nil {
        return err
    }
    err = json.Unmarshal(file, &faqs)
    if err != nil {
        return err
    }
    return nil
}

func saveFAQs() error {
    data, err := json.MarshalIndent(faqs, "", "  ")
    if err != nil {
        return err
    }
    err = ioutil.WriteFile(filePath, data, 0644)
    if err != nil {
        return err
    }
    return nil
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
    // Load FAQs and settings
    err := loadFAQs()
    if err != nil {
        http.Error(w, "Error loading FAQs", http.StatusInternalServerError)
        return
    }
    err = loadSettings()
    if err != nil {
        http.Error(w, "Error loading settings", http.StatusInternalServerError)
        return
    }

    // Data to pass to the template
    data := struct {
        FAQs     []FAQ
        Settings Settings
    }{
        FAQs:     faqs,
        Settings: settings,
    }

    // Parse and execute the template
    tmpl, err := template.ParseFiles("templates/index.html")
    if err != nil {
        http.Error(w, "Error parsing template", http.StatusInternalServerError)
        return
    }
    tmpl.Execute(w, data)
}

func faqHandler(w http.ResponseWriter, r *http.Request) {
    tmpl, err := template.ParseFiles("templates/faq.html")
    if err != nil {
        http.Error(w, "Error parsing template", http.StatusInternalServerError)
        return
    }
    tmpl.Execute(w, faqs)
}

func addFAQHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method == "POST" {
        question := r.FormValue("question")
        answer := r.FormValue("answer")
        id := len(faqs) + 1
        faq := FAQ{
            ID:       id,
            Question: question,
            Answer:   answer,
        }
        faqs = append(faqs, faq)
        saveFAQs()
    }
    http.Redirect(w, r, "/faq", http.StatusSeeOther)
}

func editFAQHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method == "POST" {
        id, _ := strconv.Atoi(r.FormValue("id"))
        question := r.FormValue("question")
        answer := r.FormValue("answer")

        for i, faq := range faqs {
            if faq.ID == id {
                faqs[i].Question = question
                faqs[i].Answer = answer
                saveFAQs()
                break
            }
        }
    }
    http.Redirect(w, r, "/faq", http.StatusSeeOther)
}

func deleteFAQHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method == "POST" {
        id, _ := strconv.Atoi(r.FormValue("id"))
        for i, faq := range faqs {
            if faq.ID == id {
                faqs = append(faqs[:i], faqs[i+1:]...)
                for j := i; j < len(faqs); j++ {
                    faqs[j].ID--
                }
                saveFAQs()
                break
            }
        }
    }
    http.Redirect(w, r, "/faq", http.StatusSeeOther)
}

// Struct for Settings
type Settings struct {
    AppName     string `json:"app_name"`
    About       string `json:"about"`
    Description string `json:"description"`
}

var settings Settings
const settingsFilePath = "settings.json"

// Load settings from settings.json
func loadSettings() error {
    file, err := ioutil.ReadFile(settingsFilePath)
    if err != nil {
        return err
    }
    err = json.Unmarshal(file, &settings)
    if err != nil {
        return err
    }
    return nil
}

// Save settings to settings.json
func saveSettings() error {
    data, err := json.MarshalIndent(settings, "", "  ")
    if err != nil {
        return err
    }
    err = ioutil.WriteFile(settingsFilePath, data, 0644)
    if err != nil {
        return err
    }
    return nil
}

// Handler to display settings page
func settingsHandler(w http.ResponseWriter, r *http.Request) {
    tmpl, err := template.ParseFiles("templates/settings.html")
    if err != nil {
        http.Error(w, "Error parsing template", http.StatusInternalServerError)
        return
    }
    tmpl.Execute(w, settings)
}

// Handler to update settings
func updateSettingsHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method == "POST" {
        settings.AppName = r.FormValue("app_name")
        settings.About = r.FormValue("about")
        settings.Description = r.FormValue("description")
        saveSettings()
    }
    http.Redirect(w, r, "/settings", http.StatusSeeOther)
}

func main() {
    // Load environment variables
    err := godotenv.Load()
    if err != nil {
        log.Fatalf("Error loading .env file")
    }

    host := os.Getenv("HOST")
    port := os.Getenv("PORT")

    if host == "" {
        host = "0.0.0.0"
    }
    if port == "" {
        port = "8080"
    }

    // Load FAQ data from JSON
    loadFAQs()
    // Load settings from settings.json
    loadSettings()

     // Define routes for the application
    http.HandleFunc("/", indexHandler)
    http.HandleFunc("/faq", faqHandler)
    http.HandleFunc("/faq/add", addFAQHandler)
    http.HandleFunc("/faq/edit", editFAQHandler)
    http.HandleFunc("/faq/delete", deleteFAQHandler)
    http.HandleFunc("/settings", settingsHandler)
    http.HandleFunc("/settings/update", updateSettingsHandler)

    // Serve static files such as CSS from the "static" folder
    fs := http.FileServer(http.Dir("./static"))
    http.Handle("/static/", http.StripPrefix("/static/", fs))

    address := host + ":" + port
    fmt.Printf("Server running on http://%s\n", address)
    log.Fatal(http.ListenAndServe(address, nil))
}
