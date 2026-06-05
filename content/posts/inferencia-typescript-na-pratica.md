---
title: "Inferência de tipos no TypeScript que ninguém te ensinou"
date: "2026-04-10"
tags: ["typescript", "tooling"]
excerpt: "const assertions, infer em conditional types e satisfies. Três ferramentas que transformam o compilador em copiloto de verdade."
readMin: 8
---

A maioria dos projetos usa 20% do sistema de tipos do TS. O resto — onde a inferência fica genuinamente inteligente — fica escondido atrás de sintaxe que parece intimidadora mas resolve problemas reais.

## satisfies: o operador subestimado

`satisfies` valida que um valor bate com um tipo **sem alargar** a inferência. Você ganha checagem e mantém o tipo literal estreito.

```typescript
const rotas = {
  home: '/',
  artigo: '/post/:slug',
} satisfies Record<string, string>;

// rotas.home tem tipo '/', não string
rotas.home; // '/'  (literal preservado)
```

## infer extrai tipos de dentro

Dentro de um conditional type, `infer` captura uma parte do tipo para você reusar. É como desestruturação, mas no nível dos tipos.

```typescript
type Elemento<T> = T extends Array<infer U> ? U : never;

type A = Elemento<string[]>; // string
type B = Elemento<number[]>; // number
```

## const assertions congelam o formato

`as const` diz ao compilador para tratar tudo como literal e readonly. Combinado com os outros dois, você descreve dados uma vez e deriva todos os tipos automaticamente.

> Boa tipagem não é escrever mais tipos. É escrever menos, e deixar o compilador derivar o resto.
