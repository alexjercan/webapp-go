package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

var myPosts []Post

type Post struct {
	Slug        uuid.UUID `json:"slug"`
	Name        string    `json:"name" binding:"required"`
	Description string    `json:"description"`
}

func DefaultPost() *Post {
	return &Post{
		Slug:        uuid.New(),
		Description: "",
	}
}

func getPost(c *gin.Context) {
	slug := c.Param("slug")
	c.String(http.StatusOK, "Hello %s", slug)
}

func getAllPosts(c *gin.Context) {
	c.String(http.StatusOK, "Hello all")
}

func createPost(c *gin.Context) {
    post := DefaultPost()

	if err := c.BindJSON(&post); err != nil {
        c.String(http.StatusBadRequest, err.Error())
	}

    myPosts = append(myPosts, *post)

	c.String(http.StatusOK, post.Slug.String())
}

func updatePost(c *gin.Context) {
	slug := c.Param("slug")
	c.String(http.StatusOK, "Update %s", slug)
}

func deletePost(c *gin.Context) {
	slug := c.Param("slug")
	c.String(http.StatusOK, "Delete %s", slug)
}

func getIndexPage(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", gin.H{
		"Posts": myPosts,
	})
}

func getPostPage(c *gin.Context) {
    var post Post
	slug, err := uuid.Parse(c.Param("slug"))
    if err != nil {
        c.String(http.StatusBadRequest, err.Error())
    }

    for _, p := range myPosts {
        if p.Slug == slug {
            post = p
        }
    }

	c.HTML(http.StatusOK, "post.html", gin.H{
		"Post": post,
	})
}

func main() {
	// Creates a gin router with default middleware:
	// logger and recovery (crash-free) middleware
	router := gin.Default()

	router.LoadHTMLGlob("templates/*")

	router.GET("/api/posts/:slug", getPost)
	router.GET("/api/posts", getAllPosts)
	router.POST("/api/posts", createPost)
	router.PUT("/api/posts/:slug", updatePost)
	router.DELETE("/api/posts/:slug", deletePost)

	router.GET("/", getIndexPage)
	router.GET("/posts/:slug", getPostPage)

	// By default it serves on :8080 unless a
	// PORT environment variable was defined.
	router.Run()
	// router.Run(":3000") for a hard coded port
}
