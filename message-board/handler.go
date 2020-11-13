package main

import (
	"fmt"
	"log"
	"net"
	"net/url"
	"strconv"
	"strings"
	"time"

	"encoding/json"

	"github.com/go-redis/redis/v8"
	"github.com/kjk/betterguid"
)

type rawReply struct {
	From    string `json:"from"`
	Content string `json:"content"`
}
type reply struct {
	From    string `json:"from"`
	Content string `json:"content"`
	Time    int64  `json:"time"`
}
type rawMessage struct {
	Author  string  `json:"author"`
	Content string  `json:"content"`
	Replies []reply `json:"replies"`
}
type message struct {
	Author  string  `json:"author"`
	Content string  `json:"content"`
	Time    int64   `json:"time"`
	Replies []reply `json:"replies"`
}

func handle(conn net.Conn, req request) error {
	resp := ""
	resp = setHeader(resp, "HTTP/1.1 200 OK")
	sid := req.getCookie(cookieKey)
	if req.Path == "/" {
		resp = setHeader(resp, "Content-Type: text/html")
		if sid == "" {
			resp = setBody(resp, readHTMLFile(redirectLoginPage))
			_, err := conn.Write([]byte(resp))
			if err != nil {
				log.Print("failed to write response contents")
				return err
			}
		} else {
			_, err := rdb.Get(ctx, sid).Result()
			if err == redis.Nil {
				resp = setBody(resp, readHTMLFile(redirectLoginPage))
				_, err := conn.Write([]byte(resp))
				if err != nil {
					log.Print("failed to write response contents")
					return err
				}
			} else if err == nil {
				resp = setBody(resp, readHTMLFile(redirectHomePage))
				_, err := conn.Write([]byte(resp))
				if err != nil {
					log.Print("failed to write response contents")
					return err
				}
			} else {
				resp = setBody(resp, readHTMLFile(redirectLoginPage))
				_, err := conn.Write([]byte(resp))
				if err != nil {
					log.Print("failed to write response contents")
					return err
				}
			}
		}
	} else if req.Path == "/home" {
		resp = setHeader(resp, "Content-Type: text/html")
		if sid == "" {
			resp = setBody(resp, readHTMLFile(redirectLoginPage))
			_, err := conn.Write([]byte(resp))
			if err != nil {
				log.Print("failed to write response contents")
				return err
			}
		} else {
			_, err := rdb.Get(ctx, sid).Result()
			if err == redis.Nil {
				resp = setBody(resp, readHTMLFile(redirectLoginPage))
				_, err := conn.Write([]byte(resp))
				if err != nil {
					log.Print("failed to write response contents")
					return err
				}
			} else if err == nil {
				resp = setBody(resp, readHTMLFile(homePage))
				_, err := conn.Write([]byte(resp))
				if err != nil {
					log.Print("failed to write response contents")
					return err
				}
			} else {
				resp = setBody(resp, readHTMLFile(redirectLoginPage))
				_, err := conn.Write([]byte(resp))
				if err != nil {
					log.Print("failed to write response contents")
					return err
				}
			}
		}
	} else if req.Path == "/logout" {
		resp = setHeader(resp, "Content-Type: text/html")
		rdb.Del(ctx, sid)
		resp = setCookie(resp, cookieKey, "")
		resp = setBody(resp, readHTMLFile(redirectLoginPage))
		_, err := conn.Write([]byte(resp))
		if err != nil {
			log.Print("failed to write response contents")
			return err
		}
	} else if strings.HasPrefix(req.Path, "/login") {
		resp = setHeader(resp, "Content-Type: text/html")
		if req.Path == "/login" || req.Path == "/login/" {
			resp = setBody(resp, readHTMLFile(loginPage))
			_, err := conn.Write([]byte(resp))
			if err != nil {
				log.Print("failed to write response contents")
				return err
			}
		} else {
			if req.Path == "/login?status=lf" {
				resp = setBody(resp, readHTMLFile(loginFailPage))
				_, err := conn.Write([]byte(resp))
				if err != nil {
					log.Print("failed to write response contents")
					return err
				}
			} else if req.Path == "/login?status=rf" {
				resp = setBody(resp, readHTMLFile(registerFailPage))
				_, err := conn.Write([]byte(resp))
				if err != nil {
					log.Print("failed to write response contents")
					return err
				}
			} else {
				resp = setBody(resp, "<h1>404 not found<h1>")
				_, err := conn.Write([]byte(resp))
				if err != nil {
					log.Print("failed to write response contents")
					return err
				}
			}
		}
	} else if strings.HasPrefix(req.Path, "/signin") {
		resp = setHeader(resp, "Content-Type: text/html")
		u, _ := url.Parse(req.Path)
		m, _ := url.ParseQuery(u.RawQuery)
		if len(m["name"]) > 0 && len(m["password"]) > 0 {
			pswd, err := rdb.Get(ctx, m["name"][0]).Result()
			if err == redis.Nil {
				resp = setBody(resp, readHTMLFile(redirectLoginFailPage))
				_, err := conn.Write([]byte(resp))
				if err != nil {
					log.Print("failed to write response contents")
					return err
				}
			} else if err == nil {
				if pswd == m["password"][0] {
					newSID := betterguid.New()
					rdb.Set(ctx, newSID, m["name"][0], 0)
					resp = setCookie(resp, cookieKey, newSID)
					resp = setBody(resp, readHTMLFile(redirectHomePage))
					_, err := conn.Write([]byte(resp))
					if err != nil {
						log.Print("failed to write response contents")
						return err
					}
				} else {
					resp = setBody(resp, readHTMLFile(redirectLoginFailPage))
					_, err := conn.Write([]byte(resp))
					if err != nil {
						log.Print("failed to write response contents")
						return err
					}
				}
			} else {
				resp = setBody(resp, readHTMLFile(redirectLoginFailPage))
				_, err := conn.Write([]byte(resp))
				if err != nil {
					log.Print("failed to write response contents")
					return err
				}
			}
		} else {
			resp = setBody(resp, readHTMLFile(redirectLoginFailPage))
			_, err := conn.Write([]byte(resp))
			if err != nil {
				log.Print("failed to write response contents")
				return err
			}
		}
	} else if strings.HasPrefix(req.Path, "/signup") {
		resp = setHeader(resp, "Content-Type: text/html")
		u, _ := url.Parse(req.Path)
		m, _ := url.ParseQuery(u.RawQuery)
		if len(m["name"]) > 0 && len(m["password"]) > 0 {
			if len(m["name"][0]) > 0 && len(m["password"][0]) > 0 {
				_, err := rdb.Get(ctx, m["name"][0]).Result()
				if err == redis.Nil {
					newSID := betterguid.New()
					rdb.Set(ctx, m["name"][0], m["password"][0], 0)
					rdb.Set(ctx, newSID, m["name"][0], 0)
					resp = setCookie(resp, cookieKey, newSID)
					resp = setBody(resp, readHTMLFile(redirectHomePage))
					_, err := conn.Write([]byte(resp))
					if err != nil {
						log.Print("failed to write response contents")
						return err
					}

				} else if err == nil {
					resp = setBody(resp, readHTMLFile(redirectRegisterFailPage))
					_, err := conn.Write([]byte(resp))
					if err != nil {
						log.Print("failed to write response contents")
						return err
					}
				} else {
					resp = setBody(resp, readHTMLFile(redirectRegisterFailPage))
					_, err := conn.Write([]byte(resp))
					if err != nil {
						log.Print("failed to write response contents")
						return err
					}
				}
			} else {
				resp = setBody(resp, readHTMLFile(redirectRegisterFailPage))
				_, err := conn.Write([]byte(resp))
				if err != nil {
					log.Print("failed to write response contents")
					return err
				}
			}
		} else {
			resp = setBody(resp, readHTMLFile(redirectRegisterFailPage))
			_, err := conn.Write([]byte(resp))
			if err != nil {
				log.Print("failed to write response contents")
				return err
			}
		}
	} else if req.Path == "/api/message" {
		resp = setHeader(resp, "Content-Type: application/json")
		l, _ := strconv.Atoi(req.Header["Content-Length"])

		res := rawMessage{}
		json.Unmarshal([]byte(req.Body[0:l]), &res)
		final := message{
			Author:  res.Author,
			Content: res.Content,
			Time:    time.Now().Unix(),
			Replies: res.Replies,
		}
		v, _ := json.Marshal(final)
		rdb.RPush(ctx, "msgdata", string(v))

		mp := map[string]string{"msg": "ok"}
		jsonReply, _ := json.Marshal(mp)
		resp = setBody(resp, string(jsonReply))
		_, err := conn.Write([]byte(resp))
		if err != nil {
			log.Print("failed to write response contents")
			return err
		}
	} else if req.Path == "/api/messages" {
		resp = setHeader(resp, "Content-Type: application/json")
		len, _ := rdb.LLen(ctx, "msgdata").Result()
		msgs, _ := rdb.LRange(ctx, "msgdata", 0, len-1).Result()

		var final []message
		for _, msg := range msgs {
			var res message
			json.Unmarshal([]byte(msg), &res)
			final = append(final, res)
		}
		jsonReply, _ := json.Marshal(final)
		resp = setBody(resp, string(jsonReply))
		_, err := conn.Write([]byte(resp))
		if err != nil {
			log.Print("failed to write response contents")
			return err
		}
	} else if strings.HasPrefix(req.Path, "/api/reply") {
		u, _ := url.Parse(req.Path)
		m, _ := url.ParseQuery(u.RawQuery)
		if len(m["msgid"]) > 0 {
			idx, _ := strconv.ParseInt(m["msgid"][0], 10, 64)
			if req.Method == "POST" {
				l, _ := strconv.Atoi(req.Header["Content-Length"])
				res := rawReply{}
				e := json.Unmarshal([]byte(req.Body[0:l]), &res)
				fmt.Println(e)
				log.Println(res)

				original, _ := rdb.LIndex(ctx, "msgdata", idx).Result()
				var res2 message
				json.Unmarshal([]byte(original), &res2)
				res2.Replies = append(res2.Replies, reply{
					From:    res.From,
					Content: res.Content,
					Time:    time.Now().Unix(),
				})

				v, _ := json.Marshal(res2)

				rdb.LSet(ctx, "msgdata", idx, string(v))

				mp := map[string]string{"msg": "ok"}
				jsonReply, _ := json.Marshal(mp)
				resp = setBody(resp, string(jsonReply))
				_, err := conn.Write([]byte(resp))
				if err != nil {
					log.Print("failed to write response contents")
					return err
				}
			} else {
				msg, _ := rdb.LIndex(ctx, "msgdata", idx).Result()
				var res message
				json.Unmarshal([]byte(msg), &res)
				jsonReply, _ := json.Marshal(res.Replies)
				resp = setBody(resp, string(jsonReply))
				_, err := conn.Write([]byte(resp))
				if err != nil {
					log.Print("failed to write response contents")
					return err
				}
			}
		}
	} else if strings.HasPrefix(req.Path, "/api/username") {
		resp = setHeader(resp, "Content-Type: application/json")
		u, _ := url.Parse(req.Path)
		m, _ := url.ParseQuery(u.RawQuery)
		if len(m["sid"]) > 0 {
			username, _ := rdb.Get(ctx, m["sid"][0]).Result()
			mp := map[string]string{"username": username}
			jsonReply, _ := json.Marshal(mp)
			resp = setBody(resp, string(jsonReply))
			_, err := conn.Write([]byte(resp))
			if err != nil {
				log.Print("failed to write response contents")
				return err
			}
		} else {
			resp = setBody(resp, "<h1>404 not found<h1>")
			_, err := conn.Write([]byte(resp))
			if err != nil {
				log.Print("failed to write response contents")
				return err
			}
		}
	} else {
		resp = setBody(resp, "<h1>404 not found<h1>")
		_, err := conn.Write([]byte(resp))
		if err != nil {
			log.Print("failed to write response contents")
			return err
		}
	}
	return nil
}
