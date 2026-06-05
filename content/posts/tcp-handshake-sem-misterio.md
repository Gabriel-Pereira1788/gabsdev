---
title: "O three-way handshake do TCP sem mistério"
date: "2026-02-14"
tags: ["networking", "systems"]
excerpt: "SYN, SYN-ACK, ACK. Três pacotes que você manda milhares de vezes por dia e que carregam mais decisão de engenharia do que parece."
readMin: 7
---

Toda conexão TCP começa com um ritual de três pacotes. Parece burocracia, mas cada passo resolve um problema concreto de comunicação sobre um canal não confiável.

## Os três passos

- **SYN** — cliente envia um número de sequência inicial aleatório.
- **SYN-ACK** — servidor confirma o do cliente e manda o seu.
- **ACK** — cliente confirma o do servidor. Conexão estabelecida.

Por que números aleatórios? Para evitar que pacotes de uma conexão antiga sejam confundidos com a nova — e para dificultar spoofing.

## Onde a latência mora

O handshake custa um *round-trip* inteiro antes de qualquer dado útil. É por isso que protocolos modernos como QUIC fundem handshake e dados — eliminar essa ida e volta é ganho direto de performance.

> Cada round-trip é um imposto que você paga na velocidade da luz. Não dá pra otimizar a física — só pra mandar menos pacotes.
