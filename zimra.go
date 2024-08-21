package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
)

func sendToZimbra(acc Acc, message string, fromdate string, todate string) {
	// URL для отправки POST запроса с данными для авторизации
	loginURL := "https://mail.ft.by/?ignoreLoginURL=1" // Замените на URL авторизации вашего Zimbra сервера

	// Создание клиента с поддержкой cookie
	jar, err := cookiejar.New(nil)
	if err != nil {
		log.Fatal(err)
	}
	client := &http.Client{
		Jar: jar,
	}

	// Шаг 1: Получение страницы входа, чтобы инициализировать cookies
	req, err := http.NewRequest("GET", loginURL, nil)
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/107.0.0.0 Safari/537.36")

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	// Шаг 2: Чтение и игнорирование тела ответа
	_, err = io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	// Данные для отправки в форме
	data := url.Values{}
	data.Set("loginOp", "login")
	data.Set("username", acc.Login)  // Замените на имя пользователя
	data.Set("password", acc.Pass)   // Замените на пароль
	data.Set("client", "advanced  ") // Используйте клиент по умолчанию
	csrf := ""
	for _, cookie := range client.Jar.Cookies(req.URL) {
		fmt.Printf("Name: %s, Value: %s\n", cookie.Name, cookie.Value)
		if cookie.Name == "ZM_LOGIN_CSRF" {
			csrf = cookie.Value
		}
	}

	data.Set("login_csrf", csrf)
	//

	// Шаг 3: Создание запроса на авторизацию
	req, err = http.NewRequest("POST", loginURL, strings.NewReader(data.Encode()))
	if err != nil {
		log.Fatal(err)
	}

	// Установка заголовков
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/107.0.0.0 Safari/537.36")

	// Отправка запроса на авторизацию
	resp, err = client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	// Проверка статуса ответа
	if resp.StatusCode != http.StatusOK {
		log.Fatalf("Ошибка авторизации: %s", resp.Status)
	}

	fmt.Println("Успешная авторизация")

	// Шаг 4: Проверка cookies
	fmt.Println("Cookies после авторизации:")

	ZM_AUTH_TOKEN := ""
	//JSESSIONID := ""
	for _, cookie := range client.Jar.Cookies(req.URL) {
		fmt.Printf("Name: %s, Value: %s\n", cookie.Name, cookie.Value)
		switch cookie.Name {
		case "ZM_AUTH_TOKEN":
			ZM_AUTH_TOKEN = cookie.Value
		case "ZM_LOGIN_COOKIES":
			//JSESSIONID = cookie.Value
		}

	}

	data.Set("login_csrf", csrf)
	//body, err := io.ReadAll(resp.Body)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//fmt.Println(string(body))
	// Пример запроса к защищенной странице после авторизации
	protectedURL := "https://mail.ft.by/?ignoreLoginURL=1" // Замените на URL защищенной страницы

	resp, err = client.Get(protectedURL)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	// Проверка доступа к защищенной странице
	if resp.StatusCode != http.StatusOK {
		log.Fatalf("Ошибка доступа к защищенной странице: %s", resp.Status)
	}

	//Чтение и вывод защищенной страницы
	protectedBody, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Содержимое защищенной страницы:")
	//fmt.Println(string(protectedBody))
	_, after, _ := strings.Cut(string(protectedBody), "localStorage.setItem")
	log.Print(after[16:58])
	TOKEN := after[16:58]
	//xml := `<soap:Envelope xmlns:soap="http://www.w3.org/2003/05/soap-envelope"><soap:Header><context xmlns="urn:zimbra"><userAgent xmlns="" name="ZimbraWebClient - FF126 (Win)" version="8.8.15_GA_4156"/><session xmlns="" id="1503980"/><account xmlns="" by="name">autotest@ft.by</account><format xmlns="" type="js"/><csrfToken xmlns="">` + TOKEN + `</csrfToken></context></soap:Header><soap:Body><BatchRequest xmlns="urn:zimbra" onerror="stop"><NoOpRequest xmlns="urn:zimbraMail" requestId="0"/><ModifyPrefsRequest xmlns="urn:zimbraAccount" requestId="1"><pref xmlns="" name="zimbraPrefOutOfOfficeFromDate">20240731210000Z</pref><pref xmlns="" name="zimbraPrefOutOfOfficeReplyEnabled">TRUE</pref><pref xmlns="" name="zimbraPrefOutOfOfficeReply">HELLO WORLD</pref><pref xmlns="" name="zimbraPrefOutOfOfficeUntilDate">20240830205900Z</pref></ModifyPrefsRequest></BatchRequest></soap:Body></soap:Envelope>`
	xml := `<soap:Envelope xmlns:soap="http://www.w3.org/2003/05/soap-envelope"><soap:Header><context xmlns="urn:zimbra"><userAgent xmlns="" name="ZimbraWebClient - FF126 (Win)" version="8.8.15_GA_4156"/><session xmlns="" id="1503980"/><account xmlns="" by="name">autotest@ft.by</account><format xmlns="" type="js"/><csrfToken xmlns="">` + TOKEN + `</csrfToken></context></soap:Header><soap:Body><BatchRequest xmlns="urn:zimbra" onerror="stop"><NoOpRequest xmlns="urn:zimbraMail" requestId="0"/><ModifyPrefsRequest xmlns="urn:zimbraAccount" requestId="1"><pref xmlns="" name="zimbraPrefOutOfOfficeFromDate">` + fromdate + `</pref><pref xmlns="" name="zimbraPrefOutOfOfficeReplyEnabled">TRUE</pref><pref xmlns="" name="zimbraPrefOutOfOfficeReply">` + message + `</pref><pref xmlns="" name="zimbraPrefOutOfOfficeUntilDate">` + todate + `</pref></ModifyPrefsRequest></BatchRequest></soap:Body></soap:Envelope>`

	log.Print(ZM_AUTH_TOKEN)
	response, err := client.Post("https://mail.ft.by/service/soap/ModifyPrefsRequest", "text/xml", bytes.NewBufferString(xml))

	if err != nil {
		fmt.Println(err)
	}
	defer response.Body.Close()
	content, _ := io.ReadAll(response.Body)
	s := strings.TrimSpace(string(content))
	log.Print(s)
}
