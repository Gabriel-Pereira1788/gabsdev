package domain

// PostMeta = interface PostMeta no TypeScript
// só os metadados, sem o conteúdo completo
type PostMeta struct {
	Slug    string   `yaml:"slug"`
	Title   string   `yaml:"title"`
	Date    string   `yaml:"date"`
	Tags    []string `yaml:"tags"`
	Excerpt string   `yaml:"excerpt"`
	ReadMin int      `yaml:"readMin"`
}

// Post = interface Post extends PostMeta no TypeScript
// tem tudo do PostMeta + o conteúdo Markdown
type Post struct {
	PostMeta
	Content string
}
