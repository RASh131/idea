package config

import (
    "errors"
    "idea/crypt"
    "log"

    "database/sql"

    _ "github.com/go-sql-driver/mysql"
)

func NewUser(username, email string) *UserModel {
    return &UserModel{
        Username: username,
        Email:    email,
    }
}

type UserModel struct {
    Username      string
    Password      string
    Email         string
    authenticated bool
}

func (u *UserModel) SetPassword(password string) error {
    hash, err := crypt.PasswordEncrypt(password)
    if err != nil {
        return err
    }
    u.Password = hash
    return nil
}

func (u *UserModel) Authenticate() {
    u.authenticated = true
}

func Exists(username string) bool {
    db, err := sql.Open("mysql", "root:root@tcp(127.0.0.1:8889)/idea")
    log.Println("Connected to mysql")
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()
    rows, err := db.Query("SELECT * FROM user WHERE name=\"" + username + "\"")
    defer rows.Close()
    if err != nil {
        panic(err.Error())
    }
    var user UserModel
    for rows.Next() {
        err := rows.Scan(&user.Username, &user.Email, &user.Password)
        if err != nil {
            panic(err.Error())
        }
    }
    println("test ; " + user.Username)
    if user.Username != "" {
        return true
    }
    return false
}

func SaveUser(username, email, password string) error {
    if Exists(username) {
        return errors.New("user \"" + username + "\" already exists")
    }
    hash, err := crypt.PasswordEncrypt(password)
    if err != nil {
        return err
    }
    db, err := sql.Open("mysql", "root:root@tcp(127.0.0.1:8889)/idea")
    log.Println("Connected to mysql")
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()
    insert, err := db.Prepare("INSERT INTO user(name,mail,password) VALUE(?,?,?)")
    if err != nil {
        log.Fatal(err)
    }
    insert.Exec(username, email, hash)
    return nil
}
func GetUser(username, password string) (*UserModel, error) {
    var user UserModel
    if !Exists(username) {
        return nil, errors.New("user \"" + username + "\" doesn't exists")
    }
    db, err := sql.Open("mysql", "root:root@tcp(127.0.0.1:8889)/idea")
    log.Println("Connected to mysql")
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()
    rows, err := db.Query("SELECT * FROM user WHERE name= \"" + username + "\"")
    defer rows.Close()
    if err != nil {
        panic(err.Error())
    }

    for rows.Next() {
        err := rows.Scan(&user.Username, &user.Email, &user.Password)
        if err != nil {
            panic(err.Error())
        }
        if err := crypt.CompareHashAndPassword(user.Password, password); err != nil {
            println("user \"" + username + "\",password \"" + password + " パスワード間違っとるよ")
        }
    }
    return &user, nil
}
