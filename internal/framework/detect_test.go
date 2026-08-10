package framework

import "testing"

func TestDetectNode(t *testing.T) {
	cases := []struct {
		cmd  string
		want string
	}{
		{"node ./node_modules/next/dist/bin/next dev", "Next.js"},
		{"node node_modules/.bin/vite", "Vite"},
		{"node node_modules/.bin/nuxt dev", "Nuxt"},
		{"node node_modules/.bin/react-scripts start", "Create React App"},
		{"node dist/main.js", ""},
	}
	for _, c := range cases {
		got := Detect("Node.js", c.cmd, "")
		if got != c.want {
			t.Errorf("Detect(Node.js, %q) = %q, want %q", c.cmd, got, c.want)
		}
	}
}

func TestDetectPython(t *testing.T) {
	cases := []struct {
		cmd  string
		want string
	}{
		{"python -m uvicorn main:app --reload", "Uvicorn / FastAPI"},
		{"python manage.py runserver", "Django"},
		{"flask run", "Flask"},
		{"gunicorn app:app", "Gunicorn"},
		{"python script.py", ""},
	}
	for _, c := range cases {
		got := Detect("Python", c.cmd, "")
		if got != c.want {
			t.Errorf("Detect(Python, %q) = %q, want %q", c.cmd, got, c.want)
		}
	}
}
