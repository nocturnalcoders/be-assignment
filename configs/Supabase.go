package configs

import (
	supa "github.com/nedpals/supabase-go"
)

func SetupSupabase() *supa.Client {
	supabaseUrl := "https://oxywjnmtwyqqtafqveny.supabase.co"
	supabaseKey := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJzdXBhYmFzZSIsInJlZiI6Im94eXdqbm10d3lxcXRhZnF2ZW55Iiwicm9sZSI6ImFub24iLCJpYXQiOjE3MTU0OTAwODIsImV4cCI6MjAzMTA2NjA4Mn0.U63K171jy6lFJnU5wszP8l_TZ1LCSujnTtwPXfdQzpM"
	supabase := supa.CreateClient(supabaseUrl, supabaseKey)

	return supabase
}
