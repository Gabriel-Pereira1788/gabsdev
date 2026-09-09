---
title: "Salve DB: SQLite offline-first para React Native"
date: "2026-09-04"
tags: ["react-native", "sqlite", "offline-first", "mobile", "arquitetura"]
excerpt: "A maioria das libs \"offline-first\" ainda depende de JavaScript rodando em background para sincronizar. E se o motor de sync inteiro, da fila até o refresh de token OAuth2, rodasse 100% em código nativo?"
readMin: 14
---

Toda vez que um app precisa funcionar offline, a solução mais comum acaba sendo a mesma: grava local, e quando der, dispara uma sincronização que roda em JavaScript em background. Funciona, até o sistema operacional decidir que aquele processo JS não merece continuar vivo. Nesse momento a sincronização simplesmente para, e ninguém percebe até o usuário reclamar que os dados sumiram.

Foi esse problema que me fez construir o [Salve DB](https://github.com/Salve-Software/react-native-salve-db), uma biblioteca de SQLite offline-first para React Native. Neste artigo vou te mostrar como ela funciona por dentro: por que o motor de sincronização inteiro roda em código nativo (C++/Swift/Kotlin) e nunca precisa subir a JS engine para sincronizar, como uma escrita local vira uma sincronização de verdade, e quais trade-offs eu assumi conscientemente para chegar nesse design.

## O problema com "offline-first" de JavaScript

Antes de falar da solução, vamos entender o que realmente acontece quando uma sincronização em background depende do JavaScript. No iOS, isso significa colar callbacks do `BGTaskScheduler` num contexto JS que pode nem estar vivo naquele momento. No Android, um job do `WorkManager` precisa subir uma JS runtime inteira antes de conseguir fazer uma única requisição HTTP.

Isso custa tempo e bateria só para *iniciar* a engine, e ainda fica refém de como o sistema operacional agenda a JS thread sob os limites de execução em background. É uma classe inteira de bug de confiabilidade que a maioria das libs de sync simplesmente aceita como custo do jogo.

**Um ponto importante a destacar é que** muita biblioteca se diz "offline-first" só por gravar no disco local antes de tentar a rede. Isso não é a mesma coisa que sincronizar de verdade sem depender do app estar aberto. No Salve DB, a captura de mudança é automática e infalível: não é uma chamada de código que alguém pode esquecer de fazer, é uma trigger de banco. E a sincronização em si não depende do app estar aberto: um job periódico do sistema operacional acorda o motor nativo sozinho.

## A virada de chave: mover o motor de sync inteiro para o nativo

A tese central do Salve DB é simples de enunciar: o TypeScript só declara *o quê* sincronizar e *como* sincronizar, mas nunca participa da execução. Orquestração, HTTP, credenciais OAuth2, agendamento em background, tudo roda em C++/Swift/Kotlin.

Isso aparece na arquitetura em quatro camadas:

<figure>
  <img src="/static/images/salve-db-architecture.png" alt="Diagrama da arquitetura em camadas do Salve DB: TypeScript se comunica via JSI com o Native Core em C++, que por sua vez usa serviços de plataforma (BGTaskScheduler/WorkManager, NWPathMonitor/ConnectivityManager, Keychain/Keystore, Swift/Kotlin) para falar com a REST API" loading="lazy">
  <figcaption>TypeScript só declara; JSI é a ponte síncrona; o Native Core em C++ concentra a lógica; Swift/Kotlin cobrem só o que o C++ não alcança.</figcaption>
</figure>

Vamos passar por cada camada, porque cada uma existe por um motivo específico e não é só "organização por organizar".

**TypeScript** é só declaração. Schemas são dado puro, o query builder monta SQL e parâmetros mas não executa nada, e os hooks/provider são a superfície de DX. A regra "sem JS" do projeto é sobre o motor de sync rodando sem supervisão em background, não sobre o builder em foreground: quando você chama `.execute()` com o app aberto, o JS continua de pé.

**JSI (via Nitro Modules)** é a ponte. Uma bridge zero-copy e síncrona, via `HybridObject`. É *por causa* dessa sincronicidade que a execução de query em foreground é síncrona na JS thread, o que por sua vez explica a existência de guard-rails que eu vou detalhar mais à frente.

**Native Core (C++)** é onde mora a maior parte da lógica: o executor de SQLite com cache LRU de prepared statements, o motor de migração, o motor de triggers, o orquestrador de sync, o provider de credenciais e o cliente HTTP. Tudo independente de plataforma.

**Swift/Kotlin** são as bordas, e só existem porque há coisas que o C++ genuinamente não alcança: `BGTaskScheduler`/`WorkManager` são APIs exclusivas de Swift/Kotlin, e fazer bridge dos callbacks de ciclo de vida via JNI/runtime Objective-C seria frágil e propenso a erro. Keychain e Keystore também não têm binding C++. E empacotar `libcurl` significaria reimplementar manualmente proxy, cert pinning e configuração de TLS que `URLSession`/`OkHttp` já resolvem corretamente.

Por que não escrever tudo em Swift/Kotlin, então? Porque a lógica core (fila de sync, interpretador de expressão, resolução de conflito) não tem dependência nenhuma de plataforma, e escrever isso duas vezes garante que as duas implementações vão divergir mais cedo ou mais tarde.

Agora que entendemos as camadas, vamos ver o exemplo que vai atravessar o resto do artigo: uma tabela `users` sincronizada.

## Um schema, do início ao fim

Todo o resto deste artigo gira em torno de um único schema. É assim que se declara uma tabela `users` sincronizável no Salve DB:

```ts
import type { ISchemaDefinition } from '@salve-software/react-native-salve-db';

export interface User {
  id: number;
  name: string;
  email: string;
  updatedAt: number;
}

export const UserSchema = {
  name: 'users',
  version: 1,
  primaryKey: 'id',
  columns: {
    id: { type: 'integer' },
    name: { type: 'text' },
    email: { type: 'text' },
    updatedAt: { type: 'datetime', nullable: false },
  },
  indexes: [
    { name: 'idx_users_updated_at', columns: ['updatedAt'] },
    { name: 'idx_users_email', columns: ['email'] },
  ],
  sync: {
    enabled: true,
    direction: 'bidirectional',
    conflict: 'lastWriteWins',
    transport: 'rest',
    endpoint: { basePath: '/users', listQueryTemplate: 'updatedAfter={since}&limit={limit}' },
    pagination: { pageSize: 50, maxPagesPerSession: 20 },
  },
} satisfies ISchemaDefinition<User>;
```

Repare que não existe SQL nenhum aqui: colunas, índices e o contrato de sync são só dados TypeScript, interpretados pelo core nativo no boot. `sync` é opcional, tabelas locais que nunca precisam de rede simplesmente omitem o bloco inteiro. E há um detalhe que passa despercebido na primeira leitura: `deletedAt` é injetado automaticamente em toda tabela, mesmo sem você declarar. Isso importa porque delete, no Salve DB, nunca é um `DELETE` de SQL de verdade, é sempre um soft delete.

## Como uma escrita local vira sincronização

Com o schema declarado, o que acontece quando o app chama `Database.insert(UserSchema).values({...}).execute()`? Esse é o coração do offline-first, e vale seguir o caminho completo:

1. A chamada vira uma única SQL parametrizada, executada sincronamente no core C++ via JSI.
2. Uma trigger SQLite, não código de aplicação, dispara automaticamente e insere uma linha em `sync_queue`. INSERT e UPDATE guardam o `payload` inteiro da linha em JSON; DELETE é sempre um soft delete (`UPDATE ... SET deletedAt = ?`), nunca um `DELETE` de verdade. Todo `select`/`count` filtra `deletedAt IS NULL` automaticamente.
3. Isso vale até para SQL raw. A trigger é definida a nível de tabela, não no query builder, então `Database.execute('INSERT INTO ...')` também populua a `sync_queue`.
4. Dentro de uma `Database.transaction(fn)`, cada write ainda dispara sua trigger normalmente, mas a fila só é populada no `COMMIT`, não por write isolado.
5. Por fim, o mesmo write dispara `requestWriteSync`, um throttle leading-edge de 5 segundos por schema que descarta silenciosamente se já existe uma sessão de sync rodando. É esse gatilho que drena a fila quase instantaneamente quando o app está aberto e online.

**Um ponto importante a destacar é que** não existe "modo silencioso" de escrita. Todo write que passa pela camada de query, ou até SQL raw, dispara a trigger e entra na fila automaticamente. É impossível esquecer de sincronizar algo, porque sincronizar nunca foi uma chamada de código, sempre foi uma trigger de banco. A única exceção intencional é quando o próprio motor de sync está aplicando dados vindos do servidor, e aí existe um mecanismo específico para não entrar em loop, que eu explico já já.

## A sessão de sync: push, depois pull

Uma sessão de sincronização, disparada manualmente com `Database.sync('users')` / `Database.syncAll()` ou automaticamente pelos gatilhos que já vimos, roda sempre duas fases em ordem: push, depois pull. Push não é estritamente necessário antes do pull (repuxar a própria linha via pull é um no-op idempotente), mas é mais intuitivo e evita uma janela onde o pull traria de volta um estado que o push está prestes a substituir.

### Push: drenar a fila local

A fase de push lê a `sync_queue` em ordem FIFO, até um teto fixo de 200 itens por sessão. Cada item vira uma requisição: insert vira `POST`, update vira `PATCH`, delete vira `DELETE`. O tratamento de resposta é onde a maioria dos detalhes interessantes mora:

- **Sucesso (2xx)**: insert/update reescreve a linha local com o id retornado pelo servidor e marca a metadata como `SYNCED`; delete confirma o soft delete.
- **Falha HTTP comum (400/409/500)**: o item é marcado `FAILED`, o `retryCount` incrementa, e a fila segue para o próximo item.
- **404 num delete**: é tratado como idempotente, o servidor já deletou aquilo mesmo, então confirma local e remove da fila.
- **404 num insert/update**: é ambíguo, o alvo sumiu no servidor enquanto havia edição local pendente. O item é marcado `BLOCKED` e para de ser tentado automaticamente, porque resolução de conflito para esse caso ficou deliberadamente fora do escopo inicial.
- **Falha de rede**: aborta o resto do push imediatamente. Itens não processados continuam `PENDING` e serão tentados de novo na próxima sessão. **A fase de pull não roda se isso acontecer.**

Cada chamada HTTP individual ainda tem seu próprio orçamento de retry, fixo em 3 tentativas com 5 segundos de delay, hardcoded no engine e não configurável por schema. É um trade-off deliberado: um scheduler que roda sem supervisão em background precisa de um piso mínimo de resiliência que não dependa de configuração.

### Pull: buscar mudanças, independente da fila

A segunda fase pagina através de `GET <basePath>?<listQueryTemplate renderizado>`. Para cada linha da resposta: se `deletedAt` vier preenchido, aplica um tombstone local; se a linha já existe localmente, resolve o conflito (por padrão via `lastWriteWins`, comparando `updatedAt`); senão, insere.

O cursor avança para o timestamp da última linha da página, mas persiste `últimoTimestamp - 1`, não o valor exato, para não perder linhas empatadas no mesmo milissegundo num limite de página. Esse overlap de 1ms é resolvido de novo pelo mesmo mecanismo de `lastWriteWins`.

**Aqui está o ponto que eu acho mais interessante do design inteiro**: o pull não depende em nada do estado da `sync_queue`. Ele é guiado só pelo cursor persistido. Isso significa que o wake periódico em background continua tendo trabalho real mesmo quando a fila local está completamente vazia, porque ele está checando o que *outros* clientes mudaram no servidor, não só drenando as próprias escritas. Se você pensar em sync só como "esvaziar uma fila", está faltando metade do quadro.

### O mecanismo anti-loop

Tem um problema óbvio escondido nisso tudo: a trigger de mudança dispara em qualquer write na tabela, inclusive quando o próprio motor de sync está aplicando dados que acabou de baixar no pull. Sem tratamento, isso reenfileiraria dado que acabou de chegar do servidor, e o sistema entraria num ping-pong infinito.

A solução é uma tabela de uma linha, `_sync_apply_lock`. Cada trigger tem uma cláusula `WHEN NOT EXISTS (SELECT 1 FROM _sync_apply_lock)`. Quando o motor aplica operações vindas do servidor, ele faz `BEGIN; INSERT INTO _sync_apply_lock; <aplica os dados>; DELETE FROM _sync_apply_lock; COMMIT`, tudo numa única transação. Um crash no meio desfaz o apply e o lock junto, então o sistema nunca fica "travado ligado".

## Os guard-rails de uma execução síncrona

Como a execução de query em foreground é síncrona na JS thread, uma consequência direta do JSI que vimos antes, um resultado sem limite bloquearia a UI de forma imprevisível. Por isso existem dois guard-rails obrigatórios:

1. **`.limit()` com teto.** Se você omitir, o padrão é 500. Se informar, não pode passar de 500, senão `.execute()` lança erro. Não se aplica a UPDATE/DELETE, onde "atualizar/deletar tudo que casa" é o comportamento esperado.
2. **Regra da coluna indexada.** Toda coluna usada em `.where()`/`.orderBy()` precisa ser a `primaryKey` ou a coluna líder de um índice declarado, seguindo o mesmo leftmost-prefix que o SQLite já usa para índices compostos. Um índice `columns: ['a', 'b']` cobre filtrar por `a`, não por `b` isoladamente.

Sem o segundo guard-rail, o `.limit()` limitaria o tamanho do resultado, mas não o custo do scan por trás dele. A mensagem de erro é direta: `Synchronous execute() requires an index covering column "X" as its leading column`. Ela aponta exatamente o que declarar em `schema.indexes` para o erro sumir.

## Dois relógios diferentes: iOS vs Android

O agendamento em background é onde a divergência entre plataformas fica mais concreta. Existe um único job nativo global por banco, não um por schema, e toda tabela com `sync.enabled` é sincronizada em cada wake. Mas o que controla *quando* esse wake acontece é bem diferente nas duas plataformas:

| | Android (`WorkManager`) | iOS (`BGTaskScheduler`) |
|---|---|---|
| `minimumInterval` | piso rígido de 15 minutos, imposto pelo próprio sistema operacional. Um valor menor é aceito, mas nunca dispara mais rápido do que isso | tratado como uma dica de `earliestBeginDate`. O sistema decide a hora real por heurística de bateria e padrão de uso, não existe piso fixo a aplicar |
| Configuração extra | nenhuma. A lib já declara `INTERNET`/`ACCESS_NETWORK_STATE` e registra o job sozinha | precisa de `Info.plist` (`BGTaskSchedulerPermittedIdentifiers`, `UIBackgroundModes`) e `*.entitlements` (`keychain-access-groups`) |
| Falha silenciosa | não se aplica | sem essas entradas, o app builda e roda normal. O scheduler simplesmente nunca dispara e o credential provider não persiste token, sem nenhum erro em runtime apontando para isso |

**Um ponto importante a destacar é que** essa última linha da tabela é o tipo de detalhe que só aparece na prática. No Android o pior cenário é um erro de build. No iOS, sem os entitlements certos, tudo parece funcionar e silenciosamente não funciona nada.

### Cold start: o caso mais forte para "sem JS de verdade"

O teste de estresse da tese inteira é o que acontece quando o sistema operacional mata o app inteiro e depois acorda o job de sync num processo novo. Nesse cenário não existe app rodando, não existe JS bundle carregado, nada.

`Database.configure()` espelha o que o motor de sync precisa (nome do banco, credenciais, contrato de sync) num arquivo JSON ao lado do próprio SQLite, e não numa tabela, porque o caminho do banco nem é conhecido até esse arquivo ser lido. Quando o job acorda, o core nativo reidrata esse estado a partir do arquivo antes de tentar qualquer coisa, e roda a sessão de sync inteira sem nenhum JS envolvido.

Isso também explica por que o wake periódico segue importando mesmo quando a fila local está quase sempre vazia. Existem quatro gatilhos independentes para uma sessão de sync: pós-escrita em foreground, abertura do app, o monitor de conectividade nativo (offline para online), e o job periódico do sistema. Os três primeiros são best-effort, descartam silenciosamente se já existe sessão rodando. O job periódico é o único caminho garantido em cenários como uma escrita feita totalmente offline ou o app morto antes de qualquer outro gatilho disparar.

## Autenticação também é nativa

Se tem um lugar onde eu esperaria ver JavaScript entrar de volta na história, seria autenticação: interceptors de axios, refresh de token, esse tipo de coisa costuma morar em JS em praticamente todo app. No Salve DB, não.

O par inicial de tokens, obtido pelo fluxo de login do próprio app (fora do escopo da lib), é passado uma vez em `Database.configure({ credentials: { tokens } })` e escrito no Keychain (iOS) ou Keystore (Android) pelo `CredentialProvider` nativo. A partir daí, o JS nunca mais lê o token de volta, não existe API pública para recuperar o access ou refresh token atual.

Quando uma requisição de sync recebe `401`, o motor nativo (nunca JS) chama o endpoint de refresh configurado com o refresh token guardado, extrai `accessToken`/`refreshToken` da resposta via um `JsonPath` que você configura, reescreve no Keychain/Keystore e repete a chamada original. Isso acontece igual numa sessão disparada de JS, no `syncOnAppOpen`, ou numa wake em background onde a JS runtime nunca chegou a subir. JS não tem hook nesse fluxo, e nenhuma forma de interceptar, atrasar ou observar um refresh individual.

## Reatividade: como a UI descobre que algo mudou

`Database.subscribeToChanges` expõe notificações de escrita a nível de tabela, vindas de qualquer origem: query builder, SQL raw, migração, ou o próprio motor de sync, seja em foreground ou background. `useQuery`/`useInfiniteQuery` consomem isso por baixo dos panos:

```tsx
function UserList() {
  const { data, isLoading, error } = useQuery({
    schema: UserSchema,
    queryFn: (q) => q.where(eq('name', search)).orderBy('updatedAt', 'desc').limit(50),
    deps: [search],
  });
  // re-renderiza automaticamente em qualquer escrita em `users`, de qualquer origem
}
```

A UI nunca sabe a diferença entre dado local e dado sincronizado. Não importa se a escrita veio da tela, do Studio (que eu mostro daqui a pouco), ou de um pull de sync rodando em background: é reatividade sobre o SQLite como fonte única de verdade, não sobre um cache JS paralelo tentando ficar sincronizado com ele.

Existe também um mecanismo mais sutil: montar um `useQuery` para um schema com `sync.enabled` dispara um sync leve de leitura, throttled, silencioso, que nunca bloqueia a leitura em si. A ideia por trás disso é simples de defender: a leitura já é o melhor sinal de "isso importa agora" que existe, é literalmente o usuário olhando para aquele dado nesse instante.

## O que isso custa: trade-offs conscientes

Nenhuma decisão de arquitetura vem de graça, e seria desonesto apresentar o Salve DB como se todas essas escolhas fossem só vantagem. Algumas que eu assumi conscientemente:

- **Migrações são só `ADD COLUMN`.** No boot, `Database.register()` compara a versão declarada com a última persistida e aplica `ALTER TABLE ADD COLUMN` para colunas novas. Não existe `DROP`/`RENAME`: uma coluna removida do schema fica órfã no SQLite, silenciosamente ignorada. `DROP COLUMN` e `RENAME` são operações destrutivas e caras no SQLite, então ficaram fora do escopo por design, não por esquecimento. Na prática, isso significa que evolução de schema é intencionalmente aditiva e de mão única: para remover ou renomear algo, você adiciona uma coluna nova e migra os dados em código de aplicação.
- **Retry fixo, não configurável.** 3 tentativas, 5 segundos de delay, hardcoded. Um scheduler que roda sem supervisão precisa de um piso mínimo de resiliência que não dependa de alguém lembrar de configurar.
- **Resolução de conflito sem UI de merge.** `lastWriteWins`, `serverWins` e `clientWins` são determinísticos e implementados; `manual` está reservado, tipado, mas não implementado. Não existe fila de conflito para o usuário resolver na mão.
- **Só REST, só bidirecional.** GraphQL, gRPC, sync `incremental`/`full` isolado, tudo isso está fora do escopo atual, reduzindo a superfície do cliente HTTP nativo.
- **`currentUser()` é conveniência, não segurança.** É um valor em memória, não um filtro de linha automático. `userId` é uma coluna comum do seu schema, e cabe a você adicionar o `where()`/`values()` certo em toda query que precisa dele. Esquecer ainda lê e escreve através de todos os usuários.

Nenhum desses pontos é acidente. São escolhas de escopo documentadas e revisitáveis, não bugs escondidos debaixo do tapete.

## Ver funcionando: o Studio

Depois de tanta teoria, o jeito mais direto de fechar essa tese é mostrar ela funcionando. O Salve DB tem um Studio companheiro, uma UI local ao estilo Prisma/Drizzle Studio, conectada por WebSocket na porta `7377` direto no SQLite real do dispositivo, não numa cópia.

Quando o app chama `Database.configure()` em modo de desenvolvimento, ele conecta sozinho, sem configuração adicional. Múltiplos dispositivos rodando viram múltiplas entradas no seletor. De lá dá para navegar toda tabela, incluindo as internas como `sync_queue` e os cursores de sync, inserir, editar e deletar linhas (disparando as mesmas triggers de uma escrita normal do app), rodar SQL raw, truncar ou dropar tabelas da própria aplicação.

É a ferramenta que torna a tese "sincroniza sem JS, sem depender do app aberto" verificável a olho nu: dá para ver a fila de sync esvaziando em tempo real conforme o motor nativo trabalha sozinho, sem nenhum código React por perto.

<figure>
  <img src="/static/images/salve-db-studio.gif" alt="Studio conectado ao SQLite real do app, mostrando a fila de sync esvaziando em tempo real" loading="lazy">
  <figcaption>Studio conectado ao SQLite real do app: escreva, dispare um sync e veja a fila esvaziar em tempo real.</figcaption>
</figure>

## Conclusão

O Salve DB ainda é uma biblioteca jovem, na versão 1.1.1, mas o recorte de escopo foi deliberado desde o primeiro commit: sync só bidirecional, só REST, resolução de conflito sem UI de merge, migração só de coluna nova. Cada um desses "nãos" tem um motivo documentado ao lado dele. Não é ausência por falta de tempo, é ausência por decisão. `incremental`/`full` como estratégia isolada, conflito `manual`, transporte GraphQL/gRPC, retry configurável por schema, compressão, criptografia e relations no query builder já estão reservados e tipados no contrato, esperando a próxima fase.

O que já está de pé resolve o problema que me fez começar a construir isso: sincronizar sem torcer para o sistema operacional acordar uma JS thread na hora certa. Da trigger que captura a escrita até o refresh de token OAuth2, passando pelo scheduler nativo em cada plataforma, o caminho inteiro roda sem depender do app estar aberto. E o Studio que você viu ali em cima é a prova disso rodando ao vivo, não só no papel.

O Salve DB é open source sob MIT, mantido pela Salve Software, publicado no npm como `@salve-software/react-native-salve-db`. Se você já bateu de frente com esse mesmo problema de sync em background, ou quiser abrir uma issue apontando onde o design ainda falha, é por ali que a conversa continua. Obrigado por ler até aqui 🙏, e como sempre: teste o seu código :)
