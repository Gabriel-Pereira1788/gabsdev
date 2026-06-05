---
title: "Anatomia de um cache miss"
date: "2026-03-22"
tags: ["performance", "systems"]
excerpt: "Por que iterar uma matriz por colunas pode ser 8x mais lento que por linhas, mesmo fazendo exatamente o mesmo número de operações."
readMin: 12
---

Duas funções, mesma complexidade de Big-O, mesmo número de somas. Uma roda em 12ms, a outra em 96ms. A diferença não está no algoritmo — está em como a memória conversa com o processador.

## A hierarquia que ninguém vê

Entre o registrador e a RAM há três níveis de cache. Acessar L1 custa ~1ns; ir até a RAM custa ~100ns. Cada vez que o dado que você quer não está no cache, o processador **para e espera**.

## Localidade espacial

Quando você lê um endereço, o CPU traz uma *linha de cache* inteira (64 bytes) junto. Se o próximo acesso estiver nessa linha, é de graça. Se estiver longe, é outro miss.

```c
// rápido: percorre na ordem em que está na memória
for (int i = 0; i < N; i++)
  for (int j = 0; j < N; j++)
    soma += m[i][j];

// lento: pula 'N' elementos a cada passo -> miss constante
for (int j = 0; j < N; j++)
  for (int i = 0; i < N; i++)
    soma += m[i][j];
```

> O algoritmo certo na ordem errada de memória ainda é lento. Hardware tem opinião.

## Como medir de verdade

Não confie no relógio sozinho — use `perf stat` e olhe a taxa de cache miss. Número de instruções igual + tempo diferente = problema de memória, não de CPU.
