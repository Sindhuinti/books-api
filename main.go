package main

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"errors"
)

type book struct {
	ID		string	`json:"id"`
	Title   string	`json:"title"`
	Author  string	`json:"author"`
	Quantity int	`json:"quantity"`
}

type user struct{
	Email   string 	`json:"email"`
	Password  string  `json:"password"`
}

var books = []book{
	{ID:"1",Title:"Ruby",Author:"Stone",Quantity: 12},
	{ID:"2",Title:"Python",Author:"Snake",Quantity: 5},
	{ID:"3",Title:"Java",Author:"Cup",Quantity: 7},
}

var users = []user{
	{Email:"hi@mail.com",Password: "123"},
	
}
func HashPassword(password string) (string, error) {
    bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
    return string(bytes), err
}
func CheckPasswordHash(password, hash string) bool {
    err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
    return err == nil
}

func bookById(c *gin.Context){
	id := c.Param("id")
	book,err:= getBookById(id)
	if err!=nil{
		c.IndentedJSON(http.StatusNotFound,gin.H{"message":"Boook not found"})
		return
	}
	c.IndentedJSON(http.StatusOK, book);
}

func createBook(c * gin.Context){
	var newBook book
	if err:= c.BindJSON(&newBook); err != nil{
		return 
	}

	books = append(books,newBook)
	c.IndentedJSON(http.StatusCreated,newBook)
}


func checkOutBook(c *gin.Context){
	id,ok:=c.GetQuery("id")
	if ok ==false{
		c.IndentedJSON(http.StatusBadRequest,gin.H{"message":"invalid query"})
		return 
	}
	book,err:=getBookById(id)
	if err !=nil{
		c.IndentedJSON(http.StatusNotFound,gin.H{"message":"Book not found"})
		return
	}

	if book.Quantity ==0{
		c.IndentedJSON(http.StatusBadRequest,gin.H{"message":"Book not available"})
		return
	}
	book.Quantity -=1
	c.IndentedJSON(http.StatusOK,book)

}

func getBookById(id string)(*book,error){
	for i,b := range books{
		if b.ID==id {
			return &books[i],nil
		}
	}
	return nil,errors.New("Book not found");

}

func getBooks (c *gin.Context){
	c.IndentedJSON(http.StatusOK,books);

}
func returnBook(c *gin.Context){

	id,ok:=c.GetQuery("id")
	if ok ==false{
		c.IndentedJSON(http.StatusBadRequest,gin.H{"message":"invalid query"})
		return 
	}
	book,err:=getBookById(id)
	if err !=nil{
		c.IndentedJSON(http.StatusNotFound,gin.H{"message":"Book not found"})
		return
	}

	book.Quantity +=1
	c.IndentedJSON(http.StatusOK,book)


}

func userSignup(c *gin.Context){

	var newUser user
	if err:= c.BindJSON(&newUser); err != nil{
		c.IndentedJSON(http.StatusBadRequest,gin.H{"message":"Cant create user"})
		return 
	}

	newUser.Password ,_ = HashPassword(newUser.Password)
	
	users = append(users,newUser)
	c.IndentedJSON(http.StatusCreated,newUser)

}

func main(){
	router:=gin.Default();
	router.GET("/books",getBooks)
	router.POST("/books",createBook)
	router.GET("/books/:id",bookById)
	router.PATCH("/books/checkout",checkOutBook)
	router.PATCH("/books/return",returnBook)
	router.POST("/signup",userSignup)
	router.Run("localhost:8080")

}

