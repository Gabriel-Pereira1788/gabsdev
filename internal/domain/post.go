package domain

type PostMeta struct {
	Slug    string   `yaml:"slug"`
	Title   string   `yaml:"title"`
	Date    string   `yaml:"date"`
	Tags    []string `yaml:"tags"`
	Excerpt string   `yaml:"excerpt"`
	ReadMin int      `yaml:"readMin"`
}

type Post struct {
	PostMeta
	Content string
}
