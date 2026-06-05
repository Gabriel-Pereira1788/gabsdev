---
title: "Por que o banco ignora seu índice (e como descobrir)"
date: "2026-04-29"
tags: ["databases", "performance"]
excerpt: "Você criou o índice, mas a query continua varrendo a tabela inteira. O EXPLAIN sabe o porquê — e quase sempre é uma de cinco razões."
readMin: 9
---

Criar um índice é fácil. Garantir que o planejador *use* ele é onde mora a frustração. A boa notícia: o banco te conta exatamente o que está pensando.

## Leia o plano, não o palpite

Antes de qualquer otimização, rode o `EXPLAIN ANALYZE`. Um **Seq Scan** onde você esperava um Index Scan é o sintoma; a causa está nos detalhes.

```sql
EXPLAIN ANALYZE
SELECT * FROM eventos
WHERE lower(email) = 'gabs@dev.io';
-- Seq Scan on eventos ... (custo alto, índice ignorado)
```

## As cinco razões clássicas

Quase todo índice ignorado cai em um destes casos:

- Função na coluna (`lower(email)`) — o índice é sobre `email`, não sobre `lower(email)`.
- Tipo divergente — comparar `bigint` com string força um cast que descarta o índice.
- Seletividade baixa — se 40% das linhas casam, o scan sequencial é genuinamente mais rápido.
- Estatísticas velhas — rode `ANALYZE`; o planejador decide com dados desatualizados.
- Ordem das colunas no índice composto não bate com os predicados da query.

## Índice funcional resolve o primeiro caso

```sql
CREATE INDEX idx_eventos_email_lower
  ON eventos (lower(email));
-- agora a query acima usa Index Scan
```

> O planejador raramente está errado. Quando ele ignora seu índice, normalmente é porque o índice realmente não ajudaria — ou porque você o indexou diferente de como consulta.

Antes de adicionar índices, sempre confirme que o gargalo é leitura e não escrita: cada índice é mais um custo em todo `INSERT` e `UPDATE`.
