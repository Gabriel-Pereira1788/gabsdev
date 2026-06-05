---
title: "Escrevendo um parser recursivo descendente do zero"
date: "2026-03-03"
tags: ["compilers", "rust"]
excerpt: "Calculadoras são chatas até você perceber que precedência de operadores é um problema lindo — e que recursão resolve quase elegantemente."
readMin: 14
---

Parser generators são ótimos para produção, mas escondem a intuição. Escrever um à mão, uma vez, faz você entender o que toda linguagem está fazendo nos bastidores.

## Tokens primeiro

Antes de parsear, transforme a string em uma sequência de tokens. O lexer não entende estrutura — só categoriza pedaços.

```rust
enum Token {
    Num(f64),
    Plus, Minus, Star, Slash,
    LParen, RParen,
}
```

## Precedência vira recursão

O truque central: cada nível de precedência é uma função que chama o nível mais alto. `expr` chama `term`, que chama `factor`. A gramática vira a pilha de chamadas.

```rust
fn expr(&mut self) -> Node {
    let mut left = self.term();
    while self.peek_is(&[Plus, Minus]) {
        let op = self.next();
        let right = self.term();
        left = Node::Bin(op, box left, box right);
    }
    left
}
```

> Recursão descendente é a gramática se lendo em voz alta. Se você consegue escrever as regras, você já escreveu o parser.

A partir daqui, adicionar potência, funções e variáveis é só mais funções no mesmo padrão. O esqueleto não muda.
