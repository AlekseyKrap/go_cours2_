package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
)

//Напишите функцию, которая будет получать на вход строку с поисковым запросом (string) и массив ссылок на страницы,
//по которым стоит произвести поиск ([]string).
//	Результатом работы функции должен быть массив строк со ссылками на страницы,
//на которых обнаружен поисковый запрос. Функция должна искать точное соответствие фразе в тексте ответа от сервера по каждой из ссылок.
//Подсказка: это задача на поиск последовательности в массиве.

func search(str string, urls []string) (answer []string) {

	for _, value := range urls {

		bodyS, err := func(value string) (string, error) {
			resp, err := http.Get(value)
			if err != nil {
				return "", err
			}
			defer resp.Body.Close()

			body, err := ioutil.ReadAll(resp.Body)

			return string(body), err
		}(value)

		if err != nil {
			fmt.Println(err)
			return
		}

		i := strings.Index(bodyS, str)

		if i >= 0 {
			answer = append(answer, value)
		}

	}

	return answer
}

//скачивание файла

func dowland(fileYandex string, path string) error {

	type Resp struct {
		Href      string
		Method    string
		Templated bool
	}

	var urlY string = "https://cloud-api.yandex.net/v1/disk/public/resources/download?public_key=" + url.QueryEscape(fileYandex)

	client := &http.Client{}
	req, _ := http.NewRequest("GET", urlY, nil)
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return errors.New(resp.Status)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var response Resp
	json.Unmarshal([]byte(string(body)), &response)

	file, err := http.Get(response.Href)
	if err != nil {
		return err
	}
	defer file.Body.Close()

	m, err := url.ParseQuery(response.Href)
	if err != nil {
		return err
	}

	out, err := os.Create(path + m.Get("filename"))
	if err != nil {
		return err
	}
	defer out.Close()
	_, err = io.Copy(out, file.Body)
	if err != nil {
		return err
	}

	return nil
}

func main() {

	urls := []string{"https://yandex.ru/", "https://golang.org/", "https://www.google.com/search?newwindow=1&sxsrf=ALeKk02lbVBhdf5Fv5ZRfVj_1_YfSC2_Bg%3A1586788921635&ei=OXqUXo6XJuevrgS3mbp4&q=golang&oq=golang&gs_lcp=CgZwc3ktYWIQAzIECCMQJzIECCMQJzIECCMQJzIECAAQQzIECAAQQzIECAAQQzIECAAQQzIECAAQQzIECAAQQzIECAAQQzoECAAQR0oOCBcSCjYtMTExZzExLTRKDAgYEgg2LTFnMTEtMlCHG1iHG2DiIWgAcAJ4AIABa4gBa5IBAzAuMZgBAKABAaoBB2d3cy13aXo&sclient=psy-ab&ved=0ahUKEwjOl4ew0eXoAhXnl4sKHbeMDg8Q4dUDCAw&uact=5"}
	answer := search("golang", urls)
	fmt.Println(answer)

	err := dowland("https://yadi.sk/i/9fcgSjyKV9lKQQ", "C:\\Users\\Алексей\\Documents\\test\\")
	if err != nil {
		fmt.Println(err)
		return
	}
}
