# Guion de vídeo para presentar el repositorio

Duración objetivo: 3–5 minutos. Idioma principal: español, con frases breves en inglés para practicar la presentación técnica. En pantalla: cámara al inicio y al cierre; terminal y código durante la demo.

## Preparación antes de grabar

Desde la raíz del repositorio, prepara la base de datos. Go 1.27.2+ y Docker Compose deben estar disponibles:

```sh
cp .env.example .env && docker compose up -d db && make migrate
```

En una primera terminal, deja la API corriendo:

```sh
make run
```

Usa una segunda terminal para las llamadas HTTP. El `.env.example` incluye una clave solo para desarrollo.

## Guion de 3–5 minutos

### 0:00–0:25 — Gancho

**En cámara:**

> “¿Qué pasa si una petición crea el registro en la base de datos, pero el proceso cae antes de publicar el evento? Y, si el cliente reintenta, ¿cómo evitamos crear dos veces el mismo recurso? Este proyecto pequeño explora esas dos preguntas con Go y PostgreSQL.”

**Frase clave en inglés:** “This project makes retries safe and keeps business state and event records consistent.”

### 0:25–0:55 — Problema y alcance

**En cámara, mostrar README y diagrama Mermaid:**

> “El ejemplo expone una API para crear y consultar Jobs. El caso es deliberadamente acotado: no ejecuta trabajos de fondo. El Job se queda en estado `pending`. Lo interesante es cómo se guarda la creación y cómo se registra el evento que la acompaña.”

**Frase clave en inglés:** “The job is a small example bounded context, not a production job-processing engine.”

### 0:55–2:00 — Demo paso a paso

**En terminal, mostrar que la API ya está ejecutándose.** En la segunda terminal, carga la clave y crea el Job:

```sh
set -a; . ./.env; set +a
IDEMPOTENCY_KEY="portfolio-demo-$(date +%s)-$$"
RESPONSE=$(curl -sS -i -H "Authorization: Bearer $API_KEY" \
  -H "Idempotency-Key: $IDEMPOTENCY_KEY" -H 'Content-Type: application/json' \
  -d '{"name":"portfolio demo"}' http://localhost:8080/v1/jobs)
printf '%s\n' "$RESPONSE"
JOB_ID=$(printf '%s\n' "$RESPONSE" | sed -n 's/.*"id":"\([^"]*\)".*/\1/p')
```

**Narración:**

> “La primera llamada devuelve `201 Created` y un ID. Ahora repito exactamente la misma solicitud con la misma clave.”

```sh
curl -i -H "Authorization: Bearer $API_KEY" \
  -H "Idempotency-Key: $IDEMPOTENCY_KEY" -H 'Content-Type: application/json' \
  -d '{"name":"portfolio demo"}' http://localhost:8080/v1/jobs
```

> “La respuesta es `200 OK` y conserva el mismo ID. Si reutilizo la clave para otro nombre, la API responde `409 Conflict`.”

```sh
curl -i -H "Authorization: Bearer $API_KEY" \
  -H "Idempotency-Key: $IDEMPOTENCY_KEY" -H 'Content-Type: application/json' \
  -d '{"name":"different job"}' http://localhost:8080/v1/jobs
```

Consulta el recurso y las sondas:

```sh
curl -i -H "Authorization: Bearer $API_KEY" "http://localhost:8080/v1/jobs/$JOB_ID"
curl -i http://localhost:8080/live
curl -i http://localhost:8080/ready
```

> “La lectura devuelve el Job persistido. Liveness indica que el proceso responde; readiness también consulta PostgreSQL y comprueba la antigüedad de eventos pendientes.”

### 2:00–3:35 — Decisiones técnicas

**Mostrar `internal/adapters/postgres/jobs.go` y contar:**

1. **Una transacción para el Job y el evento.**

   > “El adaptador inserta el Job y su evento en la tabla outbox dentro de la misma transacción. Si el commit no ocurre, ninguno queda guardado. Así evito la ventana en la que queda el Job sin su evento.”

   **English:** “The job row and the outbox row commit together.”

2. **Idempotencia guardada por la base de datos.**

   > “La clave tiene una restricción única. Guardo un hash del nombre normalizado: repetir clave y contenido devuelve el resultado existente; cambiar el contenido con la misma clave da conflicto.”

   **English:** “Same key and same request returns the existing job; a different request conflicts.”

3. **Entrega al menos una vez, con deduplicación en el receptor.**

   **Mostrar `internal/adapters/events/worker.go`:**

   > “El worker reclama filas con `FOR UPDATE SKIP LOCKED`, publica por HTTP y marca la fila como entregada tras un 2xx. Si falla, la reintenta. Hay una ventana en la que el receptor puede aceptar el evento y el proceso caer antes de marcarlo; por eso el receptor debe deduplicar usando `X-Event-Id`. No afirmo entrega exactamente una vez.”

   **English:** “Delivery is at least once, so the receiver must deduplicate the stable event ID.”

### 3:35–4:15 — Tests y límites

**En terminal:**

```sh
make test
```

> “Los tests Go corren con el detector de carreras. También hay un gate completo con integración PostgreSQL y controles de dependencias. El resultado depende de las herramientas y la base local disponibles.”

**En cámara o README:**

> “Es una plantilla pequeña, no un sistema desplegado a escala. No trae autenticación por usuario, autorización de producto, broker, motor de ejecución de Jobs, métricas de rendimiento ni un SLO medido. Esas decisiones dependen del producto real.”

### 4:15–4:35 — Cierre

> “Lo que quería mostrar es cómo hago explícitos los límites de consistencia, los reintentos y los fallos, y cómo dejo esas decisiones verificables en código, tests y documentación. El siguiente paso sería reemplazar el contexto Job por un caso de negocio real y especificar sus requisitos.”

**English:** “I want the failure behavior to be visible in the design, tests, and operational documentation.”

## Versión de 60 segundos

> “Este repositorio es una pequeña API backend en Go para crear y consultar Jobs. La pregunta que explora es cómo tolerar reintentos del cliente sin duplicar la creación ni perder el evento asociado.
>
> La creación del Job y el registro `JobCreated` se guardan en una sola transacción de PostgreSQL. La clave de idempotencia devuelve el mismo Job cuando se repite la misma solicitud y responde conflicto si cambia el contenido. Un worker publica la outbox con semántica de al menos una vez y reintenta errores, así que el receptor debe deduplicar por ID estable.
>
> En la demo se ve una creación, su reintento con el mismo ID, el conflicto por reutilizar la clave y las sondas de salud. Es una plantilla enfocada: el Job permanece `pending`; no incluye motor de ejecución, autenticación de producto ni métricas de escala. Quise dejar visibles las decisiones y sus límites en el código y los tests.”

**English phrases to retain:** “The job and outbox event commit together.” / “Delivery is at least once.” / “The receiver deduplicates by event ID.”
