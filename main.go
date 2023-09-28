package main

func main() {
	app := App{}
	err := app.Initialize(DbUser, DbPassword, DbName)
	if err != nil {
		panic(err)
	}
	app.Run("localhost:9000")
}
