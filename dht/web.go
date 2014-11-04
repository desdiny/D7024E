package dht

import (
	//"bufio"
	"fmt"
	"net/http"
	//"os"
	"github.com/boltdb/bolt"
	"log"
	"time"
)

func Chord(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "<h1>Upload a new key-value pair!</h1>"+
		"<form action=\"/chord/post/\" method=\"POST\">"+
		"Key:"+
		"<textarea name=\"Post_insertkey\"></textarea><br>"+
		"Value:"+
		"<textarea name=\"Post_insertvalue\"></textarea><br>"+
		"<input type=\"submit\" value=\"Submit\">"+
		"</form>"+

		"<h1>Returns the value for a specific key!</h1>"+
		"<form action=\"/chord/get/\" method=\"POST\">"+
		"Key:"+
		"<textarea name=\"Get_insertkey\"></textarea><br>"+
		"<input type=\"submit\" value=\"Submit\">"+
		"</form>"+

		"<h1>Update the value for a specific key!</h1>"+
		"<form action=\"/chord/put/\" method=\"POST\">"+
		"Key:"+
		"<textarea name=\"Put_insertkey\"></textarea><br>"+
		"Value:"+
		"<textarea name=\"Put_insertvalue\"></textarea><br>"+
		"<input type=\"submit\" value=\"Submit\">"+
		"</form>"+

		"<h1>Delete a key-value pair with key!</h1>"+
		"<form action=\"/chord/delete/\" method=\"POST\">"+
		"Key:"+
		"<textarea name=\"Delete_insertkey\"></textarea><br>"+
		"<input type=\"submit\" value=\"Submit\">"+

		"</form>")
}

func Post(w http.ResponseWriter, r *http.Request) {
	db, err := bolt.Open("node.db", 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		log.Fatal(err)
	}

	key := r.FormValue("Post_insertkey")
	value := r.FormValue("Post_insertvalue")

	//fmt.Fprintf(w, "Print key-value post: ", key, value)
	db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("KeyValuePair"))
		v := b.Get([]byte(key))

		if v != nil {

			fmt.Fprintf(w, "<p><a href=\"/chord/\">go back</a></p>"+
				"<p>The key: "+key+" has already been used</p>")
			return nil
		} else {

			db.Update(func(tx *bolt.Tx) error {
				b, err := tx.CreateBucketIfNotExists([]byte("KeyValuePair"))
				if err != nil {
					return fmt.Errorf("create bucket: ", err)
				}

				b = tx.Bucket([]byte("KeyValuePair"))
				err = b.Put([]byte(key), []byte(value))
				fmt.Fprintf(w, "<p><a href=\"/chord/\">go back</a></p>"+
					"<p>The key: "+key+" and value: "+value+" has been upladed </p>")
				return err
			})

		}

		return nil
	})

	defer db.Close()

}

func List(w http.ResponseWriter, r *http.Request) {
	db, err := bolt.Open("node.db", 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("KeyValuePair"))
		b.ForEach(func(k, v []byte) error {
			fmt.Fprintf(w, "key=%s, value=%s\n", k, v)
			return nil
		})
		return nil
	})

}

func Get(w http.ResponseWriter, r *http.Request) {
	db, err := bolt.Open("node.db", 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	key := r.FormValue("Get_insertkey")
	fmt.Fprintf(w, "Print something get: ", key)

	db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("KeyValuePair"))
		v := b.Get([]byte(key))
		fmt.Fprintf(w, "The answer is: %s", v)
		return nil
	})

}

func Put(w http.ResponseWriter, r *http.Request) {
	db, err := bolt.Open("node.db", 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	key := r.FormValue("Put_insertkey")
	value := r.FormValue("Put_insertvalue")

	db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("KeyValuePair"))
		v := b.Get([]byte(key))

		if v != nil {

			db.Update(func(tx *bolt.Tx) error {
				b := tx.Bucket([]byte("KeyValuePair"))
				b.Delete([]byte(key))
				//err = b.Put([]byte(key), []byte(value))
				return nil
			})

			db.Update(func(tx *bolt.Tx) error {
				b, err := tx.CreateBucketIfNotExists([]byte("KeyValuePair"))
				if err != nil {
					return fmt.Errorf("create bucket: ", err)
				}

				b = tx.Bucket([]byte("KeyValuePair"))
				err = b.Put([]byte(key), []byte(value))
				fmt.Fprintf(w, "<p><a href=\"/chord/\">go back</a></p>"+
					"<p>The key: "+key+" and value: "+value+" has been updated </p>")
				return err
			})

		} else {

			fmt.Fprintf(w, "<p><a href=\"/chord/\">go back</a></p>"+
				"<p>There are no key :%s to remove</p>", key)
			return nil

		}

		return nil

	})

}

func Del(w http.ResponseWriter, r *http.Request) {
	db, err := bolt.Open("node.db", 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	key := r.FormValue("Delete_insertkey")

	db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("KeyValuePair"))
		v := b.Get([]byte(key))

		if v != nil {

			db.Update(func(tx *bolt.Tx) error {
				b := tx.Bucket([]byte("KeyValuePair"))
				b.Delete([]byte(key))
				//err = b.Put([]byte(key), []byte(value))
				fmt.Fprintf(w, "<p><a href=\"/chord/\">go back</a></p>"+
					"<p>The key:%s is removed </p>", key)
				return err
			})

		} else {

			fmt.Fprintf(w, "<p><a href=\"/chord/\">go back</a></p>"+
				"<p>There are no key :%s to remove</p>", key)
			return nil

		}

		return nil

	})

}
