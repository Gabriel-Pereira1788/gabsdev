---
title: "Salve DB: Offline-First SQLite for React Native"
date: "2026-09-04"
tags: ["react-native", "sqlite", "offline-first", "mobile", "architecture"]
excerpt: "Most \"offline-first\" libraries still depend on JavaScript running in the background to sync. What if the entire sync engine, from the queue down to OAuth2 token refresh, ran 100% in native code?"
readMin: 14
---

Every time an app needs to work offline, the solution tends to be the same: write locally, and whenever possible, fire off a sync that runs in JavaScript in the background. It works, until the operating system decides that JS process doesn't deserve to stay alive. At that point sync simply stops, and nobody notices until the user complains their data went missing.

That problem is exactly why I built [Salve DB](https://github.com/Salve-Software/react-native-salve-db), an offline-first SQLite library for React Native. In this article I'll show you how it works under the hood: why the entire sync engine runs in native code (C++/Swift/Kotlin) and never needs to boot the JS engine to sync, how a local write turns into a real sync session, and which trade-offs I made on purpose to get there.

## The problem with JavaScript-driven "offline-first"

Before we talk about the solution, let's understand what really happens when a background sync depends on JavaScript. On iOS, that means gluing `BGTaskScheduler` callbacks to a JS context that might not even be alive at that moment. On Android, a `WorkManager` job has to boot an entire JS runtime before it can make a single HTTP request.

That costs time and battery just to *start* the engine, and it's still at the mercy of how the OS schedules the JS thread under background execution limits. It's a whole class of reliability bug that most sync libraries simply accept as the cost of doing business.

**One important point to highlight is that** a lot of libraries call themselves "offline-first" just for writing to local disk before hitting the network. That's not the same thing as syncing for real without depending on the app being open. In Salve DB, change capture is automatic and infallible: it's not a function call someone can forget to make, it's a database trigger. And syncing itself never depends on the app being open: a periodic OS job wakes the native engine on its own.

## The turning point: moving the entire sync engine to native

Salve DB's central thesis is easy to state: TypeScript only declares *what* to sync and *how* to sync, and never participates in execution. Orchestration, HTTP, OAuth2 credentials, background scheduling, all of it runs in C++/Swift/Kotlin.

That shows up in a four-layer architecture:

<figure>
  <img src="/static/images/salve-db-architecture.png" alt="Diagram of Salve DB's layered architecture: TypeScript talks to the C++ Native Core over JSI, which in turn uses platform services (BGTaskScheduler/WorkManager, NWPathMonitor/ConnectivityManager, Keychain/Keystore, Swift/Kotlin) to reach the REST API" loading="lazy">
  <figcaption>TypeScript only declares; JSI is the synchronous bridge; the C++ Native Core holds most of the logic; Swift/Kotlin cover only what C++ can't reach.</figcaption>
</figure>

Let's go through each layer, because each one exists for a specific reason and isn't just organization for organization's sake.

**TypeScript** is pure declaration. Schemas are plain data, the query builder assembles SQL and parameters but never executes anything, and the hooks/provider are the DX surface. The project's "no JS" rule is about the sync engine running unsupervised in the background, not about the builder in the foreground: when you call `.execute()` with the app open, JS is very much alive.

**JSI (via Nitro Modules)** is the bridge. A zero-copy, synchronous bridge through a `HybridObject`. It's *because* of that synchronicity that foreground query execution is synchronous on the JS thread, which in turn explains the guard-rails I'll get to further down.

**Native Core (C++)** is where most of the logic lives: the SQLite executor with an LRU cache of prepared statements, the migration engine, the trigger engine, the sync orchestrator, the credential provider, and the HTTP client. All of it platform-independent.

**Swift/Kotlin** are the edges, and they exist only for what C++ genuinely can't reach: `BGTaskScheduler`/`WorkManager` are Swift/Kotlin-only APIs, and bridging their lifecycle callbacks through JNI/Objective-C runtime would be fragile and error-prone. Keychain and Keystore have no C++ bindings either. And bundling `libcurl` would mean manually reimplementing proxy handling, cert pinning, and TLS configuration that `URLSession`/`OkHttp` already get right.

So why not write everything in Swift/Kotlin instead? Because the core logic (sync queue, expression interpreter, conflict resolution) has zero platform dependency, and writing it twice guarantees the two implementations will eventually drift apart.

Now that the layers are clear, let's look at the example that carries the rest of this article: a synced `users` table.

## One schema, start to finish

Everything else in this article revolves around a single schema. Here's how you declare a syncable `users` table in Salve DB:

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

Notice there's no SQL anywhere here: columns, indexes, and the sync contract are just TypeScript data, interpreted by the native core at boot. `sync` is optional, local-only tables that never need the network simply skip that block entirely. And there's a detail easy to miss on first read: `deletedAt` gets injected into every table automatically, even without you declaring it. That matters because delete, in Salve DB, is never a real SQL `DELETE`, it's always a soft delete.

## How a local write becomes a sync

With the schema declared, what happens when the app calls `Database.insert(UserSchema).values({...}).execute()`? This is the heart of offline-first, and it's worth following the full path:

1. The call becomes a single parameterized SQL statement, executed synchronously in the C++ core via JSI.
2. A SQLite trigger, not application code, fires automatically and inserts a row into `sync_queue`. INSERT and UPDATE store the row's full `payload` as JSON; DELETE is always a soft delete (`UPDATE ... SET deletedAt = ?`), never a real `DELETE`. Every `select`/`count` automatically filters `deletedAt IS NULL`.
3. This holds even for raw SQL. The trigger is defined at the table level, not in the query builder, so `Database.execute('INSERT INTO ...')` populates `sync_queue` too.
4. Inside a `Database.transaction(fn)`, each write still fires its trigger normally, but the queue is only populated on `COMMIT`, not per isolated write.
5. Finally, that same write fires `requestWriteSync`, a 5-second leading-edge throttle per schema that silently discards if a sync session is already running. That's what drains the queue almost instantly when the app is open and online.

**One important point to highlight is that** there's no "silent write mode". Every write that goes through the query layer, or even raw SQL, fires the trigger and enters the queue automatically. It's impossible to forget to sync something, because syncing was never a function call, it was always a database trigger. The one intentional exception is when the sync engine itself is applying data coming from the server, and there's a specific mechanism to avoid looping in that case, which I'll get to next.

## The sync session: push, then pull

A sync session, triggered manually with `Database.sync('users')` / `Database.syncAll()` or automatically by the triggers we just saw, always runs two phases in order: push, then pull. Push isn't strictly necessary before pull (pulling your own row back is an idempotent no-op), but it's more intuitive and avoids a window where pull would bring back a state push is about to replace.

### Push: draining the local queue

The push phase reads `sync_queue` in FIFO order, up to a fixed cap of 200 items per session. Each item becomes a request: insert becomes `POST`, update becomes `PATCH`, delete becomes `DELETE`. The response handling is where most of the interesting detail lives:

- **Success (2xx)**: insert/update rewrites the local row with the id the server returned and marks the metadata `SYNCED`; delete confirms the soft delete.
- **Common HTTP failure (400/409/500)**: the item is marked `FAILED`, `retryCount` increments, and the queue moves to the next item.
- **404 on a delete**: treated as idempotent, the server already deleted that thing too, so it confirms locally and drops the item from the queue.
- **404 on an insert/update**: ambiguous, the target vanished on the server while a local edit was still pending. The item gets marked `BLOCKED` and stops being retried automatically, since conflict resolution for that case was deliberately kept out of the initial scope.
- **Network failure**: aborts the rest of push immediately. Unprocessed items stay `PENDING` and get retried in the next session. **The pull phase does not run if this happens.**

Each individual HTTP call still has its own retry budget, fixed at 3 attempts with a 5-second delay, hardcoded in the engine and not configurable per schema. That's a deliberate trade-off: a scheduler running unsupervised in the background needs a minimum floor of resilience that doesn't depend on anyone remembering to configure it.

### Pull: fetching changes, independent of the queue

The second phase pages through `GET <basePath>?<rendered listQueryTemplate>`. For each row in the response: if `deletedAt` is set, it applies a local tombstone; if the row already exists locally, it resolves the conflict (by default via `lastWriteWins`, comparing `updatedAt`); otherwise, it inserts.

The cursor advances to the last page row's timestamp, but persists `lastTimestamp - 1`, not the exact value, so it doesn't drop rows tied at the same millisecond right at a page boundary. That 1ms overlap gets resolved again by the same `lastWriteWins` mechanism.

**Here's the point I find most interesting about the whole design**: pull doesn't depend on the state of `sync_queue` at all. It's driven purely by the persisted cursor. That means the periodic background wake still has real work to do even when the local queue is completely empty, because it's checking what *other* clients changed on the server, not just draining your own writes. If you think of sync as just "empty a queue", you're missing half the picture.

### The anti-loop mechanism

There's an obvious problem hiding in all of this: the change trigger fires on any write to the table, including when the sync engine itself is applying data it just pulled down. Without handling, that would re-queue data that just arrived from the server, and the system would ping-pong forever.

The fix is a single-row table, `_sync_apply_lock`. Every trigger has a `WHEN NOT EXISTS (SELECT 1 FROM _sync_apply_lock)` clause. When the engine applies operations coming from the server, it does `BEGIN; INSERT INTO _sync_apply_lock; <apply the data>; DELETE FROM _sync_apply_lock; COMMIT`, all in a single transaction. A crash in the middle undoes the apply and the lock together, so the system never ends up "stuck on".

## Guard-rails for synchronous execution

Since foreground query execution is synchronous on the JS thread, a direct consequence of JSI we saw earlier, an unbounded result would block the UI unpredictably. That's why there are two mandatory guard-rails:

1. **`.limit()` with a cap.** Omit it and the default is 500. Provide one and it can't exceed 500, or `.execute()` throws. This doesn't apply to UPDATE/DELETE, where "update/delete everything that matches" is the expected behavior.
2. **The indexed-column rule.** Every column used in `.where()`/`.orderBy()` must be the `primaryKey` or the leading column of a declared index, following the same leftmost-prefix rule SQLite already uses for composite indexes. An index on `columns: ['a', 'b']` covers filtering by `a`, not by `b` alone.

Without the second guard-rail, `.limit()` would cap the size of the result, but not the cost of the scan behind it. The error message is direct: `Synchronous execute() requires an index covering column "X" as its leading column`. It points at exactly what to declare in `schema.indexes` to make the error go away.

## Two different clocks: iOS vs Android

Background scheduling is where the divergence between platforms gets concrete. There's a single global native job per database, not one per schema, and every table with `sync.enabled` gets synced on each wake. But what controls *when* that wake happens is quite different between the two:

| | Android (`WorkManager`) | iOS (`BGTaskScheduler`) |
|---|---|---|
| `minimumInterval` | hard floor of 15 minutes, imposed by the OS itself. A smaller value is accepted, but never fires faster than that | treated as a hint via `earliestBeginDate`. The system decides the actual time by battery and usage heuristics, there's no fixed floor to apply |
| Extra setup | none. The library already declares `INTERNET`/`ACCESS_NETWORK_STATE` and registers the job on its own | needs `Info.plist` (`BGTaskSchedulerPermittedIdentifiers`, `UIBackgroundModes`) and `*.entitlements` (`keychain-access-groups`) |
| Silent failure | not applicable | without those entries, the app builds and runs fine. The scheduler just never fires and the credential provider can't persist tokens, with no runtime error pointing you at it |

**One important point to highlight is that** the last row of that table is the kind of detail that only shows up in practice. On Android the worst case is a build error. On iOS, without the right entitlements, everything looks like it works and silently nothing does.

### Cold start: the strongest case for "no JS, for real"

The stress test for the whole thesis is what happens when the OS kills the entire app and later wakes the sync job in a fresh process. In that scenario there's no app running, no JS bundle loaded, nothing.

`Database.configure()` mirrors what the sync engine needs (database name, credentials, sync contract) into a JSON file sitting next to the SQLite file itself, not a table, because the database path isn't even known until that file gets read. When the job wakes up, the native core rehydrates that state from the file before attempting anything, and runs the entire sync session with no JS involved.

That also explains why the periodic wake keeps mattering even when the local queue is almost always empty. There are four independent triggers for a sync session: post-write in the foreground, app open, the native connectivity monitor (offline to online), and the periodic OS job. The first three are best-effort, silently discarding if a session is already running. The periodic job is the only guaranteed path in scenarios like a write made entirely offline, or the app being dead before any other trigger fires.

## Authentication is native too

If there's one place I'd expect to see JavaScript sneak back into the story, it's authentication: axios interceptors, token refresh, that sort of thing usually lives in JS in practically every app. In Salve DB, it doesn't.

The initial token pair, obtained through the app's own login flow (out of the library's scope), is passed once via `Database.configure({ credentials: { tokens } })` and written to the Keychain (iOS) or Keystore (Android) by the native `CredentialProvider`. From then on, JS never reads the token back, there's no public API to retrieve the current access or refresh token.

When a sync request comes back with `401`, the native engine (never JS) calls the configured refresh endpoint with the stored refresh token, extracts `accessToken`/`refreshToken` from the response via a `JsonPath` you configure, rewrites the Keychain/Keystore, and repeats the original call. That happens the same way whether the session was triggered from JS, from `syncOnAppOpen`, or from a background wake where the JS runtime never even booted. JS has no hook into this flow, and no way to intercept, delay, or observe an individual refresh.

## Reactivity: how the UI finds out something changed

`Database.subscribeToChanges` exposes table-level write notifications from any origin: query builder, raw SQL, migration, or the sync engine itself, foreground or background. `useQuery`/`useInfiniteQuery` consume that under the hood:

```tsx
function UserList() {
  const { data, isLoading, error } = useQuery({
    schema: UserSchema,
    queryFn: (q) => q.where(eq('name', search)).orderBy('updatedAt', 'desc').limit(50),
    deps: [search],
  });
  // re-renders automatically on any write to `users`, from any source
}
```

The UI never knows the difference between local data and synced data. It doesn't matter whether the write came from the screen, from the Studio (which I'll show in a moment), or from a background sync pull: it's reactivity over SQLite as the single source of truth, not over a parallel JS cache trying to stay in sync with it.

There's also a subtler mechanism at play: mounting a `useQuery` for a schema with `sync.enabled` fires a lightweight, throttled read sync that never blocks the read itself. The idea behind it is easy to defend: a read is already the best signal that "this matters right now" that exists, it's literally the user looking at that data in this exact instant.

## What it costs: conscious trade-offs

No architecture decision comes for free, and it would be dishonest to present Salve DB as if every one of these choices were pure upside. A few I made on purpose:

- **Migrations are `ADD COLUMN` only.** On boot, `Database.register()` compares the declared version against the last persisted one and applies `ALTER TABLE ADD COLUMN` for new columns. There's no `DROP`/`RENAME`: a column removed from the schema stays orphaned in SQLite, silently ignored. `DROP COLUMN` and `RENAME` are destructive, expensive SQLite operations, so they stayed out of scope by design, not by oversight. In practice this means schema evolution is intentionally additive and one-way: to remove or rename something, you add a new column and migrate the data in application code.
- **Fixed retry, not configurable.** 3 attempts, 5-second delay, hardcoded. A scheduler running unsupervised needs a minimum floor of resilience that doesn't depend on someone remembering to configure it.
- **Conflict resolution with no merge UI.** `lastWriteWins`, `serverWins`, and `clientWins` are deterministic and implemented; `manual` is reserved and typed, but not implemented. There's no conflict queue for a user to resolve by hand.
- **REST only, bidirectional only.** GraphQL, gRPC, an isolated `incremental`/`full` sync strategy, all of that is out of scope for now, keeping the native HTTP client's surface small.
- **`currentUser()` is convenience, not security.** It's an in-memory value, not an automatic row filter. `userId` is an ordinary column of your own schema, and it's on you to add the right `where()`/`values()` to every query that needs it. Forgetting still reads and writes across every user's rows.

None of these are accidents. They're documented, revisitable scope decisions, not bugs swept under the rug.

## Seeing it work: the Studio

After all this theory, the most direct way to close this thesis out is to show it working. Salve DB ships with a companion Studio, a local UI in the style of Prisma/Drizzle Studio, connected over WebSocket on port `7377` directly to the device's real SQLite, not a copy.

When the app calls `Database.configure()` in development mode, it connects on its own, no extra setup needed. Multiple running devices show up as separate entries in the selector. From there you can browse every table, including internal ones like `sync_queue` and the sync cursors, insert, edit, and delete rows (firing the same triggers a normal app write would), run raw SQL, and truncate or drop tables belonging to your own app.

It's the tool that makes the "syncs without JS, without the app being open" thesis verifiable with your own eyes: you can watch the sync queue drain in real time as the native engine works on its own, with no React code anywhere in sight.

<figure>
  <img src="/static/images/salve-db-studio.gif" alt="Studio connected to the app's real SQLite, showing the sync queue draining in real time" loading="lazy">
  <figcaption>Studio connected to the app's real SQLite: write, trigger a sync, and watch the queue drain in real time.</figcaption>
</figure>

## Wrapping up

Salve DB is still a young library, sitting at version 1.1.1, but the scope cuts were deliberate from the first commit: sync only bidirectional, only REST, conflict resolution with no merge UI, migrations that only ever add a column. Every one of those "nos" has a documented reason attached to it. It's not missing because there wasn't time, it's missing because it was a decision. An isolated `incremental`/`full` strategy, `manual` conflict, GraphQL/gRPC transport, per-schema configurable retry, compression, encryption, and relations in the query builder are already reserved and typed in the contract, waiting for the next phase.

What's already standing solves the problem that made me start building this in the first place: syncing without crossing your fingers that the OS wakes a JS thread at the right time. From the trigger that captures a write to the OAuth2 token refresh, through the native scheduler on each platform, the whole path runs without depending on the app being open. And the Studio clip you just watched is proof of that running live, not just on paper.

Salve DB is open source under MIT, maintained by Salve Software, published on npm as `@salve-software/react-native-salve-db`. If you've already run into this exact background sync problem, or want to open an issue pointing at where the design still falls short, that's where the conversation continues. Thanks for reading this far 🙏, and as always: test your code :)
