# Guia de Estilo de Escrita — Gabriel Andrade

> **Propósito.** Este documento descreve como o Gabriel escreve artigos técnicos, de forma que qualquer IA (ou ele mesmo no futuro) consiga produzir um texto novo no mesmo estilo. Foi extraído da análise de três artigos publicados no Medium e validado escrevendo um artigo real (`content/posts/next-para-go-templ-htmx.md`).
>
> **Como usar.** Antes de escrever um artigo "no estilo do Gabriel": (1) escolha um dos três Jeitos conforme o tipo de conteúdo; (2) siga o DNA comum em qualquer caso; (3) rode o checklist final antes de entregar.

---

## TL;DR para a IA

- A escrita é **guiada por jornada**: ancora no concreto → promete a entrega → leva o leitor pela mão → explica o porquê antes do como → fecha com calor humano.
- Sempre existe um **exemplo prático único** que percorre o artigo inteiro (não exemplos soltos).
- Tom **semi-informal** em PT-BR: técnico e preciso, mas conversa com o leitor.
- Há **três Jeitos** (formatos). Eles mudam só a proporção de código e o gatilho de abertura. O DNA é o mesmo.

---

## DNA comum (vale para TODO artigo)

| Traço | Regra prática |
|---|---|
| **Abertura ancorada** | Nunca abra genérico. Comece em algo concreto: uma versão recém-lançada, um desafio recebido, um debate da comunidade, uma pergunta retórica. |
| **Promessa explícita** | No 1º ou 2º parágrafo, diga o que o leitor vai levar. Ex.: "vou te mostrar como...", "vou detalhar como foi...". |
| **"Vamos" como motor** | Use a 1ª pessoa do plural para transicionar entre seções: *"Agora que entendemos X, vamos Y"*. Puxa o leitor junto. |
| **Conceito antes da prática** | Nunca jogue código ou solução sem antes explicar o porquê. Sempre teoria → mão na massa. |
| **Código embrulhado** | Todo bloco de código é **introduzido antes** (o que é / por que) e **explicado depois** (o que reparar). Nunca um bloco solto. |
| **Marcador de ênfase** | Sinalize o ponto crítico com a frase-assinatura: **"Um ponto importante a destacar é que..."**. |
| **Comparação como didática** | Ensine por contraste: velho vs novo (Bridge vs JSI), iOS vs Android, SPA vs HTML no servidor. |
| **Perguntas diretas ao leitor** | Faça perguntas retóricas que o leitor faria ("Tudo bem, mas e a busca?") e responda em seguida. |
| **Imperativo nos passos** | Em instruções práticas use comandos diretos: "Crie", "Importe", "Execute", "Repare". |
| **Fechamento humano** | Recapitule a jornada + emoji (🚀🙏😊) + encorajamento/gratidão. Nunca termine seco. |
| **Honestidade sobre trade-offs** | Não venda mágica. Sempre admita as limitações da solução. Reforça credibilidade. |

### Frases-assinatura (reutilizáveis)

- Abertura de jornada: *"Vou te guiar no passo a passo..."* / *"Neste artigo vou detalhar..."*
- Transição: *"Antes de começarmos..., vamos entender..."* / *"Agora que entendemos..., vamos finalmente..."* / *"Com tudo pronto, vamos..."*
- Ênfase: *"Um ponto importante a destacar é que..."*
- Mão na massa: *"Mãos a Obra!"*
- Fechamento: *"Fizemos um trajeto e tanto..."* / *"Espero que... 🚀"* / *"Agradeço a todos que... 🙏"* / *"lembre-se: teste o seu código :)"*

---

## Os três Jeitos

Escolha **um** por artigo, baseado no tipo de conteúdo.

### Jeito 1 — Tutorial técnico profundo
> Quando: ensinar a construir algo passo a passo. Muito código.
> Referência: *"Explorando a Nova Arquitetura do React Native com Turbo Modules"*.

- **Abertura:** ancora em contexto atual/versão. Ex.: *"Com a chegada da versão 0.76 do React Native..."*.
- **Arco:** teoria → conceitos fundamentais → "Mãos a Obra!" → implementação incremental → conclusão.
- **Código:** intenso (15+ blocos), complexidade crescente.
- **Voz:** 1ª pessoa do plural convidativa ("vamos criar", "vamos entender").
- **Fechamento:** recapitula a jornada com emoji. *"Fizemos um trajeto e tanto explorando... 🚀!"*.

### Jeito 2 — Case study narrativo (storytelling de engenharia)
> Quando: contar uma entrega real de produto. **Zero código** — foco em decisão, desafio, impacto.
> Referência: *"Implementando Live Activities... no Ploomes"*.

- **Abertura:** o desafio enquadrado como conquista. Ex.: *"Recentemente, recebi o desafio de implementar... transformando nosso CRM no primeiro a..."*.
- **Arco:** o que é a tecnologia → objetivo → protótipo → restrições de plataforma → desafios técnicos → comparação entre plataformas → posicionamento de inovação → conclusão com gratidão.
- **Voz:** 1ª pessoa do singular, com responsabilidade pessoal ("recebi", "vou detalhar", "vou destacar", "Agradeço").
- **Padrão recorrente:** estrutura *"O desafio foi..."* repetida; comparações paralelas (iOS vs Android); ênfase em UX e consistência.
- **Fechamento:** gratidão + inspiração. *"Agradeço a todos... 🙏 Espero que... inspirem outros desenvolvedores."*.

### Jeito 3 — Conceitual-educativo
> Quando: explicar um conceito/decisão e ancorá-lo num exemplo prático. Código médio.
> Referência: *"Testes em aplicações Front End"* e *"Reescrevendo este blog: de Next.js para Go + templ + HTMX"*.

- **Abertura:** pergunta retórica / debate que estabelece relevância. Ex.: *"...afinal, para que servem os testes em uma aplicação front end?"* / *"Será que todo site de conteúdo precisa mesmo de um SPA?"*.
- **Frameworks mentais como âncora:** dê nomes a modelos mentais ("Pirâmide de testes", "Troféu de testes", "pirâmide invertida"). O conceito vira o esqueleto do artigo.
- **Arco:** conceito → por quê → como fazer bem → **um exemplo prático único reaproveitado** em vários cenários → conclusão que desarma expectativa.
- **Voz:** instrucional + casual ("Você já deve ter ouvido falar", "pode parecer simples — e de fato é").
- **Fechamento:** desarma a expectativa ("este não é um tutorial para iniciantes" / "isso não é um manifesto anti-React") + encoraja ("lembre-se: teste o seu código :)").

---

## Estrutura de um artigo (esqueleto reutilizável)

1. **Abertura** (2 parágrafos): gatilho ancorado + promessa explícita.
2. **Seção de conceito/contexto**: o problema ou o framework mental, antes de qualquer código.
3. **Corpo** (3–6 seções com `##`): teoria → prática, cada bloco de código embrulhado, um exemplo único como fio condutor.
4. **Seção de trade-offs honestos**: vantagens e, com franqueza, as limitações.
5. **Conclusão**: desarma expectativa + recapitula a jornada + emoji + encorajamento.

> Cabeçalhos de seção usam `##` (no blog deste projeto, `##` alimenta o sumário/TOC lateral automaticamente — ver `internal/utils/toc.go`).

---

## Tom & mecânica

- **Idioma:** PT-BR. Quando pedido, produzir também versão em inglês — **mesmo arco e voz**, não tradução literal (manter a pergunta de abertura, os "let's", o "One important point to highlight is that...", o fecho caloroso).
- **Pessoa:** "vamos/nós" para guiar (Jeitos 1 e 3); "eu" para narrar a própria entrega (Jeito 2).
- **Parágrafos:** curtos a médios (2–6 frases). Seções introdutórias mais curtas; trechos técnicos podem ser mais longos.
- **Listas vs prosa:** mistura. Bullets para enumerar tipos/passos; prosa para explicar conceitos.
- **Emojis:** com parcimônia e quase sempre no fechamento ou em marcos (🚀🙏😊).
- **Profundidade:** assume familiaridade com a stack; explica o conceito novo a fundo, sem subestimar o leitor.

### Pontuação — NÃO usar travessão "—"

**Regra dura: nunca usar travessão (em dash `—`) como pausa.** Esse caractere denuncia texto gerado por IA. Em vez dele, reescreva a frase de uma destas formas:

- Quebre em duas frases com ponto. Ex.: *"...o problema não é o React. É usar a ferramenta pesada..."*
- Use vírgula quando a ligação for leve. Ex.: *"...o `for` é o `.map()` que você já conhece, só que é um `for` de Go..."*
- Use dois-pontos quando o segundo trecho explica o primeiro. Ex.: *"Repare no `{slug}`: rota com parâmetro, sem dependência externa."*

O mesmo vale para a versão em inglês.

### Ortografia — atenção

A escrita espontânea do autor no Medium tem desvios ("oque", "trás" por "traz", "ja" sem acento). Isso é traço de informalidade, **não um padrão a reproduzir**. Em material publicado/documentado, **manter português correto e legível**. Não imitar os erros de propósito.

---

## Checklist final (rodar antes de entregar)

- [ ] Abertura ancorada em algo concreto (não genérica)?
- [ ] Promessa explícita no 1º/2º parágrafo?
- [ ] Um Jeito escolhido e seguido de ponta a ponta?
- [ ] Existe um exemplo prático único como fio condutor (Jeitos 1 e 3)?
- [ ] Todo bloco de código está embrulhado (introduzido antes, explicado depois)?
- [ ] Pelo menos um marcador "Um ponto importante a destacar é que..."?
- [ ] Transições com "vamos" entre seções?
- [ ] Seção de trade-offs honestos presente?
- [ ] Fechamento humano: recapitulação + emoji + encorajamento?
- [ ] **Zero travessões "—"** no texto (PT e EN)?
- [ ] Português correto (sem reproduzir os desvios ortográficos do autor)?
- [ ] Se bilíngue: versão EN mantém arco e voz, não é tradução literal?

---

## Artigos de referência

| Artigo | Jeito | Link |
|---|---|---|
| Turbo Modules / Nova Arquitetura RN | 1 — Tutorial profundo | https://medium.com/@gabriel.andrade1788/explorando-a-nova-arquitetura-do-react-native-com-turbo-modules-74c88a52acb3 |
| Live Activities no Ploomes | 2 — Case study narrativo | https://medium.com/@gabriel.andrade1788/implementando-live-activities-transformando-o-check-in-de-visitas-com-live-activities-no-ploomes-51acf304fb7c |
| Testes em aplicações Front End | 3 — Conceitual-educativo | https://medium.com/@gabriel.andrade1788/testes-em-aplica%C3%A7%C3%B5es-front-end-81d55b85c16c |
| Reescrevendo este blog (Next → Go) | 3 — Conceitual-educativo | `content/posts/next-para-go-templ-htmx.md` (escrito com este guia) |
