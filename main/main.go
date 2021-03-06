package main

import (
	"fmt"
	"net/http"

	"github.com/gregpechiro/cookieManager/cook"
	"github.com/gregpechiro/cookieManager/csrf"
	"github.com/gregpechiro/cookieManager/flash"
	"github.com/gregpechiro/cookieManager/sess"
)

var count int

func main() {

	http.HandleFunc("/get", GetCook)
	http.HandleFunc("/put", PutCook)
	http.HandleFunc("/del", DelCook)
	http.HandleFunc("/getall", GetAllCook)

	http.HandleFunc("/getflash", GetFlash)
	http.HandleFunc("/setflash", SetFlash)

	http.HandleFunc("/login", Login)
	http.HandleFunc("/logout", Logout)

	http.HandleFunc("/secure", Secure)
	http.HandleFunc("/secure/put", SecurePut)
	http.HandleFunc("/secure/all", SecureAll)

	http.HandleFunc("/admin", Admin)
	http.HandleFunc("/admin/put", AdminPut)
	http.HandleFunc("/admin/all", AdminAll)

	http.HandleFunc("/csrf", Csrf)
	http.HandleFunc("/csrf-test", CsrfTest)

	http.ListenAndServe(":8080", nil)

}

func GetCook(w http.ResponseWriter, r *http.Request) {
	c := r.FormValue("cook")
	var val string
	if c == "" {
		val = ""
	} else {
		val = cook.Get(r, r.FormValue("cook"))
	}
	fmt.Fprintf(w, "Cookie: %v", val)
	return
}

func PutCook(w http.ResponseWriter, r *http.Request) {
	count++
	key := fmt.Sprintf("KEY %d", count)
	val := fmt.Sprintf("VALUE %d", count)
	cook.Put(w, key, val)
	fmt.Fprintf(w, "Set Cookie: %s, %s", key, val)
	return
}

func DelCook(w http.ResponseWriter, r *http.Request) {
	cook.Delete(w, r, r.FormValue("cook"))
	http.Redirect(w, r, "/getall", 303)
	return
}

func GetAllCook(w http.ResponseWriter, r *http.Request) {
	c := cook.GetAll(r)
	fmt.Fprintf(w, "Cookies:    %v", c)
	return
}

func GetFlash(w http.ResponseWriter, r *http.Request) {
	msgk, msgv := flash.GetFlash(w, r)
	fmt.Fprintf(w, "Message in cookie: %s, %s", msgk, msgv)
	return
}

func SetFlash(w http.ResponseWriter, r *http.Request) {
	switch r.FormValue("kind") {
	case "success":
		flash.SetSuccessRedirect(w, r, "/getflash", r.FormValue("msg"))
	case "error":
		flash.SetErrorRedirect(w, r, "/getflash", r.FormValue("msg"))
	default:
		flash.SetMsgRedirect(w, r, "/getflash", r.FormValue("msg"))
	}
	return
}

func Login(w http.ResponseWriter, r *http.Request) {
	role := r.FormValue("role")
	if role == "" {
		role = "user"
	}
	sess.Login(w, r, role)
	fmt.Fprintf(w, "You are now logged in")
	return
}

func Logout(w http.ResponseWriter, r *http.Request) {
	sess.Logout(w, r)
	fmt.Fprintf(w, "You are now logged out")
	return
}

func Secure(w http.ResponseWriter, r *http.Request) {
	_, ok := sess.Authorized(w, r)
	if !ok {
		fmt.Fprintf(w, "You are not autorized. Please visit the login page")
		return
	}
	fmt.Fprintf(w, "You are now viewing secure data")
	return
}

func SecurePut(w http.ResponseWriter, r *http.Request) {
	_, ok := sess.Authorized(w, r)
	if !ok {
		fmt.Fprintf(w, "You are not autorized. Please visit the login page")
		return
	}
	count++
	key := fmt.Sprintf("KEY %d", count)
	val := fmt.Sprintf("VALUE %d", count)
	sess.Put(w, r, key, val)
	fmt.Fprintf(w, "Set Cookie: %s, %s", key, val)
	return
}

func SecureAll(w http.ResponseWriter, r *http.Request) {
	_, ok := sess.Authorized(w, r)
	if !ok {
		fmt.Fprintf(w, "You are not autorized. Please visit the login page")
		return
	}
	c := sess.GetAll(r)
	fmt.Fprintf(w, "Session Cookies:    %v", c)
	return
}

func Admin(w http.ResponseWriter, r *http.Request) {
	role, ok := sess.Authorized(w, r)
	if !ok || role != "admin" {
		fmt.Printf("ROLE: %q0\n", role)
		fmt.Fprintf(w, "You are not autorized as admin. Please visit the login page")
		return
	}
	fmt.Fprintf(w, "You are now viewing admin data")
	return
}

func AdminPut(w http.ResponseWriter, r *http.Request) {
	role, ok := sess.Authorized(w, r)
	if !ok || role != "admin" {
		fmt.Fprintf(w, "You are not autorized as admin. Please visit the login page")
		return
	}
	count++
	key := fmt.Sprintf("KEY %d", count)
	val := fmt.Sprintf("VALUE %d", count)
	sess.Put(w, r, key, val)
	fmt.Fprintf(w, "Set Cookie: %s, %s", key, val)
	return
}

func AdminAll(w http.ResponseWriter, r *http.Request) {
	role, ok := sess.Authorized(w, r)
	if !ok || role != "admin" {
		fmt.Fprintf(w, "You are not autorized as admin. Please visit the login page")
		return
	}
	c := sess.GetAll(r)
	fmt.Fprintf(w, "Session Cookies:    %v", c)
	return
}

func Csrf(w http.ResponseWriter, r *http.Request) {
	c := csrf.SetCSRF(w)
	link := fmt.Sprintf("<a href=\"/csrf-test?_csrf=%s\">Click to test</a>", c)
	fmt.Fprintf(w, "<html><body>Your CSRF token is: %s<br>%s</body></html>", c, link)
	return
}

func CsrfTest(w http.ResponseWriter, r *http.Request) {
	if !csrf.ValidCSRF(r) {
		fmt.Fprintf(w, "Invalid CSRF token")
		return
	}
	fmt.Fprintf(w, "Successfull CSRF token match")
	return
}
