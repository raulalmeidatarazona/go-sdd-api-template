# ADR 0001: Transactional outbox

Status: accepted

Business writes and event records use one PostgreSQL transaction. The worker publishes later and acknowledges only after a successful receiver response. This avoids lost events on process failure. Delivery can repeat after a crash between publish and acknowledgement, so consumers must deduplicate the stable event ID. We accept a single database and modest polling before introducing a broker or separate read store. Revisit only with measured throughput or isolation needs.
