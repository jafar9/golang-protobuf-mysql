package main

import (
	"bufio"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/golang/protobuf/proto"
	"io"
	"log"
	"os"
	pb "protos"
	"strings"
)

func promptForAddress(r io.Reader) (*pb.Person, error) {
	p := &pb.Person{}

	rd := bufio.NewReader(r)
	fmt.Print("Enter person ID number: ")
	if _, err := fmt.Fscanf(rd, "%d\n", &p.Id); err != nil {
		return p, err
	}

	fmt.Print("Enter name: ")
	name, err := rd.ReadString('\n')
	if err != nil {
		return p, err
	}
	p.Name = strings.TrimSpace(name)

	phone := &pb.PhoneNumber{}

	fmt.Print("Enter phone number: ")
	email, err := rd.ReadString('\n')
	if err != nil {
		return p, err
	}
	phone.Number = strings.TrimSpace(email)

	address := &pb.Address{}
	fmt.Print("enter state")
	state, err := rd.ReadString('\n')
	if err != nil {
		return p, err
	}
	address.Street = strings.TrimSpace(state)
	fmt.Print("enter pincode")
	ptype, err := rd.ReadString('\n')
	if err != nil {
		return p, err
	}
	address.Zipcode = ptype
	info := &pb.Info{}
	info.Phone = phone
	info.Address = address
	p.Info = info
	fmt.Println("%+v\n", p)
	return p, nil
}

func main() {

	addr, err := promptForAddress(os.Stdin)
	if err != nil {
		log.Fatalln("Error with address:", err)
	}

	out, err := proto.Marshal(addr)
	if err != nil {
		log.Fatalln("Failed to encode address book:", err)
	}

	db, err := sql.Open("mysql", "root:root123@tcp(127.0.0.1:3306)/test")

	if err != nil {
		log.Print(err.Error())
	}
	defer db.Close()

	stmt, err := db.Prepare("INSERT INTO test3 SET Info=? ")
	if err != nil {
		panic(err.Error())
	}

	res, err := stmt.Exec(out)
	if err != nil {
		panic(err.Error())
	}
	id, err := res.LastInsertId()
	if err != nil {
		panic(err.Error())
	}

	fmt.Println(id)
	defer stmt.Close()

	rows, err := db.Query("SELECT * FROM test3")

	for rows.Next() {
		p := &pb.Person{}
		var info []byte
		err = rows.Scan(&info)
		if err != nil {
			panic(err.Error())
		}
		if err := proto.Unmarshal(info, p); err != nil {
			log.Fatalln("Failed to parse address book:", err)
		}
		fmt.Println(p.Info)
	}
}
