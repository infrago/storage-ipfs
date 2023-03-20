package store_ipfs

import (
	"bytes"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"

	. "github.com/infrago/base"
	"github.com/infrago/infra"
)

// 发送HTTP GET请求
func HttpGet(url string, args ...Map) string {
	req, err := http.NewRequest("GET", url, nil)
	// fmt.Println("get", err)
	if err == nil {

		//处理头
		headers := Map{}
		if len(args) > 0 {
			for k, v := range args[0] {
				headers[k] = v
			}
		}

		for k, v := range headers {
			req.Header.Set(k, fmt.Sprintf("%v", v))
		}

		//发送请求
		res, err := http.DefaultClient.Do(req)
		// fmt.Println("get 2", err)
		if err == nil {
			defer res.Body.Close()
			resBody, err := ioutil.ReadAll(res.Body)
			// fmt.Println("get 3444", err, string(resBody))
			if err == nil {
				return string(resBody)
			}
		}

	}

	return ""

	/*
	   res, err := http.Get(url)
	   if err == nil {
	       defer res.Body.Close()
	       body, err := ioutil.ReadAll(res.Body)
	       if err == nil {
	           return string(body)
	       }
	   }
	   return ""
	*/
}

// 发送HTTP DELETE请求
func HttpDelete(url string, args ...Map) string {

	req, err := http.NewRequest("DELETE", url, nil)
	if err == nil {

		//处理头
		headers := Map{}
		if len(args) > 0 {
			for k, v := range args[0] {
				headers[k] = v
			}
		}

		for k, v := range headers {
			req.Header.Set(k, fmt.Sprintf("%v", v))
		}

		//发送请求
		res, err := http.DefaultClient.Do(req)
		if err == nil {
			defer res.Body.Close()
			resBody, err := ioutil.ReadAll(res.Body)
			if err == nil {
				return string(resBody)
			}
		}

	}

	return ""

	/*
	   res, err := http.Get(url)
	   if err == nil {
	       defer res.Body.Close()
	       body, err := ioutil.ReadAll(res.Body)
	       if err == nil {
	           return string(body)
	       }
	   }
	   return ""
	*/
}

func HttpGetBody(url string, args ...Map) (io.ReadCloser, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	//处理头
	headers := Map{}
	if len(args) > 0 {
		for k, v := range args[0] {
			headers[k] = v
		}
	}

	for k, v := range headers {
		req.Header.Set(k, fmt.Sprintf("%v", v))
	}

	//发送请求
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	return res.Body, nil
}

func HttpForm(uuu string, data Map) string {

	values := url.Values{}
	for k, v := range data {
		values.Add(k, fmt.Sprintf("%v", v))
	}
	res, err := http.PostForm(uuu, values)
	if err == nil {
		defer res.Body.Close()

		body, err := ioutil.ReadAll(res.Body)
		if err == nil {
			return string(body)
		}
	}

	return ""
}

// args是HEADER自定的
func HttpPostBytes(url string, bodyType string, body io.Reader, args ...Map) []byte {

	req, err := http.NewRequest("POST", url, body)
	if err == nil {

		//处理头
		headers := Map{}
		if len(args) > 0 {
			for k, v := range args[0] {
				headers[k] = v
			}
		}

		req.Header.Set("Content-Type", bodyType)
		for k, v := range headers {
			req.Header.Set(k, fmt.Sprintf("%v", v))
		}

		//发送请求
		res, err := http.DefaultClient.Do(req)
		//fmt.Println("http", url, err)
		if err == nil {
			defer res.Body.Close()
			resBody, err := ioutil.ReadAll(res.Body)

			if err == nil {
				return resBody
			}
		}

	}

	return nil

	/*
	   res, err := http.Post(url, bodyType, body)
	   if err != nil {
	       return ""
	   } else {
	       defer res.Body.Close()
	       body, err := ioutil.ReadAll(res.Body)
	       if err != nil {
	           return ""
	       } else {
	           return string(body)
	       }
	   }
	*/
}
func HttpPost(url string, bodyType string, body io.Reader, args ...Map) (string, error) {

	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return "", err
	}

	//处理头
	headers := Map{}
	if len(args) > 0 {
		for k, v := range args[0] {
			headers[k] = v
		}
	}

	req.Header.Set("Content-Type", bodyType)
	for k, v := range headers {
		req.Header.Set(k, fmt.Sprintf("%v", v))
	}

	//发送请求
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()
	resBody, err := ioutil.ReadAll(res.Body)

	if err != nil {
		return "", err
	}

	return string(resBody), nil
}

func HttpGetJson(url string) Map {
	body := HttpGet(url)
	if body != "" {
		m := Map{}
		err := infra.UnmarshalJSON([]byte(body), &m)
		if err == nil {
			return m
		}
	}
	return nil
}

func HttpDeleteJson(url string, args ...Map) Map {
	body := HttpDelete(url, args...)
	if body != "" {
		m := Map{}
		err := infra.UnmarshalJSON([]byte(body), &m)
		if err == nil {
			return m
		}
	}
	return nil
}

func HttpFormJson(url string, data Map) Map {
	body := HttpForm(url, data)
	if body != "" {
		m := Map{}
		err := infra.UnmarshalJSON([]byte(body), &m)
		if err == nil {
			return m
		}
	}
	return nil
}

func HttpPostJson(url string, bodyType string, body io.Reader, args ...Map) (Map, error) {
	resBody, err := HttpPost(url, bodyType, body, args...)
	if err != nil {
		return nil, err
	}
	if resBody == "" {
		return nil, errors.New("返回空")
	}
	m := Map{}
	err = infra.UnmarshalJSON([]byte(resBody), &m)
	if err != nil {
		return nil, err
	}
	return m, nil
}

func HttpPostXml(url string, bodyType string, body io.Reader, args ...Map) (Map, error) {
	resBody, err := HttpPost(url, bodyType, body, args...)
	if err != nil {
		return nil, err
	}
	if resBody == "" {
		return nil, errors.New("返回空")
	}

	m := Map{}
	err = xml.Unmarshal([]byte(resBody), &m)
	if err != nil {
		return nil, err
	}

	return m, nil
}

func HttpUpload(url string, name, file string, args ...Map) (string, error) {
	bodyBuf := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuf)

	//关键的一步操作
	fileWriter, err := bodyWriter.CreateFormFile(name, file)
	if err != nil {
		return "", err
	}

	//打开文件句柄操作
	fh, err := os.Open(file)
	if err != nil {
		return "", err
	}
	defer fh.Close()

	//iocopy
	_, err = io.Copy(fileWriter, fh)
	if err != nil {
		return "", err
	}

	//字段
	if len(args) > 0 {
		for key, val := range args[0] {
			bodyWriter.WriteField(key, fmt.Sprintf("%v", val))
		}
	}

	contentType := bodyWriter.FormDataContentType()
	bodyWriter.Close()

	resp, err := http.Post(url, contentType, bodyBuf)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	resp_body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(resp_body), nil
}

func HttpUploadReader(url string, field, name string, file io.Reader, args ...Map) ([]byte, error) {
	bodyBuf := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuf)

	//关键的一步操作
	fileWriter, err := bodyWriter.CreateFormFile(field, name)
	if err != nil {
		return nil, err
	}

	//iocopy
	_, err = io.Copy(fileWriter, file)
	if err != nil {
		return nil, err
	}

	//字段
	if len(args) > 0 {
		for key, val := range args[0] {
			bodyWriter.WriteField(key, fmt.Sprintf("%v", val))
		}
	}

	contentType := bodyWriter.FormDataContentType()
	bodyWriter.Close()

	resp, err := http.Post(url, contentType, bodyBuf)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	resp_body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return resp_body, nil
}

// 发送HTTP download
func HttpDownload(url, file string) error {

	res, err := http.Get(url)
	if err != nil {
		return err
	} else {
		defer res.Body.Close()

		f, err := os.Create(file)
		if err != nil {
			return err
		} else {
			defer f.Close()
			_, err := io.Copy(f, res.Body)
			return err
		}
	}
}

// 发送HTTP download
func HttpDownloadFile(url string, file *os.File) error {

	res, err := http.Get(url)
	if err != nil {
		return err
	} else {
		defer res.Body.Close()

		_, err := io.Copy(file, res.Body)
		return err

	}
}
