package runtime

import "testing"

func TestDetect(t *testing.T) {
	cases := []struct {
		name    string
		cmdLine string
		want    string
	}{
		{"node", "node ./node_modules/next/dist/bin/next dev", "Node.js"},
		{"node.exe", "node server.js", "Node.js"},
		{"python3", "python3 manage.py runserver", "Python"},
		{"python", "python -m uvicorn main:app", "Python"},
		{"java", "java -jar app.jar", "Java"},
		{"dotnet", "dotnet run", ".NET"},
		{"ruby", "ruby app.rb", "Ruby"},
		{"php", "php artisan serve", "PHP"},
		{"docker-proxy", "", "Docker"},
		{"nginx", "", ""},
	}

	for _, c := range cases {
		got := Detect(c.name, c.cmdLine)
		if got != c.want {
			t.Errorf("Detect(%q, %q) = %q, want %q", c.name, c.cmdLine, got, c.want)
		}
	}
}
