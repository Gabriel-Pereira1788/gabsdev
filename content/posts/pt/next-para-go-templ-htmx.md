---
title: "Reescrevendo este blog: de Next.js para Go + templ + HTMX"
date: "2026-06-11"
tags: ["go", "templ", "htmx", "nextjs", "web"]
excerpt: "Será que todo site precisa mesmo de um Next.js da vida, ou dá para simplificar as coisas? Reescrevi este próprio blog em Go puro, templ e HTMX para responder essa pergunta na prática."
readMin: 10
---

Se há uma decisão que virou quase automática no desenvolvimento web moderno, é começar todo projeto com o mesmo suspeito de sempre. Você quer um blog? Next. Um portfólio? Next. Uma landing page? Next. Mas afinal, será que todo site precisa mesmo de um Next.js da vida, ou dá para simplificar as coisas?

Esse blog que você está lendo agora começou exatamente assim: um projeto Next.js. Neste artigo vou te mostrar como eu o reescrevi do zero em **Go puro, templ e HTMX**. E, mais importante que o "como", por que esse caminho faz sentido para um site que é, no fim das contas, majoritariamente conteúdo.

## O peso do default moderno

Antes de falar da solução, vamos entender o problema. Um site como esse tem uma natureza simples: ele renderiza uma lista de artigos, mostra um artigo, filtra por tags. Quase nada disso é *interativo* de verdade. É conteúdo que precisa chegar rápido e ser lido.

O Next renderiza no servidor, sim, ninguém está dizendo o contrário. Mas ele ainda manda o React para o navegador e hidrata a página no cliente. E esse pacote cobra um preço por padrão, mesmo quando você não precisa dele:

- Um bundle de JavaScript que o navegador precisa baixar, parsear e executar antes da página ficar útil.
- Hidratação: o servidor manda HTML, e o cliente reexecuta tudo para "dar vida" a esse HTML.
- Um pipeline de build com dezenas de dependências que envelhecem rápido.

Para uma aplicação cheia de estado no cliente, esse custo se paga. Para um blog, você está carregando um caminhão para entregar uma carta. **Um ponto importante a destacar é que** o problema não é o React. É usar a ferramenta de interatividade pesada num lugar onde a interatividade é a exceção, não a regra.

## A simplicidade como escolha

Vale separar dois tipos de complexidade. Existe a essencial, que é o problema em si: renderizar artigos, buscar, filtrar por tag. E existe a acidental, que é a que a ferramenta arrasta junto sem o problema pedir. O Next resolve a parte essencial muito bem, mas cobra uma boa dose da acidental: um runtime para hidratar, um pipeline de build, um grafo de dependências para manter de pé. Para este blog, quase todo esse peso era acidental.

Trocar de stack foi, antes de tudo, jogar peça fora. Não tem `node_modules` para instalar, não tem toolchain de Node para versionar, não tem build de frontend para quebrar numa atualização de dependência. O que sobra é um binário Go e alguns arquivos `.templ`. Menos peças significam menos superfície que envelhece, menos coisa para atualizar quando uma biblioteca resolve mudar a API.

E tem um ganho que não aparece em nenhum bundle: o projeto inteiro cabe na cabeça. O caminho de qualquer requisição é sempre o mesmo, Rota → Handler → Service → Componente, e dá para segui-lo de ponta a ponta sem abrir a documentação de três bibliotecas. Simplicidade aqui não é fazer menos. **É tirar o que não estava ajudando, para que o que sobra seja fácil de entender e difícil de quebrar.**

Agora que está claro por que simplificar compensa, vamos ver o modelo mental que torna tudo isso possível.

## A pirâmide invertida: HTML no servidor

A ideia da reescrita é virar essa pirâmide de cabeça para baixo. Em vez de "JavaScript no cliente por padrão, HTML como detalhe", a regra passa a ser **HTML no servidor por padrão, JavaScript só onde dói**.

Esse modelo mental tem três camadas, e é bom fixá-las antes de ver qualquer código:

- **Go** é o servidor. Recebe a requisição, busca o dado, decide o que renderizar.
- **templ** é o HTML, mas com checagem de tipos. Componentes que viram funções Go de verdade.
- **HTMX** é a interatividade pontual. Pequenos pedaços de HTML trocados sob demanda, sem escrever JavaScript.

Agora que entendemos o modelo, vamos finalmente ver como ele se materializa neste projeto.

## Mãos à obra: o esqueleto

A primeira surpresa de quem vem do Next é descobrir que não precisa de framework nenhum para o servidor. A biblioteca padrão do Go já entrega roteamento, e o `cmd/server/main.go` deste blog inteiro cabe numa tela:

```go
mux := http.NewServeMux()

mux.HandleFunc("/", handlers.Home)
mux.HandleFunc("/articles", handlers.Articles)
mux.HandleFunc("/articles/search", handlers.ArticlesSearch)
mux.HandleFunc("/articles/{slug}", handlers.ArticleDetail)
mux.HandleFunc("/about", handlers.About)
mux.HandleFunc("/tags", handlers.Tags)
mux.HandleFunc("/rss.xml", handlers.RSS)

mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

log.Println("Server on http://localhost:8080")
http.ListenAndServe(":8080", mux)
```

Repare no `{slug}` na linha do `/articles/{slug}`: rota com parâmetro, sem dependência externa. Cada rota aponta para um *handler*, e o padrão que se repete o projeto inteiro é simples: **Rota → Handler → Service → Componente**. O handler busca o dado num service e manda renderizar um componente. Só isso.

## templ no lugar do JSX

Aqui mora a parte que mais aproxima quem vem do React. Com **templ**, você escreve componentes que se parecem muito com JSX, mas que são compilados para funções Go com tipos checados. Veja a página de artigos:

```go
templ ArticlesPage(pathname string, posts []domain.Post, tags []string, activeTag string) {
	@layouts.Layout("Articles", pathname) {
		<div class="page enter">
			@components.SearchBar(activeTag)
			<div class="post-list" id="post-list">
				for _, post := range posts {
					@components.PostRow(post, activeTag)
				}
			</div>
		</div>
	}
}
```

Se você já escreveu React, isso parece familiar: o layout recebe filhos (o `{ children... }` do templ), os componentes se compõem com `@components.PostRow(...)`, e o `for` é o `.map()` que você já conhece, só que é um `for` de Go de verdade.

**Um ponto importante a destacar é que** esse `.templ` não é interpretado em runtime. Antes de compilar, você roda `templ generate` e cada arquivo vira um `_templ.go` correspondente, que é código Go normal. Isso significa duas coisas: se você passar o tipo errado para um componente, o build quebra (não o navegador do usuário); e não há reflexão nem template parsing em produção, é só concatenação de strings já compilada.

## HTMX: interatividade sem framework

"Tudo bem, mas e a busca? E o filtro por tag? Isso é interativo." É aqui que entra o HTMX, e é aqui que a reescrita prova que dá para ter interatividade sem trazer um framework inteiro junto.

O HTMX entra por um único `<script>` no layout (carregado via CDN) e funciona estendendo o HTML com atributos. Veja o campo de busca deste blog:

```go
templ SearchBar(activeTag string) {
	<div class="search-wrap">
		<input
			class="search-input"
			type="text"
			name="q"
			placeholder={ searchPlaceholder(activeTag) }
			hx-get="/articles/search"
			hx-trigger="input changed delay:300ms, search"
			hx-target="#post-list"
			hx-swap="innerHTML"
			hx-vals={ `{"tag":"` + activeTag + `"}` }
		/>
	</div>
}
```

Vamos ler esses atributos como uma frase: quando o input *mudar* (`hx-trigger`, com 300ms de debounce de brinde), faça um `GET` em `/articles/search` (`hx-get`), pegue o HTML que voltar e troque o conteúdo de `#post-list` por ele (`hx-target` + `hx-swap`). Não escrevi uma linha de JavaScript para isso.

E o que o servidor responde? Aqui está o detalhe mais bonito da abordagem. O handler **não devolve a página inteira**. Devolve só o pedaço que mudou:

```go
func ArticlesSearch(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	tag := r.URL.Query().Get("tag")

	posts := services.SearchPosts(query, tag)

	components.PostList(posts, tag).Render(r.Context(), w)
}
```

Compare com o handler da página cheia: aquele renderiza `pages.ArticlesPage(...)`, que vem embrulhada no `@layouts.Layout`. Este renderiza só `components.PostList(...)`, o mesmo componente da lista, sem o layout em volta. **Mesma rota mental, dois tamanhos de resposta:** página completa na navegação, fragmento na interação. O HTMX cuida de costurar o fragmento de volta na tela.

O clique numa tag segue exatamente a mesma ideia, com um bônus de `hx-push-url` para a URL continuar compartilhável:

```go
<span
	class={ getTagClass(tag, activeTag) }
	hx-get={ "/articles/search?tag=" + tag }
	hx-target="#post-list"
	hx-swap="innerHTML"
	hx-push-url={ "/articles?tag=" + tag }
>{ tag }</span>
```

## Conteúdo como arquivo, não banco

Tem mais uma decisão que a reescrita deixou cair por terra, pelo menos por enquanto: o banco de dados. E aqui vale uma honestidade sobre o estado atual do projeto. Eu decidi, de propósito, não subir banco nenhum agora. A ideia foi a mais simples possível: escrever cada artigo num Markdown e guardar localmente, versionado junto com o código.

**Um ponto importante a destacar é que** essa é uma escolha temporária, não um dogma. Lá na frente eu pretendo migrar todo esse conteúdo para um banco de verdade e construir um dashboard administrativo, só para ter um controle razoável dos dados. Mas para o estágio em que o blog está hoje, banco seria mais peça para manter do que problema resolvido. Então vamos com o caminho mais leve enquanto ele dá conta.

Na prática, o conteúdo de um blog são os artigos, e artigos são texto. Cada post vira um arquivo Markdown em `content/posts/`, com um frontmatter YAML no topo (inclusive este que você está lendo).

O service lê a pasta, separa o frontmatter do corpo num `SplitN` por `---`, e passa o Markdown pelo goldmark com GFM e syntax highlighting:

```go
md := goldmark.New(
	goldmark.WithExtensions(
		extension.GFM,
		highlighting.NewHighlighting(
			highlighting.WithStyle("nord"),
		),
	),
	goldmark.WithParserOptions(
		parser.WithAutoHeadingID(),
	),
)
```

O `WithAutoHeadingID` ainda gera os IDs dos cabeçalhos automaticamente, o que alimenta o sumário lateral que você vê na lateral deste artigo. A vantagem dessa escolha é difícil de exagerar: o conteúdo vive no git, versionado junto com o código, sem migration, sem painel de admin, sem banco para subir no deploy.

## Vantagens (e os trade-offs honestos)

Reescrever esse blog rendeu ganhos concretos frente ao Next:

- **Menos JavaScript no cliente.** O que chega no navegador é HTML pronto. O único JS de peso é o próprio HTMX, pequeno e estável.
- **Sem build de frontend.** O passo de build é `templ generate && go build`. Os `_templ.go` gerados ficam commitados, o que torna o deploy no Railway trivial, sem toolchain de Node e sem `node_modules`.
- **Type-safety de ponta a ponta.** O mesmo compilador que valida o service valida o template.
- **Simplicidade que cabe na cabeça.** Rota → Handler → Service → Componente. Dá para ler o projeto inteiro numa tarde.

Mas eu não estaria sendo honesto se vendesse isso como mágica. O ecossistema de Go + HTMX é menor que o do React. Você vai achar menos componentes prontos e menos respostas no Stack Overflow. E interações realmente ricas (um editor de texto, um canvas, drag-and-drop complexo) ainda pedem JavaScript de verdade. O HTMX brilha na troca de fragmentos, não em estado complexo de cliente. Existem boas soluções para cobrir essa lacuna de reatividade, como o Alpine.js, que se encaixa bem ao lado do HTMX para o estado local que mora no navegador. Ainda assim, a questão nunca foi "qual é melhor", e sim **qual o peso certo para o problema certo**.

## Conclusão

Foi um trajeto e tanto sair do Next e chegar num servidor que devolve HTML pronto 🚀. Demos uma olhada em como o Go segura o roteamento sozinho, como o templ traz a ergonomia do JSX com checagem de tipos, e como o HTMX dá interatividade sem arrastar um framework inteiro junto.

Que fique claro: isso não é um manifesto anti-React. O Next continua sendo uma escolha excelente para aplicações de verdade cheias de estado no cliente. O ponto é outro. Antes de aceitar o default, vale perguntar de quanta interatividade o seu projeto realmente precisa. No caso deste blog, a resposta era "quase nenhuma", e o stack passou a refletir isso.

Espero que essa jornada te faça olhar para o próximo projeto e considerar o caminho mais simples antes de pegar o caminhão. Obrigado por ler até aqui 🙏. E, claro, esse blog inteiro é open source: vá lá fuçar o código.
