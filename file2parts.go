package main

import (
	"database/sql"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-fsnotify/fsnotify"
	_ "github.com/go-sql-driver/mysql"
	redis "gopkg.in/redis.v5"
)

type config struct {
	FileName string
	Settings settings
}

type settings struct {
	MySQL mySQLSettings `json:"mysql"`
	Redis redisSettings `json:"redis"`
}

type mySQLSettings struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	DBName   string `json:"db_name"`
}

type redisSettings struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Password string `json:"password"`
	DB       int    `json:"db"`
}

type parts struct {
	Path string
	HTML string
}

const (
	partExt     = ".volt"
	redisPrefix = "daisy-cms:parts-html-"
)

var (
	configFileName  = flag.String("config", "file2parts.json", "config file name")
	dump            = flag.Bool("d", false, "dump parts that exist as files from database")
	dumpAll         = flag.Bool("da", false, "dump all parts from database")
	restore         = flag.Bool("r", false, "restore parts to database")
	watch           = flag.Bool("w", false, "watch and restore modified part file")
	restoreAndWatch = flag.Bool("rw", false, "args r and w")
)

// start endpoint
func main() {
	flag.Usage = func() {
		fmt.Fprint(os.Stderr, `Usage of file2parts:
  -config: config file name
  -d:      dump parts that exist as files from database
  -da:     dump all parts from database
  -r:      restore parts to database
  -w:      watch and restore modified part file
  -rw:     alias of "-r -w"
`)
	}
	flag.Parse()

	conf, err := loadConfig()
	if err != nil {
		log.Fatal("failed to load configuration file: ", err)
	}

	err = saveConfig(conf)
	if err != nil {
		log.Fatal("failed to save configuration file: ", err)
	}

	dir, err := os.Getwd()
	if err != nil {
		log.Fatal("failed to get current directory: ", err)
	}

	if *dump {
		db, err := getDB(conf.Settings.MySQL)
		if err != nil {
			log.Fatal("failed to connect MySQL: ", err)
		}

		defer db.Close()

		ps, err := fetchAllParts(db)
		if err != nil {
			log.Fatal("failed to connect MySQL: ", err)
		}

		err = saveExistingParts(dir, ps)
		if err != nil {
			log.Fatal("failed to save parts: ", err)
		}
	} else if *dumpAll {
		db, err := getDB(conf.Settings.MySQL)
		if err != nil {
			log.Fatal("failed to connect MySQL: ", err)
		}

		defer db.Close()

		ps, err := fetchAllParts(db)
		if err != nil {
			log.Fatal("failed to connect MySQL: ", err)
		}

		err = saveAllParts(dir, ps)
		if err != nil {
			log.Fatal("failed to save parts: ", err)
		}
	} else if *restore || *restoreAndWatch || *watch {
		db, err := getDB(conf.Settings.MySQL)
		if err != nil {
			log.Fatal("failed to connect MySQL: ", err)
		}

		defer db.Close()

		rc := getRedis(conf.Settings.Redis)
		_, err = rc.Ping().Result()
		if err != nil {
			log.Fatal("failed to connect redis: ", err)
		}
		defer rc.Close()

		if *restore || *restoreAndWatch {
			err = restorePartsFiles(db, rc, dir)
			if err != nil {
				log.Fatal("failed to restore parts: ", err)
			}
		}
		if *watch || *restoreAndWatch {
			err = watchParts(db, rc, dir)
			if err != nil {
				log.Fatal("failed to restore parts: ", err)
			}
		}
	} else {
		flag.Usage()
	}
}

// create default configuration
func createConfig() *config {
	return &config{
		FileName: *configFileName,
		Settings: settings{
			MySQL: mySQLSettings{
				Host:     "127.0.0.1",
				Port:     3306,
				User:     "daisy",
				Password: "team-lab",
				DBName:   "daisy_cms",
			},
			Redis: redisSettings{
				Host:     "127.0.0.1",
				Port:     6379,
				Password: "",
				DB:       0,
			},
		},
	}
}

// loading configuration file
func loadConfig() (*config, error) {
	conf := createConfig()
	filename := conf.FileName

	b, err := ioutil.ReadFile(filename)
	if err != nil && !os.IsNotExist(err) {
		return nil, err
	}

	if b != nil {
		err = json.Unmarshal(b, &conf.Settings)
		if err != nil {
			return nil, fmt.Errorf("could not deconde json: %v", err)
		}
	}
	return conf, nil
}

// save configuration file
func saveConfig(conf *config) error {
	filename := conf.FileName
	b, err := json.MarshalIndent(conf.Settings, "", "\t")
	if err != nil {
		return err
	}
	return ioutil.WriteFile(filename, b, 0700)
}

// create db
func getDB(mysql mySQLSettings) (*sql.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", mysql.User, mysql.Password, mysql.Host, mysql.Port, mysql.DBName)
	return sql.Open("mysql", dsn)
}

// create redis
func getRedis(rs redisSettings) *redis.Client {
	addr := fmt.Sprintf("%s:%d", rs.Host, rs.Port)
	return redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: rs.Password,
		DB:       rs.DB,
	})
}

// request all parts from DB
func fetchAllParts(db *sql.DB) ([]parts, error) {
	rows, err := db.Query("SELECT path, html FROM parts")
	if err != nil {
		return nil, fmt.Errorf("failed to request query: %v", err)
	}
	ps := make([]parts, 0)
	for rows.Next() {
		p := parts{}
		err = rows.Scan(&p.Path, &p.HTML)
		if err != nil {
			return nil, fmt.Errorf("failed to scan response: %v", err)
		}
		ps = append(ps, p)
	}
	return ps, nil
}

func part2file(filename string, p parts) error {
	dir := filepath.Dir(filename)
	err := os.MkdirAll(dir, 0700)
	if err != nil {
		return fmt.Errorf("failed to save parts file: %v", err)
	}
	b := []byte(p.HTML)
	return ioutil.WriteFile(filename, b, 0700)
}

//　save parts file ps to dir
func saveAllParts(dir string, ps []parts) error {
	for _, p := range ps {
		filename := filepath.Join(dir, p.Path+".volt")
		err := part2file(filename, p)
		if err != nil {
			return fmt.Errorf("failed to save parts file \"%s\": %v", filename, err)
		}
	}
	return nil
}

func saveExistingParts(dir string, ps []parts) error {
	for _, p := range ps {
		filename := filepath.Join(dir, p.Path+".volt")
		_, err := os.Stat(filename)
		if err == nil {
			err = part2file(filename, p)
			if err != nil {
				return fmt.Errorf("failed to save parts file \"%s\": %v", filename, err)
			}
		}
	}
	return nil
}

func restorePartsFiles(db *sql.DB, rc *redis.Client, dir string) error {
	files, err := findFileParts(dir)
	if err != nil {
		return err
	}

	err = restoreParts(db, rc, dir, files)
	if err != nil {
		return err
	}

	return nil
}

// get all volt file in dir
func findFileParts(dir string) ([]string, error) {
	ps := make([]string, 0)
	fis, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("failed to read file or directory \"%s\": %v", dir, err)
	}
	for _, fi := range fis {
		fn := fi.Name()
		filename := filepath.Join(dir, fn)
		// ignore hidden directory
		if fn[:1] == "." {
			continue
		}

		if fi.IsDir() {
			// find subdirectory
			dps, err := findFileParts(filename)
			if err != nil {
				return nil, err
			}
			ps = append(ps, dps...)
		} else {
			ext := filepath.Ext(filename)
			// get volt file only
			if ext != partExt {
				continue
			}

			ps = append(ps, filename)
		}
	}

	return ps, nil
}

// restore parts from part file list
func restoreParts(db *sql.DB, rc *redis.Client, dir string, partsfiles []string) error {
	fmt.Println(len(partsfiles))
	for _, file := range partsfiles {
		p, err := file2Parts(dir, file)
		if err != nil {
			return err
		}
		fmt.Println(file)

		err = updateParts(db, rc, p)
		if err != nil {
			return err
		}
	}
	return nil
}

// make parts from files
func file2Parts(dir string, file string) (*parts, error) {
	p := &parts{}
	if len(file) < len(dir) && file[:len(dir)] != dir {
		return nil, fmt.Errorf("invaild file path")
	}
	partsPath := file[len(dir)+1:]

	if len(dir) < len(partExt) && partsPath[:len(partExt)] != partExt {
		return nil, fmt.Errorf("invaild filename")
	}

	pp := partsPath[:(len(partsPath) - len(partExt))]
	// windows directory delimiter is "\"
	p.Path = strings.Join(filepath.SplitList(pp), "/")

	b, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, fmt.Errorf("cannot read file %v: %v", file, err)
	}

	p.HTML = string(b)
	return p, nil
}

// update database to parts
func updateParts(db *sql.DB, rc *redis.Client, p *parts) error {
	_, err := db.Exec("UPDATE parts SET path = ?, html = ? WHERE path = ?", p.Path, p.HTML, p.Path)
	if err != nil {
		return fmt.Errorf("failed to update parts: %v", err)
	}
	fmt.Println(p.Path)

	ic := rc.Del(redisPrefix + p.Path)
	if err = ic.Err(); err != nil {
		return fmt.Errorf("failed to clear redis cache: %v", err)
	}
	return nil
}

// watch and restore
func watchParts(db *sql.DB, rc *redis.Client, dir string) error {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return fmt.Errorf("failed to watch parts files: %v", err)
	}
	defer watcher.Close()

	done := make(chan bool)
	go func() {
		for {
			select {
			case event := <-watcher.Events:
				if event.Op&fsnotify.Write == fsnotify.Write {
					log.Println("Modified file: ", event.Name)
					ok, err := isWatchFile(dir, event.Name)
					if err != nil {
						log.Printf("failed to update parts: %v", err)
						done <- true
						return
					}
					if ok {
						p, err := file2Parts(dir, event.Name)
						if err != nil {
							log.Println("failed to load parts")
							done <- true
							return
						}
						err = updateParts(db, rc, p)
						if err != nil {
							log.Printf("failed to update parts: %v", err)
							done <- true
							return
						}
					}
				}
			case err := <-watcher.Errors:
				log.Printf("error: %v", err)
				done <- true
				return
			}
		}
	}()

	err = watcher.Add(dir)
	if err != nil {
		return fmt.Errorf("cannot watch directory: %v", err)
	}

	fmt.Println("watching midifed parts files...")
	<-done

	return nil
}

func isWatchFile(dir string, file string) (bool, error) {
	if len(file) < len(dir) && file[:len(dir)] != dir {
		return false, fmt.Errorf("invaild file path")
	}
	partsPath := file[len(dir)+1:]

	if len(dir) < len(partExt) && partsPath[:len(partExt)] != partExt {
		return false, nil
	}

	dirs := filepath.SplitList(partsPath)
	for _, d := range dirs {
		if d[:1] == "." {
			return false, nil
		}
	}
	return true, nil
}
