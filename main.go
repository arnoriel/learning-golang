package main

import (
    "encoding/json"
    "fmt"
    "html/template"
    "io/ioutil"
    "net/http"
    "strconv"
)

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

// Render index.html with FAQ list
func indexHandler(w http.ResponseWriter, r *http.Request) {
    tmpl, err := template.ParseFiles("templates/index.html")
    if err != nil {
        http.Error(w, "Error parsing template", http.StatusInternalServerError)
        return
    }
    tmpl.Execute(w, faqs)
}

// Render faq.html for CRUD operations
func faqHandler(w http.ResponseWriter, r *http.Request) {
    tmpl, err := template.ParseFiles("templates/faq.html")
    if err != nil {
        http.Error(w, "Error parsing template", http.StatusInternalServerError)
        return
    }
    tmpl.Execute(w, faqs)
}

// Add a new FAQ
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


// Delete FAQ by ID
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

func main() {
    // Load FAQ data from JSON
    loadFAQs()

    // Define routes for the application
    http.HandleFunc("/", indexHandler)
    http.HandleFunc("/faq", faqHandler)
    http.HandleFunc("/faq/add", addFAQHandler)
    http.HandleFunc("/faq/delete", deleteFAQHandler)

    // Serve static files such as CSS from the "static" folder
    fs := http.FileServer(http.Dir("./static"))
    http.Handle("/static/", http.StripPrefix("/static/", fs))

    fmt.Println("Server running on http://localhost:8080")
    http.ListenAndServe(":8080", nil)
}
