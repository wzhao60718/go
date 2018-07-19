// appserver.go
package main

import (
  "encoding/json"
  "fmt"
  "net/http"
	"os"
	"strconv"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type LocalConfig struct {
	MongoDbURL    string `json:"mongodburl"`
	MongoDbName   string `json:"mongodbname"`
	HttpServerDir string `json:"httpserverdir"`
}

type Blog struct {
  ID           int64  `json:"id,string"`           //format: yyyymmdd
  Author       string `json:"author"`
  Headline     string `json:"headline"`
  Subtitle     string `json:"subtitle"`
  Introduction string `json:"introduction"`
  Body         string `json:"body"`
  Photo        string `json:"photo"`
  Video        string `json:"video"`
  CreatedDate  string `json:"createddate"` // format yyyy-mm-dd
  UpdatedDate  string `json:"updateddate"` // format yyyy-mm-dd
  IsDraft      bool   `json:"isdraft"`
}

type Photo struct {
  ID           int64      `json:"id,string"`           //format: yyyymmdd
  Author       string     `json:"author"`
  Headline     string     `json:"headline"`
  Subtitle     string     `json:"subtitle"`
  Introduction string     `json:"introduction"`
  Body         string     `json:"body"`
  Photos       []PhotoOne `json:"photos"`
  CreatedDate  string     `json:"createddate"` // format yyyy-mm-dd
  UpdatedDate  string     `json:"updateddate"` // format yyyy-mm-dd
  IsDraft      bool       `json:"isdraft"`
}

type PhotoOne struct {
  ID          string `json:"id"`
  Photo       string `json:"photo"`
  Thumbnail   string `json:"thumbnail"`
  Description string `json:"description"`
}

var local_config LocalConfig

func main() {

	configFile, err := os.Open("config.json")
  if err != nil {
	  fmt.Println("open config file", err)
  }
	defer configFile.Close()
	
	jsonParser := json.NewDecoder(configFile)
	if err = jsonParser.Decode(&local_config); err != nil {
			fmt.Println("parsing config file", err)
	}
	
	session, err := mgo.Dial(local_config.MongoDbURL)
	if err != nil {
    panic(err)
	}
	defer session.Close()

	http.HandleFunc("/json", func(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Access-Control-Allow-Origin", "*")
    w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT")
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Accept", "application/json")
		
		actions, a_ok := r.URL.Query()["a"]    
		collections, c_ok := r.URL.Query()["c"]    
		ids, id_ok := r.URL.Query()["id"]    

    if !a_ok || len(actions[0]) < 1 {
			json.NewEncoder(w).Encode("[{}]")
    } else if(actions[0] == "get") {
      if !c_ok || len(collections[0]) < 1 {
				json.NewEncoder(w).Encode("{}")
			} else if(collections[0] == "blog") {
  			if !id_ok || len(ids[0]) < 1 {
					blogs, err := GetBlogs(session)
					if err != nil {
						fmt.Println("Err:", err)
					}
					json.NewEncoder(w).Encode(blogs)
				} else {
					b_id, b_id_err := strconv.ParseInt(ids[0], 10, 64)
					if b_id_err != nil {
							panic(b_id_err)
					}
					blog, err := GetBlog(session, b_id)
					if err != nil {
						fmt.Println("Err:", err)
					}
					json.NewEncoder(w).Encode(blog)
				}
			} else if(collections[0] == "photo") {
  			if !id_ok || len(ids[0]) < 1 {
					photos, err := GetPhotos(session)
					if err != nil {
						fmt.Println("Err:", err)
					}
					json.NewEncoder(w).Encode(photos)
				} else {
					p_id, p_id_err := strconv.ParseInt(ids[0], 10, 64)
					if p_id_err != nil {
							panic(p_id_err)
					}
					photo, err := GetPhoto(session, p_id)
					if err != nil {
						fmt.Println("Err:", err)
					}	
					json.NewEncoder(w).Encode(photo)
				}
			}
		} else if(actions[0] == "add") {
      if !c_ok || len(collections[0]) < 1 {
				json.NewEncoder(w).Encode("{}")
			} else if(collections[0] == "blog") {
				decoder := json.NewDecoder(r.Body)
        var blog Blog
				err := decoder.Decode(&blog)
				if err != nil {
					panic(err)
				}
				err = AddBlog(session, &blog)
				if err != nil {
					panic(err)
				}
				json.NewEncoder(w).Encode(&blog)
			} else if(collections[0] == "photo") {
				decoder := json.NewDecoder(r.Body)
        var photo Photo
				err := decoder.Decode(&photo)
				if err != nil {
					panic(err)
				}
				err = AddPhoto(session, &photo)
				if err != nil {
					panic(err)
				}
				json.NewEncoder(w).Encode(&photo)
			}
    } else if(actions[0] == "update") {
      if !c_ok || len(collections[0]) < 1 {
				json.NewEncoder(w).Encode("{}")
			} else if(collections[0] == "blog") {
				decoder := json.NewDecoder(r.Body)
        var blog Blog
				err := decoder.Decode(&blog)
				if err != nil {
					panic(err)
				}
				err = UpdateBlog(session, &blog)
				if err != nil {
					panic(err)
				}
				json.NewEncoder(w).Encode(&blog)
			} else if(collections[0] == "photo") {
				decoder := json.NewDecoder(r.Body)
        var photo Photo
				err := decoder.Decode(&photo)
				if err != nil {
					panic(err)
				}
				err = UpdatePhoto(session, &photo)
				if err != nil {
					panic(err)
				}
				json.NewEncoder(w).Encode(&photo)
			}
    } else if(actions[0] == "delete") {
      if !c_ok || len(collections[0]) < 1 {
				json.NewEncoder(w).Encode("{}")
			} else if(collections[0] == "blog") {
				decoder := json.NewDecoder(r.Body)
        var blog Blog
				err := decoder.Decode(&blog)
				if err != nil {
					panic(err)
				}
				err = DeleteBlog(session, &blog)
				if err != nil {
					panic(err)
				}
				json.NewEncoder(w).Encode(&blog)
			} else if(collections[0] == "photo") {
				decoder := json.NewDecoder(r.Body)
        var photo Photo
				err := decoder.Decode(&photo)
				if err != nil {
					panic(err)
				}
				err = DeletePhoto(session, &photo)
				if err != nil {
					panic(err)
				}
				json.NewEncoder(w).Encode(&photo)
			}
    }
  })

	http.HandleFunc("/welcome", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, " website, J & W !")
  })
	
	fs := http.FileServer(http.Dir(local_config.HttpServerDir))
	http.Handle("/", http.StripPrefix("/", fs))
	
  http.ListenAndServe(":8080", nil)
}

func GetBlogs(s *mgo.Session) (*[]Blog, error) {
	var blogs []Blog
	err := s.DB(local_config.MongoDbName).C("Blogs").Find(bson.M{}).All(&blogs)
	if err != nil {
		return nil, err
	}
	return &blogs, err
}

func GetBlog(s *mgo.Session, id int64) (*Blog, error) {
	var blog Blog
	err := s.DB(local_config.MongoDbName).C("Blogs").Find(bson.M{"id": id}).One(&blog)
	if err != nil {
		return nil, err
	}
	return &blog, err
}

func AddBlog(s *mgo.Session, b *Blog) error {
	return s.DB(local_config.MongoDbName).C("Blogs").Insert(b)
}

func UpdateBlog(s *mgo.Session, b *Blog) error {
	return s.DB(local_config.MongoDbName).C("Blogs").Update(bson.M{"id": b.ID}, b)	
}

func DeleteBlog(s *mgo.Session, b *Blog) error {
	return DeleteBlogByID(s, b.ID)	
}

func DeleteBlogByID(s *mgo.Session, b_id int64) error {
	return s.DB(local_config.MongoDbName).C("Blogs").Remove(bson.M{"id": b_id})	
}

func GetPhotos(s *mgo.Session) (*[]Photo, error) {
	var photos []Photo
	err := s.DB(local_config.MongoDbName).C("Photos").Find(bson.M{}).All(&photos)
	if err != nil {
		return nil, err
	}
	return &photos, err
}

func GetPhoto(s *mgo.Session, id int64) (*Photo, error) {
	var photo Photo
	err := s.DB(local_config.MongoDbName).C("Photos").Find(bson.M{"id": id}).One(&photo)
	if err != nil {
		return nil, err
	}
	return &photo, err
}

func AddPhoto(s *mgo.Session, p *Photo) error {
	return s.DB(local_config.MongoDbName).C("Photos").Insert(p)
}

func UpdatePhoto(s *mgo.Session, p *Photo) error {
	return s.DB(local_config.MongoDbName).C("Photos").Update(bson.M{"id": p.ID}, p)	
}

func DeletePhoto(s *mgo.Session, p *Photo) error {
	return DeletePhotoByID(s, p.ID)
}

func DeletePhotoByID(s *mgo.Session, p_id int64) error {
	return s.DB(local_config.MongoDbName).C("Photos").Remove(bson.M{"id": p_id})	
}
