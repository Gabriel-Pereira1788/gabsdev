---
title: "Um modelo mental pro borrow checker do Rust"
date: "2026-05-18"
tags: ["rust", "systems"]
excerpt: "Parar de brigar com o compilador começa quando você entende que ownership não é uma regra arbitrária — é um grafo de quem pode tocar em quê, e quando."
readMin: 11
---

Todo mundo que aprende Rust passa pela mesma fase: o compilador recusa código que *obviamente* deveria funcionar. A virada de chave não é decorar regras — é internalizar o modelo que essas regras protegem.

## Ownership é um grafo

Pense em cada valor como um nó com exatamente **um** dono em qualquer instante. Mover um valor reescreve a aresta de propriedade; o dono antigo deixa de existir para o compilador. Não há cópia silenciosa, não há dois donos.

```rust
let s = String::from("gabs");
let t = s;            // move: s deixa de ser válido
println!("{}", t);    // ok
// println!("{}", s); // erro: value borrowed after move
```

A primeira vez que isso te morde parece punição. Mas é o compilador te impedindo de criar um *use-after-free* que em C passaria batido até produção.

## Empréstimos são leitura/escrita

Um `&T` é acesso de leitura compartilhado; um `&mut T` é acesso de escrita exclusivo. A regra de ouro: você pode ter N leitores OU 1 escritor, nunca os dois. É literalmente o RWLock que você implementaria à mão — só que verificado em tempo de compilação e com custo zero.

> O borrow checker não te impede de fazer coisas. Ele te impede de fazer coisas que você não conseguiria provar corretas sozinho.

## Lifetimes são só rótulos

Lifetimes assustam pela sintaxe, mas são apenas nomes que conectam a duração de uma referência à do dado que ela aponta. O compilador infere quase todos; você só anota quando há ambiguidade real sobre qual entrada sobrevive até a saída.

```rust
fn maior<'a>(a: &'a str, b: &'a str) -> &'a str {
    if a.len() > b.len() { a } else { b }
}
```

### O atalho prático

Quando travar: pergunte "quem é o dono desse dado e por quanto tempo?". 90% dos erros somem quando você decide o dono conscientemente em vez de espalhar referências.
